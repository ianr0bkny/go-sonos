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

	// ????
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
	UpdateID                        uint32
	NumberOfURIs                    uint32
	EnqueuedURIs                    string
	EnqueuedURIsMetaData            string
	ContainerURI                    string
	ContainerMetaData               string
	DesiredFirstTrackNumberEnqueued uint32
	EnqueueAsNext                   bool
}

//
// Output parameters for AddMultipleURIsToQueue.
//
type AddMultipleURIsToQueueOut struct {
	FirstTrackNumberEnqueued uint32
	NumTracksAdded           uint32
	NewQueueLength           uint32
	NewUpdateID              uint32
}

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

func (this *AVTransport) ReorderTracksInQueue(instanceId, startingIndex, numberOfTracks, insertBefore, updateId uint32) (err error) {
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
	err = doc.Error()
	return
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

func (this *AVTransport) RemoveTrackRangeFromQueue(instanceId, updateId, startingIndex, numberOfTracks uint32) (newUpdateId uint32, err error) {
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
	newUpdateId = doc.NewUpdateID
	err = doc.Error()
	return
}

//
// Remove all tracks from the queue (Q:0).  For Sonos instanceId will
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

func (this *AVTransport) SaveQueue(instanceId uint32, title, objectId string) (assignedObjectId string, err error) {
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
	assignedObjectId = doc.AssignedObjectID
	err = doc.Error()
	return
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

type MediaInfo struct {
	NrTracks           uint32
	MediaDuration      string
	CurrentURI         string
	CurrentURIMetaData string
	NextURI            string
	NextURIMetaData    string
	PlayMedium         string
	RecordMedium       string
	WriteStatus        string
}

func (this *AVTransport) GetMediaInfo(instanceId uint32) (mediaInfo *MediaInfo, err error) {
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
	mediaInfo = &doc.MediaInfo
	err = doc.Error()
	return
}

type TransportInfo struct {
	CurrentTransportState  string
	CurrentTransportStatus string
	CurrentSpeed           string
}

func (this *AVTransport) GetTransportInfo(instanceId uint32) (transportInfo *TransportInfo, err error) {
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
	transportInfo = &doc.TransportInfo
	err = doc.Error()
	return
}

type PositionInfo struct {
	Track         uint32
	TrackDuration string
	TrackMetaData string
	TrackURI      string
	RelTime       string
	AbsTime       string
	RelCount      uint32
	AbsCount      uint32
}

func (this *AVTransport) GetPositionInfo(instanceId uint32) (positionInfo *PositionInfo, err error) {
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
	positionInfo = &doc.PositionInfo
	err = doc.Error()
	return
}

type DeviceCapabilities struct {
	PlayMedia       string
	RecMedia        string
	RecQualityModes string
}

func (this *AVTransport) GetDeviceCapabilities(instanceId uint32) (deviceCapabilities *DeviceCapabilities, err error) {
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
	deviceCapabilities = &doc.DeviceCapabilities
	err = doc.Error()
	return
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

func (this *AVTransport) GetCrossfadeMode(instanceId uint32) (crossfadeMode bool, err error) {
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
	crossfadeMode = doc.CrossfadeMode
	err = doc.Error()
	return
}

func (this *AVTransport) Stop(instanceId uint32) (err error) {
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
	err = doc.Error()
	return
}

const PlaySpeed_1 = "1"

func (this *AVTransport) Play(instanceId uint32, speed string) (err error) {
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
	err = doc.Error()
	return
}

func (this *AVTransport) Pause(instanceId uint32) (err error) {
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
	err = doc.Error()
	return
}

const (
	SeekMode_TRACK_NR = "TRACK_NR"
	SeekMode_REL_TIME = "REL_TIME"
	SeekMode_SECTION  = "SECTION"
)

func (this *AVTransport) Seek(instanceId uint32, unit, target string) (err error) {
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
	err = doc.Error()
	return
}

func (this *AVTransport) Next(instanceId uint32) (err error) {
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
	err = doc.Error()
	return
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

func (this *AVTransport) Previous(instanceId uint32) (err error) {
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
	err = doc.Error()
	return
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

func (this *AVTransport) GetCurrentTransportActions(instanceId uint32) (actions string, err error) {
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
	actions = doc.Actions
	err = doc.Error()
	return
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
