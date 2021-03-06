// Package person is responsible for querying ldap to get a person's data, and
// pinging an internal service to get a 95-number for said person
package person

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/ldapservice"
	"github.com/uoregon-libraries/student-course-integrator/src/roles"
	"github.com/uoregon-libraries/student-course-integrator/src/service"
	ldap "gopkg.in/ldap.v2"
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

// serviceCaller is the interface used for calling the UO services which tell
// us about a user's two main IDs
type serviceCaller interface {
	Call() error
	Response() *service.Response
}

type ldapSearcher interface {
	Search(string) (*ldap.Entry, error)
}

// Find searches LDAP for the given id (either duckid or bannerid) and returns a Person record
// filled in with the details needed for SCI
func Find(stringID string) (*Person, error) {
	var c, err = ldapservice.Connect()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	var s = service.DuckID(stringID)
	if isBannerID(stringID) {
		s = service.BannerID(stringID)
	}

	return find(stringID, s, c)
}

func find(stringID string, s serviceCaller, ls ldapSearcher) (*Person, error) {
	var err = s.Call()
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

	var entry *ldap.Entry
	entry, err = ls.Search(r.User.DuckID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &Person{
		BannerID:     r.User.BannerID,
		DuckID:       r.User.DuckID,
		DisplayName:  entry.GetAttributeValue("displayName"),
		Affiliations: entry.GetAttributeValues("UOAD-UoPersonAffiliation"),
	}, nil
}

// CanBeRole returns true if this person's affiliations allow being assigned
// the given role on a course
func (p *Person) CanBeRole(role string) bool {
	switch role {
	case roles.TA:
		return p.canBeTA()
	case roles.Grader:
		return p.canBeGrader()
	}

	return false
}

// validTAAffiliations stores our list of which UO LDAP affiliations are valid
// for determining if somebody is allowed to be assigned as a TA
var validTAAffiliations = map[string]bool{
	"gtf": true,
}

// canBeTA is true if the person has at least one affiliation from our
// hard-coded map for TAs
func (p *Person) canBeTA() bool {
	for _, aff := range p.Affiliations {
		if validTAAffiliations[aff] {
			return true
		}
	}
	return false
}

// canBeGrader is always true, as the logic for permitting this is that the
// faculty member must have a form on file: not something we can validate with
// code.  If this changes in any way, at least we have the placeholder here to
// add logic / restrictions.
func (p *Person) canBeGrader() bool {
	return true
}

// BannerIDRegex pattern for use in match in isBannerID()
var BannerIDRegex = regexp.MustCompile("^95[0-9]{7}$")

// isBannerID returns true if id is a 95 number
func isBannerID(stringID string) bool {
	clean := strings.Replace(stringID, "-", "", -1)
	return BannerIDRegex.MatchString(clean)
}
