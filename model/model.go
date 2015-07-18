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

//
// A collection of object classes used in message passing in go-sonos.
//
package model

import (
	"github.com/ianr0bkny/go-sonos/didl"
	_ "log"
)

// An abstraction of a DIDL-Lite <Container> or <Item> block.
type Object interface {
	// The ObjectID of this item or container
	ID() string

	// The ObjectID of the parent container of this item or container
	ParentID() string

	// When true, the ability to change or delete this item or container is restricted
	Restricted() bool

	// The URI of the resource described by this item or container.
	// For a music track this could be the URI of the disk file on
	// the storage share.  For a playlist the URI may refer to the
	// queue's ObjectID.
	Res() string

	// The display name of the container or item.
	Title() string

	// A string giving the type of resource described by this Object, e.g.:
	//
	//  Containers:
	//	* object.container
	//	* object.container.albumlist
	//	* object.container.album.musicAlbum
	//	* object.container.genre.musicGenre
	//	* object.container.person.musicArtist
	//	* object.container.playlistContainer
	//	* object.container.playlistContainer.sameArtist
	//
	//  Items:
	//	* object.item.audioItem.musicTrack
	Class() string

	// The URI to use to access the artwork for this container or item.
	AlbumArtURI() string

	// The display name of the Artist or Album Artist.
	Creator() string

	// The display name of the containing album. This field is Valid
	// for Items only, not Containers.
	Album() string

	// The track number of this item in original album sort
	// order. This field is Valid for Items only, not Containers.
	OriginalTrackNumber() string

	// True, if this Object represents a container; false otherwise.
	IsContainer() bool
}

//
// A flattened structure of exported fields to allow Objects to be passed
// via XML, JSON, or other encoding relying on reflection.  Fields in this
// struct mirror the usage of like-named methods in the Object interface.
//
type ObjectMessage struct {
	ID          string
	ParentID    string
	URI         string
	Title       string
	Class       string
	AlbumArtURI string
	Creator     string
	Album       string
}

func makeObjectMessage(obj Object) *ObjectMessage {
	return &ObjectMessage{
		ID:          obj.ID(),
		ParentID:    obj.ParentID(),
		URI:         obj.Res(),
		Title:       obj.Title(),
		Class:       obj.Class(),
		AlbumArtURI: obj.AlbumArtURI(),
		Creator:     obj.Creator(),
		Album:       obj.Album(),
	}
}

type modelObjectImpl struct {
	id                  string
	parentId            string
	restricted          bool
	res                 string
	title               string
	class               string
	albumArtURI         string
	creator             string
	album               string
	originalTrackNumber string
	isContainer         bool
}

func (this modelObjectImpl) ID() string {
	return this.id
}

func (this modelObjectImpl) ParentID() string {
	return this.parentId
}

func (this modelObjectImpl) Restricted() bool {
	return this.restricted
}

func (this modelObjectImpl) Res() string {
	return this.res
}

func (this modelObjectImpl) Title() string {
	return this.title
}

func (this modelObjectImpl) Class() string {
	return this.class
}

func (this modelObjectImpl) AlbumArtURI() string {
	return this.albumArtURI
}

func (this modelObjectImpl) Creator() string {
	return this.creator
}

func (this modelObjectImpl) Album() string {
	return this.album
}

func (this modelObjectImpl) OriginalTrackNumber() string {
	return this.originalTrackNumber
}

func (this modelObjectImpl) IsContainer() bool {
	return this.isContainer
}

func makeContainer(in *didl.Container) Object {
	obj := modelObjectImpl{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	if 0 < len(in.Res) {
		obj.res = in.Res[0].Value
	}
	if 0 < len(in.Title) {
		obj.title = in.Title[0].Value
	}
	if 0 < len(in.Class) {
		obj.class = in.Class[0].Value
	}
	if 0 < len(in.AlbumArtURI) {
		obj.albumArtURI = in.AlbumArtURI[0].Value
	}
	if 0 < len(in.Creator) {
		obj.creator = in.Creator[0].Value
	}
	obj.isContainer = true
	return obj
}

func makeItem(in *didl.Item) Object {
	obj := modelObjectImpl{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	if 0 < len(in.Res) {
		obj.res = in.Res[0].Value
	}
	if 0 < len(in.Title) {
		obj.title = in.Title[0].Value
	}
	if 0 < len(in.Class) {
		obj.class = in.Class[0].Value
	}
	if 0 < len(in.AlbumArtURI) {
		obj.albumArtURI = in.AlbumArtURI[0].Value
	}
	if 0 < len(in.Creator) {
		obj.creator = in.Creator[0].Value
	}
	if 0 < len(in.Album) {
		obj.album = in.Album[0].Value
	}
	if 0 < len(in.OriginalTrackNumber) {
		obj.originalTrackNumber = in.OriginalTrackNumber[0].Value
	}
	return obj
}

// Create a list of Objects from a didl.Lite document.
func ObjectStream(in *didl.Lite) (objects []Object) {
	for _, container := range in.Container {
		objects = append(objects, makeContainer(&container))
	}
	for _, item := range in.Item {
		objects = append(objects, makeItem(&item))
	}
	return
}

// Create a list of ObjectMessages from a list of Objects.
func ObjectMessageStream(objs []Object) []*ObjectMessage {
	var out []*ObjectMessage
	for _, obj := range objs {
		out = append(out, makeObjectMessage(obj))
	}
	return out
}
