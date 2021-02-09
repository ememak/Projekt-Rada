package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"math/big"
	"reflect"
	"strconv"
	"testing"

	"github.com/ememak/Projekt-Rada/store"
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

func TestSignBallot(t *testing.T) {
	in := testsSignBallot
	for i, test := range in {

		s, _ := serverInit("testSB" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			ctx := context.Background()
			s.PollInit(ctx, test.schema)
			store.SaveToken(s.data, "Good token", 1)
			se, err := s.SignBallot(ctx, test.envelope)
			if !reflect.DeepEqual(err, test.exp_err) {
				t.Errorf("Error %v, want error %v", err, test.exp_err)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(se.Envelope, test.envelope.Envelope) {
				t.Errorf("Envelope %v, want envelope %v", se.Envelope, test.envelope.Envelope)
			}

		})
		s.data.Close()
	}
}

func TestPollVote(t *testing.T) {
	in := testsPollVote
	for i, test := range in {

		s, _ := serverInit("testPV" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			ctx := context.Background()
			s.PollInit(ctx, test.schema)
			store.SaveToken(s.data, "Good token", 1)
			se, _ := s.SignBallot(ctx, test.envelope)
			test.votereq.Sign.Sign = se.Sign
			vr, err := s.PollVote(ctx, test.votereq)
			if !(proto.Equal(vr, test.exp_out) && reflect.DeepEqual(err, test.exp_err)) {
				t.Errorf("Output %v, want output %v", vr, test.exp_out)
				t.Errorf("Error %v, want error %v", err, test.exp_err)
				return
			}
		})
		s.data.Close()
	}
}

func TestEntireProtocol(t *testing.T) {
	test := testsEntireProtocol
	t.Run("Full Test", func(t *testing.T) {
		s, err := serverInit("testAP.db")
		if err != nil {
			t.Errorf("ServerInit failed, error: %v", err)
			return
		}
		ctx := context.Background()
		_, err = s.PollInit(ctx, test.schema)
		if err != nil {
			t.Errorf("PollInit failed, error: %v", err)
			return
		}
		pwk, err := s.GetPoll(ctx, test.pollreq)
		if err != nil {
			t.Errorf("GetPoll failed, error: %v", err)
			return
		}
		key, err := x509.ParsePKCS1PublicKey(pwk.Key.Key)
		if err != nil {
			t.Errorf("Key parsing failed, error: %v", err)
			return
		}
		err = store.SaveToken(s.data, "Good token", 1)
		if err != nil {
			t.Errorf("SaveToken failed, error: %v", err)
			return
		}
		// Generate ballot to be signed.
		ballot, err := rand.Int(rand.Reader, key.N)
		if err != nil {
			t.Errorf("Rand.Int failed, error: %v", err)
			return
		}
		// We are hashing ballot.
		hash := sha256.Sum256(ballot.Bytes())
		m := new(big.Int).SetBytes(hash[:])
		// Get random blinding factor.
		r, err := rand.Int(rand.Reader, key.N)
		if err != nil {
			t.Errorf("Rand.Int failed, error: %v", err)
			return
		}
		// We want to send m*r^e mod N to server.
		bfactor := new(big.Int).Exp(r, big.NewInt(int64(key.E)), key.N)
		// blinded = m*(r^e) mod N
		blinded := bfactor.Mod(bfactor.Mul(bfactor, m), key.N)
		test.envelope.Envelope = blinded.Bytes()

		se, err := s.SignBallot(ctx, test.envelope)
		if err != nil {
			t.Errorf("SignBallot failed, error: %v", err)
			return
		}

		// Having (m^d)*r mod N we are removing blinding factor r,
		smi := new(big.Int).SetBytes(se.Sign)
		revr := new(big.Int).ModInverse(r, key.N)
		smirevr := new(big.Int).Mul(revr, smi)
		// Now we can calculate second part of sign.
		// sign = smirevr mod N = m^d mod N
		sign := new(big.Int).Mod(smirevr, key.N)
		test.votereq.Sign.Ballot = ballot.Bytes()
		test.votereq.Sign.Sign = sign.Bytes()
		_, err = s.PollVote(ctx, test.votereq)
		if err != nil {
			t.Errorf("PollVote failed, error: %v", err)
			return
		}
	})
}
