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

type AlarmClock struct {
	Svc *upnp.Service
}

func (this *AlarmClock) SetFormat(desiredTimeFormat, desiredDateFormat string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"DesiredTimeFormat", desiredTimeFormat},
		{"DesiredDateFormat", desiredDateFormat},
	}
	response := upnp.Call(this.Svc, "SetFormat", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) GetFormat() (currentTimeFormat, currentDateFormat string, err error) {
	type Response struct {
		XMLName           xml.Name
		CurrentTimeFormat string
		CurrentDateFormat string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetFormat")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentTimeFormat = doc.CurrentTimeFormat
	currentDateFormat = doc.CurrentDateFormat
	err = doc.Error()
	return
}

func (this *AlarmClock) SetTimeZone(index int32, autoAdjustDst bool) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"Index", index},
		{"AutoAdjustDst", autoAdjustDst},
	}
	response := upnp.Call(this.Svc, "SetTimeZone", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) GetTimeZone() (index int32, autoAdjustDst bool, err error) {
	type Response struct {
		XMLName       xml.Name
		Index         int32
		AutoAdjustDst bool
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetTimeZone")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	index = doc.Index
	autoAdjustDst = doc.AutoAdjustDst
	err = doc.Error()
	return
}

func (this *AlarmClock) GetTimeZoneAndRule() (index int32, autoAdjustDst bool, timeZone string, err error) {
	type Response struct {
		XMLName       xml.Name
		Index         int32
		AutoAdjustDst bool
		TimeZone      string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetTimeZoneAndRule")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	index = doc.Index
	autoAdjustDst = doc.AutoAdjustDst
	timeZone = doc.TimeZone
	err = doc.Error()
	return
}

func (this *AlarmClock) GetTimeZoneRule(index int32) (timeZone string, err error) {
	type Response struct {
		XMLName  xml.Name
		TimeZone string
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"Index", index},
	}
	response := upnp.Call(this.Svc, "GetTimeZoneRule", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	timeZone = doc.TimeZone
	err = doc.Error()
	return
}

func (this *AlarmClock) SetTimeServer(desiredTimeServer string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"DesiredTimeServer", desiredTimeServer},
	}
	response := upnp.Call(this.Svc, "SetTimeServer", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) GetTimeServer() (currentTimeServer string, err error) {
	type Response struct {
		XMLName           xml.Name
		CurrentTimeServer string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetTimeServer")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentTimeServer = doc.CurrentTimeServer
	err = doc.Error()
	return
}

func (this *AlarmClock) SetTimeNow(desiredTime, timeZoneForDesiredTime string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"DesiredTime", desiredTime},
		{"TimeZoneForDesiredTime", timeZoneForDesiredTime},
	}
	response := upnp.Call(this.Svc, "SetTimeNow", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) GetHouseholdTimeAtStamp(timeStamp string) (householdUTCTime string, err error) {
	type Response struct {
		XMLName          xml.Name
		HouseholdUTCTime string
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"TimeStamp", timeStamp},
	}
	response := upnp.Call(this.Svc, "GetHouseholdTimeAtStamp", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	householdUTCTime = doc.HouseholdUTCTime
	err = doc.Error()
	return
}

type GetTimeNowResponse struct {
	CurrentUTCTime        string
	CurrentLocalTime      string
	CurrentTimeZone       string
	CurrentTimeGeneration uint32
}

func (this *AlarmClock) GetTimeNow() (getTimeNowResponse *GetTimeNowResponse, err error) {
	type Response struct {
		XMLName xml.Name
		GetTimeNowResponse
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetTimeNow")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	getTimeNowResponse = &doc.GetTimeNowResponse
	err = doc.Error()
	return
}

const (
	Recurrence_ONCE     = "ONCE"
	Recurrence_WEEKDAYS = "WEEKDAYS"
	Recurrence_WEEKENDS = "WEEKENDS"
	Recurrence_DAILY    = "DAILY"
)

const (
	AlarmPlayMode_NORMAL           = "NORMAL"
	AlarmPlayMode_REPEAT_ALL       = "REPEAT_ALL"
	AlarmPlayMode_SHUFFLE_NOREPEAT = "SHUFFLE_NOREPEAT"
	AlarmPlayMode_SHUFFLE          = "SHUFFLE"
)

type CreateAlarmRequest struct {
	StartLocalTime     string
	Duration           string
	Recurrence         string
	Enabled            bool
	RoomUUID           string
	ProgramURI         string
	ProgramMetaData    string
	PlayMode           string
	Volume             uint16
	IncludeLinkedZones bool
}

func (this *AlarmClock) CreateAlarm(req *CreateAlarmRequest) (assignedId uint32, err error) {
	type Response struct {
		XMLName    xml.Name
		AssignedID uint32
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"StartLocalTime", req.StartLocalTime},
		{"Duration", req.Duration},
		{"Recurrence", req.Recurrence},
		{"Enabled", req.Enabled},
		{"RoomUUID", req.RoomUUID},
		{"ProgramURI", req.ProgramURI},
		{"ProgramMetaData", req.ProgramMetaData},
		{"PlayMode", req.PlayMode},
		{"Volume", req.Volume},
		{"IncludeLinkedZones", req.IncludeLinkedZones},
	}
	response := upnp.Call(this.Svc, "CreateAlarm", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	assignedId = doc.AssignedID
	err = doc.Error()
	return
}

type UpdateAlarmRequest CreateAlarmRequest

func (this *AlarmClock) UpdateAlarm(id uint32, req *UpdateAlarmRequest) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"ID", id},
		{"StartLocalTime", req.StartLocalTime},
		{"Duration", req.Duration},
		{"Recurrence", req.Recurrence},
		{"Enabled", req.Enabled},
		{"RoomUUID", req.RoomUUID},
		{"ProgramURI", req.ProgramURI},
		{"ProgramMetaData", req.ProgramMetaData},
		{"PlayMode", req.PlayMode},
		{"Volume", req.Volume},
		{"IncludeLinkedZones", req.IncludeLinkedZones},
	}
	response := upnp.Call(this.Svc, "UpdateAlarm", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) DestroyAlarm(id uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"ID", id},
	}
	response := upnp.Call(this.Svc, "DestroyAlarm", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) ListAlarms() (currentAlarmList, currentAlarmListVersion string, err error) {
	type Response struct {
		XMLName                 xml.Name
		CurrentAlarmList        string
		CurrentAlarmListVersion string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "ListAlarms")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentAlarmList = doc.CurrentAlarmList
	currentAlarmListVersion = doc.CurrentAlarmListVersion
	err = doc.Error()
	return
}

func (this *AlarmClock) SetDailyIndexRefreshTime(desiredDailyIndexRefreshTime string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"DesiredDailyIndexRefreshTime", desiredDailyIndexRefreshTime},
	}
	response := upnp.Call(this.Svc, "SetDailyIndexRefreshTime", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AlarmClock) GetDailyIndexRefreshTime() (currentDailyIndexRefreshTime string, err error) {
	type Response struct {
		XMLName                      xml.Name
		CurrentDailyIndexRefreshTime string
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "GetDailyIndexRefreshTime")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	currentDailyIndexRefreshTime = doc.CurrentDailyIndexRefreshTime
	err = doc.Error()
	return
}
