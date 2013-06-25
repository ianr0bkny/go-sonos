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

package reciva

import (
	"encoding/xml"
	"github.com/ianr0bkny/go-sonos/upnp"
)

type RecivaRadio struct {
	Svc *upnp.Service
}

func (this *RecivaRadio) BeginSet(svc *upnp.Service, channel chan upnp.Event) {
}

func (this *RecivaRadio) HandleProperty(svc *upnp.Service, value string, channel chan upnp.Event) error {
	return nil
}

func (this *RecivaRadio) EndSet(svc *upnp.Service, channel chan upnp.Event) {
}

func (this *RecivaRadio) GetDateTime() (retDataTimeValue string, err error) {
	type Response struct {
		XMLName          xml.Name
		RetDateTimeValue string
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("GetDateTime")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetDateTimeValue, doc.Error()
}

func (this *RecivaRadio) GetTimeZone() (retTimeZoneValue string, err error) {
	type Response struct {
		XMLName          xml.Name
		RetTimeZoneValue string
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("GetTimeZone")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetTimeZoneValue, doc.Error()
}

func (this *RecivaRadio) GetNumberOfPresets() (retNumberOfPresetsValue uint32, err error) {
	type Response struct {
		XMLName                 xml.Name
		RetNumberOfPresetsValue uint32
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("GetNumberOfPresets")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetNumberOfPresetsValue, doc.Error()
}

func (this *RecivaRadio) GetDisplayLanguages() (retLanguageListValue, retIsoCodeListValue string, err error) {
	type Response struct {
		XMLName              xml.Name
		RetLanguageListValue string
		RetIsoCodeListValue  string
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("GetDisplayLanguages")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetLanguageListValue, doc.RetIsoCodeListValue, doc.Error()
}

func (this *RecivaRadio) GetCurrentDisplayLanguage() (retLanguageValue, retIsoCodeValue string, err error) {
	type Response struct {
		XMLName          xml.Name
		RetLanguageValue string
		RetIsoCodeValue  string
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("GetCurrentDisplayLanguage")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetLanguageValue, doc.RetIsoCodeValue, doc.Error()
}

func (this *RecivaRadio) GetPowerState() (retPowerStateValue string, err error) {
	type Response struct {
		XMLName            xml.Name
		RetPowerStateValue string
		upnp.ErrorResponse
	}
	response := this.Svc.CallVa("GetPowerState")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetPowerStateValue, doc.Error()
}

func (this *RecivaRadio) SetPowerState(newPowerStateValue string) (retPowerStateValue string, err error) {
	type Response struct {
		XMLName            xml.Name
		RetPowerStateValue string
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"NewPowerStateValue", newPowerStateValue},
	}
	response := this.Svc.Call("SetPowerState", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RetPowerStateValue, doc.Error()
}
