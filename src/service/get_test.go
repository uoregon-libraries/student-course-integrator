package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/uoregon-libraries/gopkg/assert"
)

func TestMockJSON(t *testing.T) {
	var r responseJSON
	var jstr = fmt.Sprintf(mockResponseTemplate, 1, 2, "no message", 200)
	var err = json.Unmarshal([]byte(jstr), &r)
	if err != nil {
		switch e := err.(type) {
		case *json.SyntaxError:
			t.Fatalf("Syntax error in JSON: %#v", e)
		default:
			t.Fatalf("Error parsing JSON: %s", e)
		}
	}

	assert.Equal("95x000001", r.User.BannerID, "parsed banner id", t)
	assert.Equal("tester2", r.User.DuckID, "parsed duck id", t)
	assert.Equal("no message", r.Message, "parsed message", t)
	assert.Equal(200, r.StatusCode, "parsed status code", t)
}
