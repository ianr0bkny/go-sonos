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

//
// A server to demonstrate controlling Sonos from a web browser.
//
// csweb := (c)ontrol (s)onos from a (web) browser
//
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
		s = sonos.Connect(dev, nil, sonos.SVC_CONTENT_DIRECTORY|sonos.SVC_AV_TRANSPORT|sonos.SVC_RENDERING_CONTROL)
	} else {
		log.Fatal("Could not create Sonos instance")
	}
	return s
}

func replyOk(w http.ResponseWriter, value interface{}) {
	reply(w, nil, value)
}

func replyError(w http.ResponseWriter, msg string) {
	reply(w, errors.New(msg), nil)
}

type Reply struct {
	Error string      `json:",omitempty"`
	Value interface{} `json:",omitempty"`
}

func reply(w http.ResponseWriter, err error, value interface{}) {
	r := Reply{}
	if nil != err {
		r.Error = fmt.Sprintf("%v", err)
	}
	if nil != value {
		r.Value = value
	}
	encoder := json.NewEncoder(os.Stdout)
	if false {
		encoder.Encode(r)
	}
	encoder = json.NewEncoder(w)
	encoder.Encode(r)
}

type handlerFunc func(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error

//
// get-position-info
//
func handle_GetPositionInfo(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if info, err := s.GetPositionInfo(0); nil != err {
		return err
	} else {
		replyOk(w, model.GetPositionInfoMessage(info))
	}
	return nil
}

//
// get-transport-info
//
func handle_GetTransportInfo(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if info, err := s.GetTransportInfo(0); nil != err {
		return err
	} else {
		replyOk(w, info)
	}
	return nil
}

//
// get-volume
//
func handle_GetVolume(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if volume, err := s.GetVolume(0, upnp.Channel_Master); nil != err {
		return err
	} else {
		replyOk(w, volume)
	}
	return nil
}

//
// next
//
func handle_Next(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.Next(0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// next-section
//
func handle_NextSection(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.NextSection(0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// pause
//
func handle_Pause(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.Pause(0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// play
//
func handle_Play(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.Play(0, "1"); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// previous
//
func handle_Previous(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.Previous(0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// previous-section
//
func handle_PreviousSection(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.PreviousSection(0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// remove-track-from-queue
//
func handle_RemoveTrackFromQueue(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	track_s := r.FormValue("track")
	if err := s.RemoveTrackFromQueue(0, fmt.Sprintf("Q:0/%s", track_s), 0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// seek
//
func handle_Seek(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	unit := r.FormValue("unit")
	target := r.FormValue("target")
	if err := s.Seek(0, unit, target); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

//
// set-volume
//
func handle_SetVolume(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	volume_s := r.FormValue("value")
	if volume, err := strconv.ParseInt(volume_s, 10, 16); nil != err {
		return errors.New(fmt.Sprintf("Invalid volume `%s' specified", volume_s))
	} else {
		if err := s.SetVolume(0, upnp.Channel_Master, uint16(volume)); nil != err {
			return err
		}
	}
	replyOk(w, true)
	return nil
}

//
// stop
//
func handle_Stop(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if err := s.Stop(0); nil != err {
		return err
	}
	replyOk(w, true)
	return nil
}

var controlHandlerMap = map[string]handlerFunc{
	"get-position-info":       handle_GetPositionInfo,
	"get-transport-info":      handle_GetTransportInfo,
	"get-volume":              handle_GetVolume,
	"next":                    handle_Next,
	"next-section":            handle_NextSection,
	"pause":                   handle_Pause,
	"play":                    handle_Play,
	"previous":                handle_Previous,
	"previous-section":        handle_PreviousSection,
	"remove-track-from-queue": handle_RemoveTrackFromQueue,
	"seek":       handle_Seek,
	"set-volume": handle_SetVolume,
	"stop":       handle_Stop,
}

func handleControl(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) {
	f := r.FormValue("method")
	if handler, has := controlHandlerMap[f]; has {
		if err := handler(s, w, r); nil != err {
			replyError(w, fmt.Sprintf("Error in call to %s: %v", f, err))
		}
		return
	} else {
		replyError(w, fmt.Sprintf("No such method control::%s", f))
	}
}

//
// get-album-tracks
//
func handle_GetAlbumTracks(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if tracks, err := s.GetAlbumTracks(r.FormValue("album")); nil != err {
		return err
	} else {
		replyOk(w, model.GetQueueContentsMessage(tracks))
	}
	return nil
}

//
// get-all-genres
//
func handle_GetAllGenres(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if list, err := s.GetAllGenres(); nil != err {
		return err
	} else {
		replyOk(w, model.GetQueueContentsMessage(list))
	}
	return nil
}

//
// get-artist-albums
//
func handle_GetArtistAlbums(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if list, err := s.GetArtistAlbums(r.FormValue("artist")); nil != err {
		return err
	} else {
		replyOk(w, model.GetQueueContentsMessage(list))
	}
	return nil
}

//
// get-direct-children
//
func handle_GetDirectChildren(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if list, err := s.GetDirectChildren(r.FormValue("root")); nil != err {
		return err
	} else {
		replyOk(w, model.GetQueueContentsMessage(list))
	}
	return nil
}

//
// get-genre-artists
//
func handle_GetGenreArtists(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if artists, err := s.GetGenreArtists(r.FormValue("genre")); nil != err {
		return err
	} else {
		replyOk(w, model.GetQueueContentsMessage(artists))
	}
	return nil
}

//
// get-queue-contents
//
func handle_GetQueueContents(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) error {
	if queue, err := s.GetQueueContents(); nil != err {
		return err
	} else {
		replyOk(w, model.GetQueueContentsMessage(queue))
	}
	return nil
}

var browseHandlerMap = map[string]handlerFunc{
	"get-album-tracks":    handle_GetAlbumTracks,
	"get-all-genres":      handle_GetAllGenres,
	"get-artist-albums":   handle_GetArtistAlbums,
	"get-direct-children": handle_GetDirectChildren,
	"get-genre-artists":   handle_GetGenreArtists,
	"get-queue-contents":  handle_GetQueueContents,
}

func handleBrowse(s *sonos.Sonos, w http.ResponseWriter, r *http.Request) {
	f := r.FormValue("method")
	if handler, has := browseHandlerMap[f]; has {
		if err := handler(s, w, r); nil != err {
			replyError(w, fmt.Sprintf("Error in call to %s: %v", f, err))
		}
		return
	} else {
		replyError(w, fmt.Sprintf("No such method browse::%s", f))
	}
}

func setupHttp(s *sonos.Sonos) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, strings.Join([]string{"web", r.RequestURI}, "/"))
	})

	http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
		handleControl(s, w, r)
	})

	http.HandleFunc("/browse", func(w http.ResponseWriter, r *http.Request) {
		handleBrowse(s, w, r)
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

	log.Printf("Starting server loop ...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", CSWEB_HTTP_PORT), nil))
}
