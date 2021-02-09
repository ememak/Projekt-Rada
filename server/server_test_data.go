package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ememak/Projekt-Rada/query"
)

var testsPollInitIn = []*query.PollSchema{
	&query.PollSchema{ // test0 - positive
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
	&query.PollSchema{ // test1 - negative, wrong type value
		Questions: []*query.PollSchema_QA{
			{
				Question: "Do you like this system? Options: yes/no",
				Type:     query.PollSchema_CLOSE,
			},
			{
				Question: "Why?",
				Type:     3, // wrong type!
			},
		},
	},
	&query.PollSchema{ // test2 - negative, wrong type value
		Questions: []*query.PollSchema_QA{
			{
				Question: "Do you like this system? Options: yes/no",
				Type:     -1, //wrong type!
			},
		},
	},
	&query.PollSchema{ // test3 - positive
		Questions: []*query.PollSchema_QA{
			{
				Question: "Check numbers you like? Options: 1;2;5;e;74",
				Type:     query.PollSchema_CHECKBOX,
			},
		},
	},
	&query.PollSchema{ // test4 - negative, wrong characters in question
		Questions: []*query.PollSchema_QA{
			{
				Question: "\x00\x01\x02\xff\xe7",
				Type:     query.PollSchema_OPEN,
			},
		},
	},
}

var testsPollInitOutEmpty = []struct {
	exp_out *query.PollQuestion
	exp_err error
}{
	{
		exp_out: &query.PollQuestion{
			Id: 1,
			Schema: &query.PollSchema{
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
		},
		exp_err: nil,
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Wrong question type.")),
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Wrong question type.")),
	},
	{
		exp_out: &query.PollQuestion{
			Id: 1,
			Schema: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Check numbers you like? Options: 1;2;5;e;74",
						Type:     query.PollSchema_CHECKBOX,
					},
				},
			},
		},
		exp_err: nil,
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Question contains invalid characters.")),
	},
}

var testsPollInitOut1Poll = []struct {
	exp_out *query.PollQuestion
	exp_err error
}{
	{
		exp_out: &query.PollQuestion{
			Id: 2,
			Schema: &query.PollSchema{
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
		},
		exp_err: nil,
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Wrong question type.")),
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Wrong question type.")),
	},
	{
		exp_out: &query.PollQuestion{
			Id: 2,
			Schema: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Check numbers you like? Options: 1;2;5;e;74",
						Type:     query.PollSchema_CHECKBOX,
					},
				},
			},
		},
		exp_err: nil,
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Question contains invalid characters.")),
	},
}

var testsGetPollIn = []struct {
	schema  *query.PollSchema
	pollreq *query.GetPollRequest
}{
	{ // test0 - positive
		schema: &query.PollSchema{
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
		pollreq: &query.GetPollRequest{
			Pollid: 1,
		},
	},
	{ // test1 - negative, wrong poll requested
		schema: &query.PollSchema{
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
		pollreq: &query.GetPollRequest{
			Pollid: 2,
		},
	},
	{ // test2 - negative, wrong poll requested
		schema: &query.PollSchema{
			Questions: []*query.PollSchema_QA{
				{
					Question: "Do you like this system? Options: yes/no",
					Type:     query.PollSchema_CLOSE,
				},
			},
		},
		pollreq: &query.GetPollRequest{
			Pollid: -10000,
		},
	},
	{ // test3 - negative, wrong type - poll won't be saved in database
		schema: &query.PollSchema{
			Questions: []*query.PollSchema_QA{
				{
					Question: "Do you like this system? Options: yes/no",
					Type:     -3,
				},
				{
					Question: "Why?",
					Type:     query.PollSchema_OPEN,
				},
			},
		},
		pollreq: &query.GetPollRequest{
			Pollid: 1,
		},
	},
}

var testsGetPollOut = []struct {
	exp_out *query.PollSchema
	exp_err error
}{
	{
		exp_out: &query.PollSchema{
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
		exp_out: nil,
		exp_err: fmt.Errorf("Error in GetPoll while retrieving key from database: %w", fmt.Errorf("No key for this poll in database.")),
	},
	{
		exp_out: nil,
		exp_err: fmt.Errorf("Error in GetPoll while retrieving key from database: %w", fmt.Errorf("No key for this poll in database.")),
	},
	{
		exp_out: nil,
		exp_err: fmt.Errorf("Error in GetPoll while retrieving key from database: %w", fmt.Errorf("No key for this poll in database.")),
	},
}

var testsSignBallot = []struct {
	schema   *query.PollSchema
	envelope *query.EnvelopeToSign
	exp_err  error
}{
	{ // test0 - positive
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: []byte{1, 3, 4, 5, 6, 7, 8, 9, 0},
			Pollid:   1,
			Token:    "Good token",
		},
		exp_err: nil,
	},
	{ // test1 - negative, wrong poll requested
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: []byte{1, 3, 4, 5, 6, 7, 8, 9, 0},
			Pollid:   2,
			Token:    "Good token",
		},
		exp_err: fmt.Errorf("No such poll: 2"),
	},
	{ // test2 - negative, wrong poll requested
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: []byte{1, 3, 4, 5, 6, 7, 8, 9, 0},
			Pollid:   -3,
			Token:    "Good token",
		},
		exp_err: fmt.Errorf("No such poll: -3"),
	},
	{ // test3 - negative, token not valid (only valid token is "Good token")
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: []byte{1, 3, 4, 5, 6, 7, 8, 9, 0},
			Pollid:   1,
			Token:    "Bad token",
		},
		exp_err: fmt.Errorf("No such token"),
	},
	{ // test4 - negative, token not valid (only valid token is "Good token")
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: nil,
			Pollid:   1,
			Token:    "Good token",
		},
		exp_err: fmt.Errorf("Error in SignBallot, envelope shouldn't be null"),
	},
	{ // test5 - positive
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: []byte("Some random value\x01\x10!@#$%^&*()_+{}"),
			Pollid:   1,
			Token:    "Good token",
		},
		exp_err: nil,
	},
}

// For simplicity there is no blinding in this tests.
// Hash of "12345678"
var hash, _ = hex.DecodeString("ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f")

// Hash of "rvrbhd54":^V(B)*TBytvw.ucq<{_@x-mzua"
var hash2, _ = hex.DecodeString("563b5c8ae85cabe26529b9857d4b503e6da069389bbb94490cf1873176d7ff94")

var testsPollVote = []struct {
	schema   *query.PollSchema
	envelope *query.EnvelopeToSign
	votereq  *query.VoteRequest
	exp_out  *query.VoteReply
	exp_err  error
}{
	{ // test0 - positive
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid:  1,
			Answers: &query.PollSchema{},
			Sign: &query.RSASignature{
				Ballot: []byte("12345678"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Thank you for your vote!",
		},
		exp_err: nil,
	},
	{ // test1 - positive
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash2,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid:  1,
			Answers: &query.PollSchema{},
			Sign: &query.RSASignature{
				Ballot: []byte("rvrbhd54\":^V(B)*TBytvw.ucq<{_@x-mzua"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Thank you for your vote!",
		},
		exp_err: nil,
	},
	{ // test2 - negative, wrong pollid in votereq
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash2,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid:  0,
			Answers: &query.PollSchema{},
			Sign: &query.RSASignature{
				Ballot: []byte("rvrbhd54\":^V(B)*TBytvw.ucq<{_@x-mzua"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Error in PollVote",
		},
		exp_err: fmt.Errorf("Error in PollVote while retrieving key from database: %w", fmt.Errorf("No key for this poll in database.")),
	},
	{ // test3 - negative, wrong pollid in votereq
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid:  1,
			Answers: &query.PollSchema{},
			Sign: &query.RSASignature{
				Ballot: []byte("Some value"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Error in PollVote",
		},
		exp_err: fmt.Errorf("Error in PollVte, Sign invalid!"),
	},
	{ // test4 - negative, wrong characters
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Question\x00\x01",
						Type:     query.PollSchema_OPEN,
					},
				},
			},
			Sign: &query.RSASignature{
				Ballot: []byte("12345678"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Error in PollVote",
		},
		exp_err: fmt.Errorf("Error in PollVote while saving key in database: %w", fmt.Errorf("Error! Question contains invalid characters.")),
	},
	{ // test5 - negative, wrong characters
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Question",
						Type:     query.PollSchema_OPEN,
						Answers:  []string{"\xe2\x03"},
					},
				},
			},
			Sign: &query.RSASignature{
				Ballot: []byte("12345678"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Error in PollVote",
		},
		exp_err: fmt.Errorf("Error in PollVote while saving key in database: %w", fmt.Errorf("Error! Answer contains invalid characters.")),
	},
	{ // test6 - negative, wrong characters
		schema: &query.PollSchema{},
		envelope: &query.EnvelopeToSign{
			Envelope: hash,
			Pollid:   1,
			Token:    "Good token",
		},
		votereq: &query.VoteRequest{
			Pollid: 1,
			Answers: &query.PollSchema{
				Questions: []*query.PollSchema_QA{
					{
						Question: "Question",
						Type:     -1,
						Answers:  []string{"\xe2\x03"},
					},
				},
			},
			Sign: &query.RSASignature{
				Ballot: []byte("12345678"),
			},
		},
		exp_out: &query.VoteReply{
			Mess: "Error in PollVote",
		},
		exp_err: fmt.Errorf("Error in PollVote while saving key in database: %w", fmt.Errorf("Error! Wrong question type.")),
	},
}

var testsEntireProtocol = struct {
	pollreq  *query.GetPollRequest
	schema   *query.PollSchema
	envelope *query.EnvelopeToSign
	votereq  *query.VoteRequest
}{
	pollreq: &query.GetPollRequest{
		Pollid: 1,
	},
	schema: &query.PollSchema{
		Questions: []*query.PollSchema_QA{
			{
				Question: "Question",
				Type:     query.PollSchema_OPEN,
			},
		},
	},
	envelope: &query.EnvelopeToSign{
		Envelope: nil,
		Pollid:   1,
		Token:    "Good token",
	},
	votereq: &query.VoteRequest{
		Pollid: 1,
		Answers: &query.PollSchema{
			Questions: []*query.PollSchema_QA{
				{
					Question: "Question",
					Type:     query.PollSchema_OPEN,
					Answers:  []string{"Answer"},
				},
			},
		},
		Sign: &query.RSASignature{
			Ballot: []byte("12345678"),
		},
	},
}
