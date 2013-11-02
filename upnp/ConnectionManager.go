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

package upnp

import (
	"encoding/xml"
	_ "log"
)

var (
	ConnectionManager_EventType = registerEventType("ConnectionManager")
)

type ConnectionManagerState struct {
	SourceProtocolInfo   string
	SinkProtocolInfo     string
	CurrentConnectionIDs string
}

type ConnectionManagerEvent struct {
	ConnectionManagerState
	Svc *Service
}

func (this ConnectionManagerEvent) Service() *Service {
	return this.Svc
}

func (this ConnectionManagerEvent) Type() int {
	return ConnectionManager_EventType
}

type ConnectionManager struct {
	ConnectionManagerState
	Svc *Service
}

func (this *ConnectionManager) BeginSet(svc *Service, channel chan Event) {
}

type connectionManagerUpdate_XML struct {
	XMLName xml.Name `xml:"ConnectionManagerState"`
	Value   string   `xml:",innerxml"`
}

func (this *ConnectionManager) HandleProperty(svc *Service, value string, channel chan Event) error {
	update := connectionManagerUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		xml.Unmarshal(bytes, &this.ConnectionManagerState)
	}
	return nil
}

func (this *ConnectionManager) EndSet(svc *Service, channel chan Event) {
	evt := ConnectionManagerEvent{ConnectionManagerState: this.ConnectionManagerState, Svc: svc}
	channel <- evt
}

func (this *ConnectionManager) GetProtocolInfo() (source, sink string, err error) {
	type Response struct {
		XMLName xml.Name
		Source  string
		Sink    string
		ErrorResponse
	}
	response := this.Svc.CallVa("GetProtocolInfo")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	source = doc.Source
	sink = doc.Sink
	err = doc.Error()
	return
}

func (this *ConnectionManager) GetCurrentConnectionIDs() (connectionIds string, err error) {
	type Response struct {
		XMLName       xml.Name
		ConnectionIDs string
		ErrorResponse
	}
	response := this.Svc.CallVa("GetCurrentConnectionIDs")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	connectionIds = doc.ConnectionIDs
	err = doc.Error()
	return
}

const (
	Direction_Input  = "Input"
	Direction_Output = "Output"
)

const (
	Status_OK                    = "OK"
	Status_ContentFormatMismatch = "ContentFormatMismatch"
	Status_InsufficientBandwidth = "InsufficientBandwidth"
	Status_UnreliableChannel     = "UnreliableChannel"
	Status_Unknown               = "Unknown"
)

type ConnectionInfo struct {
	RcsID                 int32
	AVTransportID         int32
	ProtocolInfo          string
	PeerConnectionManager string
	PeerConnectionID      int32
	Direction             string
	Status                string
}

func (this *ConnectionManager) GetCurrentConnectionInfo(connectionId int32) (connectionInfo *ConnectionInfo, err error) {
	type Response struct {
		XMLName xml.Name
		ConnectionInfo
		ErrorResponse
	}
	args := []Arg{
		{"ConnectionID", connectionId},
	}
	response := this.Svc.Call("GetCurrentConnectionInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	connectionInfo = &doc.ConnectionInfo
	err = doc.Error()
	return
}
