package ldapservice

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uoregon-libraries/student-course-integrator/src/config"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
	ldap "gopkg.in/ldap.v2"
)

func TestSanitize(t *testing.T) {
	var str = "this is (not) a t*est:\\"
	var expected = "this is \\28not\\29 a t\\2aest:\\5c"

	var out = sanitizeSearch(str)
	if out != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, out)
	}
}

// This is a very specific test which is only run if config is set up
func TestConnect(t *testing.T) {
	var wd string
	var err error

	wd, err = os.Getwd()
	global.Conf, err = config.Parse(filepath.Join(wd, "..", "..", "sci.conf"))
	if err != nil {
		t.Fatalf("Cannot test LDAP: config is invalid: %s", err)
	}

	var conn *Connection
	var entry *ldap.Entry
	if strings.Contains(global.Conf.LDAPBaseDN, "uoregon") {
		conn, err = Connect()
		if err != nil {
			t.Fatalf("Unable to connect to ldap: %s", err)
		}

		entry, err = conn.Search("jechols")
		if err != nil {
			t.Fatalf("Unable to perform an LDAP search: %s", err)
		}
		var dn = entry.GetAttributeValue("displayName")
		if dn != "Jeremy Echols" {
			t.Fatalf("Expected displayName to be Jeremy Echols, but got %q", dn)
		}
	}
}
