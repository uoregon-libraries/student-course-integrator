package translator

import (
	"hash/crc32"
)

// get is a simple wrapper for REST API Get hits
func get(url string) (content []byte, err error) {
	var responses = [][]byte{
		[]byte(`{"banner_id": "95x000001", "duckid": "tester1", "message": "", "code": 0}`),
		[]byte(`{"banner_id": "95x000002", "duckid": "tester2", "message": "", "code": 0}`),
		[]byte(`{"banner_id": "95x000003", "duckid": "tester3", "message": "", "code": 0}`),
		[]byte(`{"banner_id": "95x000004", "duckid": "tester4", "message": "", "code": 0}`),
		[]byte(`{"banner_id": "95x000005", "duckid": "tester5", "message": "", "code": 0}`),
	}
	return responses[crc32.ChecksumIEEE([]byte(url))%uint32(len(responses))], nil
}
