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

// Conf is the global configuration exposed to the entire app, and as such
// should be built precisely once, and never be modified
var Conf = &Config{}

// Parse reads the given settings file and returns a parsed Config.  File paths
// are parsed and verified as they are used by most subsystems.  The database
// connection string is built, but is not tested.
func Parse(filename string) error {
	// Read settings and store them into Conf
	var bc = bashconf.New()
	bc.EnvironmentPrefix("SCI_")
	var err = bc.ParseFile(filename)
	if err != nil {
		return err
	}
	err = bc.Store(Conf)
	if err != nil {
		return err
	}

	return nil
}
