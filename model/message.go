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

package model

import (
	"encoding/xml"
	"github.com/ianr0bkny/go-sonos/didl"
	"github.com/ianr0bkny/go-sonos/upnp"
	"strings"
	"time"
)

type PositionInfo struct {
	Track               uint32
	TrackDuration       time.Duration
	TrackURI            string
	RelTime             time.Duration
	ProtocolInfo        string
	Title               string
	Class               string
	Creator             string
	Album               string
	OriginalTrackNumber string
}

func getDuration(in string) (d time.Duration, err error) {
	in = strings.Replace(in, ":", "h", 1)
	in = strings.Replace(in, ":", "m", 1)
	in += "s"
	return time.ParseDuration(in)
}

func GetPositionInfoMessage(in *upnp.PositionInfo) *PositionInfo {
	var trackDuration, relTime time.Duration
	trackDuration, err := getDuration(in.TrackDuration)
	if nil == err {
		trackDuration /= time.Second
	}

	relTime, err = getDuration(in.RelTime)
	if nil == err {
		relTime /= time.Second
	}

	out := &PositionInfo{
		Track:         in.Track,
		TrackDuration: trackDuration,
		TrackURI:      in.TrackURI,
		RelTime:       relTime,
	}

	metadata := &didl.Lite{}
	xml.Unmarshal([]byte(in.TrackMetaData), metadata)
	metadata.Validate()

	for _, item := range metadata.Item {
		for _, res := range item.Res {
			out.ProtocolInfo = res.ProtocolInfo
			break
		}
		for _, title := range item.Title {
			out.Title = title.Value
			break
		}
		for _, class := range item.Class {
			out.Class = class.Value
			break
		}
		for _, creator := range item.Creator {
			out.Creator = creator.Value
			break
		}
		for _, album := range item.Album {
			out.Album = album.Value
			break
		}
		for _, originalTrackNumber := range item.OriginalTrackNumber {
			out.OriginalTrackNumber = originalTrackNumber.Value
			break
		}
		break
	}
	return out
}

type QueueElement struct {
	ID                  string
	ParentID            string
	TrackURI            string
	Title               string
	Class               string
	AlbumArtURI         string
	Creator             string
	Album               string
	OriginalTrackNumber string
}

func protectEncoding(s string) string {
	return strings.Replace(s, "%", "%25", -1)
}

func GetQueueContentsMessage(in []Object) []QueueElement {
	var out []QueueElement
	for _, obj := range in {
		out = append(out, QueueElement{
			ID:                  protectEncoding(obj.ID()),
			ParentID:            obj.ParentID(),
			TrackURI:            obj.Res(),
			Title:               obj.Title(),
			Class:               obj.Class(),
			AlbumArtURI:         obj.AlbumArtURI(),
			Creator:             obj.Creator(),
			Album:               obj.Album(),
			OriginalTrackNumber: obj.OriginalTrackNumber(),
		})
	}
	return out
}

type TransportInfo *upnp.TransportInfo
