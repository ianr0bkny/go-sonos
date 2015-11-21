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
	DeviceProperties_EventType = registerEventType("DeviceProperties")
)

type DevicePropertiesState struct {
	SettingsReplicationState string
	ZoneName                 string
	Icon                     string
	Configuration            string
	Invisible                bool
	IsZoneBridge             bool
	ChannelMapSet            string
	HTSatChanMapSet          string
	HTFreq                   uint32
}

type DevicePropertiesEvent struct {
	DevicePropertiesState
	Svc *Service
}

func (this DevicePropertiesEvent) Service() *Service {
	return this.Svc
}

func (this DevicePropertiesEvent) Type() int {
	return DeviceProperties_EventType
}

type DeviceProperties struct {
	DevicePropertiesState
	Svc *Service
}

func (this *DeviceProperties) BeginSet(svc *Service, channel chan Event) {
}

type devicePropertiesUpdate_XML struct {
	XMLName xml.Name `xml:"DevicePropertiesState"`
	Value   string   `xml:",innerxml"`
}

func (this *DeviceProperties) HandleProperty(svc *Service, value string, channel chan Event) error {
	update := devicePropertiesUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		xml.Unmarshal(bytes, &this.DevicePropertiesState)
	}
	return nil
}

func (this *DeviceProperties) EndSet(svc *Service, channel chan Event) {
	evt := DevicePropertiesEvent{DevicePropertiesState: this.DevicePropertiesState, Svc: svc}
	channel <- evt
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
	response := this.Svc.Call("SetLEDState", args)
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
	response := this.Svc.CallVa("GetLEDState")
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
	response := this.Svc.Call("SetInvisible", args)
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
	response := this.Svc.CallVa("GetInvisible")
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
	response := this.Svc.Call("AddBondedZones", args)
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
	response := this.Svc.Call("RemoveBondedZones", args)
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
	response := this.Svc.Call("CreateStereoPair", args)
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
	response := this.Svc.Call("SeparateStereoPair", args)
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
	response := this.Svc.Call("SetZoneAttributes", args)
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
	response := this.Svc.CallVa("GetZoneAttributes")
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
	response := this.Svc.CallVa("GetHouseholdID")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentHouseholdId = doc.CurrentHouseholdID
	err = doc.Error()
	return
}

//
// The return value for the GetZoneInfo method
//
type ZoneInfo struct {
	// Appliance serial number
	SerialNumber string
	// Software version string
	SoftwareVersion string
	// Display software version string
	DisplaySoftwareVersion string
	// Hardware version
	HardwareVersion string
	// the IP address of the appliance
	IPAddress string
	// The hardware MAC address of the appliance
	MACAddress string
	// The Sonos Copyright statement
	CopyrightInfo string
	// ???
	ExtraInfo string
}

//
// Fetches basic properties of the appliance including IP address,
// MAC address, and relevant hardware and software version.
//
func (this *DeviceProperties) GetZoneInfo() (*ZoneInfo, error) {
	type Response struct {
		XMLName xml.Name
		ZoneInfo
		ErrorResponse
	}
	response := this.Svc.CallVa("GetZoneInfo")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.ZoneInfo, doc.Error()
}

func (this *DeviceProperties) SetAutoplayLinkedZones(includeLinkedZones bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"IncludeLinkedZones", includeLinkedZones},
	}
	response := this.Svc.Call("SetAutoplayLinkedZones", args)
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
	response := this.Svc.CallVa("GetAutoplayLinkedZones")
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
	response := this.Svc.Call("SetAutoplayRoomUUID", args)
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
	response := this.Svc.CallVa("GetAutoplayRoomUUID")
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
	response := this.Svc.Call("SetAutoplayVolume", args)
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
	response := this.Svc.CallVa("GetAutoplayVolume")
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
	response := this.Svc.Call("ImportSettings", args)
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
	response := this.Svc.Call("SetUseAutoplayVolume", args)
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
	response := this.Svc.CallVa("GetUseAutoplayVolume")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	useVolume = doc.UseVolume
	err = doc.Error()
	return
}

func (this *DeviceProperties) AddHTSatellite(htSatChanMapSet string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"HTSatChanMapSet", htSatChanMapSet},
	}
	response := this.Svc.Call("AddHTSatellite", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *DeviceProperties) RemoveHTSatellite(satRoomUUID string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"SatRoomUUID", satRoomUUID},
	}
	response := this.Svc.Call("RemoveHTSatellite", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *DeviceProperties) EnterConfigMode(mode, options string) (state string, err error) {
	type Response struct {
		XMLName xml.Name
		State string
		ErrorResponse
	}
	args := []Arg{
		{"Mode", mode},
		{"Options", options},
	}
	response := this.Svc.Call("EnterConfigMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	state = doc.State
	err = doc.Error()
	return
}

func (this *DeviceProperties) ExitConfigMode(options string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"Options", options},
	}
	response := this.Svc.Call("ExitConfigMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *DeviceProperties) GetButtonState() (state string, err error) {
	type Response struct {
		XMLName xml.Name
		State string
		ErrorResponse
	}
	response := this.Svc.CallVa("GetButtonState")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	state = doc.State
	err = doc.Error()
	return
}
