package store

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestDBInit(t *testing.T) {
	tests := testsDBInit
	for i, test := range tests {

		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			out, err := DBInit(test.in)
			if !reflect.DeepEqual(err, test.exp_err) {
				t.Errorf("Error %v, want error %v", err, test.exp_err)
			}
			if out != nil {
				if !reflect.DeepEqual(out.GoString(), test.exp_out) {
					t.Errorf("Output %v, want output %v", out.GoString(), test.exp_out)
				}
			} else {
				if test.exp_out != "" {
					t.Errorf("Ouput is nil, wanted output %v", test.exp_out)
				}
			}
		})
	}
}

func TestNewPoll(t *testing.T) {
	in := testsNewPoll
	for i, test := range in {

		data, _ := DBInit("testNP" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			p, err := NewPoll(data, test.in)
			if !reflect.DeepEqual(err, test.exp_err) {
				t.Errorf("Error %v, want error %v", err, test.exp_err)
				return
			} else {
				if err != nil {
					return
				}
			}
			o, err := GetPoll(data, p.Id)
			if !(proto.Equal(o.Schema, test.in) && reflect.DeepEqual(err, nil)) {
				t.Errorf("Output %v, want output %v", o.Schema, test.in)
				t.Errorf("Error %v, want nil error", err)
			}
		})
		data.Close()
	}
}

func TestSaveKey(t *testing.T) {
	for i := 0; i < 3; i++ { // Three random tests, all should be positive
		key, _ := rsa.GenerateKey(rand.Reader, 2048)

		data, _ := DBInit("testSK" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			err := SaveKey(data, i, key)
			if err != nil {
				t.Errorf("Error %v, want nil error", err)
				return
			}
			keyret, err := GetKey(data, int32(i))
			if err != nil || !reflect.DeepEqual(keyret, key) {
				t.Errorf("Output %v, want output %v", keyret, key)
				t.Errorf("Error %v, want nil error", err)
			}
		})
	}
	tests := testsSaveKey
	for i, test := range tests {
		data, _ := DBInit("testSK" + strconv.Itoa(i+3) + ".db")

		t.Run("Test "+strconv.Itoa(i+3), func(t *testing.T) {
			err := SaveKey(data, 1, test.in)
			if !reflect.DeepEqual(err, test.exp_err) {
				t.Errorf("Error %v, want error %v", err, test.exp_err)
				return
			} else {
				if err != nil {
					return
				}
			}
			key, err := GetKey(data, 1)
			if !reflect.DeepEqual(key, test.in) || !reflect.DeepEqual(err, nil) {
				t.Errorf("Output %v, want output %v", key, test.in)
				t.Errorf("Error %v, want nil error", err)
			}
		})
	}
}

func TestAcceptToken(t *testing.T) {
	in := testsAcceptToken
	for i, test := range in {

		data, _ := DBInit("testAT" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			NewPoll(data, testsNewPoll[0].in)
			err := SaveToken(data, test.token, test.pollid)
			if !reflect.DeepEqual(err, test.st_err) {
				t.Errorf("Error %v, want nil error", err, test.st_err)
			}

			pq, err := GetPoll(data, test.pollid)
			if !reflect.DeepEqual(err, test.gp_err) {
				t.Errorf("Error %v, want error %v", err, test.gp_err)
			}
			if test.st_err == nil {
				found := false
				for _, tok := range pq.Tokens {
					if reflect.DeepEqual(tok, test.token) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Output %v, want output %v", pq.Tokens[0], test.token)
				}
			}

			err = AcceptToken(data, test.token, test.pollid)
			if !reflect.DeepEqual(err, test.at_err) {
				t.Errorf("Error %v, want error %v", err, test.at_err)
				return
			}
		})
		data.Close()
	}
	test := testsAcceptToken[0]
	rep_err := fmt.Errorf("Token was used before")
	data, _ := DBInit("testAT" + strconv.Itoa(len(in)) + ".db")
	t.Run("Test "+strconv.Itoa(len(in)), func(t *testing.T) {
		NewPoll(data, testsNewPoll[0].in)
		err := SaveToken(data, test.token, test.pollid)
		if !reflect.DeepEqual(err, test.st_err) {
			t.Errorf("Error %v, want error %v", err, test.st_err)
		}

		pq, err := GetPoll(data, test.pollid)
		if !reflect.DeepEqual(err, test.gp_err) {
			t.Errorf("Error %v, want error %v", err, test.gp_err)
		}
		if test.st_err == nil {
			found := false
			for _, tok := range pq.Tokens {
				if reflect.DeepEqual(tok, test.token) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Output %v, want output %v", pq.Tokens[0], test.token)
			}
		}
		err = AcceptToken(data, test.token, test.pollid)
		if !reflect.DeepEqual(err, test.at_err) {
			t.Errorf("Error %v, want error %v", err, test.at_err)
		}
		err = AcceptToken(data, test.token, test.pollid)
		if !reflect.DeepEqual(err, rep_err) {
			t.Errorf("Error %v, want error %v", err, rep_err)
		}
	})
	data.Close()
}

func TestSaveVote(t *testing.T) {
	in := testsSaveVote
	for i, test := range in {

		data, _ := DBInit("testSV" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			NewPoll(data, testsNewPoll[0].in)
			vr, err := SaveVote(data, test.in)
			if !(reflect.DeepEqual(err, test.sv_err) && proto.Equal(vr, test.reply)) {
				t.Errorf("Output %v, want output %v", vr, test.reply)
				t.Errorf("Error %v, want error %v", err, test.sv_err)
			}

			pq, err := GetPoll(data, test.in.Pollid)
			if !reflect.DeepEqual(err, test.gp_err) {
				t.Errorf("Error %v, want error %v", err, test.gp_err)
			}
			if test.sv_err == nil {
				if !(proto.Equal(pq.Votes[0].Answers, test.in.Answers) && proto.Equal(pq.Votes[0].Sign, test.in.Sign)) {
					t.Errorf("Answers %v, want output %v", pq.Votes[0].Answers, test.in.Answers)
					t.Errorf("Sign %v, want output %v", pq.Votes[0].Sign, test.in.Sign)
				}
			}
		})
		data.Close()
	}
}

func TestGetSummary(t *testing.T) {
	in := testsGetSummary
	for i, test := range in {

		data, _ := DBInit("testGS" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			NewPoll(data, test.schema)
			vr, err := SaveVote(data, test.in)
			if !(reflect.DeepEqual(err, test.sv_err) && proto.Equal(vr, test.sv_out)) {
				t.Errorf("Output %v, want output %v", vr, test.sv_out)
				t.Errorf("Error %v, want error %v", err, test.sv_err)
			}

			ps, err := GetSummary(data, test.in.Pollid)
			if !reflect.DeepEqual(err, test.gs_err) || !reflect.DeepEqual(ps, test.gs_out) {
				t.Errorf("Output %v, want output %v", ps, test.gs_out)
				t.Errorf("Error %v, want error %v", err, test.gs_err)
			}
		})
		data.Close()
	}
}
