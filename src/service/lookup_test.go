package service

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/uoregon-libraries/gopkg/assert"
	"github.com/uoregon-libraries/student-course-integrator/src/config"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

func TestMain(m *testing.M) {
	global.Conf = &config.Config{
		LookupByBannerIDURL: "https://example.org/bannerid={{bannerid}}",
		LookupByDuckIDURL:   "https://example.org/duckid={{duckid}}",
	}
	os.Exit(m.Run())
}

func TestCallServiceSuccess(t *testing.T) {
	var s = DuckID("test")
	s.get = _getMock
	var err = s.Call()
	assert.NilError(err, "DuckIDToBannerID", t)
	assert.True(s.Response().User.BannerID != "", "bannerID isn't empty", t)

	var s2 = BannerID("test")
	s2.get = _getMock
	err = s2.Call()
	assert.NilError(err, "BannerIDToDuckID", t)
	assert.True(s2.Response().User.DuckID != "", "duckID isn't empty", t)
}

func TestCallServiceFailure(t *testing.T) {
	var s = DuckID("test")
	s.get = func(url string) (content []byte, err error) {
		return []byte(fmt.Sprintf(mockResponseTemplate, 0, 0, "baby fall down go boom", 500)), nil
	}
	var err = s.Call()

	assert.NilError(err, "Service failure shouldn't be an error", t)
	assert.Equal(500, s.Response().StatusCode, "service status", t)
	assert.Equal("baby fall down go boom", s.Response().Message, "service message", t)
}

func TestCallServiceError(t *testing.T) {
	var s = DuckID("test")
	s.get = func(url string) (content []byte, err error) {
		return nil, errors.New("foo")
	}
	var err = s.Call()
	assert.True(err != nil, "explicit error return should propagate", t)
	assert.Equal("foo", err.Error(), "expected error text", t)
}

func translatefn(translated *string) getter {
	return func(url string) (content []byte, err error) {
		*translated = url
		return _getMock(url)
	}
}

func TestTranslatedURIDuckIDToBannerID(t *testing.T) {
	var translated string
	var s = DuckID("jechols")
	s.get = translatefn(&translated)
	s.Call()
	assert.Equal("https://example.org/duckid=jechols", translated, "translated URI for duckid-to-bannerid", t)
}

func TestTranslatedURIBannerIDToDuckID(t *testing.T) {
	var translated string
	var s = BannerID("95x000000")
	s.get = translatefn(&translated)
	s.Call()
	assert.Equal("https://example.org/bannerid=95x000000", translated, "translated URI for bannerid-to-duckid", t)
}
