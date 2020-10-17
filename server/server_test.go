package main

import (
	"context"
	"reflect"
	"strconv"
	"testing"
)

func TestPollInit(t *testing.T) {
	in := testsPollInitIn
	out := testsPollInitOutEmpty
	for i, test := range in {

		s, _ := serverInit("test" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			ctx := context.Background()
			o, err := s.PollInit(ctx, test)
			if !(reflect.DeepEqual(o, out[i].exp_out) && reflect.DeepEqual(err, out[i].exp_err)) {
				t.Errorf("Output %v, want output %v", o, out[i].exp_out)
				t.Errorf("Error %d, want error %d", err, out[i].exp_err)
			}
		})
		s.data.Close()
	}
	out = testsPollInitOut1Poll
	for i, test := range in {

		s, _ := serverInit("test" + strconv.Itoa(len(in)+i) + ".db")
		t.Run("Test "+strconv.Itoa(len(in)+i), func(t *testing.T) {
			ctx := context.Background()
			s.PollInit(ctx, test)
			o, err := s.PollInit(ctx, test)
			if !(reflect.DeepEqual(o, out[i].exp_out) && reflect.DeepEqual(err, out[i].exp_err)) {
				t.Errorf("Output %v, want output %v", o, out[i].exp_out)
				t.Errorf("Error %d, want error %d", err, out[i].exp_err)
			}
		})
	}
}
