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
	_ "log"
	"net/http"
	_ "os"
	"strings"
)

const (
	soapEncodingStyle  = "http://schemas.xmlsoap.org/soap/encoding/"
	soapEnvelopeSchema = "http://schemas.xmlsoap.org/soap/envelope/"
	soapContentType    = "text/xml; charset=\"utf-8\""
	soapUserAgent      = "Linux UPnP/1.1 crash"
)

type Arg struct {
	Key   string
	Value interface{}
}

type Args []Arg

type ErrorResponse struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
	Detail      struct {
		XMLName   xml.Name
		UPnPError struct {
			XMLName   xml.Name
			ErrorCode string `xml:"errorCode"`
		} `xml:"urn:schemas-upnp-org:control-1-0 UPnPError"`
	} `xml:"detail"`
}

func (this *ErrorResponse) Error() (err error) {
	if 0 < len(this.FaultCode) {
		err = errors.New(this.Detail.UPnPError.ErrorCode)
	}
	return
}

type soapRequestAction struct {
	XMLName  xml.Name
	XMLNS_u  string `xml:"xmlns:u,attr"`
	Argument []soapRequestArgument
}

func soapBuildNamespace(svc *Service) (ns string) {
	return fmt.Sprintf("urn:%s:service:%s:%s", svc.serviceURI, svc.serviceType, svc.serviceVersion)
}

func soapNewRequestAction(action string, svc *Service, args Args) (ra soapRequestAction) {
	ns := soapBuildNamespace(svc)
	ra = soapRequestAction{}
	ra.XMLName.Local = "u:" + action
	ra.XMLNS_u = ns
	for _, arg := range args {
		ra.Argument = append(ra.Argument, soapNewRequestArgument(arg.Key, fmt.Sprintf("%v", arg.Value)))
	}
	return
}

type soapRequestArgument struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func soapNewRequestArgument(name, value string) (ar soapRequestArgument) {
	ar = soapRequestArgument{}
	ar.XMLName.Local = name
	ar.Value = value
	return
}

type soapRequestBody struct {
	XMLName xml.Name `xml:"s:Body"`
	Action  soapRequestAction
}

type soapResponseBody struct {
	Data string `xml:",innerxml"`
}

type soapRequestEnvelope struct {
	XMLName       xml.Name `xml:"s:Envelope"`
	Body          soapRequestBody
	XMLNS_s       string `xml:"xmlns:s,attr"`
	EncodingStyle string `xml:"s:encodingStyle,attr"`
}

type soapResponseEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    soapResponseBody
}

type soapRequest struct {
	Envelope soapRequestEnvelope
}

func soapNewRequest(action string, svc *Service, args Args) (req *soapRequest) {
	req = &soapRequest{}
	req.Envelope.XMLNS_s = soapEnvelopeSchema
	req.Envelope.EncodingStyle = soapEncodingStyle
	req.Envelope.Body.Action = soapNewRequestAction(action, svc, args)
	return
}

func upnpBuildRequest(svc *Service, action string, args Args) (msg []byte) {
	if act, err := svc.findAction(action); nil != err {
		panic(err)
	} else {
		req := soapNewRequest(act.name, svc, args)
		var err error
		if msg, err = xml.MarshalIndent(req.Envelope, "", "  "); nil != err {
			panic(err)
		}
	}
	return
}

func (this *Service) Call(action string, args Args) (response string) {
	client := &http.Client{}
	r := upnpBuildRequest(this, action, args)
	body := strings.NewReader(xml.Header + string(r))
	req, err := http.NewRequest("POST", this.controlURL.String(), body)
	req.Header.Set("CONTENT-TYPE", soapContentType)
	req.Header.Set("USER-AGENT", soapUserAgent)
	req.Header.Set("SOAPACTION", fmt.Sprintf("\"%s#%s\"", soapBuildNamespace(this), action))
	req.Header.Set("CONNECTION", "KEEP-ALIVE")
	//req.Write(os.Stdout)
	//body.Seek(0, 0)
	if nil != err {
		panic(err)
	}
	if resp, err := client.Do(req); nil != err {
		panic(err)
	} else {
		defer resp.Body.Close()
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); nil != err {
			panic(err)
		}
		doc := soapResponseEnvelope{}
		//log.Printf("%v", string(body))
		xml.Unmarshal(body, &doc)
		response = doc.Body.Data
	}
	return
}

func (this *Service) CallVa(action string, va_list ...interface{}) (response string) {
	var args Args
	for i := 0; i < len(va_list); i += 2 {
		args = append(args, Arg{va_list[i].(string), va_list[i+1]})
	}
	return this.Call(action, args)
}
