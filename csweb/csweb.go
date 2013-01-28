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
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	"github.com/ianr0bkny/go-sonos/model"
	"log"
	"net/http"
)

func getCurrentQueue(sonos *sonos.Sonos, writer http.ResponseWriter, request *http.Request) {
	if result, err := sonos.GetQueueContents(); nil != err {
		log.Fatal(err)
	} else {
		encoder := json.NewEncoder(writer)
		encoder.Encode(model.ObjectMessageStream(result))
	}
}

func pressPlay(sonos *sonos.Sonos, writer http.ResponseWriter, request *http.Request) {
	if err := sonos.Play(0, "1"); nil != err {
		log.Fatal(err)
	}
}

func pressPause(sonos *sonos.Sonos, writer http.ResponseWriter, request *http.Request) {
	if err := sonos.Pause(0); nil != err {
		log.Fatal(err)
	}
}

func pressStop(sonos *sonos.Sonos, writer http.ResponseWriter, request *http.Request) {
	if err := sonos.Stop(0); nil != err {
		log.Fatal(err)
	}
}

const (
	CSWEB_CONFIG        = "/home/ianr/.go-sonos"
	CSWEB_DEVICE        = "kitchen"
	CSWEB_DISCOVER_PORT = "13104"
	CSWEB_EVENTING_PORT = "13105"
	CSWEB_NETWORK       = "eth0"
)

var testSonos *sonos.Sonos

func initSonos(flags int) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	c := config.MakeConfig(CSWEB_CONFIG)
	c.Init()
	if dev := c.Lookup(CSWEB_DEVICE); nil != dev {
		reactor := sonos.MakeReactor(CSWEB_NETWORK, CSWEB_EVENTING_PORT)
		testSonos = sonos.Connect(dev, reactor, flags)
	} else {
		log.Fatal("Could not create test instance")
	}
}

func getSonos(flags int) *sonos.Sonos {
	if nil == testSonos {
		initSonos(flags)
	}
	return testSonos
}

func main() {
	s := getSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)

	http.HandleFunc("/queue.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "res/queue.html")
	})
	http.HandleFunc("/current_queue", func(w http.ResponseWriter, r *http.Request) {
		getCurrentQueue(s, w, r)
	})
	http.HandleFunc("/press_play", func(w http.ResponseWriter, r *http.Request) {
		pressPlay(s, w, r)
	})
	http.HandleFunc("/press_pause", func(w http.ResponseWriter, r *http.Request) {
		pressPause(s, w, r)
	})
	http.HandleFunc("/press_stop", func(w http.ResponseWriter, r *http.Request) {
		pressStop(s, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
