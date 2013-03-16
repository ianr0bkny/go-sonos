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

type ZoneGroupTopology struct {
	Svc *Service
}

const (
	ALL      = "All"
	SOFTWARE = "Software"
)

func (this *ZoneGroupTopology) BeginSoftwareUpdate(updateURL string, flags uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		upnpErrorResponse
	}
	args := []upnpArg{
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
		upnpErrorResponse
	}
	args := []upnpArg{
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
		upnpErrorResponse
	}
	args := []upnpArg{
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
		upnpErrorResponse
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
		upnpErrorResponse
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
		upnpErrorResponse
	}
	args := []upnpArg{
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
		upnpErrorResponse
	}
	response := this.Svc.CallVa("GetZoneGroupAttributes")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.ZoneGroupAttributes, doc.Error()
}
