package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
	fn()
}
