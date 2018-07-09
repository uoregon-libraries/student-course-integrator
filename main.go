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

var cmdMap = map[string]func(){"server": sciserver.Run}

func usageErr(e string) {
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: sci <subcommand>\n\n", e)
	var commands []string
	for key, _ := range cmdMap {
		commands = append(commands, key)
	}
	sort.Strings(commands)
	fmt.Fprintf(os.Stderr, "Valid subcommands: %s\n", strings.Join(commands, ","))
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usageErr("You must specify a subcommand")
	}

	var fn, ok = cmdMap[os.Args[1]]
	if !ok {
		usageErr(fmt.Sprintf("%q is not a valid subcommand", os.Args[1]))
	}

	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	getConf()
	fn()
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
