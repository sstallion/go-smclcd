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

	"github.com/sstallion/go-smclcd"
	"github.com/sstallion/go-tools/command"
)

type cursorCmd struct {
	flags *flag.FlagSet
	state smclcd.Cursor
}

var ()

func init() {
	cmd := &cursorCmd{flags: flag.NewFlagSet("cursor", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	command.Add(cmd)
}

func (cmd *cursorCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *cursorCmd) Description() string {
	return "Cursor control"
}

func (cmd *cursorCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }} <state>

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *cursorCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	args := cmd.flags.Args()
	if len(args) != 1 {
		return command.ErrNArg
	}
	return cmd.parseState(args[0])
}

func (cmd *cursorCmd) parseState(s string) error {
	switch s {
	case "off":
		cmd.state = smclcd.CursorOff
	case "block":
		cmd.state = smclcd.CursorBlock
	case "underline":
		cmd.state = smclcd.CursorUnderline
	case "both":
		cmd.state = smclcd.CursorBoth
	default:
		return errors.New("invalid argument: " + s)
	}
	return nil
}

func (cmd *cursorCmd) Run() error {
	l, err := openLCD()
	if err != nil {
		return err
	}
	defer l.Close()

	return l.SetCursor(cmd.state)
}
