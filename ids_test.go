package relay_test

import (
	"encoding/base64"
	"github.com/oreqizer/go-relay"
	"testing"
)

var tableIDs = []struct {
	ttype  string
	id     string
	global string
}{
	{ttype: "Kek", id: "1337", global: "Kek:1337"},
	{ttype: "Lol", id: "420", global: "Lol:420"},
	{ttype: "Bur", id: "69", global: "Bur:69"},
}

func TestToGlobalID(t *testing.T) {
	for _, e := range tableIDs {
		out := relay.ToGlobalID(e.ttype, e.id)
		bytes, err := base64.StdEncoding.DecodeString(out)
		if err != nil {
			t.Error(err)
		}

		if got := string(bytes); got != e.global {
			t.Errorf("Got %s, want %s", got, e.global)
		}
	}
}

func TestFromGlobalID(t *testing.T) {
	if out := relay.FromGlobalID(base64.StdEncoding.EncodeToString([]byte("Lol:Kek:Bur"))); out != nil {
		t.Error("Expected nil for multiple separators")
	}

	if out := relay.FromGlobalID("XXXXXaGVsbG8="); out != nil {
		t.Error("Expected nil for corrupt input")
	}

	for _, e := range tableIDs {
		out := relay.FromGlobalID(base64.StdEncoding.EncodeToString([]byte(e.global)))
		if out == nil {
			t.Error("Unexpected nil")
			return
		}

		if out.Type != e.ttype {
			t.Errorf("Got %s, want %s", out.Type, e.ttype)
		}

		if out.ID != e.id {
			t.Errorf("Got %s, want %s", out.ID, e.id)
		}
	}
}
