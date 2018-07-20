package translator

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// DuckIDToBannerID returns the banner id (95 number) for the given
// duckid, or an error if the service can't be reached to do the lookup
func DuckIDToBannerID(duckid string) (string, error) {
	// Simulate the cost of an API hit
	time.Sleep(time.Millisecond * 50)

	var u, err = callService(lookupByDuckID, duckid)
	return u.BannerID, err
}

// BannerIDToDuckID returns the duckid for the given banner id (95
// number), or an error if the service can't be reached to do the lookup
func BannerIDToDuckID(uid string) (string, error) {
	// Simulate the cost of an API hit
	time.Sleep(time.Millisecond * 50)

	var u, err = callService(lookupByBannerID, uid)
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

// lookupTypes represent the operation when calling the IS service
type lookupType string

// valid lookup types - the value indicates how a person is being looked up
const (
	// TODO: change this if the service endpoint swap problem gets fixed
	//lookupByDuckID   lookupType = "duckid"
	//lookupByBannerID lookupType = "banner_id"
	lookupByDuckID   lookupType = "banner_id"
	lookupByBannerID lookupType = "duckid"
)

// callService is a common wrapper to call the central translation service and
// return a user response.
func callService(lookup lookupType, val string) (user userJSON, err error) {
	var content []byte
	var url = global.Conf.TranslatorHost + "/" + path.Join("person", "uo", string(lookup), val)
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
