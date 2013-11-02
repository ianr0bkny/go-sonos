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
	MusicServices_EventType = registerEventType("MusicServices")
)

type MusicServicesState struct {
	ServiceListVersion string
}

type MusicServicesEvent struct {
	MusicServicesState
	Svc *Service
}

func (this MusicServicesEvent) Service() *Service {
	return this.Svc
}

func (this MusicServicesEvent) Type() int {
	return MusicServices_EventType
}

type MusicServices struct {
	MusicServicesState
	Svc *Service
}

func (this *MusicServices) BeginSet(svc *Service, channel chan Event) {
}

type musicServicesUpdate_XML struct {
	XMLName xml.Name `xml:"MusicServicesState"`
	Value   string   `xml:",innerxml"`
}

func (this *MusicServices) HandleProperty(svc *Service, value string, channel chan Event) error {
	update := musicServicesUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		xml.Unmarshal(bytes, &this.MusicServicesState)
	}
	return nil
}

func (this *MusicServices) EndSet(svc *Service, channel chan Event) {
	evt := MusicServicesEvent{MusicServicesState: this.MusicServicesState, Svc: svc}
	channel <- evt
}

type msPolicy_XML struct {
	XMLName      xml.Name
	Auth         string `xml:"Auth,attr"`
	PollInterval string `xml:"PollInterval,attr"`
}

type msStrings_XML struct {
	XMLName xml.Name
	Version string `xml:"Version,attr"`
	Uri     string `xml:"Uri,attr"`
}

type msPresentationMap_XML struct {
	XMLName xml.Name
	Version string `xml:"Version,attr"`
	Uri     string `xml:"Uri,attr"`
}

type msPresentation_XML struct {
	XMLName         xml.Name
	Strings         []msStrings_XML         `xml:"Strings"`
	PresentationMap []msPresentationMap_XML `xml:"PresentationMap"`
}

type msService_XML struct {
	XMLName           xml.Name
	Id                string               `xml:"Id,attr"`
	Name              string               `xml:"Name,attr"`
	Version           string               `xml:"Version,attr"`
	Uri               string               `xml:"Uri,attr"`
	SecureUri         string               `xml:"SecureUri,attr"`
	ContainerType     string               `xml:"ContainerType,attr"`
	Capabilities      string               `xml:"Capabilities,attr"`
	MaxMessagingChars string               `xml:"MaxMessagingChars,attr"`
	Policy            []msPolicy_XML       `xml:"Policy"`
	Presentation      []msPresentation_XML `xml:"Presentation"`
}

type msServices_XML struct {
	XMLName xml.Name
	Service []msService_XML `xml:"Service"`
}

func (this *MusicServices) GetSessionId(serviceId int16, username string) (sessionId string, err error) {
	type Response struct {
		XMLName   xml.Name
		SessionId string
		ErrorResponse
	}
	args := []Arg{
		{"ServiceId", serviceId},
		{"Username", username},
	}
	response := this.Svc.Call("GetSessionId", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	sessionId = doc.SessionId
	err = doc.Error()
	return
}

func (this *MusicServices) ListAvailableServices() (err error) {
	type Response struct {
		XMLName                        xml.Name
		AvailableServiceDescriptorList string
		AvailableServiceTypeList       string
		AvailableServiceListVersion    string
		ErrorResponse
	}
	response := this.Svc.CallVa("ListAvailableServices")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	services := msServices_XML{}
	xml.Unmarshal([]byte(doc.AvailableServiceDescriptorList), &services)
	err = doc.Error()
	// TODO: Return value
	return
}

func (this *MusicServices) UpdateAvailableServices() (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	response := this.Svc.CallVa("UpdateAvailableServices")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
