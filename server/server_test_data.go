package main

import (
	"fmt"
	"github.com/ememak/Projekt-Rada/query"
)

var testsPollInitIn = []*query.PollSchema{
	&query.PollSchema{ // test0 - positive
		Questions: []*query.PollSchema_QA{
			{
				Question: "Do you like this system? Options: yes/no",
				Type:     query.PollSchema_CLOSE,
				Answer:   "",
			},
			{
				Question: "Why?",
				Type:     query.PollSchema_OPEN,
				Answer:   "",
			},
		},
	},
	&query.PollSchema{ // test1 - negative, wrong type value
		Questions: []*query.PollSchema_QA{
			{
				Question: "Do you like this system? Options: yes/no",
				Type:     query.PollSchema_CLOSE,
				Answer:   "",
			},
			{
				Question: "Why?",
				Type:     3, // wrong type!
				Answer:   "",
			},
		},
	},
	&query.PollSchema{ // test2 - negative, wrong type value
		Questions: []*query.PollSchema_QA{
			{
				Question: "Do you like this system? Options: yes/no",
				Type:     -1, //wrong type!
				Answer:   "",
			},
		},
	},
	&query.PollSchema{ // test3 - positive
		Questions: []*query.PollSchema_QA{
			{
				Question: "Check numbers you like? Options: 1;2;5;e;74",
				Type:     query.PollSchema_CHECKBOX,
				Answer:   "",
			},
		},
	},
	&query.PollSchema{ // test4 - negative, wrong characters in question
		Questions: []*query.PollSchema_QA{
			{
				Question: "\x00\x01\x02\xff\xe7",
				Type:     query.PollSchema_OPEN,
				Answer:   "",
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
						Answer:   "",
					},
					{
						Question: "Why?",
						Type:     query.PollSchema_OPEN,
						Answer:   "",
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
						Answer:   "",
					},
				},
			},
		},
		exp_err: nil,
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Question contains non valid characters.")),
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
						Answer:   "",
					},
					{
						Question: "Why?",
						Type:     query.PollSchema_OPEN,
						Answer:   "",
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
						Answer:   "",
					},
				},
			},
		},
		exp_err: nil,
	},
	{
		exp_out: &query.PollQuestion{},
		exp_err: fmt.Errorf("Error in PollInit while creating new poll in database: %w", fmt.Errorf("Error! Question contains non valid characters.")),
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
					Answer:   "",
				},
				{
					Question: "Why?",
					Type:     query.PollSchema_OPEN,
					Answer:   "",
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
					Answer:   "",
				},
				{
					Question: "Why?",
					Type:     query.PollSchema_OPEN,
					Answer:   "",
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
					Answer:   "",
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
					Answer:   "",
				},
				{
					Question: "Why?",
					Type:     query.PollSchema_OPEN,
					Answer:   "",
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
					Answer:   "",
				},
				{
					Question: "Why?",
					Type:     query.PollSchema_OPEN,
					Answer:   "",
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
