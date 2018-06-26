// Package config is the project-specific configuration reader / parser /
// validator.  This uses the more generalized bashconf but adds our
// app-specific logic.
package config

import (
	"github.com/uoregon-libraries/gopkg/bashconf"
)

// Config holds the configuration needed for this application to work
type Config struct {
	DatabaseConnect string `setting:"DB"`
	Debug           bool   `setting:"DEBUG" type:"bool"`
	BindAddress     string `setting:"BIND_ADDRESS"`
	SessionSecret   string `setting:"SESSION_SECRET"`
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
