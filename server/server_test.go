package main

import (
	"context"
	"crypto/x509"
	"reflect"
	"strconv"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestPollInit(t *testing.T) {
	in := testsPollInitIn
	out := testsPollInitOutEmpty
	for i, test := range in {

		s, _ := serverInit("testPI" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			ctx := context.Background()
			o, err := s.PollInit(ctx, test)
			if !(proto.Equal(o, out[i].exp_out) && reflect.DeepEqual(err, out[i].exp_err)) {
				t.Errorf("Output %v, want output %v", o, out[i].exp_out)
				t.Errorf("Error %v, want error %v", err, out[i].exp_err)
			}
		})
		s.data.Close()
	}
	out = testsPollInitOut1Poll
	for i, test := range in {

		s, _ := serverInit("testPI" + strconv.Itoa(len(in)+i) + ".db")
		t.Run("Test "+strconv.Itoa(len(in)+i), func(t *testing.T) {
			ctx := context.Background()
			s.PollInit(ctx, test)
			poll, err := s.PollInit(ctx, test)
			if !(proto.Equal(poll, out[i].exp_out) && reflect.DeepEqual(err, out[i].exp_err)) {
				t.Errorf("Output %v, want output %v", poll, out[i].exp_out)
				t.Errorf("Error %v, want error %v", err, out[i].exp_err)
			}
		})
	}
}

func TestGetPoll(t *testing.T) {
	in := testsGetPollIn
	out := testsGetPollOut
	for i, test := range in {

		s, _ := serverInit("testGP" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			ctx := context.Background()
			s.PollInit(ctx, test.schema)
			pwk, err := s.GetPoll(ctx, test.pollreq)
			if !(proto.Equal(pwk.Poll, out[i].exp_out) && reflect.DeepEqual(err, out[i].exp_err)) {
				t.Errorf("Output %v, want output %v", pwk.Poll, out[i].exp_out)
				t.Errorf("Error %v, want error %v", err, out[i].exp_err)
				return
			}
			if err != nil {
				return
			}

			_, err = x509.ParsePKCS1PublicKey(pwk.Key.Key)
			if err != nil {
				t.Errorf("Key parsing. Error %v, want nil error", err)
				return
			}
		})
		s.data.Close()
	}
}
