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
	"net"
	"net/http"
	"time"
)

type upnpEventProperty_XML struct {
	Content string `xml:",innerxml"`
}

type upnpEvent_XML struct {
	XMLName    xml.Name                `xml:"urn:schemas-upnp-org:event-1-0 propertyset"`
	Properties []upnpEventProperty_XML `xml:"urn:schemas-upnp-org:event-1-0 property"`
}

type EventCallback func(svc *Service, value string)

type Reactor interface {
	Init(ifiname, port string)
	Subscribe(svc *Service, callback EventCallback) error
}

type upnpEvent struct {
	sid   string
	value string
}

type upnpEventRecord struct {
	svc      *Service
	callback EventCallback
}

type upnpEventMap map[string]*upnpEventRecord

type upnpDefaultReactor struct {
	ifiname     string
	port        string
	initialized bool
	server      *http.Server
	localAddr   string
	eventMap    upnpEventMap
	subscrChan  chan *upnpEventRecord
	eventChan   chan *upnpEvent
}

func (this *upnpDefaultReactor) Init(ifiname, port string) {
	if this.initialized {
		panic("Attempt to reinitialize reactor")
	}

	ifi, err := net.InterfaceByName(ifiname)
	if err != nil {
		panic(err)
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		panic(err)
	}

	this.initialized = true
	this.port = port
	this.ifiname = ifiname
	this.localAddr = net.JoinHostPort(addrs[0].(*net.IPNet).IP.String(), port)
	this.server = &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/eventSub", this)
	log.Printf("Listening for events on %s", this.localAddr)
	go this.run()
	go this.server.ListenAndServe()
}

func (this *upnpDefaultReactor) handleAck(svc *Service, resp *http.Response) (sid string, err error) {
	sid_key := http.CanonicalHeaderKey("sid")
	if sid_list, has := resp.Header[sid_key]; has {
		sid = sid_list[0]
	} else {
		err = errors.New("Subscription ack missing sid")
	}
	return
}

func (this *upnpDefaultReactor) Subscribe(svc *Service, callback EventCallback) (err error) {
	rec := upnpEventRecord{
		svc:      svc,
		callback: callback,
	}
	this.subscrChan <- &rec
	return
}

func (this *upnpDefaultReactor) subscribeImpl(rec *upnpEventRecord) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest("SUBSCRIBE", rec.svc.eventSubURL.String(), nil)
	if nil != err {
		return
	}
	req.Header.Add("CALLBACK", fmt.Sprintf("<http://%s/eventSub>", this.localAddr))
	req.Header.Add("HOST", rec.svc.eventSubURL.Host)
	req.Header.Add("USER-AGENT", "unix/5.1 UPnP/1.1 sonos.go/1.0")
	req.Header.Add("NT", "upnp:event")
	req.Header.Add("TIMEOUT", "900")
	var resp *http.Response
	if resp, err = client.Do(req); nil == err {
		var sid string
		if sid, err = this.handleAck(rec.svc, resp); nil == err {
			this.eventMap[sid] = rec
		}
	}
	return
}

func (this *upnpDefaultReactor) maybePostEvent(event *upnpEvent) {
	if rec, has := this.eventMap[event.sid]; has {
		rec.callback(rec.svc, event.value)
	}
}

func (this *upnpDefaultReactor) run() {
	for {
		select {
		case subscr := <-this.subscrChan:
			this.subscribeImpl(subscr)
		case event := <-this.eventChan:
			this.maybePostEvent(event)
		}
	}
}

func (this *upnpDefaultReactor) sendAck(writer http.ResponseWriter) {
	writer.Write(nil)
}

func (this *upnpDefaultReactor) notify(sid, value string) {
	event := &upnpEvent{
		sid:   sid,
		value: value,
	}
	this.eventChan <- event
}

func (this *upnpDefaultReactor) unpack(sid string, doc *upnpEvent_XML) {
	for _, prop := range doc.Properties {
		this.notify(sid, prop.Content)
	}
}

func (this *upnpDefaultReactor) handle(request *http.Request) {
	defer request.Body.Close()
	if body, err := ioutil.ReadAll(request.Body); nil != err {
		panic(err)
	} else {
		sid_key := http.CanonicalHeaderKey("sid")
		var sid string
		if sid_list, has := request.Header[sid_key]; has {
			sid = sid_list[0]
			doc := &upnpEvent_XML{}
			xml.Unmarshal(body, doc)
			this.unpack(sid, doc)
		}
	}
}

func (this *upnpDefaultReactor) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer this.sendAck(writer)
	this.handle(request)
}

func MakeReactor() Reactor {
	reactor := &upnpDefaultReactor{}
	reactor.eventMap = make(upnpEventMap)
	reactor.subscrChan = make(chan *upnpEventRecord)
	reactor.eventChan = make(chan *upnpEvent)
	return reactor
}
