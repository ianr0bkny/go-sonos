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
	ZoneGroupTopology_EventType = registerEventType("ZoneGroupTopology")
)

type ZoneGroupTopologyState struct {
	ZoneGroupState          string
	ThirdPartyMediaServersX string
	AvailableSoftwareUpdate string // TODO: Unpack
	AlarmRunSequence        string
	ZoneGroupName           string
	ZoneGroupID             string
	ZonePlayerUUIDsInGroup  string
}

type ZoneGroupTopologyEvent struct {
	ZoneGroupTopologyState
	Svc *Service
}

func (this ZoneGroupTopologyEvent) Service() *Service {
	return this.Svc
}

func (this ZoneGroupTopologyEvent) Type() int {
	return ZoneGroupTopology_EventType
}

type ZoneGroupTopology struct {
	ZoneGroupTopologyState
	Svc *Service
}

func (this *ZoneGroupTopology) BeginSet(svc *Service, channel chan Event) {
}

type zoneGroupTopologyUpdate_XML struct {
	XMLName xml.Name `xml:"ZoneGroupTopologyState"`
	Value   string   `xml:",innerxml"`
}

func (this *ZoneGroupTopology) HandleProperty(svc *Service, value string, channel chan Event) error {
	update := zoneGroupTopologyUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		xml.Unmarshal(bytes, &this.ZoneGroupTopologyState)
	}
	return nil
}

func (this *ZoneGroupTopology) EndSet(svc *Service, channel chan Event) {
	evt := ZoneGroupTopologyEvent{ZoneGroupTopologyState: this.ZoneGroupTopologyState, Svc: svc}
	channel <- evt
}

const (
	ALL      = "All"
	SOFTWARE = "Software"
)

func (this *ZoneGroupTopology) BeginSoftwareUpdate(updateURL string, flags uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"UpdateURL", updateURL},
		{"Flags", flags},
	}
	response := this.Svc.Call("BeginSoftwareUpdate", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type UpdateItem struct {
	Type         string `xml:"Type,attr"`
	Version      string `xml:"Version,attr"`
	UpdateURL    string `xml:"UpdateURL,attr"`
	DownloadSize string `xml:"DownloadSize,attr"`
	ManifestURL  string `xml:"ManifestURL,attr"`
}

type UpdateType string

func (this *ZoneGroupTopology) CheckForUpdate(updateType UpdateType, cachedOnly bool, version string) (updateItem *UpdateItem, err error) {
	type UpdateItemHolder struct {
		XMLName xml.Name
		UpdateItem
	}
	type UpdateItemText struct {
		XMLName xml.Name
		Text    string `xml:",chardata"`
	}
	type Response struct {
		XMLName    xml.Name
		UpdateItem UpdateItemText
		ErrorResponse
	}
	args := []Arg{
		{"UpdateType", updateType},
		{"CachedOnly", cachedOnly},
		{"Version", version},
	}
	response := this.Svc.Call("CheckForUpdate", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	rec := UpdateItemHolder{}
	xml.Unmarshal([]byte(doc.UpdateItem.Text), &rec)
	updateItem = &rec.UpdateItem
	err = doc.Error()
	return
}

const (
	REMOVE                        = "Remove"
	VERIFY_THEN_REMOVE_SYSTEMWIDE = "VerifyThenRemoveSystemwide"
)

func (this *ZoneGroupTopology) ReportUnresponsiveDevice(deviceUUID string, desiredAction string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"DeviceUUID", deviceUUID},
		{"DesiredAction", desiredAction},
	}
	response := this.Svc.Call("ReportUnresponsiveDevice", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ZoneGroupTopology) ReportAlarmStartedRunning() (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	response := this.Svc.CallVa("ReportAlarmStartedRunning")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ZoneGroupTopology) SubmitDiagnostics() (diagnosticId string, err error) {
	type Response struct {
		XMLName      xml.Name
		DiagnosticID string
		ErrorResponse
	}
	response := this.Svc.CallVa("SubmitDiagnostics")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	diagnosticId = doc.DiagnosticID
	err = doc.Error()
	return
}

func (this *ZoneGroupTopology) RegisterMobileDevice(deviceName, deviceUDN, deviceAddress string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"MobileDeviceName", deviceName},
		{"MobileDeviceUDN", deviceUDN},
		{"MobileIPAndPort", deviceAddress},
	}
	response := this.Svc.Call("RegisterMobileDevice", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type ZoneGroupAttributes struct {
	CurrentZoneGroupName   string
	CurrentZoneGroupID     string
	ZonePlayerUUIDsInGroup string
}

func (this *ZoneGroupTopology) GetZoneGroupAttributes() (*ZoneGroupAttributes, error) {
	type Response struct {
		XMLName xml.Name
		ZoneGroupAttributes
		ErrorResponse
	}
	response := this.Svc.CallVa("GetZoneGroupAttributes")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.ZoneGroupAttributes, doc.Error()
}

type ZoneGroupMember struct {
	XMLName              xml.Name
	UUID                 string `xml:"UUID,attr"`
	Location             string `xml:"Location,attr"`
	ZoneName             string `xml:"ZoneName,attr"`
	Icon                 string `xml:"Icon,attr"`
	Configuration        string `xml:"Configuration,attr"`
	Invisible            string `xml:"Invisible,attr"`
	IsZoneBridge         string `xml:"IsZoneBridge,attr"`
	SoftwareVersion      string `xml:"SoftwareVersion,attr"`
	MinCompatibleVersion string `xml:"MinCompatibleVersion,attr"`
	BootSeq              string `xml:"BootSeq,attr"`
}

type ZoneGroup struct {
	XMLName         xml.Name
	Coordinator     string `xml:"Coordinator,attr"`
	ID              string `xml:"ID,attr"`
	ZoneGroupMember []ZoneGroupMember
}

type ZoneGroups struct {
	XMLName   xml.Name
	ZoneGroup []ZoneGroup
}

func (this *ZoneGroupTopology) GetZoneGroupState() (*ZoneGroups, error) {
	type Response struct {
		XMLName        xml.Name
		ZoneGroupState string
		ErrorResponse
	}
	response := this.Svc.CallVa("GetZoneGroupState")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	state := ZoneGroups{}
	xml.Unmarshal([]byte(doc.ZoneGroupState), &state)
	return &state, doc.Error()
}
