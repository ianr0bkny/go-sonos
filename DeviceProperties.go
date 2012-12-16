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
	"log"
	"sonos/upnp"
)

type DeviceProperties struct {
	svc *upnp.Service
}

func (this *DeviceProperties) SetLEDState(state string) {
	response := upnp.CallVa(this.svc, "SetLEDState", "DesiredLEDState", state)
	log.Printf("%#v", response)
}

func (this *DeviceProperties) GetLEDState() string {
	type Response struct {
		XMLName         xml.Name
		CurrentLEDState string
	}
	response := upnp.CallVa(this.svc, "GetLEDState")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.CurrentLEDState
}

func (this *DeviceProperties) GetHouseholdID() string {
	type Response struct {
		XMLName            xml.Name
		CurrentHouseholdID string
	}
	response := upnp.CallVa(this.svc, "GetHouseholdID")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.CurrentHouseholdID
}

func (this *DeviceProperties) GetZoneAttributes() (name, icon string) {
	type Response struct {
		XMLName         xml.Name
		CurrentZoneName string
		Icon            string
	}
	response := upnp.CallVa(this.svc, "GetZoneAttributes")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.CurrentZoneName, doc.Icon
}
