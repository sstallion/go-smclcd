// Copyright (c) 2023 Steven Stallion <sstallion@gmail.com>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this cmd of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this cmd of conditions and the following disclaimer in the
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
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/sstallion/go-hid"
	"github.com/sstallion/go-smclcd"
	"github.com/sstallion/go-tools/command"
)

type listCmd struct {
	flags *flag.FlagSet
}

func init() {
	cmd := &listCmd{flags: flag.NewFlagSet("list", flag.ExitOnError)}
	cmd.flags.Usage = cmd.Usage
	command.Add(cmd)
}

func (cmd *listCmd) Name() string {
	return cmd.flags.Name()
}

func (cmd *listCmd) Description() string {
	return "List compatible displays"
}

func (cmd *listCmd) Usage() {
	command.PrintUsage(cmd.flags, `
TODO.

Usage:

  {{ .Program }} [global flags] {{ .Name }}

Use "{{ .Program }} help" for more information about global flags.
`)
}

func (cmd *listCmd) Parse(arguments []string) error {
	if err := cmd.flags.Parse(arguments); err != nil {
		return err
	}
	args := cmd.flags.Args()
	if len(args) != 0 {
		return command.ErrNArg
	}
	return nil
}

func (cmd *listCmd) Run() error {
	var seen bool
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	hid.Enumerate(smclcd.VendorID, smclcd.ProductID,
		func(info *hid.DeviceInfo) (err error) {
			if !seen {
				fmt.Fprintf(w, "Path\tManufacturer\tProduct\tSerial Number\n")
				seen = true
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				info.Path,
				strings.TrimSpace(info.MfrStr),
				strings.TrimSpace(info.ProductStr),
				strings.TrimSpace(info.SerialNbr))
			return nil
		})
	w.Flush()
	return nil
}
