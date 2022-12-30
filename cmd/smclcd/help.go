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

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/sstallion/go-tools/command"
)

type helpCmd struct {
	flags *flag.FlagSet
	usage func()
}

type helpAllFlag struct{}

func (helpAllFlag) IsBoolFlag() bool { return true }
func (helpAllFlag) String() string   { return "" }
func (helpAllFlag) Set(s string) error {
	flag.Usage()
	command.Visit(func(cmd command.Command) {
		if desc := cmd.Description(); desc != "" {
			fmt.Fprintf(os.Stderr, "\n%s\n\n", desc)
			cmd.Usage()
		}
	})
	os.Exit(0)
	return nil
}

func init() {
	cmd := &helpCmd{flags: flag.NewFlagSet("help", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	cmd.flags.Var(helpAllFlag{}, "all", "TODO")
	command.Add(cmd)
}

func (cmd *helpCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *helpCmd) Description() string { return "" }

func (cmd *helpCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }} <command>

Flags:

  {{ call .PrintDefaults }}

Commands:

  {{ call .PrintCommands }}

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *helpCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	args := cmd.flags.Args()
	if len(args) > 1 {
		return command.ErrNArg
	}
	if len(args) == 0 {
		cmd.usage = flag.Usage
	} else {
		c := command.Lookup(args[0])
		if c == nil {
			return errors.New("invalid command: " + args[0])
		}
		cmd.usage = c.Usage
	}
	return nil
}

func (cmd *helpCmd) Run() error {
	cmd.usage()
	return nil
}
