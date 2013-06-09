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
	GroupManagement_EventType = registerEventType("GroupManagement")
)

type GroupManagementState struct {
	GroupCoordinatorIsLocal bool
	LocalGroupUUID          string
	ResetVolumeAfter        bool
	VolumeAVTransportURI    string
}

type GroupManagementEvent struct {
	GroupManagementState
	Svc *Service
}

func (this GroupManagementEvent) Service() *Service {
	return this.Svc
}

func (this GroupManagementEvent) Type() int {
	return GroupManagement_EventType
}

type GroupManagement struct {
	GroupManagementState
	Svc *Service
}

func (this *GroupManagement) BeginSet(svc *Service, channel chan Event) {
}

type groupManagementUpdate_XML struct {
	XMLName xml.Name `xml:"GroupManagementState"`
	Value   string   `xml:",innerxml"`
}

func (this *GroupManagement) HandleProperty(svc *Service, value string, channel chan Event) error {
	update := groupManagementUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		xml.Unmarshal(bytes, &this.GroupManagementState)
	}
	return nil
}

func (this *GroupManagement) EndSet(svc *Service, channel chan Event) {
	evt := GroupManagementEvent{GroupManagementState: this.GroupManagementState, Svc: svc}
	channel <- evt
}

type MemberInfo struct {
	CurrentTransportSettings string
	GroupUUIDJoined          string
	ResetVolumeAfter         bool
	VolumeAVTransportURI     string
}

func (this *GroupManagement) AddMember(memberId string) (memberInfo *MemberInfo, err error) {
	type Response struct {
		XMLName xml.Name
		MemberInfo
		ErrorResponse
	}
	args := []Arg{
		{"MemberID", memberId},
	}
	response := this.Svc.Call("AddMember", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	memberInfo = &doc.MemberInfo
	err = doc.Error()
	return
}

func (this *GroupManagement) RemoveMember(memberId string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"MemberID", memberId},
	}
	response := this.Svc.Call("RemoveMember", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *GroupManagement) ReportTrackBufferingResult(memberId string, resultCode int32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"MemberID", memberId},
		{"ResultCode", resultCode},
	}
	response := this.Svc.Call("ReportTrackBufferingResult", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
