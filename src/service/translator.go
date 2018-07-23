package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// DuckIDToBannerID returns the banner id (95 number) for the given
// duckid, or an error if the service can't be reached to do the lookup
func DuckIDToBannerID(duckid string) (string, error) {
	var u, err = callService(strings.Replace(global.Conf.LookupByDuckIDURL, "{{duckid}}", duckid, -1))
	return u.BannerID, err
}

// BannerIDToDuckID returns the duckid for the given banner id (95
// number), or an error if the service can't be reached to do the lookup
func BannerIDToDuckID(id string) (string, error) {
	var u, err = callService(strings.Replace(global.Conf.LookupByBannerIDURL, "{{bannerid}}", id, -1))
	return u.DuckID, err
}

type userJSON struct {
	BannerID string `json:"bannerID"`
	DuckID   string `json:"duckID"`
}

type responseJSON struct {
	User       userJSON `json:"data"`
	Message    string   `json:"message"`
	StatusCode int      `json:"statusCode"`
}

// callService is a common wrapper to call the central translation service and
// return a user response.
func callService(url string) (user userJSON, err error) {
	var content []byte
	content, err = get(url)
	if err != nil {
		return user, err
	}

	var r responseJSON
	err = json.Unmarshal(content, &r)
	if err != nil {
		return user, err
	}

	if r.StatusCode != 200 {
		return user, fmt.Errorf("translator: service error (%q; status code %d)", r.Message, r.StatusCode)
	}

	user = r.User
	if user.BannerID == "" || user.DuckID == "" {
		return user, fmt.Errorf("invalid response data")
	}
	return user, err
}
