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
	"github.com/ianr0bkny/go-sonos/ssdp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"
)

var upnpOtherDeviceRegex *regexp.Regexp
var upnpOtherServiceRegex *regexp.Regexp

func init() {
	upnpOtherDeviceRegex = regexp.MustCompile("^urn:([^:]+):device:([^:]+)(:(.+))?$")
	upnpOtherServiceRegex = regexp.MustCompile("^urn:([^:]+):service:([^:]+)(:(.+))?$")
}

type upnpSpecVersion_XML struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

type upnpDevice_XML struct {
	DeviceType           string              `xml:"deviceType"`
	FriendlyName         string              `xml:"friendlyName"`
	Manufacturer         string              `xml:"manufacturer"`
	ManufacturerURL      string              `xml:"manufacturerURL"`
	ModelNumber          string              `xml:"modelNumber"`
	ModelDescription     string              `xml:"modelDescription"`
	ModelName            string              `xml:"modelName"`
	ModelURL             string              `xml:"modelURL"`
	SoftwareVersion      string              `xml:"softwareVersion"`
	HardwareVersion      string              `xml:"hardwareVersion"`
	SerialNum            string              `xml:"serialNum"`
	UDN                  string              `xml:"UDN"`
	IconList             upnpIconList_XML    `xml:"iconList"`
	MinCompatibleVersion string              `xml:"minCompatibleVersion"`
	DisplayVersion       string              `xml:"displayVersion"`
	ExtraVersion         string              `xml:"extraVersion"`
	RoomName             string              `xml:"roomName"`
	DisplayName          string              `xml:"displayName"`
	ZoneType             string              `xml:"zoneType"`
	Feature1             string              `xml:"feature1"`
	Feature2             string              `xml:"feature2"`
	Feature3             string              `xml:"feature3"`
	InternalSpeakerSize  string              `xml:"internalSpeakerSize"`
	ServiceList          upnpServiceList_XML `xml:"serviceList"`
	DeviceList           upnpDeviceList_XML  `xml:"deviceList"`
}

type upnpIcon_XML struct {
	Id       string `xml:"id"`
	Mimetype string `xml:"mimetype"`
	Width    string `xml:"width"`
	Height   string `xml:"height"`
	Depth    string `xml:"depth"`
	Url      string `xml:"url"`
}

type upnpIconList_XML struct {
	Icon []upnpIcon_XML `xml:"icon"`
}

type upnpService_XML struct {
	ServiceType string        `xml:"serviceType"`
	ServiceId   string        `xml:"serviceId"`
	ControlURL  ssdp.Location `xml:"controlURL"`
	EventSubURL ssdp.Location `xml:"eventSubURL"`
	SCPDURL     ssdp.Location `xml:"SCPDURL"`
}

type upnpServiceList_XML struct {
	Service []upnpService_XML `xml:"service"`
}

type upnpDeviceList_XML struct {
	Device []upnpDevice_XML `xml:"device"`
}

type upnpDescribeDevice_XML struct {
	XMLNamespace string                `xml:"xmlns,attr"`
	SpecVersion  []upnpSpecVersion_XML `xml:"specVersion"`
	Device       []upnpDevice_XML      `xml:"device"`
}

type upnpDescribeDeviceJob struct {
	result     chan []*Service
	err_result chan error
	response   *http.Response
	doc        upnpDescribeDevice_XML
	uri        ssdp.Location
}

func upnpMakeDescribeDeviceJob(uri ssdp.Location) (job *upnpDescribeDeviceJob) {
	job = &upnpDescribeDeviceJob{}
	job.result = make(chan []*Service)
	job.err_result = make(chan error)
	job.uri = uri
	job.doc = upnpDescribeDevice_XML{}
	return
}

func (this *upnpDescribeDeviceJob) BuildURL(servicePath ssdp.Location) (url *url.URL, err error) {
	if url, err = url.Parse(string(this.uri)); nil != err {
		return
	}
	if len(servicePath) > 0 && servicePath[0] == '/' {
		// We have an absolute path
		url.Path = string(servicePath)
	} else {
		// We have a path relative to the location of the description
		basePath, _ := path.Split(url.Path)
		url.Path = path.Join(basePath, string(servicePath))
	}
	return
}

func (this *upnpDescribeDeviceJob) UnpackService(dev *upnpDevice_XML, svc_doc *upnpService_XML) (svc *Service) {
	svc = upnpMakeService()
	if m := upnpOtherDeviceRegex.FindStringSubmatch(dev.DeviceType); 0 < len(m) {
		svc.deviceURI = m[1]
		svc.deviceType = m[2]
		svc.deviceVersion = m[4]
	} else {
		this.err_result <- errors.New(fmt.Sprintf("Malformed device type string `%s'", dev.DeviceType))
		return
	}
	svc.udn = dev.UDN
	if m := upnpOtherServiceRegex.FindStringSubmatch(svc_doc.ServiceType); 0 < len(m) {
		svc.serviceURI = m[1]
		svc.serviceType = m[2]
		svc.serviceVersion = m[4]
	} else {
		this.err_result <- errors.New(fmt.Sprintf("Malformed service type string `%s'", svc_doc.ServiceType))
		return
	}
	svc.serviceId = svc_doc.ServiceId
	var err error
	if svc.controlURL, err = this.BuildURL(svc_doc.ControlURL); nil != err {
		this.err_result <- err
	} else if svc.eventSubURL, err = this.BuildURL(svc_doc.EventSubURL); nil != err {
		this.err_result <- err
	} else if svc.scpdURL, err = this.BuildURL(svc_doc.SCPDURL); nil != err {
		this.err_result <- err
	}
	return
}

func (this *upnpDescribeDeviceJob) UnpackDevice(dev *upnpDevice_XML) (svc_list []*Service) {
	for _, svc := range dev.ServiceList.Service {
		svc_list = append(svc_list, this.UnpackService(dev, &svc))
	}
	for _, sub_dev := range dev.DeviceList.Device {
		svc_list = append(svc_list, this.UnpackDevice(&sub_dev)...)
	}
	return
}

func (this *upnpDescribeDeviceJob) Unpack() (svc_list []*Service) {
	for _, dev := range this.doc.Device {
		svc_list = append(svc_list, this.UnpackDevice(&dev)...)
	}
	return
}

func (this *upnpDescribeDeviceJob) Parse() {
	defer this.response.Body.Close()
	if body, err := ioutil.ReadAll(this.response.Body); nil == err {
		xml.Unmarshal(body, &this.doc)
		this.result <- this.Unpack()
	} else {
		this.err_result <- err
	}
}

func (this *upnpDescribeDeviceJob) Describe() {
	var err error
	log.Printf("Loading %s", string(this.uri))
	if this.response, err = http.Get(string(this.uri)); nil == err {
		this.Parse()
	} else {
		this.err_result <- err
	}
}

type ServiceMap map[string][]*Service

func Describe(uri ssdp.Location) (svc_map ServiceMap, err error) {
	job := upnpMakeDescribeDeviceJob(uri)
	go job.Describe()
	timeout := time.NewTimer(time.Duration(3) * time.Second)
	select {
	case svc_list := <-job.result:
		svc_map = make(ServiceMap)
		for _, svc := range svc_list {
			svc_map[svc.serviceType] = append(svc_map[svc.serviceType], svc)
		}
	case err = <-job.err_result:
	case <-timeout.C:
	}
	return
}
