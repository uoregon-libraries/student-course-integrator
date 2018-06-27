package sciserver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/config"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

var opts struct {
	Approot    string `short:"r" long:"app-root" description:"path to app root if not current working directory"`
	ConfigFile string `short:"c" long:"config" description:"path to SCI config file if not /etc/sci.conf or ./sci.conf"`
}

func Run() {
	getConf()

	var s = &server{
		Approot: opts.Approot,
		Bind:    ":8080",
		Debug:   global.Conf.Debug,
	}
	s.Listen()
}

func getConf() {
	var p = flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	var _, err = p.Parse()

	if err != nil {
		var flagErr, ok = err.(*flags.Error)
		if !ok || flagErr.Type != flags.ErrHelp {
			fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		}
		p.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if opts.Approot == "" {
		opts.Approot, err = os.Getwd()
		if err != nil {
			logger.Fatalf("Error trying to read current directory: %s", err)
		}
	}

	if opts.ConfigFile == "" {
		var files = []string{"/etc/sci.conf", filepath.Join(opts.Approot, "sci.conf")}
		for _, file := range files {
			if fileutil.Exists(file) {
				opts.ConfigFile = file
			}
		}
	}

	global.Conf, err = config.Parse(opts.ConfigFile)
	if err != nil {
		logger.Fatalf("Config error: %s", err)
	}

	if global.Conf.Debug {
		logger.Warnf("Debug mode has been enabled")
	}

	initRootTemplates(filepath.Join(opts.Approot, "templates"), global.Conf.Debug)
}