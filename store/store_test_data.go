package store

import (
	"crypto/rsa"
	"fmt"
	"github.com/ememak/Projekt-Rada/query"
	"math/big"
)

var testsDBInit = []struct {
	in      string
	exp_out string
	exp_err error
}{
	{ // test0 - positive
		in:      "data.db",
		exp_out: "bolt.DB{path:\"data.db\"}",
		exp_err: nil,
	},
	{ // test1 - positive
		in:      "data",
		exp_out: "bolt.DB{path:\"data\"}",
		exp_err: nil,
	},
	{ // test2 - negative
		in:      "",
		exp_out: "",
		exp_err: fmt.Errorf("Database name invalid\n"),
	},
}

var testsNewPoll = []struct {
	in      *query.PollSchema
	exp_err error
}{
	{
		in: &query.PollSchema{ // test0 - positive
			Questions: []*query.PollSchema_QA{
				{
					Question: "Do you like this system? Options: yes/no",
					Type:     query.PollSchema_CLOSE,
				},
				{
					Question: "Why?",
					Type:     query.PollSchema_OPEN,
				},
			},
		},
		exp_err: nil,
	},
	{
		in: &query.PollSchema{ // test1 - negative, wrong characters in question
			Questions: []*query.PollSchema_QA{
				{
					Question: "Wrong\x01characters\x07\x00",
					Type:     query.PollSchema_CLOSE,
				},
			},
		},
		exp_err: fmt.Errorf("Error! Question contains non valid characters."),
	},
	{
		in: &query.PollSchema{ // test2 - negative, wrong question type
			Questions: []*query.PollSchema_QA{
				{
					Question: "Valid question\n!@#$%",
					Type:     5,
				},
			},
		},
		exp_err: fmt.Errorf("Error! Wrong question type."),
	},
	{
		in: &query.PollSchema{ // test3 - positive
			Questions: []*query.PollSchema_QA{
				{
					Question: "^&*()'\"\\",
					Type:     0,
				},
				{
					Question: "{}:<>?,./;[]+=-_",
					Type:     1,
				},
				{
					Question: "1234567890\x09", // \x09 - Tab
					Type:     2,
				},
			},
		},
		exp_err: nil,
	},
}

var testsSaveKey = []struct {
	in      *rsa.PrivateKey
	exp_err error
}{
	{ // test3 - nil test, negative
		in:      nil,
		exp_err: fmt.Errorf("Error! Private key is nil!"),
	},
	{ // test4 - negative, wrong key
		in: &rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: new(big.Int).SetBytes([]byte{188, 55, 229, 23, 136, 225, 241, 143, 196, 59, 114, 80, 198, 19, 73, 125}),
				E: 3, // should be 65537 here
			},
			D:      new(big.Int).SetBytes([]byte{114, 248, 49, 152, 14, 164, 245, 72, 46, 250, 250, 63, 55, 189, 76, 1}),
			Primes: []*big.Int{new(big.Int).SetBytes([]byte{239, 234, 64, 132, 237, 120, 37, 181}), new(big.Int).SetBytes([]byte{200, 214, 90, 142, 87, 23, 241, 169})},
		},
		exp_err: fmt.Errorf("crypto/rsa: invalid exponents"),
	},
	{ // test5 - positive
		in: &rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: new(big.Int).SetBytes([]byte{188, 55, 229, 23, 136, 225, 241, 143, 196, 59, 114, 80, 198, 19, 73, 125}),
				E: 65537,
			},
			D:      new(big.Int).SetBytes([]byte{114, 248, 49, 152, 14, 164, 245, 72, 46, 250, 250, 63, 55, 189, 76, 1}),
			Primes: []*big.Int{new(big.Int).SetBytes([]byte{239, 234, 64, 132, 237, 120, 37, 181}), new(big.Int).SetBytes([]byte{200, 214, 90, 142, 87, 23, 241, 169})},
		},
		exp_err: nil,
	},
}

var testsAcceptToken = []struct {
	token  []byte
	pollid int32
	st_err error
	gp_err error
	at_err error
}{
	{ // test0 - positive
		token:  []byte("GoodToken"),
		pollid: 1,
		st_err: nil,
		gp_err: nil,
		at_err: nil,
	},
	{ // test1 - negative, wrong poll requested
		token:  []byte("GoodToken"),
		pollid: 2,
		st_err: fmt.Errorf("Poll ID does not exist in database. SaveToken: 2"),
		gp_err: fmt.Errorf("Poll ID does not exist in database. GetPoll: 2"),
		at_err: fmt.Errorf("No such poll: 2"),
	},
	{ // test2 - negative, token can't be empty
		token:  []byte(""),
		pollid: 1,
		st_err: fmt.Errorf("key required"),
		gp_err: nil,
		at_err: fmt.Errorf("No such token"),
	},
	{ // test3 - positive
		token:  []byte{200, 235, 239, 61, 75, 253, 155, 22, 139, 151, 158, 172, 35, 68, 192, 61, 65, 159, 101, 47, 90, 93, 218, 42, 49, 50, 32, 3, 71, 10, 28, 202, 158, 245, 51, 8, 194, 26, 105, 179, 209, 157, 190, 20, 55, 190, 129, 244, 10, 92, 130, 72, 151, 74, 111, 135, 8, 244, 121, 145, 217, 85, 152, 167, 94, 21, 16, 176, 197, 195, 82, 194, 230, 184, 73, 111, 44, 28, 32, 194, 26, 108, 93, 120, 214, 9, 85, 11, 120, 250, 157, 252, 74, 19, 53, 40, 11, 191, 208, 153, 103}, //random 100 bytes
		pollid: 1,
		st_err: nil,
		gp_err: nil,
		at_err: nil,
	},
}

var testsSaveVote = []struct {
	in     *query.VoteRequest
	reply  *query.VoteReply
	sv_err error
	gp_err error
}{
	{ // test0 - positive
		in: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Do you like this system? Options: yes/no Answer: yes",
						Type:     query.PollSchema_CLOSE,
					},
					{
						Question: "Why? Answer: Its cool.",
						Type:     query.PollSchema_OPEN,
					},
				}},
			Sign: &query.RSASignature{
				Ballot: []byte{1}, // Here RSA signature is not checked, but it can't be empty
				Sign:   []byte{1},
			},
		},
		reply: &query.VoteReply{
			Mess: "Thank you for your vote!",
		},
		sv_err: nil,
		gp_err: nil,
	},
	{ // test1 - negative, wrong pollid
		in: &query.VoteRequest{
			Pollid: 2,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Do you like this system? Options: yes/no Answer: yes",
						Type:     query.PollSchema_CLOSE,
					},
					{
						Question: "Why? Answer: Its cool.",
						Type:     query.PollSchema_OPEN,
					},
				}},
			Sign: &query.RSASignature{
				Ballot: []byte{1},
				Sign:   []byte{1},
			},
		},
		reply:  &query.VoteReply{},
		sv_err: fmt.Errorf("No such poll: 2"),
		gp_err: fmt.Errorf("Poll ID does not exist in database. GetPoll: 2"),
	},
	{ // test2 - negative, wrong characters in answer
		in: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Do you like this system? Options: yes/no",
						Type:     query.PollSchema_CLOSE,
						Answers:  []string{"yes"},
					},
					{
						Question: "Why?",
						Type:     query.PollSchema_OPEN,
						Answers:  []string{"\x01\x21\xae"}, // Non valid characters
					},
				}},
			Sign: &query.RSASignature{
				Ballot: []byte{1},
				Sign:   []byte{1},
			},
		},
		reply:  &query.VoteReply{},
		sv_err: fmt.Errorf("Error! Answer contains non valid characters."),
		gp_err: nil,
	},
	{ // test3 - positive
		in: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "!@#$%^&*()_+:\"<>?,./;'[]-=",
						Type:     0,
						Answers:  []string{"qwertyuiopasdfghjk"},
					},
				}},
			Sign: &query.RSASignature{
				Ballot: []byte{123, 34, 56, 4, 19},
				Sign:   []byte{23, 6, 3, 0, 8},
			},
		},
		reply: &query.VoteReply{
			Mess: "Thank you for your vote!",
		},
		sv_err: nil,
		gp_err: nil,
	},
	{ // test4 - negative, wrong characters in answer
		in: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Do you like this system? Options: yes/no",
						Type:     query.PollSchema_CLOSE,
						Answers:  []string{"yes"},
					},
					{
						Question: "\x01\x21\xae", // Non valid characters
						Type:     query.PollSchema_OPEN,
						Answers:  []string{"whatever"},
					},
				}},
			Sign: &query.RSASignature{
				Ballot: []byte{1},
				Sign:   []byte{1},
			},
		},
		reply:  &query.VoteReply{},
		sv_err: fmt.Errorf("Error! Question contains non valid characters."),
		gp_err: nil,
	},
}
