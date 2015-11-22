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
	RenderingControl_EventType = registerEventType("RenderingControl")
)

type renderingControl_Value_XML struct {
	XMLName xml.Name
	Channel string `xml:"channel,attr"`
	Val     string `xml:"val,attr"`
}

type renderingControl_InstanceID_XML struct {
	renderingControl_Value_XML
	Volume,
	Mute,
	Bass,
	Treble,
	Loudness,
	OutputFixed,
	HeadphoneConnected,
	SpeakerSize,
	SubGain,
	SubCrossover,
	SubPolarity,
	SubEnabled []renderingControl_Value_XML
	PresetNameList string
}

type renderingControl_Event_XML struct {
	XMLName    xml.Name
	InstanceID renderingControl_InstanceID_XML
}

type RenderingControlState struct {
	LastChange renderingControl_Event_XML
}

type RenderingControlEvent struct {
	RenderingControlState
	Svc *Service
}

func (this RenderingControlEvent) Service() *Service {
	return this.Svc
}

func (this RenderingControlEvent) Type() int {
	return RenderingControl_EventType
}

type RenderingControl struct {
	RenderingControlState
	Svc *Service
}

func (this *RenderingControl) BeginSet(svc *Service, channel chan Event) {
}

type renderingControlUpdate_XML struct {
	XMLName xml.Name `xml:"RenderingControlState"`
	Value   string   `xml:",innerxml"`
}

func (this *RenderingControl) HandleProperty(svc *Service, value string, channel chan Event) error {
	type Response struct {
		XMLName    xml.Name
		LastChange string
	}
	update := renderingControlUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		doc := Response{}
		xml.Unmarshal(bytes, &doc)
		xml.Unmarshal([]byte(doc.LastChange), &this.RenderingControlState.LastChange)
	}
	return nil
}

func (this *RenderingControl) EndSet(svc *Service, channel chan Event) {
	evt := RenderingControlEvent{RenderingControlState: this.RenderingControlState, Svc: svc}
	channel <- evt
}

const (
	Channel_Master = "Master"
	Channel_RF     = "RF"
	Channel_LF     = "LF"
)

func (this *RenderingControl) GetMute(instanceId uint32, channel string) (currentMute bool, err error) {
	type Response struct {
		XMLName     xml.Name
		CurrentMute bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
	}
	response := this.Svc.Call("GetMute", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentMute = doc.CurrentMute
	err = doc.Error()
	return
}

func (this *RenderingControl) SetMute(instanceId uint32, channel string, desiredMute bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"DesiredMute", desiredMute},
	}
	response := this.Svc.Call("SetMute", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type BasicEQ struct {
	Bass        int16
	Treble      int16
	Loudness    bool
	LeftVolume  uint16
	RightVolume uint16
}

func (this *RenderingControl) ResetBasicEQ(instanceId uint32) (basicEQ *BasicEQ, err error) {
	type Response struct {
		XMLName xml.Name
		BasicEQ
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("ResetBasicEQ", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	basicEQ = &doc.BasicEQ
	err = doc.Error()
	return
}

func (this *RenderingControl) ResetExtEQ(instanceId uint32, eqType string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"EQType", eqType},
	}
	response := this.Svc.Call("ResetExtEQ", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetVolume(instanceId uint32, channel string) (currentVolume uint16, err error) {
	type Response struct {
		XMLName       xml.Name
		CurrentVolume uint16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
	}
	response := this.Svc.Call("GetVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentVolume = doc.CurrentVolume
	err = doc.Error()
	return
}

//
// Set the playback volume.  For Sonos @instanceId will always be 0;
// @channel is one of the constants given in this file (e.g. Channel_Master);
// @volume is an integer between 0 and 100, where 100 is the loudest.
//
func (this *RenderingControl) SetVolume(instanceId uint32, channel string, volume uint16) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"DesiredVolume", volume},
	}
	response := this.Svc.Call("SetVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *RenderingControl) SetRelativeVolume(instanceId uint32, channel string, adjustment int32) (newVolume uint16, err error) {
	type Response struct {
		XMLName   xml.Name
		NewVolume uint16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"Adjustment", adjustment},
	}
	response := this.Svc.Call("SetRelativeVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	newVolume = doc.NewVolume
	err = doc.Error()
	return
}

func (this *RenderingControl) GetVolumeDB(instanceId uint32, channel string) (currentVolume int16, err error) {
	type Response struct {
		XMLName       xml.Name
		CurrentVolume int16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
	}
	response := this.Svc.Call("GetVolumeDB", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentVolume = doc.CurrentVolume
	err = doc.Error()
	return
}

func (this *RenderingControl) SetVolumeDB(instanceId uint32, channel string, volume int16) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"DesiredVolume", volume},
	}
	response := this.Svc.Call("SetVolumeDB", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetVolumeDBRange(instanceId uint32, channel string) (min, max int16, err error) {
	type Response struct {
		XMLName  xml.Name
		MinValue int16
		MaxValue int16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
	}
	response := this.Svc.Call("GetVolumeDBRange", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	min = doc.MinValue
	max = doc.MaxValue
	err = doc.Error()
	return
}

func (this *RenderingControl) GetBass(instanceId uint32) (currentBass int16, err error) {
	type Response struct {
		XMLName     xml.Name
		CurrentBass int16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("GetBass", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentBass = doc.CurrentBass
	err = doc.Error()
	return
}

func (this *RenderingControl) SetBass(instanceId, desiredBass int16) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"DesiredBass", desiredBass},
	}
	response := this.Svc.Call("SetBass", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetTreble(instanceId uint32) (currentTreble int16, err error) {
	type Response struct {
		XMLName       xml.Name
		CurrentTreble int16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("GetTreble", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentTreble = doc.CurrentTreble
	err = doc.Error()
	return
}

func (this *RenderingControl) SetTreble(instanceId, desiredTreble int16) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"DesiredTreble", desiredTreble},
	}
	response := this.Svc.Call("SetTreble", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetEQ(instanceId uint32, eqType string) (currentValue int16, err error) {
	type Response struct {
		XMLName      xml.Name
		CurrentValue int16
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"EQType", eqType},
	}
	response := this.Svc.Call("GetEQ", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentValue = doc.CurrentValue
	err = doc.Error()
	return
}

func (this *RenderingControl) SetEQ(instanceId uint32, eqType string, desiredValue int16) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"EQType", eqType},
		{"DesiredValue", desiredValue},
	}
	response := this.Svc.Call("SetEQ", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetLoudness(instanceId uint32, channel string) (loudness bool, err error) {
	type Response struct {
		XMLName         xml.Name
		CurrentLoudness bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
	}
	response := this.Svc.Call("GetLoudness", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	loudness = doc.CurrentLoudness
	err = doc.Error()
	return
}

func (this *RenderingControl) SetLoudness(instanceId uint32, channel string, loudness bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"DesiredLoudness", loudness},
	}
	response := this.Svc.Call("SetLoudness", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetSupportsOutputFixed(instanceId uint32) (currentSupportsFixed bool, err error) {
	type Response struct {
		XMLName              xml.Name
		CurrentSupportsFixed bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("GetSupportsOutputFixed", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentSupportsFixed = doc.CurrentSupportsFixed
	err = doc.Error()
	return
}

func (this *RenderingControl) GetOutputFixed(instanceId uint32) (currentFixed bool, err error) {
	type Response struct {
		XMLName      xml.Name
		CurrentFixed bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("GetOutputFixed", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentFixed = doc.CurrentFixed
	err = doc.Error()
	return
}

func (this *RenderingControl) SetOutputFixed(instanceId uint32, desiredFixed bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"DesiredFixed", desiredFixed},
	}
	response := this.Svc.Call("SetOutputFixed", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) GetHeadphoneConnected(instanceId uint32) (currentHeadphoneConnected bool, err error) {
	type Response struct {
		XMLName                   xml.Name
		CurrentHeadphoneConnected bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("GetHeadphoneConnected", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentHeadphoneConnected = doc.CurrentHeadphoneConnected
	err = doc.Error()
	return
}

const (
	RampType_SleepTimer = "SLEEP_TIMER_RAMP_TYPE"
	RampType_Alarm      = "ALARM_RAMP_TYPE"
	RampType_Autoplay   = "AUTOPLAY_RAMP_TYPE"
)

type RampRequest struct {
	RampType         string
	DesiredVolume    uint16
	ResetVolumeAfter bool
	ProgramURI       string
}

func (this *RenderingControl) RampToVolume(instanceId uint32, channel, req RampRequest) (rampTime uint32, err error) {
	type Response struct {
		XMLName  xml.Name
		RampTime uint32
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"RampType", req.RampType},
		{"DesiredVolume", req.DesiredVolume},
		{"ResetVolumeAfter", req.ResetVolumeAfter},
		{"ProgramURI", req.ProgramURI},
	}
	response := this.Svc.Call("RampToVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	rampTime = doc.RampTime
	err = doc.Error()
	return
}

func (this *RenderingControl) RestoreVolumePriorToRamp(instanceId uint32, channel string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
	}
	response := this.Svc.Call("RestoreVolumePriorToRamp", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *RenderingControl) SetChannelMap(instanceId uint32, channelMap string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"ChannelMap", channelMap},
	}
	response := this.Svc.Call("SetChannelMap", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

/* Reciva */
/* NSZ-GS7 */
func (this *RenderingControl) ListPresets(instanceId uint32) (presets string, err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
		CurrentPresetNameList string
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := this.Svc.Call("ListPresets", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.CurrentPresetNameList, doc.Error()
}

/* Reciva */
/* NSZ-GS7 */
func (this *RenderingControl) SelectPreset(instanceId uint32, presetName string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"PresetName", presetName},
	}
	response := this.Svc.Call("SelectPreset", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *RenderingControl) SetSonarCalibrationX(instanceId uint32, calibrationId, sonarCoefficients string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"CalibrationID", calibrationId},
		{"SonarCoefficients", sonarCoefficients},
	}
	response := this.Svc.Call("SetSonarCalibrationX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *RenderingControl) GetSonarStatus(instanceId uint32) (sonarEnabled, sonarCalibrationAvailable bool, err error) {
	type Response struct {
		XMLName xml.Name
		SonarEnabled bool
		SonarCalibrationAvailable bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response:= this.Svc.Call("GetSonarStatus", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	sonarEnabled = doc.SonarEnabled
	sonarCalibrationAvailable = doc.SonarCalibrationAvailable
	err = doc.Error()
	return
}

func (this *RenderingControl) SetSonarStatus(instanceId uint32, sonarEnabled bool) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"SonarEnabled", sonarEnabled},
	}
	response:= this.Svc.Call("SetSonarStatus", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}
