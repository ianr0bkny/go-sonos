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

type ZoneGroupTopology struct {
	Svc *upnp.Service
}

const (
	ALL      = "All"
	SOFTWARE = "Software"
)

func (this *ZoneGroupTopology) BeginSoftwareUpdate(updateURL string, flags uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"UpdateURL", updateURL},
		{"Flags", flags},
	}
	response := upnp.Call(this.Svc, "BeginSoftwareUpdate", args)
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
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"UpdateType", updateType},
		{"CachedOnly", cachedOnly},
		{"Version", version},
	}
	response := upnp.Call(this.Svc, "CheckForUpdate", args)
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
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"DeviceUUID", deviceUUID},
		{"DesiredAction", desiredAction},
	}
	response := upnp.Call(this.Svc, "ReportUnresponsiveDevice", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ZoneGroupTopology) ReportAlarmStartedRunning() (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "ReportAlarmStartedRunning")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ZoneGroupTopology) SubmitDiagnostics() (diagnosticId string, err error) {
	type Response struct {
		XMLName      xml.Name
		DiagnosticID string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "SubmitDiagnostics")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	diagnosticId = doc.DiagnosticID
	err = doc.Error()
	return
}

func (this *ZoneGroupTopology) RegisterMobileDevice(deviceName, deviceUDN, deviceAddress string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"MobileDeviceName", deviceName},
		{"MobileDeviceUDN", deviceUDN},
		{"MobileIPAndPort", deviceAddress},
	}
	response := upnp.Call(this.Svc, "RegisterMobileDevice", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
