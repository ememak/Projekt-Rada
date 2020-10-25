package main

import (
	"fmt"
	"github.com/ememak/Projekt-Rada/query"
)

var testsPollInitIn = []*query.PollSchema{
	&query.PollSchema{ // test1 - positive
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
	&query.PollSchema{ // test2 - negative, wrong type value
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
	&query.PollSchema{ // test3 - negative, wrong type value
		Questions: []*query.PollSchema_QA{
			{
				Question: "Do you like this system? Options: yes/no",
				Type:     -1, //wrong type!
				Answer:   "",
			},
		},
	},
	&query.PollSchema{ // test4 - positive
		Questions: []*query.PollSchema_QA{
			{
				Question: "Check numbers you like? Options: 1;2;5;e;74",
				Type:     query.PollSchema_CHECKBOX,
				Answer:   "",
			},
		},
	},
	&query.PollSchema{ // test5 - negative, wrong characters in question
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
