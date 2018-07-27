package service

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

type getter func(url string) (content []byte, err error)

func applyHeaders(req *http.Request) error {
	var h = global.Conf.TranslatorAPIHeaders
	var list = strings.Split(h, "\x1e")
	for _, fv := range list {
		var parts = strings.SplitN(fv, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header declaration %q", fv)
		}
		req.Header.Set(parts[0], parts[1])
	}

	return nil
}

// _getReal is a simple wrapper for REST API Get hits
func _getReal(url string) (content []byte, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	applyHeaders(req)

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
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
		"data": {
			"bannerID": "95x00000%d",
			"duckID": "tester%d"
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
