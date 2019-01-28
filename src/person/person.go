// Package person is responsible for querying ldap to get a person's data, and
// pinging an internal service to get a 95-number for said person
package person

import (
	"fmt"
    "strings"
    "regexp"

	"github.com/uoregon-libraries/student-course-integrator/src/service"
)

// A Person is different from a user in that this represents anybody at UO, not
// just those who can log into the application.  This is currently not in any
// local database; we query LDAP and a custom duckid-to-95x endpoint to get the
// data we need to populate this structure.
type Person struct {
	BannerID     string // "95 number"
	DuckID       string
	Affiliations []string // Affiliations: things like faculty, staff, GTF, etc
	DisplayName  string   // Banner's display name for the individual
}

// FindByDuckID searches LDAP for the given duckid and returns a Person record
// filled in with the details needed for SCI
func FindByDuckID(duckid string) (*Person, error) {
	var c, err = connect()
	if err != nil {
		return nil, err
	}
	defer c.lc.Close()

    var s = service.DuckID(duckid)
    if IsBannerID(duckid) {
        s = service.BannerID(duckid)
    }
	var serviceErr = s.Call()
	if serviceErr != nil {
		return nil, fmt.Errorf("unable to look up Banner ID for %s: %s", duckid, serviceErr)
	}
	var r = s.Response
	if r.StatusCode == 404 {
		return nil, nil
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("service: status %d looking up %s: %s", r.StatusCode, duckid, r.Message)
	}
	if r.User.BannerID == "" {
		return nil, fmt.Errorf("lookup for duckid %s: response contains empty Banner ID", duckid)
	}
	var p *Person
	p, err = c.find(r.User.DuckID)
	if err != nil {
		return nil, err
	}
	p.BannerID = r.User.BannerID
	return p, nil

}

// validGEAffiliations stores our list of which UO LDAP affiliations are valid
// for determining if somebody is allowed to be assigned as a GE
var validGEAffiliations = map[string]bool{
	"gtf": true,
}

// CanBeGE returns true if this person's affiliations allow being assigned as a
// GE on a course
func (p *Person) CanBeGE() bool {
	for _, aff := range p.Affiliations {
		if validGEAffiliations[aff] {
			return true
		}
	}
	return false
}

// IsBannerID returns true if id is a 9 digit number
func IsBannerID(stringId string) bool {
    clean := strings.Replace(stringId, "-", "", -1)
    re := regexp.MustCompile("[0-9]{9}")
    return re.MatchString(clean)
}