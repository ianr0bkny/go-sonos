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

package didl

import (
	"encoding/xml"
)

type Title struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Class struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Container struct {
	XMLName    xml.Name
	ID         string  `xml:"id,attr"`
	ParentID   string  `xml:"parentID,attr"`
	Restricted string  `xml:"restricted,attr"`
	Title      []Title `xml:"title"`
	Class      []Class `xml:"class"`
}

type Res struct {
	XMLName      xml.Name
	ProtocolInfo string `xml:"protocolInfo,attr"`
	Value        string `xml:",chardata"`
}

type AlbumArtURI struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Creator struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Album struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type OriginalTrackNumber struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type Item struct {
	XMLName             xml.Name
	ID                  string                `xml:"id,attr"`
	ParentID            string                `xml:"parentID,attr"`
	Restricted          string                `xml:"restricted,attr"`
	Res                 []Res                 `xml:"res"`
	AlbumArtURI         []AlbumArtURI         `xml:"albumArtURI"`
	Title               []Title               `xml:"title"`
	Class               []Class               `xml:"class"`
	Creator             []Creator             `xml:"creator"`
	Album               []Album               `xml:"album"`
	OriginalTrackNumber []OriginalTrackNumber `xml:"originalTrackNumber"`
}

type Lite struct {
	XMLName   xml.Name
	Container []Container `xml:"container"`
	Item      []Item      `xml:"item"`
}
