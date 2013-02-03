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
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 

// csweb := (c)ontrol (s)onos from a (web) browser

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	"github.com/ianr0bkny/go-sonos/model"
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
func getCurrentQueue(sonos *sonos.Sonos, writer http.ResponseWriter, request *http.Request) {
	if result, err := sonos.GetQueueContents(); nil != err {
		log.Fatal(err)
	} else {
		encoder := json.NewEncoder(writer)
		encoder.Encode(model.ObjectMessageStream(result))
	}
}
*/

const (
	CSWEB_CONFIG        = "/home/ianr/.go-sonos"
	CSWEB_DEVICE        = "kitchen"
	CSWEB_DISCOVER_PORT = "13104"
	CSWEB_EVENTING_PORT = "13105"
	CSWEB_NETWORK       = "eth0"
	CSWEB_HTTP_PORT     = 8080
)

func initSonos(config *config.Config) *sonos.Sonos {
	var s *sonos.Sonos
	if dev := config.Lookup(CSWEB_DEVICE); nil != dev {
		reactor := sonos.MakeReactor(CSWEB_NETWORK, CSWEB_EVENTING_PORT)
		s = sonos.Connect(dev, reactor, sonos.SVC_CONTENT_DIRECTORY|sonos.SVC_AV_TRANSPORT|sonos.SVC_RENDERING_CONTROL)
	} else {
		log.Fatal("Could not create Sonos instance")
	}
	return s
}

func replyOk(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	encoder.Encode(true)
}

func replyError(w http.ResponseWriter, msg string) {
	reply(w, errors.New(msg), nil)
}

type Reply struct {
	Error string      `json:",omitempty"`
	Value interface{} `json:",omitempty"`
}

func reply(w http.ResponseWriter, err error, value interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	r := Reply{}
	if nil != err {
		r.Error = fmt.Sprintf("%v", err)
	}
	if nil != value {
		r.Value = value
	}
	encoder.Encode(r)
	encoder = json.NewEncoder(w)
	encoder.Encode(r)
}

func handleControl(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) {
	f := r.FormValue("method")
	switch f {
	case "previous":
		s.Previous(0)
	case "previous-section":
		s.PreviousSection(0)
	case "play":
		s.Play(0, "1")
	case "pause":
		s.Pause(0)
	case "stop":
		s.Stop(0)
	case "next-section":
		s.NextSection(0)
	case "next":
		s.Next(0)
	//
	case "volume":
		volume_s := r.FormValue("value")
		if volume, err := strconv.ParseInt(volume_s, 10, 16); nil != err {
			replyError(w, fmt.Sprintf("Invalid volume `%s' specified", volume_s))
			return
		} else {
			s.SetVolume(0, upnp.Channel_Master, uint16(volume))
			replyOk(w)
		}
		return
	case "get-volume":
		if volume, err := s.GetVolume(0, upnp.Channel_Master); nil != err {
			replyError(w, fmt.Sprintf("Error in call to %s: %v", f, err))
		} else {
			reply(w, nil, volume)
		}
		return
	case "get-position-info":
		if info, err := s.GetPositionInfo(0); nil != err {
			replyError(w, fmt.Sprintf("Error in call to %s: %v", f, err))
		} else {
			reply(w, nil, model.GetPositionInfoMessage(info))
		}
		return
	default:
		replyError(w, fmt.Sprintf("No such method `%s'", f))
		return
	}
	replyOk(w)
}

func setupHttp(s *sonos.Sonos) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, strings.Join([]string{"res", r.RequestURI}, "/"))
	})

	http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
		handleControl(s, w, r)
	})
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	config := config.MakeConfig(CSWEB_CONFIG)
	config.Init()

	s := initSonos(config)
	if nil != s {
		setupHttp(s)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", CSWEB_HTTP_PORT), nil))
}
