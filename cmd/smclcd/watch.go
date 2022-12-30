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
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/sstallion/go-smclcd"
	"github.com/sstallion/go-tools/command"
)

type watchBacklight int

const (
	watchBacklightOn watchBacklight = iota
	watchBacklightOff
	watchBacklightAuto
)

func (b *watchBacklight) Run(l *smclcd.LCD) (err error) {
	switch *b {
	case watchBacklightOn:
		err = l.SetBacklight(smclcd.BacklightOn)
	case watchBacklightOff, watchBacklightAuto:
		err = l.SetBacklight(smclcd.BacklightOff)
		if *b == watchBacklightOff {
			break
		}

		var keys = make(chan smclcd.Key)
		go func() {
			for {
				_ = <-keys
				l.SetBacklight(smclcd.BacklightOn)
				time.Sleep(10 * time.Second)
				l.SetBacklight(smclcd.BacklightOff)
			}
		}()
		go func() {
			for {
				key, err := l.GetInput()
				if err != nil {
					continue
				}
				select {
				case keys <- key:
				default:
					// ignore spurious events
				}
			}
		}()
	}
	return
}

type watchCmd struct {
	flags     *flag.FlagSet
	n         time.Duration
	backlight watchBacklight
	name      string
	args      []string
}

var ()

func init() {
	cmd := &watchCmd{flags: flag.NewFlagSet("watch", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	cmd.flags.Func("backlight", "`state`", cmd.parseBacklight)
	cmd.flags.DurationVar(&cmd.n, "n", 30*time.Second, "`interval`")
	command.Add(cmd)
}

func (cmd *watchCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *watchCmd) Description() string {
	return "Write periodic command output to display"
}

func (cmd *watchCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }} [flags] <command> [arguments...]

Flags:

  {{ call .PrintDefaults }}

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *watchCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	args := cmd.flags.Args()
	if len(args) < 1 {
		return command.ErrNArg
	}
	cmd.name = args[0]
	cmd.args = args[1:]
	return nil
}

func (cmd *watchCmd) parseBacklight(s string) error {
	switch s {
	case "on":
		cmd.backlight = watchBacklightOn
	case "off":
		cmd.backlight = watchBacklightOff
	case "auto":
		cmd.backlight = watchBacklightAuto
	default:
		return errors.New("invalid argument: " + s)
	}
	return nil
}

func (cmd *watchCmd) Run() error {
	l, err := openLCD()
	if err != nil {
		return err
	}
	defer l.Close()

	if err = cmd.backlight.Run(l); err != nil {
		return err
	}

	var b bytes.Buffer
	for {
		c := exec.Command(cmd.name, cmd.args...)
		c.Stdout = &b
		c.Stderr = os.Stderr
		if err = c.Run(); err != nil {
			return fmt.Errorf("failed to execute command: %v", cmd)
		}
		if err = l.Clear(); err != nil {
			return err
		}
		if _, err = b.WriteTo(l); err != nil {
			if err != io.EOF {
				return err
			}
		}
		b.Reset()
		time.Sleep(cmd.n)
	}
}
