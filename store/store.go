// Package for database management.
//
// Database structure:
// KeysBucket is storing keys for queries.
// Each key is stored in pair (keyid, key), where id is number of query
// and key is PKCS1 encoding of key.
// * KeysBucket
//   - (keyid, key)
//
// QueriesBucket is storing data of queries: fields, tokens and votes.
// * QueriesBucket
//
//	 Each query is stored in separate bucket.
//	 Id in name is its number (starting with 1!).
//   - QueryidBucket
//
//		 FieldsBucket stores query options represented as pair with its number and name (string).
// 		 + FieldsBucket
// 				- (id, name)
//
//		 TokensBucket is storing tokens to query.
//		 Each is stored as its value as key and empty value.
// 		 + TokensBucket
// 				- (token, _)
//
// 		 VotesBucket stores votes for query.
//		 + VotesBucket
//
//			  Each vote is stored in a bucket named after ballot used to sign it.
//				Answer to each field is pair consisting this field number and a number (actual answer).
// 				- Ballot
// 					* (id, answer)
// Each number value is stored as its decimal representation.
package store

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ememak/Projekt-Rada/query"
	bolt "go.etcd.io/bbolt"
)

// DBInit Opens database and create buckets for data.
//
// This function have to be called before any other database related function.
func DBInit(filename string) (*bolt.DB, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("KeyBucket"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("QueriesBucket"))
		return err
	})
	return db, err
}

// GetKey reads key for specific query from database.
//
// Key should be stored in bucket KeyBucket with label keyid, where id is number of query.
// If keyid is not in database, nil is returned.
// Key is stored in PKCS1 format.
func GetKey(db *bolt.DB, queryid int) *rsa.PrivateKey {
	var bkeycpy []byte
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		kbuck := tx.Bucket([]byte("KeyBucket"))
		bkey := kbuck.Get([]byte("key" + strconv.Itoa(queryid)))
		bkeycpy = make([]byte, len(bkey))
		copy(bkeycpy, bkey)
		return nil
	})

	var key *rsa.PrivateKey
	if len(bkeycpy) != 0 {
		key, err = x509.ParsePKCS1PrivateKey(bkeycpy)
		if err != nil {
			fmt.Printf("Failed to convert key from binary: %v", err)
		}
	} else {
		key = nil
	}
	return key
}

// SaveKey saves query key to database.
//
// Key is saved in bucket KeyBucket with label keyid, where id is number of query.
// Key is stored in PKCS1 format.
func SaveKey(db *bolt.DB, queryid int, key *rsa.PrivateKey) error {
	bkey := x509.MarshalPKCS1PrivateKey(key)
	return db.Update(func(tx *bolt.Tx) error {
		keybuck := tx.Bucket([]byte("KeyBucket"))
		return keybuck.Put([]byte("key"+strconv.Itoa(queryid)), bkey)
	})
}

// NewQuery creates bucket for new query.
//
// Return values is an id of query in database and error returned by database.
func NewQuery(db *bolt.DB) (qid int, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		// All queries are stored in QueriesBucket.
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))
		id, _ := queriesbuck.NextSequence()
		qid = int(id)

		// Each query is contained in bucket named by its number.
		qbuck, err := queriesbuck.CreateBucketIfNotExists([]byte("Query" + strconv.Itoa(int(id)) + "Bucket"))
		if err != nil {
			return err
		}

		// Inside of a query bucket there are three buckets:
		// One for fields, one for tokens, one for votes.
		_, err = qbuck.CreateBucketIfNotExists([]byte("FieldsBucket"))
		if err != nil {
			return err
		}

		_, err = qbuck.CreateBucketIfNotExists([]byte("TokensBucket"))
		if err != nil {
			return err
		}

		_, err = qbuck.CreateBucketIfNotExists([]byte("VotesBucket"))
		return err
	})
	return
}

// ModifyQueryField is editing or adding new field to query.
//
// Input is number of modified query, number of modified field (-1 is new field) and
// string with name of this field.
func ModifyQueryField(db *bolt.DB, queryid int, fieldid int32, name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))

		// Each query have its own bucket inside QueriesBucket.
		qbuck := queriesbuck.Bucket([]byte("Query" + strconv.Itoa(queryid) + "Bucket"))
		if qbuck == nil {
			fmt.Printf("Wrong Query number: %v", queryid)
			return nil
		}
		fbuck := qbuck.Bucket([]byte("FieldsBucket"))

		id := fieldid
		// Edit field in query, -1 is a signal of new field.
		if fieldid == -1 {
			id64, _ := fbuck.NextSequence()
			id = int32(id64)
		} else {
			field := fbuck.Get([]byte(strconv.Itoa(int(fieldid))))
			if field == nil {
				// Not existing field is requested.
				fmt.Printf("Wrong Query field number\n")
				return nil
			}
		}

		return fbuck.Put([]byte(strconv.Itoa(int(id))), []byte(name))
	})
}

// GetQuery reads query from database.
func GetQuery(db *bolt.DB, id int) query.PollQuestion {
	q := query.PollQuestion{
		Id: int32(id),
	}
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))

		qbuck := queriesbuck.Bucket([]byte("Query" + strconv.Itoa(id) + "Bucket"))
		if qbuck == nil {
			fmt.Printf("Wrong Query number: %v", id)
			return nil
		}

		// Fields are stored as pairs (nr, name), where name is a name of this choice.
		fbuck := qbuck.Bucket([]byte("FieldsBucket"))
		c := fbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			f := query.PollQuestion_QueryField{}
			f.Name = string(v)
			q.Fields = append(q.Fields, &f)
		}

		// Tokens are stored as pairs (token, _), where token is sha256 hash of a random int.
		tbuck := qbuck.Bucket([]byte("TokensBucket"))
		c = tbuck.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			t := query.VoteToken{}
			t.Token = k
			q.Tokens = append(q.Tokens, &t)
		}

		// Votes are stored in VotesBucket.
		// Each vote is a different bucket inside VotesBucket, with name Vote+nr.
		vbuck := qbuck.Bucket([]byte("VotesBucket"))

		c = vbuck.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ansbuck := vbuck.Bucket(k)
			vt := query.PollQuestion_StoredVote{}

			cins := ansbuck.Cursor()
			for ki, vi := cins.First(); ki != nil; ki, vi = cins.Next() {
				val, errins := strconv.Atoi(string(vi))
				if errins != nil {
					return errins
				}
				vt.Answer = append(vt.Answer, int32(val))
			}
			q.Votes = append(q.Votes, &vt)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to read from database in GetQuery: %v", err)
	}
	return q
}

// NewToken checks if token request is valid and return new token.
func NewToken(db *bolt.DB, in *query.TokenRequest) (*query.VoteToken, error) {
	t := query.VoteToken{}
	err := db.Update(func(tx *bolt.Tx) error {
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))

		qbuck := queriesbuck.Bucket([]byte("Query" + strconv.Itoa(int(in.Nr)) + "Bucket"))
		if qbuck == nil {
			fmt.Printf("Wrong Query number in NewToken: %v\n", in.Nr)
			return nil
		}
		tbuck := qbuck.Bucket([]byte("TokensBucket"))

		// Token is a sha256 hash of a random number from range [0, 2^1024).
		max := big.NewInt(2)
		max = max.Exp(max, big.NewInt(1024), big.NewInt(0))
		val, err := rand.Int(rand.Reader, max)
		if err != nil {

		}
		token := sha256.Sum256(val.Bytes())

		err = tbuck.Put(token[:], []byte{})
		if err != nil {
			return err
		}

		t.Token = token[:]
		return err
	})
	return &t, err
}

// AcceptToken checks if MessageToSign is matching server informations.
//
// Function returns true if token is a token of data[qNum] (provided
// such query exists). In other case it returns false.
func AcceptToken(db *bolt.DB, token *query.VoteToken, queryid int32) (res bool, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))

		qbuck := queriesbuck.Bucket([]byte("Query" + strconv.Itoa(int(queryid)) + "Bucket"))

		// Check if query of number queryid exists.
		if qbuck == nil {
			fmt.Printf("No such query: %v\n", queryid)
			res = false
			return nil
		}
		tbuck := qbuck.Bucket([]byte("TokensBucket"))

		// We check if requested token exists. If so, v will have one element, else 0.
		v := tbuck.Get(token.Token)
		if v != nil {
			res = true
			// After use we remove token from database.
			tbuck.Delete(token.Token)
		} else {
			fmt.Print("No such token\n")
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
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))

		qbuck := queriesbuck.Bucket([]byte("Query" + strconv.Itoa(int(v.Nr)) + "Bucket"))

		// Check if query of number v.Nr exists.
		if qbuck == nil {
			vr.Mess = "Vote error"
			return fmt.Errorf("No such query: %v\n", v.Nr)
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
			fmt.Printf("Wrong vote size: %v, wanted: %v\n", len(v.Answer), nF)
			vr.Mess = "Vote error"
			return nil
		}

		// Save vote. Vote is stored as a bucket.
		// Name of this bucket is ballot used for signing it.
		newvbuck, errins := vbuck.CreateBucket(sv.Signm)

		if errins != nil {
			fmt.Printf("Failed to create bucket\n")
			vr.Mess = "Vote error"
			return errins
		}

		// Inside vote bucket are pairs (nr, vote), where vote is int value representing choice.
		for i := 0; i < nF; i++ {
			errins := newvbuck.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(int(v.Answer[i]))))
			if errins != nil {
				fmt.Printf("Failed to put part of the vote: %v\n", i)
				vr.Mess = "Vote error"
				return errins
			}
		}

		vr.Mess = "Thank you for your vote!\n"
		return nil
	})
	fmt.Printf("Response: %v\n", vr)
	return
}
