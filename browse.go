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
	"github.com/ianr0bkny/go-sonos/model"
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
	ObjectID_Attribute_Genres    = "A:GENRE"
	ObjectID_Attribute_Album     = "A:ALBUM"
	ObjectID_Attribute_Artist    = "A:ARTIST"
	ObjectID_Attribute_Composers = "A:COMPOSER"
)

func (this *Sonos) GetRootLevelChildren() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      upnp.BrowseObjectID_Root,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) ListQueues() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_Queues,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) ListSavedQueues() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_SavedQueues,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) ListInternetRadio() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_InternetRadio,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) ListAttributes() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_Attributes,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) ListMusicShares() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_MusicShares,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) GetAllGenres() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_Attribute_Genres,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) GetAllComposers() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_Attribute_Composers,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func objectIDForGenre(genre string) string {
	return strings.Join([]string{ObjectID_Attribute_Genres, genre}, "/")
}

func objectIDForAlbum(album string) string {
	return strings.Join([]string{ObjectID_Attribute_Album, album}, "/")
}

func objectIDForArtist(artist string) string {
	return strings.Join([]string{ObjectID_Attribute_Artist, artist}, "/")
}

func (this *Sonos) GetGenreArtists(genre string) ([]model.Object, error) {
	req := &upnp.BrowseRequest{
		ObjectID:      objectIDForGenre(genre),
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err := this.Browse(req); nil != err {
		log.Printf("Could not browse artists for genre `%s': %v", genre, err)
		return nil, err
	} else {
		return model.ObjectStream(result.Doc), nil
	}
}

func (this *Sonos) ListChildren(objectId string) (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      objectId,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) GetMetadata(objectId string) (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      objectId,
		BrowseFlag:    upnp.BrowseFlag_BrowseMetadata,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) GetDirectChildren(objectId string) (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      objectId,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) GetQueueContents() (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      ObjectID_Queue_AVT_Instance_0,
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}

func (this *Sonos) GetAlbumTracks(album string) ([]model.Object, error) {
	req := &upnp.BrowseRequest{
		ObjectID:      objectIDForAlbum(album),
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	log.Printf("Browsing tracks for album `%s'", album)
	if result, err := this.Browse(req); nil != err {
		log.Printf("Could not browse tracks for album `%s': %v", album, err)
		return nil, err
	} else {
		return model.ObjectStream(result.Doc), nil
	}
}

func (this *Sonos) GetTrackFromAlbum(album, track string) ([]model.Object, error) {
	if tracks, err := this.GetAlbumTracks(album); nil != err {
		return nil, err
	} else {
		var track_objs []model.Object
		for _, track_obj := range tracks {
			if track_obj.Title() == track {
				track_objs = append(track_objs, track_obj)
			}
		}
		return track_objs, nil
	}
}

func (this *Sonos) GetArtistAlbums(artist string) (objects []model.Object, err error) {
	var result *upnp.BrowseResult
	req := &upnp.BrowseRequest{
		ObjectID:      objectIDForArtist(artist),
		BrowseFlag:    upnp.BrowseFlag_BrowseDirectChildren,
		Filter:        upnp.BrowseFilter_All,
		StartingIndex: 0,
		RequestCount:  0,
		SortCriteria:  upnp.BrowseSortCriteria_None,
	}
	if result, err = this.Browse(req); nil != err {
		return
	} else {
		objects = model.ObjectStream(result.Doc)
	}
	return
}
