//
// go-sonos
// ========
//
// Copyright (c) 2012, Ian T. Richards <ianr@panix.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright notice,
//     this list of conditions and the following disclaimer.
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENTING SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package main

import (
	"errors"
	"fmt"
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	_ "log"
)

type State struct {
	config *config.Config
	sonos  *sonos.Sonos
}

type Interp struct {
	State
}

type HandlerFn func(this *Interp, args []Arg) error

type Handler struct {
	fn  HandlerFn
	sig string
}

type HandlerMap map[string]Handler

var handlerMap HandlerMap

func init() {
	handlerMap = map[string]Handler{
		"config": Handler{(*Interp).do_config, "s"},
		"mute":   Handler{fn: (*Interp).do_mute},
		"unmute": Handler{fn: (*Interp).do_unmute},
		"select": Handler{(*Interp).do_select, "s"},
	}
}

func (this *Interp) do_config(args []Arg) (err error) {
	this.config = config.MakeConfig(args[0].sv)
	this.config.Init()
	return
}

func (this *Interp) do_mute(args []Arg) (err error) {
	this.sonos.SetMute(0, "Master", true)
	return
}

func (this *Interp) do_unmute(args []Arg) (err error) {
	this.sonos.SetMute(0, "Master", false)
	return
}

func (this *Interp) do_select(args []Arg) (err error) {
	if dev := this.config.Lookup(args[0].sv); nil != dev {
		this.sonos = sonos.Connect(dev, nil, sonos.SVC_RENDERING_CONTROL)
	} else {
		err = errors.New("unknown device " + args[0].sv)
	}
	return
}

var argTypeMap = map[rune]int{
	's': STRING,
	'i': NUMBER,
	'f': NUMBER,
}

func (this *Interp) checkArgs(cmd Cmd, handler Handler) {
	if len(cmd.args) != len(handler.sig) {
		panic("wrong number of arguments")
	} else if 0 < len(handler.sig) {
		for i, c := range handler.sig {
			if id, exists := argTypeMap[c]; exists {
				if id != cmd.args[i].id {
					panic(fmt.Sprintf("wrong type for argument %d", i+1))
				}
			} else {
				panic(fmt.Sprintf("unknown argument type character %c", c))
			}
		}
	}
}

func (this *Interp) execute(cmd Cmd) {
	if handler, has := handlerMap[cmd.name]; has {
		this.checkArgs(cmd, handler)
		if err := handler.fn(this, cmd.args); nil != err {
			panic(err)
		}
	} else {
		panic("unknown command " + cmd.name)
	}
}
