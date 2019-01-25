// +build ignore

package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	var err error

	var cmd = exec.Command("git", "rev-parse", "--short=8", "--verify", "HEAD")
	var out []byte
	out, err = cmd.CombinedOutput()
	if err != nil {
		panic("Unable to run `git describe`: " + err.Error())
	}
	var build = string(out)
	build = strings.TrimSpace(build)

	var f *os.File
	f, err = os.Create("commit.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(`// commit.go is a generated file and should not be modified by hand

package version

func init() {
	commit = "` + build + `"
}
`)
	if err != nil {
		panic(err)
	}
}
