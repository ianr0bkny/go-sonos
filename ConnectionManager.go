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
	"github.com/ianr0bkny/go-sonos/upnp"
	_ "log"
)

type ConnectionManager struct {
	Svc *upnp.Service
}

func (this *ConnectionManager) GetProtocolInfo() (source, sink string, err error) {
	type Response struct {
		XMLName xml.Name
		Source  string
		Sink    string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetProtocolInfo")
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
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetCurrentConnectionIDs")
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
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"ConnectionID", connectionId},
	}
	response := upnp.Call(this.Svc, "GetCurrentConnectionInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	connectionInfo = &doc.ConnectionInfo
	err = doc.Error()
	return
}
