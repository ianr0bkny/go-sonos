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
	"github.com/ianr0bkny/go-sonos/didl"
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
	"strings"
)

const (
	ObjectID_Attributes    = "A:"
	ObjectID_MusicShares   = "S:"
	ObjectID_Queues        = "Q:"
	ObjectID_SavedQueues   = "SQ:"
	ObjectID_InternetRadio = "R:"
	ObjectID_EntireNetwork = "EN:"
	//
	ObjectID_Queue_AVT_Instance_0 = "Q:0"
	//
	ObjectID_Attribute_Genres = "A:GENRE"
)

type Container struct {
	ID         string
	ParentID   string
	Restricted bool
	Title      string
	Class      string
}

func makeContainer(in *didl.Container) *Container {
	return &Container{
		in.ID,
		in.ParentID,
		in.Restricted,
		in.Title[0].Value,
		in.Class[0].Value,
	}
}

func makeContainers(in []didl.Container) (out []*Container) {
	for _, container := range in {
		out = append(out, makeContainer(&container))
	}
	return
}

func (this *Sonos) GetRootLevelChildren() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		upnp.BrowseObjectID_Root,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListQueues() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_Queues,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListSavedQueues() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_SavedQueues,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListInternetRadio() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_InternetRadio,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListAttributes() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_Attributes,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListMusicShares() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_MusicShares,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListGenres() (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_Attribute_Genres,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func objectIDForGenre(genre string) string {
	return strings.Join([]string{ObjectID_Attribute_Genres, genre}, "/")
}

func (this *Sonos) ListGenre(genre string) (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		objectIDForGenre(genre),
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) ListChildren(objectId string) (containers []*Container, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		objectId,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		containers = makeContainers(result.Doc.Container)
	}
	return
}

func (this *Sonos) GetQueueContents() (err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID_Queue_AVT_Instance_0,
		upnp.BrowseFlag_BrowseDirectChildren,
		upnp.BrowseFilter_All,
		0, /*StartingIndex*/
		0, /*RequestCount*/
		upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		/*
			for _, container := range result.Doc.Container {
				c := Container{}
				c.ID = container.ID
				c.Title = container.Title[0].Value
				c.Class = container.Class[0].Value
				containers = append(containers, c)
			}
		*/
		log.Printf("%#v", result.Doc)
	}
	return
}
