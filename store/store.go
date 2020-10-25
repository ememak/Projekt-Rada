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
//       SchemaBucket stores poll questions.
//			 Each QA structure have its own bucket containing it.
//       + SchemaBucket
//         - QAid
//					 * ("Question", question)
//					 * ("Type", type)
//					 * ("Answer", answer)
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
//         PollAnswer structure is stored similarily as PollSchema, but
//				 apart from QA structures, this bucket also contains second part of RSA sign.
//         - Ballot
//					 + ("Sign", sign)
//					 + QAid
//					 	 * ("Question", question)
//					 	 * ("Type", type)
//					 	 * ("Answer", answer)
// Each number value is stored using strconv.Itoa function.
package store

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strconv"

	"github.com/ememak/Projekt-Rada/query"
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
func GetKey(db *bolt.DB, pollid int) (*rsa.PrivateKey, error) {
	var bkeycpy []byte
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		kbuck := tx.Bucket([]byte("KeyBucket"))
		bkey := kbuck.Get([]byte("key" + strconv.Itoa(pollid)))
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
	if key == nil || key.Validate() != nil {
		if key == nil {
			return fmt.Errorf("Error! Private key is nil!")
		}
		return key.Validate()
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

		// Inside of a poll bucket there are three buckets:
		// One for schema, one for tokens, one for votes.
		sbuck, err := pbuck.CreateBucketIfNotExists([]byte("SchemaBucket"))
		if err != nil {
			return err
		}

		for i, qa := range sch.Questions {
			qbuck, err := sbuck.CreateBucketIfNotExists([]byte("QA" + string(i)))
			if err != nil {
				return err
			}

			if !query.IsStringPrintable(qa.Question) {
				return fmt.Errorf("Error! Question contains non valid characters.")
			}
			err = qbuck.Put([]byte("Question"), []byte(qa.Question))
			if err != nil {
				return err
			}

			if !qa.Type.IsValid() {
				return fmt.Errorf("Error! Wrong question type.")
			}
			qbuck.Put([]byte("Type"), []byte(strconv.Itoa(int(qa.Type))))
			if err != nil {
				return err
			}

			qbuck.Put([]byte("Answer"), []byte(qa.Answer))
			if err != nil {
				return err
			}
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
func GetPoll(db *bolt.DB, id int) (query.PollQuestion, error) {
	q := query.PollQuestion{
		Id:     int32(id),
		Schema: &query.PollSchema{},
	}
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		pbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(id) + "Bucket"))
		if pbuck == nil {
			return fmt.Errorf("Wrong Poll number in GetPoll: %v", id)
		}

		// Read Schema stored as an array of buckets with QA structures.
		sbuck := pbuck.Bucket([]byte("SchemaBucket"))
		c := sbuck.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			qa := query.PollSchema_QA{}
			qabuck := sbuck.Bucket(k)
			qa.Question = string(qabuck.Get([]byte("Question")))
			inttype, errins := strconv.Atoi(string(qabuck.Get([]byte("Type"))))
			if errins != nil {
				return fmt.Errorf("Failed to convert answer to number in GetPoll: %w", errins)
			}
			qa.Type = query.PollSchema_QuestionType(inttype)
			q.Schema.Questions = append(q.Schema.Questions, &qa)
		}
		/* not tested or can be tested yet
		// Tokens are stored as pairs (token, bool), where bool is true if token was not used.
		// If used == false, token is not read from database.
		tbuck := pbuck.Bucket([]byte("TokensBucket"))
		c = tbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := query.VoteToken{}
			b, errins := strconv.Atoi(string(v))
			if errins != nil {
				return fmt.Errorf("Failed to convert answer to number in GetPoll: %w", errins)
			}
			if b == 1 {
				t.Token = k
			}
			q.Tokens = append(q.Tokens, &t)
		}

		// Votes are stored in VotesBucket.
		// Each vote is a different bucket inside VotesBucket, with name Vote+nr.
		vbuck := pbuck.Bucket([]byte("VotesBucket"))

		c = vbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ansbuck := vbuck.Bucket(k)
			pa := query.PollAnswer{}
			pa.Sign = quert.RSASignature{
				Ballot: k,
			}

			cins := ansbuck.Cursor()
			for k, v := cins.First(); k != nil; k, v = cins.Next() {
				if string(k) == "Sign" {
					pa.Sign.Sign = v
				} else {
					qa := query.PollSchema_QA{}
					qabuck := sbuck.Bucket(k)
					qa.Question = string(qabuck.Get([]byte("Question")))
					inttype, errins := strconv.Atoi(string(qabuck.Get([]byte("Type"))))
					if errins != nil {
						return fmt.Errorf("Failed to convert answer to number in GetPoll: %w", errins)
					}
					qa.Type = query.PollSchema_QuestionType(inttype)
					qa.Answer = string(qabuck.Get([]byte("Answer")))
					qa. = string(qabuck.Get([]byte("Answer")))
					pa.Answer.Questions = append(pa.Answer.Questions, &qa)
				}
			}
			q.Votes = append(q.Votes, &pa)
		}
		*/
		return nil
	})
	return q, err
}

/*
// AcceptToken checks if BallotToSign is matching server informations.
//
// Function returns true if token is a token of data[qNum] (provided
// such poll exists). In other case it returns false.
func AcceptToken(db *bolt.DB, token *query.VoteToken, pollid int32) (res bool, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		qbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(pollid)) + "Bucket"))

		// Check if poll of number pollid exists.
		if qbuck == nil {
			res = false
			return fmt.Errorf("No such poll: %v", pollid)
		}
		tbuck := qbuck.Bucket([]byte("TokensBucket"))

		// We check if requested token exists. If so, v will have one element, else 0.
		v := tbuck.Get(token.Token)
		if v != nil {
			res = true
			// After use we remove token from database.
			tbuck.Delete(token.Token)
		} else {
			fmt.Print("No such token")
			res = false
		}
		return nil
	})
	return
}

// AcceptVote is saving properly signed vote to database.
func AcceptVote(db *bolt.DB, sv *query.SignedVote) (vr *query.VoteReply, err error) {
	vr = &query.VoteReply{}
	v := sv.Vote
	err = db.Update(func(tx *bolt.Tx) error {
		pollsbuck := tx.Bucket([]byte("PollsBucket"))

		qbuck := pollsbuck.Bucket([]byte("Poll" + strconv.Itoa(int(v.Pollid)) + "Bucket"))

		// Check if poll of number v.Nr exists.
		if qbuck == nil {
			vr.Mess = "Vote error"
			return fmt.Errorf("No such poll: %v", v.Pollid)
		}
		vbuck := qbuck.Bucket([]byte("VotesBucket"))
		fbuck := qbuck.Bucket([]byte("FieldsBucket"))

		// Check if answer have good number of fields.
		c := fbuck.Cursor()
		nF := 0
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			nF += 1
		}

		if len(v.Answer) != nF {
			vr.Mess = "Vote error"
			return fmt.Errorf("Wrong vote size: %v, wanted: %v", len(v.Answer), nF)
		}

		// Save vote. Vote is stored as a bucket.
		// Name of this bucket is ballot used for signing it.
		newvbuck, errins := vbuck.CreateBucketIfNotExists(sv.Signm)

		if errins != nil {
			vr.Mess = "Vote error"
			return errins
		}

		// Inside vote bucket are pairs (nr, vote), where vote is int value representing choice.
		for i := 0; i < nF; i++ {
			errins := newvbuck.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(int(v.Answer[i]))))
			if errins != nil {
				vr.Mess = "Vote error"
				return errins
			}
		}

		vr.Mess = "Thank you for your vote!"
		return nil
	})
	fmt.Printf("Response: %v", vr)
	return
}*/
