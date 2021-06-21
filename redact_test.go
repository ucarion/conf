package conf_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ucarion/conf"
)

func ExampleRedact() {
	cfg := struct {
		Username      string
		Password      string `redact:"true"`
		InnocentBytes []byte `redact:"false"`
		SecretBytes   []byte `redact:"1"`
	}{
		Username:      "jdoe",
		Password:      "iloveyou",
		InnocentBytes: []byte{1, 2, 3},
		SecretBytes:   []byte{4, 5, 6},
	}

	fmt.Println(conf.Redact(cfg)) // returns a deep copy with redacted fields set to zero
	fmt.Println(cfg)              // original copy unmodified

	// Output:
	// {jdoe  [1 2 3] []}
	// {jdoe iloveyou [1 2 3] [4 5 6]}
}

func TestRedact(t *testing.T) {
	type subConfig struct {
		S1 string `redact:"true"`
		S2 string `redact:"t"`
		S3 string `redact:"false"`
		S4 string `redact:"f"`
		S5 string
	}

	type config struct {
		S1 string `redact:"true"`
		S2 string `redact:"t"`
		S3 string `redact:"false"`
		S4 string `redact:"f"`
		S5 string

		Bool      bool          `redact:"true"`
		Int       int           `redact:"true"`
		Array     [64]byte      `redact:"true"`
		Chan      chan string   `redact:"true"`
		Func      func()        `redact:"true"`
		Interface interface{}   `redact:"true"`
		Map       map[bool]bool `redact:"true"`

		StringPtr *string      `redact:"true"`
		ChanPtr   *chan string `redact:"true"`

		SubConfig       subConfig
		SubConfigRedact subConfig `redact:"true"`

		unexported string
	}

	s := "a"
	c := make(chan string)
	got := conf.Redact(config{
		S1:        "a",
		S2:        "a",
		S3:        "a",
		S4:        "a",
		S5:        "a",
		Bool:      true,
		Int:       1,
		Array:     [64]byte{1},
		Chan:      make(chan string),
		Func:      func() {},
		Interface: 1,
		Map:       map[bool]bool{true: true},
		StringPtr: &s,
		ChanPtr:   &c,
		SubConfig: subConfig{
			S1: "a",
			S2: "a",
			S3: "a",
			S4: "a",
			S5: "a",
		},
		SubConfigRedact: subConfig{
			S1: "a",
			S2: "a",
			S3: "a",
			S4: "a",
			S5: "a",
		},
	})

	want := config{
		S1: "",
		S2: "",
		S3: "a",
		S4: "a",
		S5: "a",
		SubConfig: subConfig{
			S1: "",
			S2: "",
			S3: "a",
			S4: "a",
			S5: "a",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want != got, want: %#v, got: %#v", want, got)
	}
}

func TestRedact_Panic_Invalid_Tag_Value(t *testing.T) {
	defer func() {
		r := recover()
		if r.(error).Error() != "conf: error parsing redact tag: strconv.ParseBool: parsing \"notbool\": invalid syntax" {
			t.Fatalf("failed to panic with expected error: %v", r)
		}
	}()

	type config struct {
		S string `redact:"notbool"`
	}

	conf.Redact(config{})
}

func TestRedact_Panic_Invalid_Kind(t *testing.T) {
	defer func() {
		r := recover()
		if r.(error).Error() != "conf: Redact called on ptr (only structs are acceptable)" {
			t.Fatalf("failed to panic with expected error: %v", r)
		}
	}()

	conf.Redact(&struct{}{})
}
