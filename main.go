package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	flags "github.com/jessevdk/go-flags"
	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/config"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
	"github.com/uoregon-libraries/student-course-integrator/src/importcsv"
	"github.com/uoregon-libraries/student-course-integrator/src/sci-server"
)

type command struct {
	desc string
	run  func()
}

var cmdMap = map[string]command{
	"server": {
		desc: "Backoffice web server which faculty log into in order to associate students with courses",
		run:  sciserver.Run,
	},
	"import-csv": {
		desc: "CSV importer for populating the database with courses and faculty",
		run:  importcsv.Run,
	},
}

func usageErr(e string) {
	fmt.Fprintf(os.Stderr, "\x1b[31;1m%s\x1b[0m\n\nUsage: sci <subcommand>\n\n", e)
	var keys []string
	var maxlen int
	for key, _ := range cmdMap {
		keys = append(keys, key)
		if len(key) > maxlen {
			maxlen = len(key)
		}
	}
	sort.Strings(keys)
	fmt.Fprintf(os.Stderr, "Valid subcommands:\n")
	for _, key := range keys {
		var cmd = cmdMap[key]
		fmt.Fprintf(os.Stderr, "  - %s%s: %s\n", key, strings.Repeat(" ", maxlen - len(key)), cmd.desc)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usageErr("You must specify a subcommand")
	}

	var cmd, ok = cmdMap[os.Args[1]]
	if !ok {
		usageErr(fmt.Sprintf("%q is not a valid subcommand", os.Args[1]))
	}

	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	getConf()
	cmd.run()
}

func getConf() {
	var p = flags.NewParser(&global.Opts, flags.HelpFlag|flags.PassDoubleDash)
	var _, err = p.Parse()

	if err != nil {
		var flagErr, ok = err.(*flags.Error)
		if !ok || flagErr.Type != flags.ErrHelp {
			fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		}
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if global.Opts.Approot == "" {
		global.Opts.Approot, err = os.Getwd()
		if err != nil {
			logger.Fatalf("Error trying to read current directory: %s", err)
		}
	}

	var files = []string{"/etc/sci.conf", filepath.Join(global.Opts.Approot, "sci.conf")}
	if global.Opts.ConfigFile == "" {
		for _, file := range files {
			if fileutil.Exists(file) {
				global.Opts.ConfigFile = file
			}
		}
	}

	if global.Opts.ConfigFile == "" {
		logger.Fatalf("Config error: no config file found in %s", strings.Join(files, ", "))
	}

	global.Conf, err = config.Parse(global.Opts.ConfigFile)
	if err != nil {
		logger.Fatalf("Config error: %s", err)
	}

	if global.Conf.Debug {
		logger.Warnf("Debug mode has been enabled")
	}

	global.DB, err = sql.Open("mysql", global.Conf.DatabaseConnect)
	if err != nil {
		logger.Fatalf("Error trying to connect to database: %s", err)
	}
	global.DB.SetConnMaxLifetime(time.Second)

	global.Store = sessions.NewCookieStore([]byte(global.Conf.SessionSecret))
}
