package translator

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
)

// get aliases the getter function.  TODO: this alias should only be the real
// getter except when testing.
var get = _getMock

// _getReal is a simple wrapper for REST API Get hits
func _getReal(url string) (content []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// mockResponseTemplate MUST match the actual service's JSON if we want to
// trust the mock service call
var mockResponseTemplate = `
	{
		"result": {
			"banner_id": "95x00000%d",
			"duckid": "tester%d"
		},
		"message": "%s",
		"statusCode": %d
	}
`

// _getMock is a simple wrapper for fake REST API Get hits
func _getMock(url string) (content []byte, err error) {
	var idnum = crc32.ChecksumIEEE([]byte(url))%uint32(5) + 1
	var response = fmt.Sprintf(mockResponseTemplate, idnum, idnum, "no message", 200)
	return []byte(response), nil
}
