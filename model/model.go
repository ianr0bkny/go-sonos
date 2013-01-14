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

package model

import (
	"fmt"
	"github.com/ianr0bkny/go-sonos/didl"
	_ "log"
)

type Object interface {
	ID() string
	ParentID() string
	Restricted() bool
	Title() string
	Class() string
}

type modelObjectImpl struct {
	id         string
	parentId   string
	restricted bool
	title      string
	class      string
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

func (this modelObjectImpl) Title() string {
	return this.title
}

func (this modelObjectImpl) Class() string {
	return this.class
}

type Container struct {
	modelObjectImpl
}

func makeContainer(in *didl.Container) Object {
	obj := Container{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

type Person struct {
	Container
}

type MusicArtist struct {
	Person
}

func makeMusicArtist(in *didl.Container) Object {
	obj := MusicArtist{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

type Album struct {
	Container
}

type AlbumList struct {
	Container
}

func makeAlbumList(in *didl.Container) Object {
	obj := AlbumList{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

type PlaylistContainer struct {
	Container
}

func makePlaylistContainer(in *didl.Container) Object {
	obj := PlaylistContainer{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

type MusicAlbum struct {
	modelObjectImpl
}

func makeMusicAlbum(in *didl.Container) Object {
	obj := MusicAlbum{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

type SameArtist struct {
	modelObjectImpl
}

func makeSameArtist(in *didl.Container) Object {
	obj := MusicAlbum{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

type Genre struct {
	Container
}

type MusicGenre struct {
	Genre
}

func makeMusicGenre(in *didl.Container) Object {
	obj := MusicGenre{}
	obj.id = in.ID
	obj.parentId = in.ParentID
	obj.restricted = in.Restricted
	obj.title = in.Title[0].Value
	obj.class = in.Class[0].Value
	return obj
}

func makeContainerObject(in *didl.Container) Object {
	if 0 == len(in.Class) {
		panic("missing object class")
	}
	switch in.Class[0].Value {
	case "object.container.album.musicAlbum":
		return makeMusicAlbum(in)
	case "object.container.playlistContainer.sameArtist":
		return makeSameArtist(in)
	case "object.container":
		return makeContainer(in)
	case "object.container.playlistContainer":
		return makePlaylistContainer(in)
	case "object.container.albumlist":
		return makeAlbumList(in)
	case "object.container.genre.musicGenre":
		return makeMusicGenre(in)
	case "object.container.person.musicArtist":
		return makeMusicArtist(in)
	default:
		panic(fmt.Sprintf("unsupported DIDL-Lite object class %s", in.Class[0].Value))
	}
	return nil
}

func makeItem(in *didl.Item) Object {
	return nil
}

func ObjectStream(in *didl.Lite) (objects []Object) {
	for _, container := range in.Container {
		objects = append(objects, makeContainerObject(&container))
	}
	for _, item := range in.Item {
		objects = append(objects, makeItem(&item))
	}
	return
}
