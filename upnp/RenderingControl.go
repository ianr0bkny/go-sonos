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

type RenderingControl struct {
	Svc *Service
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
	response := Call(this.Svc, "GetMute", args)
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
	response := Call(this.Svc, "SetMute", args)
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
	response := Call(this.Svc, "ResetBasicEQ", args)
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
	response := Call(this.Svc, "ResetExtEQ", args)
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
	response := Call(this.Svc, "GetVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentVolume = doc.CurrentVolume
	err = doc.Error()
	return
}

func (this *RenderingControl) SetVolume(instanceId uint32, channel string, volume uint16) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Channel", channel},
		{"DesiredVolume", volume},
	}
	response := Call(this.Svc, "SetVolume", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
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
	response := Call(this.Svc, "SetRelativeVolume", args)
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
	response := Call(this.Svc, "GetVolumeDB", args)
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
	response := Call(this.Svc, "SetVolumeDB", args)
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
	response := Call(this.Svc, "GetVolumeDBRange", args)
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
	response := Call(this.Svc, "GetBass", args)
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
	response := Call(this.Svc, "SetBass", args)
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
	response := Call(this.Svc, "GetTreble", args)
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
	response := Call(this.Svc, "SetTreble", args)
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
	response := Call(this.Svc, "GetEQ", args)
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
	response := Call(this.Svc, "SetEQ", args)
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
	response := Call(this.Svc, "GetLoudness", args)
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
	response := Call(this.Svc, "SetLoudness", args)
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
	response := Call(this.Svc, "GetSupportsOutputFixed", args)
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
	response := Call(this.Svc, "GetOutputFixed", args)
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
	response := Call(this.Svc, "SetOutputFixed", args)
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
	response := Call(this.Svc, "GetHeadphoneConnected", args)
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
	response := Call(this.Svc, "RampToVolume", args)
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
	response := Call(this.Svc, "RestoreVolumePriorToRamp", args)
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
	response := Call(this.Svc, "SetChannelMap", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
