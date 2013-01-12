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
	"github.com/ianr0bkny/go-sonos/didl"
	_ "log"
)

type ContentDirectory struct {
	Svc *Service
}

func (this *ContentDirectory) GetSearchCapabilities() (searchCaps string, err error) {
	type Response struct {
		XMLName    xml.Name
		SearchCaps string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetSearchCapabilities")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	searchCaps = doc.SearchCaps
	err = doc.Error()
	return
}

func (this *ContentDirectory) GetSortCapabilities() (sortCaps string, err error) {
	type Response struct {
		XMLName  xml.Name
		SortCaps string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetSortCapabilities")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	sortCaps = doc.SortCaps
	err = doc.Error()
	return
}

func (this *ContentDirectory) GetSystemUpdateID() (id string, err error) {
	type Response struct {
		XMLName xml.Name
		Id      string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetSystemUpdateID")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	id = doc.Id
	err = doc.Error()
	return
}

func (this *ContentDirectory) GetAlbumArtistDisplayOption() (albumArtistDisplayOption string, err error) {
	type Response struct {
		XMLName                  xml.Name
		AlbumArtistDisplayOption string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetAlbumArtistDisplayOption")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	albumArtistDisplayOption = doc.AlbumArtistDisplayOption
	err = doc.Error()
	return
}

func (this *ContentDirectory) GetLastIndexChange() (lastIndexChange string, err error) {
	type Response struct {
		XMLName         xml.Name
		LastIndexChange string
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetLastIndexChange")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	lastIndexChange = doc.LastIndexChange
	err = doc.Error()
	return
}

const (
	BrowseObjectID_Root = "0"
)

const (
	BrowseFlag_BrowseMetadata       = "BrowseMetadata"
	BrowseFlag_BrowseDirectChildren = "BrowseDirectChildren"
)

const (
	BrowseFilter_All = "*"
)

const (
	BrowseSortCriteria_None = ""
)

type BrowseRequest struct {
	ObjectID      string
	BrowseFlag    string
	Filter        string
	StartingIndex uint32
	RequestCount  uint32
	SortCriteria  string
}

type BrowseResult struct {
	NumberReturned int32
	TotalMatches   int32
	UpdateID       int32
	Doc            *didl.Lite
}

func (this *ContentDirectory) Browse(req *BrowseRequest) (browseResult *BrowseResult, err error) {
	type Response struct {
		XMLName xml.Name
		Result  string
		BrowseResult
		ErrorResponse
	}
	args := []Arg{
		{"ObjectID", req.ObjectID},
		{"BrowseFlag", req.BrowseFlag},
		{"Filter", req.Filter},
		{"StartingIndex", req.StartingIndex},
		{"RequestedCount", req.RequestCount},
		{"SortCriteria", req.SortCriteria},
	}
	response := Call(this.Svc, "Browse", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	doc.Doc = &didl.Lite{}
	xml.Unmarshal([]byte(doc.Result), doc.Doc)
	browseResult = &doc.BrowseResult
	err = doc.Error()
	return
}

func (this *ContentDirectory) FindPrefix(objectId, prefix string) (startingIndex, updateId uint32, err error) {
	type Response struct {
		XMLName       xml.Name
		StartingIndex uint32
		UpdateID      uint32
		ErrorResponse
	}
	args := []Arg{
		{"ObjectID", objectId},
		{"StartingIndex", startingIndex},
		{"UpdateID", updateId},
	}
	response := Call(this.Svc, "FindPrefix", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	startingIndex = doc.StartingIndex
	updateId = doc.UpdateID
	err = doc.Error()
	return
}

type PrefixLocations struct {
	TotalPrefixes     uint32
	PrefixAndIndexCSV string
	UpdateID          uint32
}

func (this *ContentDirectory) GetAllPrefixLocations(objectId string) (prefixLocations *PrefixLocations, err error) {
	type Response struct {
		XMLName xml.Name
		PrefixLocations
		ErrorResponse
	}
	args := []Arg{
		{"ObjectID", objectId},
	}
	response := Call(this.Svc, "GetAllPrefixLocations", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	prefixLocations = &doc.PrefixLocations
	err = doc.Error()
	return
}

func (this *ContentDirectory) CreateObject(container, elements string) (objectId, result string, err error) {
	type Response struct {
		XMLName  xml.Name
		ObjectID string
		Result   string
		ErrorResponse
	}
	args := []Arg{
		{"Container", container},
		{"Elements", elements},
	}
	response := Call(this.Svc, "CreateObject", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	objectId = doc.ObjectID
	result = doc.Result
	err = doc.Error()
	return
}

func (this *ContentDirectory) UpdateObject(objectId, currentTagValue, newTagValue string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"ObjectID", objectId},
		{"CurrentTagValue", currentTagValue},
		{"NewTagValue", newTagValue},
	}
	response := Call(this.Svc, "UpdateObject", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ContentDirectory) DestroyObject(objectId string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"ObjectID", objectId},
	}
	response := Call(this.Svc, "DestroyObject", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ContentDirectory) RefreshShareList() (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	response := CallVa(this.Svc, "RefreshShareList")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ContentDirectory) RefreshShareIndex(albumArtistDisplayOption string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AlbumArtistDisplayOption", albumArtistDisplayOption},
	}
	response := Call(this.Svc, "RefreshShareIndex", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ContentDirectory) RequestResort(sortOrder string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"SortOrder", sortOrder},
	}
	response := Call(this.Svc, "RequestResort", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *ContentDirectory) GetShareIndexInProgress() (isIndexing bool, err error) {
	type Response struct {
		XMLName    xml.Name
		IsIndexing bool
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetShareIndexInProgress")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	isIndexing = doc.IsIndexing
	err = doc.Error()
	return
}

func (this *ContentDirectory) GetBrowseable() (isBrowseable bool, err error) {
	type Response struct {
		XMLName      xml.Name
		IsBrowseable bool
		ErrorResponse
	}
	response := CallVa(this.Svc, "GetBrowseable")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	isBrowseable = doc.IsBrowseable
	err = doc.Error()
	return
}

func (this *ContentDirectory) SetBrowseable(browseable bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"Browseable", browseable},
	}
	response := Call(this.Svc, "SetBrowseable", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
