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
	"flag"
	"io"
	"strings"
	"text/template"

	"github.com/sstallion/go-tools/command"
)

const bannerTmpl = `
(@ @) go-smclcd
(.".) {{ .Version }}
`

type bannerCmd struct {
	flags *flag.FlagSet
}

func init() {
	cmd := &bannerCmd{flags: flag.NewFlagSet("banner", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	command.Add(cmd)
}

func (cmd *bannerCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *bannerCmd) Description() string { return "" }

func (cmd *bannerCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }}

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *bannerCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	args := cmd.flags.Args()
	if len(args) != 0 {
		return command.ErrNArg
	}
	return nil
}

func (cmd *bannerCmd) Run() error {
	l, err := openLCD()
	if err != nil {
		return err
	}
	defer l.Close()

	var b bytes.Buffer
	text := strings.TrimSpace(bannerTmpl) + "\n"
	t := template.Must(template.New("").Parse(text))
	t.Execute(&b, map[string]interface{}{
		"Version": version,
	})
	if err = l.Clear(); err != nil {
		return err
	}
	if _, err = l.Write(b.Bytes()); err != nil {
		if err != io.EOF {
			return err
		}
	}
	return nil
}
