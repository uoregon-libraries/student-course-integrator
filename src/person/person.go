// Package person is responsible for querying ldap to get a person's data, and
// pinging an internal service to get a 95-number for said person
package person

import (
	"fmt"

	"github.com/uoregon-libraries/student-course-integrator/src/translator"
)

// A Person is different from a user in that this represents anybody at UO, not
// just those who can log into the application.  This is currently not in any
// local database; we query LDAP and a custom duckid-to-95x endpoint to get the
// data we need to populate this structure.
type Person struct {
	UniversityID string // "95 number"
	DuckID       string
	Affiliations []string
	DisplayName  string
}

// FindByDuckID searches LDAP for the given duckid and returns a Person record
// filled in with the details needed for SCI
func FindByDuckID(duckid string) (*Person, error) {
	var c, err = connect()
	if err != nil {
		return nil, err
	}
	defer c.lc.Close()

	var p *Person
	p, err = c.find(duckid)
	if err != nil {
		return nil, err
	}

	if p != nil {
		p.UniversityID, err = translator.DuckIDToUniversityID(p.DuckID)
		if err != nil {
			return nil, fmt.Errorf("unable to look up university id for duckid %s: %s", p.DuckID, err)
		}
	}
	return p, nil
}

// IsGTF returns true if this person's affiliations allow being assigned as
// a GTF on a course
func (p *Person) IsGTF() bool {
	var validAffiliations = map[string]bool{
		"gtf": true,
	}
	for _, aff := range p.Affiliations {
		if validAffiliations[aff] {
			return true
		}
	}
	return false
}
