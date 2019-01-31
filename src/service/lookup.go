package service

import (
	"encoding/json"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// Lookup implements the person lookups
type Lookup struct {
	get      getter
	uri      string
	response *Response
}

// DuckID creates a Lookup for getting a user from a given duckid
func DuckID(id string) *Lookup {
	return &Lookup{
		get: _getReal,
		uri: strings.Replace(global.Conf.LookupByDuckIDURL, "{{duckid}}", id, -1),
	}
}

// BannerID creates a Lookup for getting a user from a given bannerid
func BannerID(id string) *Lookup {
	return &Lookup{
		get: _getReal,
		uri: strings.Replace(global.Conf.LookupByBannerIDURL, "{{bannerid}}", id, -1),
	}
}

// User is the relevant data returned by the service
type User struct {
	BannerID string `json:"bannerID"`
	DuckID   string `json:"duckID"`
}

// Response contains all data returned by the API call
type Response struct {
	User       User   `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// Call requests data from the service, storing the response and returning any
// errors which prevented the service endpoint from being reached.  Errors
// within the service are not checked here.
func (l *Lookup) Call() error {
	var content, err = l.get(l.uri)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &l.response)
}

// Response returns the result of the most recent service call
func (l *Lookup) Response() *Response {
	return l.response
}
