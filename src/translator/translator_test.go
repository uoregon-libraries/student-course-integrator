package translator

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/uoregon-libraries/gopkg/assert"
	"github.com/uoregon-libraries/student-course-integrator/src/config"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

func TestMain(m *testing.M) {
	global.Conf = &config.Config{TranslatorHost: "https://example.org"}
	os.Exit(m.Run())
}

func TestCallServiceSuccess(t *testing.T) {
	get = _getMock
	var bannerID, err = DuckIDToBannerID("test")
	assert.NilError(err, "DuckIDToBannerID", t)
	assert.True(bannerID != "", "bannerID isn't empty", t)

	var duckID string
	duckID, err = BannerIDToDuckID("test")
	assert.NilError(err, "BannerIDToDuckID", t)
	assert.True(duckID != "", "duckID isn't empty", t)
}

func TestCallServiceFailure(t *testing.T) {
	get = func(url string) (content []byte, err error) {
		return []byte(fmt.Sprintf(mockResponseTemplate, 0, 0, "baby fall down go boom", 500)), nil
	}

	var _, err = DuckIDToBannerID("test")
	assert.True(err != nil, "500 status code from json should produce an error", t)
	assert.True(strings.Contains(err.Error(), "baby fall down go boom"), "expected error text", t)
}

func TestCallServiceError(t *testing.T) {
	get = func(url string) (content []byte, err error) {
		return nil, errors.New("foo")
	}
	var _, err = DuckIDToBannerID("test")
	assert.True(err != nil, "explicit error return should propagate", t)
	assert.Equal("foo", err.Error(), "expected error text", t)
}