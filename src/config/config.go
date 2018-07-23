// Package config is the project-specific configuration reader / parser /
// validator.  This uses the more generalized bashconf but adds our
// app-specific logic.
package config

import (
	"github.com/uoregon-libraries/gopkg/bashconf"
)

// Config holds the configuration needed for this application to work
type Config struct {
	DatabaseConnect      string `setting:"DB"`
	Debug                bool   `setting:"DEBUG" type:"bool"`
	BindAddress          string `setting:"BIND_ADDRESS"`
	SessionSecret        string `setting:"SESSION_SECRET"`
	AuthHeader           string `setting:"AUTH_HEADER"`
	LDAPServer           string `setting:"LDAP_SERVER"`
	LDAPUser             string `setting:"LDAP_BIND_USER"`
	LDAPPass             string `setting:"LDAP_BIND_PASS"`
	LDAPBaseDN           string `setting:"LDAP_BASE_DN"`
	BannerCSVPath        string `setting:"BANNER_CSV_PATH" type:"path"`
	CanvasCSVPath        string `setting:"CANVAS_CSV_PATH" type:"path"`
	LookupByDuckIDURL    string `setting:"LOOKUP_BY_DUCKID_URL" type:"url"`
	LookupByBannerIDURL  string `setting:"LOOKUP_BY_BANNERID_URL" type:"url"`
	TranslatorAPIHeaders string `setting:"TRANSLATOR_API_HEADERS"`
}

// Parse reads the given settings file and returns a parsed Config.  File paths
// are parsed and verified as they are used by most subsystems.  The database
// connection string is built, but is not tested.
func Parse(filename string) (*Config, error) {
	var bc = bashconf.New()
	bc.EnvironmentPrefix("SCI_")
	var err = bc.ParseFile(filename)
	if err != nil {
		return nil, err
	}

	var conf = new(Config)
	err = bc.Store(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
