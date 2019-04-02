package person

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/uoregon-libraries/gopkg/assert"
	"github.com/uoregon-libraries/student-course-integrator/src/roles"
	"github.com/uoregon-libraries/student-course-integrator/src/service"
	ldap "gopkg.in/ldap.v2"
)

type FakeLookup struct {
	User       service.User
	Message    string
	StatusCode int
	err        string
}

func (f *FakeLookup) Call() error {
	if f.err != "" {
		return errors.New(f.err)
	}
	return nil
}

func (f *FakeLookup) Response() *service.Response {
	return &service.Response{User: f.User, Message: f.Message, StatusCode: f.StatusCode}
}

type FakeLdap struct {
	dn         string
	attributes []*ldap.EntryAttribute
	err        string
}

func (f *FakeLdap) Search(id string) (*ldap.Entry, error) {
	if f.err != "" {
		return nil, errors.New(f.err)
	}
	return &ldap.Entry{DN: f.dn, Attributes: f.attributes}, nil
}

var attrs = []*ldap.EntryAttribute{
	{Name: "displayName", Values: []string{"Sam Smith"}, ByteValues: nil},
	{Name: "UOAD-UoPersonAffiliation", Values: []string{"gtf"}, ByteValues: nil},
}

type testvars struct {
	bannerID      string
	duckID        string
	status        int
	lookupMessage string
	lookupErr     string
}

func initTvars() testvars {
	return testvars{"950123456", "ssmith", 200, "", ""}
}

var c = FakeLdap{initTvars().duckID, attrs, ""}
var user = service.User{BannerID: initTvars().bannerID, DuckID: initTvars().duckID}

func TestFindSuccess(t *testing.T) {
	var tvars = initTvars()
	var s = FakeLookup{User: user, Message: tvars.lookupMessage, StatusCode: tvars.status, err: tvars.lookupErr}
	var response, err = find(tvars.duckID, &s, &c)
	assert.True(response != nil, "returns a person", t)
	assert.True(err == nil, "no errors", t)
}

func TestFindError(t *testing.T) {
	var tvars = initTvars()
	tvars.lookupErr = "something is busted"
	var s = FakeLookup{User: user, Message: tvars.lookupMessage, StatusCode: tvars.status, err: tvars.lookupErr}
	var response, err = find(tvars.duckID, &s, &c)
	assert.True(response == nil, "should return nil", t)
	assert.True(strings.Contains(err.Error(), "unable to look up Banner ID"), "expected error text", t)
	assert.True(strings.Contains(err.Error(), tvars.duckID), "expected error text", t)
	assert.True(strings.Contains(err.Error(), tvars.lookupErr), "expected error text", t)
}

func TestFind404(t *testing.T) {
	var tvars = initTvars()
	tvars.status = 404
	var s = FakeLookup{User: user, Message: tvars.lookupMessage, StatusCode: tvars.status, err: tvars.lookupErr}
	var response, err = find(tvars.duckID, &s, &c)
	assert.True(response == nil, "should return nil", t)
	assert.True(err == nil, "but not an error", t)
}

func TestFindNot200(t *testing.T) {
	var tvars = initTvars()
	tvars.status = 418
	tvars.lookupMessage = "time flies"
	var s = FakeLookup{User: user, Message: tvars.lookupMessage, StatusCode: tvars.status, err: tvars.lookupErr}
	var response, err = find(tvars.duckID, &s, &c)
	assert.True(response == nil, "should return nil", t)
	assert.True(strings.Contains(err.Error(), "service: status "+strconv.Itoa(tvars.status)), "expected error text", t)
	assert.True(strings.Contains(err.Error(), tvars.lookupMessage), "expected error text", t)
}

func TestFindNoBannerID(t *testing.T) {
	var tvars = initTvars()
	tvars.bannerID = ""
	var user = service.User{BannerID: tvars.bannerID, DuckID: tvars.duckID}
	var s = FakeLookup{User: user, Message: tvars.lookupMessage, StatusCode: tvars.status, err: tvars.lookupErr}
	var response, err = find(tvars.duckID, &s, &c)
	assert.True(response == nil, "find response should be nil", t)
	assert.True(strings.Contains(err.Error(), "lookup for duckid "+tvars.duckID), "expected error text", t)
}

func TestFindSearchFail(t *testing.T) {
	var tvars = initTvars()
	var ldapErr = "something went wrong"
	var s = FakeLookup{User: user, Message: tvars.lookupMessage, StatusCode: tvars.status, err: tvars.lookupErr}
	var c = FakeLdap{tvars.duckID, attrs, ldapErr}
	var response, err = find(tvars.duckID, &s, &c)
	assert.True(response == nil, "response should be nil", t)
	assert.True(strings.Contains(err.Error(), ldapErr), "expected error text", t)
}

func TestCanBeRoleTA(t *testing.T) {
	var pTA = Person{"950123456", "ssmith", []string{"gtf"}, "Sam Smith"}
	var resultTA = pTA.CanBeRole(roles.TA)
	assert.True(resultTA, "person can be a TA", t)
}

func TestCanBeRoleGrader(t *testing.T) {
	var pGr = Person{"950123456", "ssmith", []string{""}, "Sam Smith"}
	var resultGr = pGr.CanBeRole(roles.Grader)
	assert.True(resultGr, "person can be a Grader", t)
}

func TestCanBeRoleTAFail(t *testing.T) {
	var pTA = Person{"950123456", "ssmith", []string{""}, "Sam Smith"}
	var resultTA = pTA.CanBeRole(roles.TA)
	assert.False(resultTA, "person is not a TA", t)
}

func TestIsBannerID(t *testing.T) {
	var idShort = "950123"
	var idNon95 = "123456789"
	var id95 = "950123456"
	assert.False(isBannerID(idShort), "not a BannerID", t)
	assert.False(isBannerID(idNon95), "not a BannerID", t)
	assert.True(isBannerID(id95), "is a BannerID", t)
}
