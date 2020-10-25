package store

import (
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"strconv"
	"testing"
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

		data, _ := DBInit("test" + strconv.Itoa(i) + ".db")
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
			o, err := GetPoll(data, int(p.Id))
			if !(reflect.DeepEqual(o.Schema, test.in) && reflect.DeepEqual(err, nil)) {
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

		data, _ := DBInit("test" + strconv.Itoa(i) + ".db")
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			err := SaveKey(data, i, key)
			if err != nil {
				t.Errorf("Error %v, want nil error", err)
				return
			}
			keyret, err := GetKey(data, i)
			if err != nil || !reflect.DeepEqual(keyret, key) {
				t.Errorf("Output %v, want output %v", keyret, key)
				t.Errorf("Error %v, want nil error", err)
			}
		})
	}
	tests := testsSaveKey
	for i, test := range tests {
		data, _ := DBInit("test" + strconv.Itoa(i+3) + ".db")

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
