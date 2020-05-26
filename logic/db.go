package logic

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

// GetKey reads key and from database.
//
// Key should be stored in bucket KeyBucket with label key.
// If key is not in database, nil is returned.
func GetKey(db *bolt.DB) *rsa.PrivateKey {
	var bkeycpy []byte
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		kbuck := tx.Bucket([]byte("KeyBucket"))
		bkey := kbuck.Get([]byte("key"))
		bkeycpy = make([]byte, len(bkey))
		copy(bkeycpy, bkey)
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to create bucket: %v", err)
	}

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

// SaveKey saves server key to database.
func SaveKey(db *bolt.DB, key *rsa.PrivateKey) error {
	bkey := x509.MarshalPKCS1PrivateKey(key)
	return db.Update(func(tx *bolt.Tx) error {
		keybuck := tx.Bucket([]byte("KeyBucket"))
		return keybuck.Put([]byte("key"), bkey)
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
		qbuck, err := queriesbuck.CreateBucketIfNotExists([]byte(strconv.FormatUint(id, 10)))
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

		qbuck := queriesbuck.Bucket([]byte(strconv.FormatInt(int64(queryid), 10)))
		fbuck := qbuck.Bucket([]byte("FieldsBucket"))

		id := fieldid
		// Edit field in query, -1 is a signal of new field.
		if fieldid == -1 {
			id64, _ := fbuck.NextSequence()
			id = int32(id64)
		} else {
			token := fbuck.Get([]byte(strconv.FormatInt(int64(fieldid), 10)))
			if token == nil {
				// Not existing field is requested.
				fmt.Printf("Wrong Query field number\n")
				// There should be some error returned here.
				// TODO
				return nil
			}
		}
		return fbuck.Put([]byte(strconv.FormatInt(int64(id), 10)), []byte(name))
	})
}

// GetKey reads query and from database.
//
func GetQuery(db *bolt.DB, id int) query.PollQuestion {
	q := query.PollQuestion{
		Id: int32(id),
	}
	// Database db should be open before this call.
	err := db.View(func(tx *bolt.Tx) error {
		queriesbuck := tx.Bucket([]byte("QueriesBucket"))

		qbuck := queriesbuck.Bucket([]byte(strconv.FormatInt(int64(id), 10)))
		fbuck := qbuck.Bucket([]byte("FieldsBucket"))
		c := fbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			f := query.PollQuestion_QueryField{}
			f.Name = string(v)
			q.Fields = append(q.Fields, &f)
		}

		tbuck := qbuck.Bucket([]byte("TokensBucket"))
		c = tbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t := query.VoteToken{}
			tokint, err := strconv.ParseInt(string(v), 10, 32)
			t.Token = int32(tokint)
			if err == nil {
				q.Tokens = append(q.Tokens, &t)
			}
		}

		//TODO
		/*vbuck := qbuck.Bucket([]byte("VotesBucket"))
		c = vbuck.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			vt := query.PollQuestion_StoredVote{}
			vt.Answer = v
			q.Votes = append(q.Votes, &vt)
		}*/
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to create bucket: %v", err)
	}

	return q
}
