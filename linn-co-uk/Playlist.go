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

package linn

import (
	"encoding/xml"
	"github.com/ianr0bkny/go-sonos/upnp"
)

type Playlist struct {
	Svc *upnp.Service
}

func (this *Playlist) BeginSet(svc *upnp.Service, channel chan upnp.Event) {
}

func (this *Playlist) HandleProperty(svc *upnp.Service, value string, channel chan upnp.Event) error {
	return nil
}

func (this *Playlist) EndSet(svc *upnp.Service, channel chan upnp.Event) {
}

func (this *Playlist) IdArray() (aIdArrayToken uint32, aIdArray string, err error) {
	type Response struct {
		XMLName       xml.Name
		AIdArrayToken uint32 `xml:"aIdArrayToken"`
		AIdArray      string `xml:"aIdArray"`
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("IdArray")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.AIdArrayToken, doc.AIdArray, doc.Error()
}

func (this *Playlist) TracksMax() (aTracksMax uint32, err error) {
	type Response struct {
		XMLName    xml.Name
		ATracksMax uint32 `xml:"aTracksMax"`
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("IdArray")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.ATracksMax, doc.Error()
}
