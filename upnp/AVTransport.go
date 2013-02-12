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
	_ "log"
	"strings"
)

type AVTransport struct {
	Svc *Service
}

func (this *AVTransport) SetAVTransportURI(instanceId uint32, currentURI, currentURIMetaData string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"CurrentURI", currentURI},
		{"CurrentURIMetaData", currentURIMetaData},
	}
	response := Call(this.Svc, "SetAVTransportURI", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

//
// Input parameters for AddURIToQueue.
//
type AddURIToQueueIn struct {
	// The URI of the track to be added to the queue, corresponding
	// the to <res> tag in a DIDL-Lite description (@see dldl,
	// @ContentDirectory, etc) e.g.:
	//     "x-file-cifs://servername/path/to/track.mp3"
	EnqueuedURI string
	// A DIDL-Lite document describing the the resource given by @EnqueuedURI
	EnqueuedURIMetaData string
	// This field should be 0 to insert the new item at the end
	// of the queue.  If non-zero the new track will be inserted at
	// this location, and later tracks will see their track numbers
	// incremented.
	DesiredFirstTrackNumberEnqueued uint32
	// ???? (possibly unsupported)
	EnqueueAsNext bool
}

//
// Output parameters for AddURIToQueue
//
type AddURIToQueueOut struct {
	// The track number of the newly added track.
	FirstTrackNumberEnqueued uint32
	// The number of tracks added by this request (always 1).
	NumTracksAdded uint32
	// The length of the queue now that this track has been added
	NewQueueLength uint32
}

//
// Add a single track to the queue (Q:0).  For Sonos @instanceId will
// always be 0.  See @AddURIToQueueIn for a discussion of the input
// parameters and @AddURIToQueueOut for a discussion of the output
// parameters.
//
func (this *AVTransport) AddURIToQueue(instanceId uint32, req *AddURIToQueueIn) (*AddURIToQueueOut, error) {
	type Response struct {
		XMLName xml.Name
		AddURIToQueueOut
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"EnqueuedURI", req.EnqueuedURI},
		{"EnqueuedURIMetaData", req.EnqueuedURIMetaData},
		{"DesiredFirstTrackNumberEnqueued", req.DesiredFirstTrackNumberEnqueued},
		{"EnqueueAsNext", req.EnqueueAsNext},
	}
	response := Call(this.Svc, "AddURIToQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.AddURIToQueueOut, doc.Error()
}

//
// Input parameters for AddMultipleURIsToQueue.
//
type AddMultipleURIsToQueueIn struct {
	// UpdateID (in), can be 0
	UpdateID uint32
	// The number of URIs to be added in this request
	NumberOfURIs uint32
	// A list of @NumberOfURIs URIs, separated by a space
	EnqueuedURIs string
	// A list of @NumberOfURIs DIDL-Lite documents, separated by a space
	EnqueuedURIsMetaData string
	// URI of a container in the ContentDirectory containing the
	// URIs to be added.  If adding tracks this should be the URI for
	// the A:TRACK entry in the directory.
	ContainerURI string
	// A DIDL-Lite document describing the resource given by @ContainerURI
	ContainerMetaData string
	// This field should be 0 to insert the new item at the end
	// of the queue.  If non-zero the new track will be inserted at
	// this location, and later tracks will see their track numbers
	// incremented.
	DesiredFirstTrackNumberEnqueued uint32
	// ???? (possibly unsupported)
	EnqueueAsNext bool
}

//
// Output parameters for AddMultipleURIsToQueue.
//
type AddMultipleURIsToQueueOut struct {
	// The starting position int the queue (Q:0) of the newly added tracks
	FirstTrackNumberEnqueued uint32
	// The number of tracks added by the request
	NumTracksAdded uint32
	// The length of the queue after the request was complete
	NewQueueLength uint32
	// The new UpdateID
	NewUpdateID uint32
}

//
// Add multiple tracks to the queue (Q:0).  This method does not seem
// to be a standard part of AVTransport:1, but rather a Sonos extension.
// As such it is not entirely clear how it should be used.
//
// For Sonos @instanceId should always be 0; @UpdateID should be 0;
// @NumberOfURIs should be the number of tracks to be added by the
// request; @EnqueuedURIs is a space-separated list of URIs (as given by
// the Res() method of the model.Object class); @EnqueuedURIMetData is a
// space-separated list of DIDL-Lite documents describing the resources
// to be added; @ContainerURI should be the ContentDirectory URI for
// A:TRACK, when adding tracks; @ContainerMetaData should be a DIDL-Lite
// document describing A:TRACK. Other arguments have the same meaning as
// in @AddURIToQueue.
//
// Note that the number of DIDL-Lite documents in @EnqueuedURIsMetaData
// must match the number of URIs in @EnqueuedURIs.  These DIDL-Lite documents
// can be empty, but must be present.  @ContainerMetaData must be a string
// of non-zero length, but need not be a valid DIDL-Lite document.
//
func (this *AVTransport) AddMultipleURIsToQueue(instanceId uint32, req *AddMultipleURIsToQueueIn) (*AddMultipleURIsToQueueOut, error) {
	type Response struct {
		XMLName xml.Name
		AddMultipleURIsToQueueOut
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"UpdateID", req.UpdateID},
		{"NumberOfURIs", req.NumberOfURIs},
		{"EnqueuedURIs", req.EnqueuedURIs},
		{"EnqueuedURIsMetaData", req.EnqueuedURIsMetaData},
		{"ContainerURI", req.ContainerURI},
		{"ContainerMetaData", req.ContainerMetaData},
		{"DesiredFirstTrackNumberEnqueued", req.DesiredFirstTrackNumberEnqueued},
		{"EnqueueAsNext", req.EnqueueAsNext},
	}
	response := Call(this.Svc, "AddMultipleURIsToQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.AddMultipleURIsToQueueOut, doc.Error()
}

//
// Move a contiguous range of tracks to a given point in the queue.
// For Sonos @instanceId will always be 0; @startingIndex is the first
// track in the range to be moved, where the first track in the queue is
// track 1; @numberOfTracks is the length of the range; @insertBefore set
// the destination position in the queue; @updateId should be 0.
//
// Note that to move tracks to the end of the queue @insertBefore must be
// set to the number of tracks in the queue plus 1.  This method fails with
// 402 if @startingndex, @numberOfTracks, or @insertBefore are out of range.
//
func (this *AVTransport) ReorderTracksInQueue(instanceId, startingIndex, numberOfTracks, insertBefore, updateId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"StartingIndex", startingIndex},
		{"NumberOfTracks", numberOfTracks},
		{"InsertBefore", insertBefore},
		{"UpdateID", updateId},
	}
	response := Call(this.Svc, "ReorderTracksInQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

//
// Remove a single track from the queue (Q:0).  For Sonos @instanceId
// will always be 0; @objectId will be the identifier of the item to be
// removed from the queue (e.g. "Q:0/5" for the fifth element in the queue);
// @updateId will always be 0.
//
func (this *AVTransport) RemoveTrackFromQueue(instanceId uint32, objectId string, updateId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"ObjectID", objectId},
		{"UpdateID", updateId},
	}
	response := Call(this.Svc, "RemoveTrackFromQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

//
// Remove a continguous range of tracks from the queue (Q:0).  For Sonos
// @instanceId will always be 0; @updateId should be 0; @startingIndex is
// the first track to remove where the first track is 1; @numberOfTracks
// is the number of tracks to remove.  Returns the new @updateId.
//
// This method fails with 402 if either @startingIndex or @numberOfTracks
// is out of range.
//
func (this *AVTransport) RemoveTrackRangeFromQueue(instanceId, updateId, startingIndex, numberOfTracks uint32) (uint32, error) {
	type Response struct {
		XMLName     xml.Name
		NewUpdateID uint32
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"UpdateID", updateId},
		{"StartingIndex", startingIndex},
		{"NumberOfTracks", numberOfTracks},
	}
	response := Call(this.Svc, "RemoveTrackRangeFromQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.NewUpdateID, doc.Error()
}

//
// Remove all tracks from the queue (Q:0).  For Sonos @instanceId will
// always be 0.  Emptying an already empty queue is not an error.
//
func (this *AVTransport) RemoveAllTracksFromQueue(instanceId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "RemoveAllTracksFromQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

//
// Create a new named queue (SQ:n) from the contents of the current
// queue (Q:0).  For Sonos @instanceId should always be 0; @title is the
// display name of the new named queue; @objectId should be left blank.
// This method returns the objectId of the newly created queue.
//
func (this *AVTransport) SaveQueue(instanceId uint32, title, objectId string) (string, error) {
	type Response struct {
		XMLName          xml.Name
		AssignedObjectID string
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Title", title},
		{"ObjectID", objectId},
	}
	response := Call(this.Svc, "SaveQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.AssignedObjectID, doc.Error()
}

func (this *AVTransport) BackupQueue(instanceId uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "BackupQueue", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

//
// The return type of the GetMediaInfo method
//
type MediaInfo struct {
	// The number of tracks in the queue (Q:0)
	NrTracks uint32
	// ???? (possibly not supported)
	MediaDuration string
	// The URI of the queue
	CurrentURI string
	// ????
	CurrentURIMetaData string
	// ???? (possibly not supported)
	NextURI string
	// ???? (possibly not supported)
	NextURIMetaData string
	// ????
	PlayMedium string
	// ???? (possibly not supported)
	RecordMedium string
	// ???? (possibly not supported)
	WriteStatus string
}

//
// Gets information about the currently selected media, its URI, length
// in tracks, and recording status, if any.  For Sonos @instanceId should
// always be 0 and most of the fields are unsupported.
//
func (this *AVTransport) GetMediaInfo(instanceId uint32) (*MediaInfo, error) {
	type Response struct {
		XMLName xml.Name
		MediaInfo
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetMediaInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.MediaInfo, doc.Error()
}

// Legal values for TransportInfo.CurrentTransportState
const (
	State_PLAYING         = "PLAYING"
	State_PAUSED_PLAYBACK = "PAUSED_PLAYBACK"
	State_STOPPED         = "STOPPED"
)

//
// The return type for the GetTransportInfo method
//
type TransportInfo struct {
	// Indicates whether the device is playing, paused, or stopped
	CurrentTransportState string
	// Indicates if an error condition exists ("OK" otherwise)
	CurrentTransportStatus string
	// Playback speed relative to normal playback speed (e.g. "1" or "1/2")
	CurrentSpeed string
}

//
// Return the current state of the transport (playing, stopped, paused),
// its status, and playback speed.  For Sonos @instanceId should always be 0.
//
func (this *AVTransport) GetTransportInfo(instanceId uint32) (*TransportInfo, error) {
	type Response struct {
		XMLName xml.Name
		TransportInfo
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetTransportInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.TransportInfo, doc.Error()
}

//
// The return type of the GetPositionInfo method
//
type PositionInfo struct {
	// Track number relative to the beginning of the queue (not the containing album).
	Track uint32
	// Total length of the track in HH:MM:SS format
	TrackDuration string
	// The DIDL-Lite document describing this item in the ContentDirectory
	TrackMetaData string
	// The URI of the track, corresponding // the to <res> tag in
	// a DIDL-Lite description (@see dldl, @ContentDirectory, etc) e.g.:
	//     "x-file-cifs://servername/path/to/track.mp3"
	TrackURI string
	// The number of elapsed seconds into the track in HH:MM:SS format
	RelTime string
	// ???? (possibly unsupported)
	AbsTime string
	// ???? (possibly unsupported)
	RelCount uint32
	// ???? (possibly unsupported)
	AbsCount uint32
}

//
// Get information about the track that is currently playing, its URI,
// position in the queue (Q:0), and elapsed time.  For Sonos @instanceId
// should always be 0.
//
func (this *AVTransport) GetPositionInfo(instanceId uint32) (*PositionInfo, error) {
	type Response struct {
		XMLName xml.Name
		PositionInfo
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetPositionInfo", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.PositionInfo, doc.Error()
}

//
// The return type of the GetDeviceCapabilities method
//
type DeviceCapabilities struct {
	// Configured sources of media
	PlayMedia string
	// ???? (possibly unsupported)
	RecMedia string
	// ???? (possibly unsupported)
	RecQualityModes string
}

//
// Return the device capabilities, sources of input media, recording
// media, and recoding quality modes.  For Sonos @instanceId should always
// be 0, and the record-related fields are unsupported.
//
func (this *AVTransport) GetDeviceCapabilities(instanceId uint32) (*DeviceCapabilities, error) {
	type Response struct {
		XMLName xml.Name
		DeviceCapabilities
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetDeviceCapabilities", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return &doc.DeviceCapabilities, doc.Error()
}

type TransportSettings struct {
	PlayMode       string
	RecQualityMode string
}

func (this *AVTransport) GetTransportSettings(instanceId uint32) (transportSettings *TransportSettings, err error) {
	type Response struct {
		XMLName xml.Name
		TransportSettings
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetTransportSettings", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	transportSettings = &doc.TransportSettings
	err = doc.Error()
	return
}

//
// Returns true if crossfade mode is active; false otherwise.  For Sonos
// @instanceId should always be 0.
//
func (this *AVTransport) GetCrossfadeMode(instanceId uint32) (bool, error) {
	type Response struct {
		XMLName       xml.Name
		CrossfadeMode bool
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetCrossfadeMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.CrossfadeMode, doc.Error()
}

//
// Stops playback and return to the beginning of the queue (Q:1).
// For Sonos @instanceId should always be 0.
//
func (this *AVTransport) Stop(instanceId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "Stop", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

// Playback at normal speed
const PlaySpeed_1 = "1"

//
// Starts or resumes playback at the given speed.  For Sonos @instanceId
// should always be 0; @speed is a fraction relative to normal speed
// (e.g. "1" or "1/2").
//
func (this *AVTransport) Play(instanceId uint32, speed string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Speed", speed},
	}
	response := Call(this.Svc, "Play", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

//
// Pause playback, prepared to resume at the current position.  For Sonos
// @instanceId should always be 0.
//
func (this *AVTransport) Pause(instanceId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "Pause", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

//
// Legal values for @unit in calls to Seek.
//
const (
	// Seek to the beginning of the specified track
	SeekMode_TRACK_NR = "TRACK_NR"
	// Seek to the given offset in the current track
	SeekMode_REL_TIME = "REL_TIME"
	// Seek to the specified section (not tested)
	SeekMode_SECTION  = "SECTION"
)

//
// A general function to seek within the playback queue (Q:0).  For Sonos
// @instanceId should always be 0; @unit should be one of the values given
// for seek mode (TRACK_NR, REL_TIME, or SECTION); and @target should
// give the track, time offset, or section where playback should resume.
//
// For TRACK_NR the integer track number relative to the start of the queue
// is supplied to @target.  For REL_TIME a duration in the format HH:MM:SS
// is given as @target.  SECTION is not tested.
//
func (this *AVTransport) Seek(instanceId uint32, unit, target string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Unit", unit},
		{"Target", target},
	}
	response := Call(this.Svc, "Seek", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

//
// Skip ahead to the next track in the queue (Q:).  For Sonos @instanceId
// should always be 0.  This method returns an error 711 if the current
// track is the last track in the queue.
//
func (this *AVTransport) Next(instanceId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "Next", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *AVTransport) NextProgrammedRadioTracks(instanceId uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "NextProgrammedRadioTracks", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

//
// Move to a previous track in the queue (Q:0).  For Sonos @instanceId
// should always be 0.  This method returns error 711 if the current track
// is the first track in the queue.
//
func (this *AVTransport) Previous(instanceId uint32) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "Previous", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *AVTransport) NextSection(instanceId uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "NextSection", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) PreviousSection(instanceId int) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "PreviousSection", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

const (
	PlayMode_NORMAL           = "NORMAL"
	PlayMode_REPEAT_ALL       = "REPEAT_ALL"
	PlayMode_SHUFFLE_NOREPEAT = "SHUFFLE_NOREPEAT"
	PlayMode_SHUFFLE          = "SHUFFLE"
)

func (this *AVTransport) SetPlayMode(instanceId uint32, newPlayMode string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"NewPlayMode", newPlayMode},
	}
	response := Call(this.Svc, "SetPlayMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) SetCrossfadeMode(instanceId uint32, crossfadeMode bool) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"CrossfadeMode", crossfadeMode},
	}
	response := Call(this.Svc, "SetCrossfadeMode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) NotifyDeletedURI(instanceId uint32, deletedURI string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"DeletedURI", deletedURI},
	}
	response := Call(this.Svc, "NotifyDeletedURI", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

//
// Returns a list of the actions that are valid at this time.  The list
// consists of human-readable strings, such as "Play", and "Stop".  For Sonos
// @instanceId will always be 0.
//
func (this *AVTransport) GetCurrentTransportActions(instanceId uint32) ([]string, error) {
	type Response struct {
		XMLName xml.Name
		Actions string
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetCurrentTransportActions", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return strings.Split(doc.Actions, ", "), doc.Error()
}

func (this *AVTransport) BecomeCoordinatorOfStandaloneGroup(instanceId uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "BecomeCoordinatorOfStandaloneGroup", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type BecomeGroupCoordinatorRequest struct {
	CurrentCoordinator    string
	CurrentGroupID        string
	OtherMembers          string
	TransportSettings     string
	CurrentURI            string
	CurrentURIMetaData    string
	SleepTimerState       string
	AlarmState            string
	StreamRestartState    string
	CurrentQueueTrackList string
}

func (this *AVTransport) BecomeGroupCoordinator(instanceId uint32, req *BecomeGroupCoordinatorRequest) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"CurrentCoordinator", req.CurrentCoordinator},
		{"CurrentGroupID", req.CurrentGroupID},
		{"OtherMembers", req.OtherMembers},
		{"TransportSettings", req.TransportSettings},
		{"CurrentURI", req.CurrentURI},
		{"CurrentURIMetaData", req.CurrentURIMetaData},
		{"SleepTimerState", req.SleepTimerState},
		{"AlarmState", req.AlarmState},
		{"StreamRestartState", req.StreamRestartState},
		{"CurrentQueueTrackList", req.CurrentQueueTrackList},
	}
	response := Call(this.Svc, "BecomeCoordinatorOfStandaloneGroup", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type BecomeGroupCoordinatorAndSourceRequest struct {
	CurrentCoordinator    string
	CurrentGroupID        string
	OtherMembers          string
	CurrentURI            string
	CurrentURIMetaData    string
	SleepTimerState       string
	AlarmState            string
	StreamRestartState    string
	CurrentAVTTrackList   string
	CurrentQueueTrackList string
	CurrentSourceState    string
	ResumePlayback        bool
}

func (this *AVTransport) BecomeGroupCoordinatorAndSource(instanceId uint32, req *BecomeGroupCoordinatorAndSourceRequest) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"CurrentCoordinator", req.CurrentCoordinator},
		{"CurrentGroupID", req.CurrentGroupID},
		{"OtherMembers", req.OtherMembers},
		{"CurrentURI", req.CurrentURI},
		{"CurrentURIMetaData", req.CurrentURIMetaData},
		{"SleepTimerState", req.SleepTimerState},
		{"AlarmState", req.AlarmState},
		{"StreamRestartState", req.StreamRestartState},
		{"CurrentAVTTrackList", req.CurrentAVTTrackList},
		{"CurrentQueueTrackList", req.CurrentQueueTrackList},
		{"CurrentSourceState", req.CurrentSourceState},
		{"ResumePlayback", req.ResumePlayback},
	}
	response := Call(this.Svc, "BecomeGroupCoordinatorAndSource", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type ChangeCoordinatorRequest struct {
	CurrentCoordinator   string
	NewCoordinator       string
	NewTransportSettings string
}

func (this *AVTransport) ChangeCoordinator(instanceId uint32, req *ChangeCoordinatorRequest) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"CurrentCoordinator", req.CurrentCoordinator},
		{"NewCoordinator", req.NewCoordinator},
		{"NewTransportSettings", req.NewTransportSettings},
	}
	response := Call(this.Svc, "ChangeCoordinator", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) ChangeTransportSettings(instanceId uint32, newTransportSettings, currentAVTransportURI string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"NewTransportSettings", newTransportSettings},
		{"CurrentTransportURI", currentAVTransportURI},
	}
	response := Call(this.Svc, "ChangeTransportSettings", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) ConfigureSleepTimer(instanceId uint32, newSleepTimerDuration string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"NewSleepTimerDuration", newSleepTimerDuration},
	}
	response := Call(this.Svc, "ConfigureSleepTimer", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) GetRemainingSleepTimerDuration(instanceId uint32) (remainingSleepTimerDuration string,
	currentSleepTimerGeneration uint32, err error) {
	type Response struct {
		XMLName                     xml.Name
		RemainingSleepTimerDuration string
		CurrentSleepTimerGeneration uint32
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetRemainingSleepTimerDuration", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	remainingSleepTimerDuration = doc.RemainingSleepTimerDuration
	currentSleepTimerGeneration = doc.CurrentSleepTimerGeneration
	err = doc.Error()
	return
}

type RunAlarmRequest struct {
	AlarmID            uint32
	LoggedStartTime    string
	Duration           string
	ProgramURI         string
	ProgramMetaData    string
	PlayMode           string
	Volume             uint32
	IncludeLinkedZones bool
}

func (this *AVTransport) RunAlarm(instanceId uint32, req *RunAlarmRequest) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"AlarmID", req.AlarmID},
		{"LoggedStartTime", req.LoggedStartTime},
		{"Duration", req.Duration},
		{"ProgramURI", req.ProgramURI},
		{"ProgramMetaData", req.ProgramMetaData},
		{"PlayMode", req.PlayMode},
		{"Volume", req.Volume},
		{"IncludeLinkedZones", req.IncludeLinkedZones},
	}
	response := Call(this.Svc, "RunAlarm", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

type StartAutoplayRequest struct {
	ProgramURI         string
	ProgramMetaData    string
	Volume             uint32
	IncludeLinkedZones bool
	ResetVolumeAfter   bool
}

func (this *AVTransport) StartAutoplay(instanceId uint32, req *StartAutoplayRequest) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"ProgramURI", req.ProgramURI},
		{"ProgramMetaData", req.ProgramMetaData},
		{"Volume", req.Volume},
		{"IncludeLinkedZones", req.IncludeLinkedZones},
		{"ResetVolumeAfter", req.ResetVolumeAfter},
	}
	response := Call(this.Svc, "StartAutoplay", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *AVTransport) GetRunningAlarmProperties(instanceId uint32) (alarmId uint32, groupId, loggedStartTime string, err error) {
	type Response struct {
		XMLName         xml.Name
		AlarmID         uint32
		GroupID         string
		LoggedStartTime string
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
	}
	response := Call(this.Svc, "GetRunningAlarmProperties", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	alarmId = doc.AlarmID
	groupId = doc.GroupID
	loggedStartTime = doc.LoggedStartTime
	err = doc.Error()
	return
}

func (this *AVTransport) SnoozeAlarm(instanceId uint32, duration string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"InstanceID", instanceId},
		{"Duration", duration},
	}
	response := Call(this.Svc, "SnoozeAlarm", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
