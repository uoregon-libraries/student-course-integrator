package translator

import (
	"fmt"
	"hash/crc32"
)

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

// get is a simple wrapper for fake REST API Get hits
func get(url string) (content []byte, err error) {
	var idnum = crc32.ChecksumIEEE([]byte(url))%uint32(5) + 1
	var response = fmt.Sprintf(mockResponseTemplate, idnum, idnum, "no message", 200)
	return []byte(response), nil
}
