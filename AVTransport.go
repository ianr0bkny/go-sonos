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

package sonos

import (
	"encoding/xml"
	_ "log"
	"sonos/upnp"
)

type AVTransport struct {
	svc *upnp.Service
}

const (
	NORMAL           = "NORMAL"
	REPEAT_ALL       = "REPEAT_ALL"
	SHUFFLE_NOREPEAT = "SHUFFLE_NOREPEAT"
	SHUFFLE          = "SHUFFLE"
)

func (this *AVTransport) SetPlayMode(instance int, mode string) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"NewPlayMode", mode},
	}
	response := upnp.Call(this.svc, "SetPlayMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) PreviousSection(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "PreviousSection", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) NextSection(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "NextSection", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) Previous(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "Previous", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) NextProgrammedRadioTracks(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "NextProgrammedRadioTracks", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) Next(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "Next", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

const (
	TRACK_NR = "TRACK_NR"
	REL_TIME = "REL_TIME"
	SECTION  = "SECTION"
)

func (this *AVTransport) Seek(instance int, unit, target string) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"Unit", unit},
		{"Target", target},
	}
	response := upnp.Call(this.svc, "Seek", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) Pause(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "Pause", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) Play(instance int, speed string) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"Speed", speed},
	}
	response := upnp.Call(this.svc, "Play", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) Stop(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "Stop", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) GetCrossfadeMode(instance int) bool {
	type Response struct {
		XMLName       xml.Name
		CrossfadeMode bool
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "GetCrossfadeMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.CrossfadeMode
}

type TransportSettings struct {
	PlayMode       string
	RecQualityMode string
}

func (this *AVTransport) GetTransportSettings(instance int) *TransportSettings {
	type Response struct {
		XMLName xml.Name
		TransportSettings
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "GetTransportSettings", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.TransportSettings
}

type DeviceCapabilities struct {
	PlayMedia       string
	RecMedia        string
	RecQualityModes string
}

func (this *AVTransport) GetDeviceCapabilities(instance int) *DeviceCapabilities {
	type Response struct {
		XMLName xml.Name
		DeviceCapabilities
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "GetDeviceCapabilities", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.DeviceCapabilities
}

type PositionInfo struct {
	Track         int
	TrackDuration string
	TrackMetaData string
	TrackURI      string
	RelTime       string
	AbsTime       string
	RelCount      int
	AbsCount      int
}

func (this *AVTransport) GetPositionInfo(instance int) *PositionInfo {
	type Response struct {
		XMLName xml.Name
		PositionInfo
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "GetPositionInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.PositionInfo
}

type TransportInfo struct {
	CurrentTransportState  string
	CurrentTransportStatus string
	CurrentSpeed           string
}

func (this *AVTransport) GetTransportInfo(instance int) *TransportInfo {
	type Response struct {
		XMLName xml.Name
		TransportInfo
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "GetTransportInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.TransportInfo
}

type MediaInfo struct {
	NrTracks           int
	MediaDuration      string
	CurrentURI         string
	CurrentURIMetaData string
	NextURI            string
	NextURIMetaData    string
	PlayMedium         string
	RecordMedium       string
	WriteStatus        string
}

func (this *AVTransport) GetMediaInfo(instance int) *MediaInfo {
	type Response struct {
		XMLName xml.Name
		MediaInfo
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "GetMediaInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.MediaInfo
}

func (this *AVTransport) BackupQueue(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "BackupQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) SaveQueue(instance int, title, object string) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"Title", title},
		{"ObjectID", object},
	}
	response := upnp.Call(this.svc, "SaveQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) RemoveAllTracksFromQueue(instance int) {
	type Response struct {
		XMLName xml.Name
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
	}
	response := upnp.Call(this.svc, "RemoveAllTracksFromQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	// TODO
}

func (this *AVTransport) RemoveTrackRangeFromQueue(instance, update, start, number int) int {
	type Response struct {
		XMLName     xml.Name
		NewUpdateID int
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"UpdateID", update},
		{"StartingIndex", start},
		{"NumberOfTracks", number},
	}
	response := upnp.Call(this.svc, "RemoveTrackRangeFromQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.NewUpdateID
}

func (this *AVTransport) RemoveTrackFromQueue(instance int, object string) int {
	type Response struct {
		XMLName  xml.Name
		UpdateID int
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"ObjectID", object},
	}
	response := upnp.Call(this.svc, "RemoveTrackFromQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.UpdateID
}

func (this *AVTransport) ReorderTracksInQueue(instance, start, num, before int) int {
	type Response struct {
		XMLName  xml.Name
		UpdateID int
	}
	args := []upnp.Arg{
		{"InstanceID", instance},
		{"StartingIndex", start},
		{"NumberOfTracks", num},
		{"InsertBefore", before},
	}
	response := upnp.Call(this.svc, "ReorderTracksInQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.UpdateID
}
