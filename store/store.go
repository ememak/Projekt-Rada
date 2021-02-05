// Package for database management.
//
// All structures are stored similar to how they are defined in query.proto file.
// Main database buckets are containing RSA keys and polls data.
// Here is more detailed sheme:
//
//   KeysBucket is storing keys for polls.
//   Each key is stored in pair (keyid, key), where id is number of poll
//   and key is PKCS1 encoding of key.
//   * KeysBucket
//     - (keyid, key)
//
//   PollsBucket is storing data of polls: schema, tokens and votes.
//   * PollsBucket
//
//     Each poll is stored in separate bucket.
//     Id in name is its number (starting with 1!).
//     - PollidBucket
//
//       Schema structure stores poll questions.
//       It is stored in database encoded using proto.Marshal function.
//       + ("Schema", struct)
//
//       TokensBucket is storing tokens to poll.
//       Each is stored as its value as key and bool value specifying if it was used.
//       + TokensBucket
//         - (token, used)
//
//       VotesBucket stores votes for poll.
//       + VotesBucket
//
//         Each vote is stored in a bucket named after ballot used to signing it.
//         Ballot and sign are first and second value of RSASignature used in voting.
//         Answer is a PollSchema containing questions and answers encoded using
//         proto.Marshal function.
//         - Ballot
//           + ("Sign", sign)
//           + ("Answer", structure)
//
// Each number value is stored using strconv.Itoa function.
package store

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strconv"

	"github.com/ememak/Projekt-Rada/query"
	"github.com/golang/protobuf/proto"
	bolt "go.etcd.io/bbolt"
)

// DBInit Opens database and create buckets for data.
//
// This function have to be called before any other database related function.
func DBInit(filename string) (*bolt.DB, error) {
	if filename == "" || !query.IsStringPrintable(filename) {
		return nil, fmt.Errorf("Database name invalid\n")
	}
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("KeyBucket"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("PollsBucket"))
		return err
	})
	return db, err
}

// GetKey reads key for specific poll from database.
//
// Key should be stored in bucket KeyBucket with label keyid, where id is number of poll.
// If keyid is not in database, nil is returned.
// Key is stored in PKCS1 format.
func GetKey(db *bolt.DB, pollid int32) (*rsa.PrivateKey, error) {
	var bkeycpy []byte
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		kbuck := tx.Bucket([]byte("KeyBucket"))
		bkey := kbuck.Get([]byte("key" + strconv.Itoa(int(pollid))))
		if bkey == nil {
			return fmt.Errorf("No key for this poll in database.")
		}

		bkeycpy = make([]byte, len(bkey))
		copy(bkeycpy, bkey)
		return nil
	})

	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(bkeycpy)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert key from binary: %w", err)
	}
	return key, nil
}

// SaveKey saves poll key to database.
//
// Key is saved in bucket KeyBucket with label keyid, where id is number of poll.
// Key is stored in PKCS1 format.
func SaveKey(db *bolt.DB, pollid int, key *rsa.PrivateKey) error {
	if key == nil {
		return fmt.Errorf("Error! Private key is nil!")
	}

	if err := key.Validate(); err != nil {
		return err
	}
	bkey := x509.MarshalPKCS1PrivateKey(key)
	return db.Update(func(tx *bolt.Tx) error {
		keybuck := tx.Bucket([]byte("KeyBucket"))
		return keybuck.Put([]byte("key"+strconv.Itoa(pollid)), bkey)
	})
}

// NewPoll creates bucket for new poll.
//
// Return values is an id of poll in database and error returned by database.
func NewPoll(db *bolt.DB, sch *query.PollSchema) (*query.PollQuestion, error) {
	poll := &query.PollQuestion{
		Schema: sch,
	}
	err := db.Update(func(tx *bolt.Tx) error {
		// All polls are stored in PollsBucket.
		pollsbuck := tx.Bucket([]byte("PollsBucket"))
		id, _ := pollsbuck.NextSequence()
		poll.Id = int32(id)

		// Each poll is contained in bucket named by its number.
		pbuck, err := pollsbuck.CreateBucketIfNotExists([]byte("Poll" + strconv.Itoa(int(id)) + "Bucket"))
		if err != nil {
			return err
		}

		// Inside of a poll bucket there are two buckets and one value:
		// Value for schema, buckets for tokens and votes.

		// Check if Schema is valid.
		if err = sch.IsValid(); err != nil {
			return err
		}
		binschema, err := proto.Marshal(sch)
		if err != nil {
			return err
		}

		err = pbuck.Put([]byte("Schema"), binschema)
		if err != nil {
			return err
		}

		_, err = pbuck.CreateBucketIfNotExists([]byte("TokensBucket"))
		if err != nil {
			return err
		}

		_, err = pbuck.CreateBucketIfNotExists([]byte("VotesBucket"))
		return err
	})
	if err != nil {
		return &query.PollQuestion{}, err
	}
	return poll, err
}

// GetPoll reads poll from database.
func GetPoll(db *bolt.DB, pollid int32) (query.PollQuestion, error) {
	q := query.PollQuestion{
		Id:     pollid,
		Schema: &query.PollSchema{},
	}
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		pbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(pollid)) + "Bucket"))
		if pbuck == nil {
			return fmt.Errorf("Poll ID does not exist in database. GetPoll: %v", pollid)
		}

		// Read Schema stored as bytes converted via proto.Marchal.
		binschema := pbuck.Get([]byte("Schema"))
		err := proto.Unmarshal(binschema, q.Schema)
		if err != nil {
			return fmt.Errorf("Failed to read schema from database in GetPoll: %w", err)
		}

		// Tokens are stored as pairs (token, bool), where bool is true if token was not used.
		// If used == false, token is not read from database.
		tbuck := pbuck.Bucket([]byte("TokensBucket"))
		c := tbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := query.VoteToken{}
			if !bytes.Equal(v, []byte{0}) {
				t.Token = k
			}
			q.Tokens = append(q.Tokens, &t)
		}

		// Votes are stored in VotesBucket.
		// Each vote is a different bucket inside VotesBucket, with name Vote+nr.
		vbuck := pbuck.Bucket([]byte("VotesBucket"))

		c = vbuck.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ansbuck := vbuck.Bucket(k)
			pa := query.PollAnswer{
				Answers: &query.PollSchema{},
			}
			sign := ansbuck.Get([]byte("Sign"))
			pa.Sign = &query.RSASignature{
				Ballot: k,
				Sign:   sign,
			}

			binans := ansbuck.Get([]byte("Answer"))
			// Read Answers stored as bytes converted via proto.Marchal.
			err = proto.Unmarshal(binans, pa.Answers)
			if err != nil {
				return fmt.Errorf("Failed to read vote from database in GetPoll: %w", err)
			}
			q.Votes = append(q.Votes, &pa)
		}
		return nil
	})
	return q, err
}

// SaveToken saves token for specified poll in database.
//
// Token is represented as byte array.
func SaveToken(db *bolt.DB, token []byte, pollid int32) error {
	return db.Update(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		pbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(pollid)) + "Bucket"))
		if pbuck == nil {
			return fmt.Errorf("Poll ID does not exist in database. SaveToken: %v", pollid)
		}

		tbuck := pbuck.Bucket([]byte("TokensBucket"))
		return tbuck.Put(token, []byte{1})
	})
}

// AcceptToken checks if token sent by client is valid.
//
// Function returns true if token is present in database and
// if this token was not used before.
// If returned error is nil, token is accepted.
func AcceptToken(db *bolt.DB, token *query.VoteToken, pollid int32) error {
	return db.Update(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		pbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(pollid)) + "Bucket"))

		// Check if poll of number pollid exists.
		if pbuck == nil {
			return fmt.Errorf("No such poll: %v", pollid)
		}
		tbuck := pbuck.Bucket([]byte("TokensBucket"))

		// We check if requested token exists. If so, v will have one element, else 0.
		v := tbuck.Get(token.Token)
		if v == nil {
			return fmt.Errorf("No such token")
		}
		if v[0] == 0 {
			return fmt.Errorf("Token was used before")
		}
		// After use we remove token from database by setting its value to 0.
		return tbuck.Put(token.Token, []byte{0})
	})
}

// SaveVote is saving properly signed vote to database.
func SaveVote(db *bolt.DB, vr *query.VoteRequest) (*query.VoteReply, error) {
	reply := &query.VoteReply{}
	err := db.Update(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		pbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(vr.Pollid)) + "Bucket"))

		// Check if poll of number vr.Pollid exists.
		if pbuck == nil {
			return fmt.Errorf("No such poll: %v", vr.Pollid)
		}
		vbuck := pbuck.Bucket([]byte("VotesBucket"))

		// Save vote. Vote is stored as a bucket.
		// Name of this bucket is ballot used for signing it.
		ansbuck, err := vbuck.CreateBucketIfNotExists(vr.Sign.Ballot)
		if err != nil {
			return err
		}

		err = ansbuck.Put([]byte("Sign"), vr.Sign.Sign)
		if err != nil {
			return err
		}

		// Check if Vote is valid (in sense of valid characters etc.).
		if err = vr.Answers.IsValid(); err != nil {
			return err
		}

		binans, err := proto.Marshal(vr.Answers)
		if err != nil {
			return err
		}

		err = ansbuck.Put([]byte("Answer"), binans)
		if err != nil {
			return err
		}

		reply.Mess = "Thank you for your vote!"
		return nil
	})
	return reply, err
}

// GetSummary reads poll's answers from database.
func GetSummary(db *bolt.DB, pollid int32) (*query.PollSummary, error) {
	s := &query.PollSummary{
		Id:     pollid,
		Schema: &query.PollSchema{},
	}
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		pbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(pollid)) + "Bucket"))
		if pbuck == nil {
			return fmt.Errorf("Poll ID does not exist in database. GetPoll: %v", pollid)
		}

		// Read Schema stored as bytes converted via proto.Marchal.
		binschema := pbuck.Get([]byte("Schema"))
		err := proto.Unmarshal(binschema, s.Schema)
		if err != nil {
			return fmt.Errorf("Failed to read schema from database in GetPoll: %w", err)
		}

		// Votes are stored in VotesBucket.
		// Each vote is a different bucket inside VotesBucket, with name Vote+nr.
		vbuck := pbuck.Bucket([]byte("VotesBucket"))

		c := vbuck.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ansbuck := vbuck.Bucket(k)
			pa := &query.PollSchema{}

			binans := ansbuck.Get([]byte("Answer"))
			// Read Answers stored as bytes converted via proto.Marchal.
			err := proto.Unmarshal(binans, pa)
			if err != nil {
				return fmt.Errorf("Failed to read vote from database in GetPoll: %w", err)
			}
			s.Votes = append(s.Votes, pa)
		}
		return nil
	})
	return s, err
}
