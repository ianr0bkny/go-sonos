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

type DeviceProperties struct {
	Svc *Service
}

const (
	LEDState_On  = "On"
	LEDState_Off = "Off"
)

func (this *DeviceProperties) SetLEDState(desiredLEDState string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"DesiredLEDState", desiredLEDState},
	}
	response := Call(this.Svc, "SetLEDState", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetLEDState() (currentLEDState string, err error) {
	type Response struct {
		XMLName         xml.Name
		CurrentLEDState string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetLEDState")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentLEDState = doc.CurrentLEDState
	err = doc.Error()
	return
}

func (this *DeviceProperties) SetInvisible(desiredInvisible bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"DesiredInvisible", desiredInvisible},
	}
	response := Call(this.Svc, "SetInvisible", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetInvisible() (currentInvisible bool, err error) {
	type Response struct {
		XMLName          xml.Name
		CurrentInvisible bool
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetInvisible")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentInvisible = doc.CurrentInvisible
	err = doc.Error()
	return
}

func (this *DeviceProperties) AddBondedZones(channelMapSet string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"ChannelMapSet", channelMapSet},
	}
	response := Call(this.Svc, "AddBondedZones", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) RemoveBondedZones(channelMapSet string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"ChannelMapSet", channelMapSet},
	}
	response := Call(this.Svc, "RemoveBondedZones", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) CreateStereoPair(channelMapSet string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"ChannelMapSet", channelMapSet},
	}
	response := Call(this.Svc, "CreateStereoPair", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) SeparateStereoPair(channelMapSet string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"ChannelMapSet", channelMapSet},
	}
	response := Call(this.Svc, "SeparateStereoPair", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) SetZoneAttributes(desiredZoneName, desiredIcon string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"DesiredZoneName,", desiredZoneName},
		{"DesiredIcon,", desiredIcon},
	}
	response := Call(this.Svc, "SetZoneAttributes", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetZoneAttributes() (currentZoneName, currentIcon string, err error) {
	type Response struct {
		XMLName         xml.Name
		CurrentZoneName string
		CurrentIcon     string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetZoneAttributes")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentZoneName = doc.CurrentZoneName
	currentIcon = doc.CurrentIcon
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetHouseholdID() (currentHouseholdId string, err error) {
	type Response struct {
		XMLName            xml.Name
		CurrentHouseholdID string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetHouseholdID")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentHouseholdId = doc.CurrentHouseholdID
	err = doc.Error()
	return
}

type ZoneInfo struct {
	SerialNumber           string
	SoftwareVersion        string
	DisplaySoftwareVersion string
	HardwareVersion        string
	IPAddress              string
	MACAddress             string
	CopyrightInfo          string
	ExtraInfo              string
}

func (this *DeviceProperties) GetZoneInfo() (zoneInfo *ZoneInfo, err error) {
	type Response struct {
		XMLName xml.Name
		ZoneInfo
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetZoneInfo")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	zoneInfo = &doc.ZoneInfo
	err = doc.Error()
	return
}

func (this *DeviceProperties) SetAutoplayLinkedZones(includeLinkedZones bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"IncludeLinkedZones", includeLinkedZones},
	}
	response := Call(this.Svc, "SetAutoplayLinkedZones", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetAutoplayLinkedZones() (includeLinkedZones bool, err error) {
	type Response struct {
		XMLName            xml.Name
		IncludeLinkedZones bool
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetAutoplayLinkedZones")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	includeLinkedZones = doc.IncludeLinkedZones
	err = doc.Error()
	return
}

func (this *DeviceProperties) SetAutoplayRoomUUID(roomUUID string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"RoomUUID", roomUUID},
	}
	response := Call(this.Svc, "SetAutoplayRoomUUID", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetAutoplayRoomUUID() (roomUUID string, err error) {
	type Response struct {
		XMLName  xml.Name
		RoomUUID string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetAutoplayRoomUUID")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	roomUUID = doc.RoomUUID
	err = doc.Error()
	return
}

func (this *DeviceProperties) SetAutoplayVolume(volume uint16) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"Volume", volume},
	}
	response := Call(this.Svc, "SetAutoplayVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetAutoplayVolume() (currentVolume uint16, err error) {
	type Response struct {
		XMLName       xml.Name
		CurrentVolume uint16
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetAutoplayVolume")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentVolume = doc.CurrentVolume
	err = doc.Error()
	return
}

func (this *DeviceProperties) ImportSetting(settingID uint32, settingURI string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"SettingID", settingID},
		{"SettingURI", settingURI},
	}
	response := Call(this.Svc, "ImportSettings", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) SetUseAutoplayVolume(useVolume bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"UseVolume", useVolume},
	}
	response := Call(this.Svc, "SetUseAutoplayVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *DeviceProperties) GetUseAutoplayVolume() (useVolume bool, err error) {
	type Response struct {
		XMLName   xml.Name
		UseVolume bool
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetUseAutoplayVolume")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	useVolume = doc.UseVolume
	err = doc.Error()
	return
}
