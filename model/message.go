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

// This module is intended to define an object model that go-sonos can
// use to refer to albums, tracks, and other structures and collections.
// This object model is based on and derived from the classes found in
// DIDL-Lite documents.
package model

import (
	"encoding/xml"
	"github.com/ianr0bkny/go-sonos/didl"
	"github.com/ianr0bkny/go-sonos/upnp"
)

type PositionInfo struct {
	Track               uint32
	TrackDuration       string
	TrackURI            string
	RelTime             string
	ProtocolInfo        string
	Title               string
	Class               string
	Creator             string
	Album               string
	OriginalTrackNumber string
}

func GetPositionInfoMessage(in *upnp.PositionInfo) *PositionInfo {
	out := &PositionInfo{
		Track:         in.Track,
		TrackDuration: in.TrackDuration,
		TrackURI:      in.TrackURI,
		RelTime:       in.RelTime,
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
