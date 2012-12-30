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
)

type ConnectionManager struct {
	Svc *upnp.Service
}

func (this *ConnectionManager) GetProtocolInfo() (source, sink string) {
	type Response struct {
		XMLName xml.Name
		Source  string
		Sink    string
	}
	response := upnp.CallVa(this.Svc, "GetProtocolInfo")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Source, doc.Sink
}

func (this *ConnectionManager) GetCurrentConnectionIDs() string {
	type Response struct {
		XMLName       xml.Name
		ConnectionIDs string
	}
	response := upnp.CallVa(this.Svc, "GetCurrentConnectionIDs")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.ConnectionIDs
}

type ConnectionInfo struct {
	RcsID                 int
	AVTransportID         int
	ProtocolInfo          string
	PeerConnectionManager string
	PeerConnectionID      int
	Direction             string
	Status                string
}

func (this *ConnectionManager) GetCurrentConnectionInfo(connection int) *ConnectionInfo {
	type Response struct {
		XMLName xml.Name
		ConnectionInfo
	}
	response := upnp.CallVa(this.Svc, "GetCurrentConnectionInfo", "ConnectionID", connection)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.ConnectionInfo
}
