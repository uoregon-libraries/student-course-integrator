// Package person is responsible for querying ldap to get a person's data, and
// pinging an internal service to get a 95-number for said person
package person

import (
	"fmt"
	"regexp"
	"strings"

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

// Find searches LDAP for the given id (either duckid or bannerid) and returns a Person record
// filled in with the details needed for SCI
func Find(stringID string) (*Person, error) {
	var c, err = connect()
	if err != nil {
		return nil, err
	}
	defer c.lc.Close()

	var s = service.DuckID(stringID)
	if isBannerID(stringID) {
		s = service.BannerID(stringID)
	}
	err = s.Call()
	if err != nil {
		return nil, fmt.Errorf("unable to look up Banner ID for %s: %s", stringID, err)
	}
	var r = s.Response()
	if r.StatusCode == 404 {
		return nil, nil
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("service: status %d looking up %s: %s", r.StatusCode, stringID, r.Message)
	}
	if r.User.BannerID == "" {
		return nil, fmt.Errorf("lookup for duckid %s: response contains empty Banner ID", stringID)
	}
	var p *Person
	p, err = c.find(r.User.DuckID)
	if err != nil {
		return nil, err
	}
	if p != nil {
		p.BannerID = r.User.BannerID
	}
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

// BannerIDRegex pattern for use in match in isBannerID()
var BannerIDRegex = regexp.MustCompile("^95[0-9]{7}$")

// isBannerID returns true if id is a 95 number
func isBannerID(stringID string) bool {
	clean := strings.Replace(stringID, "-", "", -1)
	return BannerIDRegex.MatchString(clean)
}
