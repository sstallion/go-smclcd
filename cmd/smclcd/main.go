// Copyright (c) 2023 Steven Stallion <sstallion@gmail.com>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

//go:generate doxxer . help -all
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sstallion/go-hid"
	"github.com/sstallion/go-smclcd"
	"github.com/sstallion/go-tools/command"
	"github.com/sstallion/go-tools/util"
)

// go build -ldflags '-X main.version=<version>'
var version string

func init() {
	util.FixVersion(&version)
}

type versionFlag struct{}

func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) String() string   { return "" }
func (versionFlag) Set(s string) error {
	fmt.Printf("%s version %s (HIDAPI version %s)\n",
		util.Program(), version, hid.GetVersionStr())
	os.Exit(0)
	return nil
}

type debugFlag struct{}

func (debugFlag) IsBoolFlag() bool { return true }
func (debugFlag) String() string   { return "" }
func (debugFlag) Set(s string) error {
	smclcd.DebugLog.SetOutput(os.Stderr)
	return nil
}

var pathFlag, serialFlag string

func usage() {
	command.PrintGlobalUsage(`
TODO.

Usage:

  {{ .Program }} [global flags] <command> [arguments...]

Global Flags:

  {{ call .PrintDefaults }}

Commands:

  {{ call .PrintCommands }}

Use "{{ .Program }} help <command>" for more information about that command.
Report issues to https://github.com/sstallion/go-smclcd/issues.
`)
}

func main() {
	flag.Usage = usage
	flag.Var(versionFlag{}, "V", "TODO")
	flag.Var(debugFlag{}, "debug", "TODO")
	flag.StringVar(&pathFlag, "path", "", "TODO")
	flag.StringVar(&serialFlag, "serial", "", "TODO")
	command.Parse()
}
