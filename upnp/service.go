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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type upnpChardataValue_XML struct {
	Chardata string `xml:",chardata"`
}

type upnpAllowedValueList_XML struct {
	AllowedValue []upnpChardataValue_XML `xml:"allowedValue"`
}

type upnpAllowedValueRange_XML struct {
	Minimum []upnpChardataValue_XML `xml:"minimum"`
	Maximum []upnpChardataValue_XML `xml:"maximum"`
	Step    []upnpChardataValue_XML `xml:"step"`
}

type upnpStateVariable_XML struct {
	SendEvents        string                      `xml:"sendEvents,attr"`
	Name              []upnpChardataValue_XML     `xml:"name"`
	DataType          []upnpChardataValue_XML     `xml:"dataType"`
	AllowedValueList  []upnpAllowedValueList_XML  `xml:"allowedValueList"`
	AllowedValueRange []upnpAllowedValueRange_XML `xml:"allowedValueRange"`
}

type upnpServiceStateTable_XML struct {
	StateVariable []upnpStateVariable_XML `xml:"stateVariable"`
}

type upnpActionArgument_XML struct {
	Name                 []upnpChardataValue_XML `xml:"name"`
	Direction            []upnpChardataValue_XML `xml:"direction"`
	RelatedStateVariable []upnpChardataValue_XML `xml:"relatedStateVariable"`
}

type upnpActionArgumentList_XML struct {
	Argument []upnpActionArgument_XML `xml:"argument"`
}

type upnpAction_XML struct {
	Name         []upnpChardataValue_XML      `xml:"name"`
	ArgumentList []upnpActionArgumentList_XML `xml:"argumentList"`
}

type upnpActionList_XML struct {
	Action []upnpAction_XML `xml:"action"`
}

type upnpDescribeService_XML struct {
	XMLNamespace      string                      `xml:"xmlns,attr"`
	SpecVersion       []upnpSpecVersion_XML       `xml:"specVersion"`
	ServiceStateTable []upnpServiceStateTable_XML `xml:"serviceStateTable"`
	ActionList        []upnpActionList_XML        `xml:"actionList"`
}

type upnpDescribeServiceJob struct {
	result     chan *Service
	err_result chan error
	response   *http.Response
	doc        upnpDescribeService_XML
	svc        *Service
}

func upnpMakeDescribeServiceJob(svc *Service) (job *upnpDescribeServiceJob) {
	job = &upnpDescribeServiceJob{}
	job.result = make(chan *Service)
	job.err_result = make(chan error)
	job.svc = svc
	job.doc = upnpDescribeService_XML{}
	return
}

func (this *upnpDescribeServiceJob) UnpackChardataValue(val *upnpChardataValue_XML) (s string) {
	return val.Chardata
}

func (this *upnpDescribeServiceJob) UnpackAllowedValueList(val_list *upnpAllowedValueList_XML) (allowed_list []string) {
	for _, value := range val_list.AllowedValue {
		allowed_list = append(allowed_list, this.UnpackChardataValue(&value))
	}
	return
}

type upnpValueRange struct {
	min  string
	max  string
	step string
}

func (this *upnpDescribeServiceJob) UnpackAllowedValueRange(val_range *upnpAllowedValueRange_XML) (vrange *upnpValueRange) {
	vrange = &upnpValueRange{}
	for _, min := range val_range.Minimum {
		vrange.min = this.UnpackChardataValue(&min)
	}
	for _, max := range val_range.Maximum {
		vrange.max = this.UnpackChardataValue(&max)
	}
	for _, step := range val_range.Step {
		vrange.step = this.UnpackChardataValue(&step)
	}
	return
}

type upnpStateVariable struct {
	name          string
	dataType      string
	allowedValues []string
	allowedRange  *upnpValueRange
}

func (this *upnpDescribeServiceJob) UnpackStateVariable(v *upnpStateVariable_XML) (sv *upnpStateVariable) {
	sv = &upnpStateVariable{}
	for _, name := range v.Name {
		sv.name = this.UnpackChardataValue(&name)
	}
	for _, datatype := range v.DataType {
		sv.dataType = this.UnpackChardataValue(&datatype)
	}
	for _, val_list := range v.AllowedValueList {
		sv.allowedValues = this.UnpackAllowedValueList(&val_list)
	}
	for _, val_range := range v.AllowedValueRange {
		sv.allowedRange = this.UnpackAllowedValueRange(&val_range)
	}
	return
}

func (this *upnpDescribeServiceJob) UnpackStateTable(tab *upnpServiceStateTable_XML) (table []*upnpStateVariable) {
	for _, v := range tab.StateVariable {
		table = append(table, this.UnpackStateVariable(&v))
	}
	return
}

type upnpActionArgument struct {
	name     string
	dir      string
	variable string
}

func (this *upnpDescribeServiceJob) UnpackActionArgument(arg *upnpActionArgument_XML) (aarg *upnpActionArgument) {
	aarg = &upnpActionArgument{}
	for _, name := range arg.Name {
		aarg.name = this.UnpackChardataValue(&name)
	}
	for _, dir := range arg.Direction {
		aarg.dir = this.UnpackChardataValue(&dir)
	}
	for _, related := range arg.RelatedStateVariable {
		aarg.variable = this.UnpackChardataValue(&related)
	}
	return
}

func (this *upnpDescribeServiceJob) UnpackArgumentList(arg_list *upnpActionArgumentList_XML) (aarg_list []*upnpActionArgument) {
	for _, arg := range arg_list.Argument {
		aarg_list = append(aarg_list, this.UnpackActionArgument(&arg))
	}
	return
}

type upnpAction struct {
	name    string
	argList []*upnpActionArgument
}

func (this *upnpDescribeServiceJob) UnpackAction(action *upnpAction_XML) (act *upnpAction) {
	act = &upnpAction{}
	for _, name := range action.Name {
		act.name = this.UnpackChardataValue(&name)
	}
	for _, arg_list := range action.ArgumentList {
		act.argList = this.UnpackArgumentList(&arg_list)
	}
	return
}

func (this *upnpDescribeServiceJob) UnpackActionList(act_list *upnpActionList_XML) (action_list []*upnpAction) {
	for _, action := range act_list.Action {
		action_list = append(action_list, this.UnpackAction(&action))
	}
	return
}

type Service struct {
	deviceURI      string
	deviceType     string
	deviceVersion  string
	udn            string
	serviceURI     string
	serviceType    string
	serviceVersion string
	serviceId      string
	controlURL     *url.URL
	eventSubURL    *url.URL
	scpdURL        *url.URL
	described      bool
	stateTable     []*upnpStateVariable
	actionList     []*upnpAction
}

func (this *Service) Actions() (actions []string) {
	for _, action := range this.actionList {
		actions = append(actions, action.name)
	}
	return
}

func upnpMakeService() (svc *Service) {
	return &Service{}
}

func (this *upnpDescribeServiceJob) Unpack() {
	for _, tab := range this.doc.ServiceStateTable {
		this.svc.stateTable = this.UnpackStateTable(&tab)
	}
	for _, act_list := range this.doc.ActionList {
		this.svc.actionList = this.UnpackActionList(&act_list)
	}
	return
}

func (this *upnpDescribeServiceJob) Parse() {
	defer this.response.Body.Close()
	if body, err := ioutil.ReadAll(this.response.Body); nil == err {
		xml.Unmarshal(body, &this.doc)
		this.Unpack()
		this.result <- this.svc
	} else {
		this.err_result <- err
	}
}

func (this *upnpDescribeServiceJob) Describe() {
	var err error
	uri := this.svc.scpdURL.String()
	log.Printf("Loading %s", string(uri))
	if this.response, err = http.Get(string(uri)); nil == err {
		this.Parse()
	} else {
		this.err_result <- err
	}
}

func (this *Service) Describe() (err error) {
	if this.described {
		return
	}
	job := upnpMakeDescribeServiceJob(this)
	go job.Describe()
	timeout := time.NewTimer(time.Duration(3) * time.Second)
	select {
	case <-job.result:
		this.described = true
	case err = <-job.err_result:
	case <-timeout.C:
	}
	return
}

func (this *Service) findAction(action string) (act *upnpAction, err error) {
	if !this.described {
		err = errors.New("Service is not described")
	} else {
		for _, act = range this.actionList {
			if action == act.name {
				return
			}
		}
		err = errors.New(fmt.Sprintf("No such method %s for service %s", action, this.serviceId))
	}
	return
}
