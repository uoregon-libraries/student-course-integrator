package ldapservice

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/uoregon-libraries/student-course-integrator/src/global"
	ldap "gopkg.in/ldap.v2"
)

// Connection wraps ldap.Conn to provide search functionality
type Connection struct {
	lc   *ldap.Conn
	open bool
}

// Connect starts a connection to the configured LDAP endpoint, which must
// support TLS, then binds with the bind name / pass
func Connect() (*Connection, error) {
	var l, err = ldap.Dial("tcp", global.Conf.LDAPServer)
	if err != nil {
		return nil, err
	}

	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		l.Close()
		return nil, err
	}

	err = l.Bind(global.Conf.LDAPUser, global.Conf.LDAPPass)
	if err != nil {
		l.Close()
		return nil, err
	}

	return &Connection{lc: l, open: true}, nil
}

// Search runs a search against LDAP for the given duckid, returning the raw
// LDAP record on success.  If LDAP fails, or the search finds more than one
// person, it is considered unsuccessful and an error is returned.  If nobody
// is found, a nil record is returned, but this is not considered an error,
// just a nonexistent duckid.
func (c *Connection) Search(duckid string) (*ldap.Entry, error) {
	// NewSearchRequest is way more painful to read than just instantiating the
	// thing with named values
	var req = &ldap.SearchRequest{
		BaseDN:       global.Conf.LDAPBaseDN,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		SizeLimit:    0,
		TimeLimit:    0,
		TypesOnly:    false,
		Filter:       "(name=" + sanitizeSearch(duckid) + ")",
		Attributes:   []string{"displayName", "UOAD-UoPersonAffiliation"},
		Controls:     nil,
	}

	var sr, err = c.lc.Search(req)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("too many entries found for %q", duckid)
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	return sr.Entries[0], nil
}

// sanitizeSearch attempts to make a duckid safe for use in an LDAP search
func sanitizeSearch(duckid string) string {
	// We have to escape backslashes before other things to avoid double-escaping of them
	duckid = strings.Replace(duckid, "\\", "\\5c", -1)

	// Escape all other dangerous characters
	var escapes = map[string]string{
		"*":      "\\2a",
		"(":      "\\28",
		")":      "\\29",
		"\u0000": "\\00",
	}
	for unsafe, safe := range escapes {
		duckid = strings.Replace(duckid, unsafe, safe, -1)
	}

	return duckid
}

// Close ends the underlying LDAP connection
func (c *Connection) Close() {
	c.lc.Close()
}
