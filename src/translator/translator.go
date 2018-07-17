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

	var u, err = callService(lookupBannerID, duckid)
	return u.BannerID, err
}

// BannerIDToDuckID returns the duckid for the given banner id (95
// number), or an error if the service can't be reached to do the lookup
func BannerIDToDuckID(uid string) (string, error) {
	// Simulate the cost of an API hit
	time.Sleep(time.Millisecond * 50)

	var u, err = callService(lookupDuckID, uid)
	return u.DuckID, err
}

type userJSON struct {
	BannerID string `json:"banner_id"`
	DuckID   string `json:"duckid"`
}

type responseJSON struct {
	User       userJSON `json:"result"`
	Message    string   `json:"message"`
	StatusCode int      `json:"statusCode"`
}

// lookupTypes represent the operation when calling the IS service
type lookupType string

// Valid lookup types - names don't match strings because the name is meant to
// convey the meaning of what you're requesting, not what you have (and the
// string is really the API path endpoint, so may not be meaningful as things
// change anyway)
const (
	lookupBannerID lookupType = "duckid"    // Look up a banner ID from a given duckid
	lookupDuckID   lookupType = "banner_id" // Look up a duckid from a given banner ID
)

// callService is a common wrapper to call the central translation service and
// return a user response.
func callService(lookup lookupType, val string) (user userJSON, err error) {
	var content []byte
	var url = global.Conf.TranslatorHost + "/" + path.Join(string(lookup), val)
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
		return user, fmt.Errorf("translator: service returned non-success status code %d and message %q", r.StatusCode, r.Message)
	}

	user = r.User
	if user.BannerID == "" || user.DuckID == "" {
		return user, fmt.Errorf("invalid response data")
	}
	return user, err
}
