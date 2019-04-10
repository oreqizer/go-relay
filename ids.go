package relay

import (
	"encoding/base64"
	"strings"
)

type Local struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ToGlobalID creates a globally unique ID
func ToGlobalID(ttype string, id string) string {
	return base64.StdEncoding.EncodeToString([]byte(ttype + ":" + id))
}

// FromGlobalID splits the global ID into a type and the original ID
func FromGlobalID(id string) *Local {
	bytes, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil
	}

	tokens := strings.Split(string(bytes), ":")
	if len(tokens) != 2 {
		return nil
	}

	return &Local{
		Type: tokens[0],
		ID:   tokens[1],
	}
}
