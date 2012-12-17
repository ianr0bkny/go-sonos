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
	"github.com/ianr0bkny/go-sonos/didl"
	"github.com/ianr0bkny/go-sonos/upnp"
	_ "log"
)

type ContentDirectory struct {
	svc *upnp.Service
}

func (this *ContentDirectory) GetSearchCapabilities() string {
	type Response struct {
		XMLName    xml.Name
		SearchCaps string
	}
	response := upnp.CallVa(this.svc, "GetSearchCapabilities")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.SearchCaps
}

func (this *ContentDirectory) GetSortCapabilities() string {
	type Response struct {
		XMLName  xml.Name
		SortCaps string
	}
	response := upnp.CallVa(this.svc, "GetSortCapabilities")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.SortCaps
}

func (this *ContentDirectory) GetLastIndexChange() string {
	type Response struct {
		XMLName         xml.Name
		LastIndexChange string
	}
	response := upnp.CallVa(this.svc, "GetLastIndexChange")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.LastIndexChange
}

func (this *ContentDirectory) Browse(object, flag, filter string, start, count int, sort string) (result *didl.Lite) {
	type Response struct {
		XMLName        xml.Name
		Result         string
		NumberReturned int
		TotalMatches   int
		UpdateID       int
	}
	args := []upnp.Arg{
		{"ObjectID", object},
		{"BrowseFlag", flag},
		{"Filter", filter},
		{"StartingIndex", start},
		{"RequestedCount", count},
		{"SortCriteria", sort},
	}
	response := upnp.Call(this.svc, "Browse", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	//log.Printf("%s", doc.Result)
	result = &didl.Lite{}
	xml.Unmarshal([]byte(doc.Result), &result)
	//log.Printf("%#v", result)
	return
}
