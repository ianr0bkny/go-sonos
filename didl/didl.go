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
