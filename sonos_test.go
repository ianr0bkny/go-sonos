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
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENTING SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package sonos_test

import (
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	"github.com/ianr0bkny/go-sonos/didl"
	"github.com/ianr0bkny/go-sonos/ssdp"
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
	"strings"
	"testing"
)

const (
	TEST_CONFIG        = "/home/ianr/.go-sonos"
	TEST_SONOS         = "kitchen"
	TEST_RECIVA        = "basement"
	TEST_DISCOVER_PORT = "13104"
	TEST_EVENTING_PORT = "13106"
	TEST_NETWORK       = "eth0"
)

var testSonos *sonos.Sonos

func initTestSonos(flags int) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	c := config.MakeConfig(TEST_CONFIG)
	c.Init()
	if dev := c.Lookup(TEST_SONOS); nil != dev {
		testSonos = sonos.Connect(dev, nil, flags)
	} else {
		log.Fatal("Could not create test instance")
	}
}

func getTestSonos(flags int) *sonos.Sonos {
	if nil == testSonos {
		initTestSonos(flags)
	}
	return testSonos
}

var testReciva *sonos.Reciva

func initTestReciva(flags int) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	c := config.MakeConfig(TEST_CONFIG)
	c.Init()
	if dev := c.Lookup(TEST_RECIVA); nil != dev {
		testReciva = sonos.ConnectReciva(dev, nil, flags)
	} else {
		log.Fatal("Could not create test instance")
	}
}

func getTestReciva(flags int) *sonos.Reciva {
	if nil == testReciva {
		initTestReciva(flags)
	}
	return testReciva
}

//
// AlarmClock
//
func TestAlarmClock(t *testing.T) {
	s := getTestSonos(sonos.SVC_ALARM_CLOCK)

	if currentTimeFormat, currentDateFormat, err := s.GetFormat(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetFormat() -> \"%s\",\"%s\"", currentTimeFormat, currentDateFormat)
	}

	if index, autoAdjustDst, err := s.GetTimeZone(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTimeZone() -> %d,%v", index, autoAdjustDst)
	}

	if index, autoAdjustDst, timeZone, err := s.GetTimeZoneAndRule(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTimeZoneAndRule() -> %d,%v,\"%s\"", index, autoAdjustDst, timeZone)
		if timeZone, err := s.GetTimeZoneRule(index); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("GetTimeZoneRule(index=%d) -> \"%v\"", index, timeZone)
		}
	}

	if currentTimeServer, err := s.GetTimeServer(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTimeServer() -> \"%s\"", currentTimeServer)
	}

	if getTimeNowResponse, err := s.GetTimeNow(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTimeNow() ->")
		t.Logf("\tCurrentUTCTime = \"%s\"", getTimeNowResponse.CurrentUTCTime)
		t.Logf("\tCurrentLocalTime = \"%s\"", getTimeNowResponse.CurrentLocalTime)
		t.Logf("\tCurrentTimeZone = \"%s\"", getTimeNowResponse.CurrentTimeZone)
		t.Logf("\tCurrenTimeGeneration = %d", getTimeNowResponse.CurrentTimeGeneration)
	}

	if currentAlarmList, currentAlarmListVersion, err := s.ListAlarms(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("ListAlarms() -> \"%s\",\"%s\"", currentAlarmList, currentAlarmListVersion)
	}

	if currentDailyIndexRefreshTime, err := s.GetDailyIndexRefreshTime(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetDailyIndexRefreshTime() -> \"%s\"", currentDailyIndexRefreshTime)
	}
}

//
// AVTransport
//
func TestAVTransport(t *testing.T) {
	s := getTestSonos(sonos.SVC_AV_TRANSPORT)

	if mediaInfo, err := s.GetMediaInfo(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetMediaInfo(0) ->")
		t.Logf("\tNrTracks = %d", mediaInfo.NrTracks)
		t.Logf("\tMediaDuration = \"%s\"", mediaInfo.MediaDuration)
		t.Logf("\tCurrentURI = \"%s\"", mediaInfo.CurrentURI)
		t.Logf("\tCurrentURIMetaData = \"%s\"", mediaInfo.CurrentURIMetaData)
		t.Logf("\tNextURI = \"%s\"", mediaInfo.NextURI)
		t.Logf("\tNextURIMetaData = \"%s\"", mediaInfo.NextURIMetaData)
		t.Logf("\tPlayMedium = \"%s\"", mediaInfo.PlayMedium)
		t.Logf("\tRecordMedium = \"%s\"", mediaInfo.RecordMedium)
		t.Logf("\tWriteStatus = \"%s\"", mediaInfo.WriteStatus)
	}

	if transportInfo, err := s.GetTransportInfo(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTransportInfo(0) ->")
		t.Logf("\tCurrentTransportState = \"%s\"", transportInfo.CurrentTransportState)
		t.Logf("\tCurrentTransportStatus = \"%s\"", transportInfo.CurrentTransportStatus)
		t.Logf("\tCurrentSpeed = \"%s\"", transportInfo.CurrentSpeed)
	}

	if positionInfo, err := s.GetPositionInfo(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetPositionInfo(0) ->")
		t.Logf("\tTrack = %d", positionInfo.Track)
		t.Logf("\tTrackDuration = \"%s\"", positionInfo.TrackDuration)
		t.Logf("\tTrackMetaData = \"%s\"", positionInfo.TrackMetaData)
		t.Logf("\tTrackURI = \"%s\"", positionInfo.TrackURI)
		t.Logf("\tRelTime = \"%s\"", positionInfo.RelTime)
		t.Logf("\tAbsTime = \"%s\"", positionInfo.AbsTime)
		t.Logf("\tRelCount = %d", positionInfo.RelCount)
		t.Logf("\tAbsCount = %d", positionInfo.AbsCount)
	}

	if deviceCapabilities, err := s.GetDeviceCapabilities(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetDeviceCapabilities() ->")
		t.Logf("\tPlayMedia = \"%s\"", deviceCapabilities.PlayMedia)
		t.Logf("\tRecMedia = \"%s\"", deviceCapabilities.RecMedia)
		t.Logf("\tRecQualityModes = \"%s\"", deviceCapabilities.RecQualityModes)
	}

	if transportSettings, err := s.GetTransportSettings(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTransportSettings() ->")
		t.Logf("\tPlayMode = \"%s\"", transportSettings.PlayMode)
		t.Logf("\tRecQualityMode = \"%s\"", transportSettings.RecQualityMode)
	}

	if crossfadeMode, err := s.GetCrossfadeMode(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetCrossfadeMode() -> %v", crossfadeMode)
	}

	if actions, err := s.GetCurrentTransportActions(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetCurrentTransportActions() -> \"%s\"", actions)
	}

	if remainingSleepTimerDuration, currentSleepTimerGeneration, err := s.GetRemainingSleepTimerDuration(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetRemainingSleepTimerDuration() -> \"%s\",%d", remainingSleepTimerDuration, currentSleepTimerGeneration)
	}

	/*
		if alarmId, groupId, loggedStartTime, err := s.GetRunningAlarmProperties(0); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("GetRunningAlarmProperties() ->")
			t.Logf("\tAlarmID = %d", alarmId)
			t.Logf("\tGroupID = \"%s\"", groupId)
			t.Logf("\tLoggedStartTime = \"%s\"", loggedStartTime)
		}
	*/
}

//
// ConnectiionManager
//
func TestConnectionManager(t *testing.T) {
	s := getTestSonos(sonos.SVC_CONNECTION_MANAGER)

	if source, sink, err := s.GetProtocolInfo(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetProtocolInfo() -> \"%s\",\"%s\"", source, sink)
	}

	if connectionIds, err := s.GetCurrentConnectionIDs(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetCurrentConnectionIds() -> \"%s\"", connectionIds)
	}
}

//
// ContentDirectory
// @see also TestBrowse
//
func TestContentDirectory(t *testing.T) {
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY)

	if searchCaps, err := s.GetSearchCapabilities(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetSearchCapabilities() -> \"%s\"", searchCaps)
	}

	if sortCaps, err := s.GetSortCapabilities(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetSortCapabilities() -> \"%s\"", sortCaps)
	}

	if id, err := s.GetSystemUpdateID(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetSystemUpdateID() -> \"%s\"", id)
	}

	if albumArtistDisplayOption, err := s.GetAlbumArtistDisplayOption(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetAlbumArtistDisplayOption() -> \"%s\"", albumArtistDisplayOption)
	}

	if lastIndexChange, err := s.GetLastIndexChange(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetLastIndexChange() -> \"%s\"", lastIndexChange)
	}

	if isIndexing, err := s.GetShareIndexInProgress(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetShareIndexInProgress() -> %v", isIndexing)
	}

	if isBrowseable, err := s.GetBrowseable(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetBrowseable() -> %v", isBrowseable)
	}
}

//
// DeviceProperties
//
func TestDeviceProperties(t *testing.T) {
	s := getTestSonos(sonos.SVC_DEVICE_PROPERTIES)

	if currentLEDState, err := s.GetLEDState(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetLEDState() -> \"%s\"", currentLEDState)
	}

	if currentInvisible, err := s.GetInvisible(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetInvisible() -> %v", currentInvisible)
	}

	if currentZoneName, currentIcon, err := s.GetZoneAttributes(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetZoneAttributes() -> \"%s\",\"%s\"", currentZoneName, currentIcon)
	}

	if currentHouseholdId, err := s.GetHouseholdID(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetHouseholdID() -> \"%s\"", currentHouseholdId)
	}

	if zoneInfo, err := s.GetZoneInfo(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetZoneInfo() ->")
		t.Logf("\tSerialNumber = \"%s\"", zoneInfo.SerialNumber)
		t.Logf("\tSoftwareVersion = \"%s\"", zoneInfo.SoftwareVersion)
		t.Logf("\tDisplaySoftwareVersion = \"%s\"", zoneInfo.DisplaySoftwareVersion)
		t.Logf("\tHardwareVersion = \"%s\"", zoneInfo.HardwareVersion)
		t.Logf("\tIPAddress = \"%s\"", zoneInfo.IPAddress)
		t.Logf("\tMACAddress = \"%s\"", zoneInfo.MACAddress)
		t.Logf("\tCopyrightInfo = \"%s\"", zoneInfo.CopyrightInfo)
		t.Logf("\tExtraInfo = \"%s\"", zoneInfo.ExtraInfo)
	}

	if includeLinkedZones, err := s.GetAutoplayLinkedZones(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetAutoplayLinkedZones() -> %v", includeLinkedZones)
	}

	if roomUUID, err := s.GetAutoplayRoomUUID(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetAutoplayRoomUUID() -> \"%s\"", roomUUID)
	}

	if currentVolume, err := s.GetAutoplayVolume(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetAutoplayVolume() -> %d", currentVolume)
	}

	if useVolume, err := s.GetUseAutoplayVolume(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetUseAutoplayVolume() -> %v", useVolume)
	}
}

//
// GroupManagement
//
func TestGroupManagement(t *testing.T) {
	// TODO
}

//
// MusicServices
//
func TestMusicServices(t *testing.T) {
	s := getTestSonos(sonos.SVC_MUSIC_SERVICES)

	if sessionId, err := s.GetSessionId(6 /*iheartradio*/, ""); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetSessionId() -> \"%s\"", sessionId)
	}

	if err := s.ListAvailableServices(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("ListAvailableServices() -> OK")
	}

	if err := s.UpdateAvailableServices(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("UpdateAvailableServices() -> OK")
	}
}

//
// RenderingControl
//
func TestRenderingControl(t *testing.T) {
	s := getTestSonos(sonos.SVC_RENDERING_CONTROL)

	if currentMute, err := s.GetMute(0, upnp.Channel_Master); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetMute() -> %v", currentMute)
	}

	if basicEQ, err := s.ResetBasicEQ(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("ResetBasicEQ() ->")
		t.Logf("\tBass = %d", basicEQ.Bass)
		t.Logf("\tTreble = %d", basicEQ.Treble)
		t.Logf("\tLoudness = %v", basicEQ.Loudness)
		t.Logf("\tLeftVolume = %d", basicEQ.LeftVolume)
		t.Logf("\tRightVolume = %d", basicEQ.RightVolume)
	}

	if currentVolume, err := s.GetVolume(0, upnp.Channel_Master); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetVolume() -> %d", currentVolume)
	}

	/*
		if currentVolume, err := s.GetVolumeDB(0, upnp.Channel_Master); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("GetVolumeDB() -> %d", currentVolume)
		}

		if min, max, err := s.GetVolumeDBRange(0, upnp.Channel_Master); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("GetVolumeDBRange() -> %d,%d", min, max)
		}
	*/

	if currentBass, err := s.GetBass(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetBass() -> %d", currentBass)
	}

	if currentTreble, err := s.GetTreble(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetTreble() -> %d", currentTreble)
	}

	if loudness, err := s.GetLoudness(0, upnp.Channel_Master); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetLoudness() -> %v", loudness)
	}

	if currentSupportsFixed, err := s.GetSupportsOutputFixed(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetSupportsOutputFixed() -> %v", currentSupportsFixed)
	}

	if currentFixed, err := s.GetOutputFixed(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetOutputFixed() -> %v", currentFixed)
	}

	if currentHeadphoneConnected, err := s.GetHeadphoneConnected(0); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("GetHeadphoneConnected() -> %v", currentHeadphoneConnected)
	}
}

//
// SystemProperties
//
func TestSystemProperties(t *testing.T) {
	// TODO
}

//
// ZoneGroupTopology
//
func TestZoneGroupTopology(t *testing.T) {
	s := getTestSonos(sonos.SVC_ZONE_GROUP_TOPOLOGY)

	if updateItem, err := s.CheckForUpdate(upnp.ALL, false, ""); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("CheckForUpdate() ->")
		t.Logf("\tType = \"%s\"", updateItem.Type)
		t.Logf("\tVersion = \"%s\"", updateItem.Version)
		t.Logf("\tUpdateURL = \"%s\"", updateItem.UpdateURL)
		t.Logf("\tDownloadSize = \"%s\"", updateItem.DownloadSize)
		t.Logf("\tManifestURL = \"%s\"", updateItem.ManifestURL)
	}
}

//
// Coverage
//
func TestCoverage(t *testing.T) {
	s := getTestSonos(sonos.SVC_ALL)
	sonos.Coverage(s)
}

//
// Discovery
//
func _TestDiscovery(t *testing.T) {
	if mgr, err := sonos.Discover(TEST_NETWORK, TEST_DISCOVER_PORT); nil != err {
		panic(err)
	} else {
		reactor := sonos.MakeReactor(TEST_NETWORK, TEST_EVENTING_PORT)
		found := sonos.ConnectAny(mgr, reactor, sonos.SVC_DEVICE_PROPERTIES)
		for _, s := range found {
			id, _ := s.GetHouseholdID()
			name, _, _ := s.GetZoneAttributes()
			t.Logf("Found device \"%s\",\"%s\"", id, name)
		}
	}
}

//
// Browse
//
func TestBrowse(t *testing.T) {
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY)

	t.Logf("Root Level Children")
	t.Logf("-------------------")
	if result, err := s.GetRootLevelChildren(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-15s %s", container.ID(), container.Title(), container.Class())
		}
	}

	t.Logf("===================")
	t.Logf("Queues")
	t.Logf("-------------------")
	if result, err := s.ListQueues(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-15s %s", container.ID(), container.Title(), container.Class())
		}
	}

	t.Logf("===================")
	t.Logf("Saved Queues")
	t.Logf("-------------------")
	if result, err := s.ListSavedQueues(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%-25s %s", container.Title(), container.Class())
		}
	}

	t.Logf("===================")
	t.Logf("Attributes")
	t.Logf("-------------------")
	if result, err := s.ListAttributes(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%-25s %s", container.Title(), container.Class())
		}
	}

	t.Logf("===================")
	t.Logf("Music Shares")
	t.Logf("-------------------")
	if result, err := s.ListMusicShares(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%-25s %s", container.Title(), container.Class())
		}
	}

	t.Logf("===================")
	t.Logf("Genres")
	t.Logf("-------------------")
	if result, err := s.GetAllGenres(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%-25s %s", container.Title(), container.Class())
		}
	}

	t.Logf("===================")
	t.Logf("R&B")
	t.Logf("-------------------")
	var target string
	if result, err := s.GetGenreArtists("R&B"); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%-25s %s", container.Title(), container.Class())
			if "John Legend" == container.Title() {
				target = container.ID()
			}
		}
	}

	t.Logf("===================")
	t.Logf("John Legend")
	t.Logf("-------------------")
	if result, err := s.ListChildren(target); nil != err {
		t.Fatal(err)
	} else {
		target = ""
		for _, container := range result {
			t.Logf("%-25s %s", container.Title(), container.Class())
			if "Get Lifted" == container.Title() {
				target = container.ID()
			}
		}
	}

	t.Logf("===================")
	t.Logf("Get Lifted")
	t.Logf("-------------------")
	if result, err := s.ListChildren(target); nil != err {
		t.Fatal(err)
	} else {
		target = ""
		for _, container := range result {
			t.Logf("%-32s %s", container.Title(), container.Class())
			if "I Can Change" == container.Title() {
				target = container.ID()
			}
		}
	}

	t.Logf("===================")
	t.Logf("I Can Change")
	t.Logf("-------------------")
	if result, err := s.GetMetadata(target); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%-32s", container.Title())
		}
	}
}

func TestRadio(t *testing.T) {
	// FAILS: It does not seem possible to Browse() R:
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY)

	t.Logf("Internet Radio")
	t.Logf("-------------------")
	if result, err := s.ListInternetRadio(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-25s %s", container.ID(), container.Title(), container.Class())
		}
	}
}

func TestQueue(t *testing.T) {
	// This test prints the current queue, removes the last track
	// in the queue, and then empties the queue completely.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)

	var lastId string
	t.Logf("Current Queue")
	t.Logf("-------------------")
	if result, err := s.GetQueueContents(); nil != err {
		t.Fatal(err)
	} else {
		for _, item := range result {
			t.Logf("%6s %s", item.ID(), item.Title())
			lastId = item.ID()
		}
	}

	if "" != lastId {
		if err := s.RemoveTrackFromQueue(0 /*instanceId*/, lastId, 0 /*updateId*/); nil != err {
			t.Fatal(err)
		}
	}

	if err := s.RemoveAllTracksFromQueue(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	}
}

func TestAddTrack(t *testing.T) {
	// This test duplicates the last track at the end of the queue.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)

	var lastTrackURI string
	t.Logf("Current Queue")
	t.Logf("-------------------")
	if result, err := s.GetQueueContents(); nil != err {
		t.Fatal(err)
	} else {
		for _, item := range result {
			t.Logf("%#v", item)
			lastTrackURI = item.Res()
		}
	}

	if "" != lastTrackURI {
		req := upnp.AddURIToQueueIn{
			EnqueuedURI: lastTrackURI,
		}
		if result, err := s.AddURIToQueue(0 /*instanceId*/, &req); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("%#v", result)
		}
	}
}

func TestAddAlbum(t *testing.T) {
	// Replace the current queue (Q:0) with George Harrison's Dark Horse
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if err := s.RemoveAllTracksFromQueue(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	}
	if objs, err := s.GetMetadata("A:ALBUM/Dark Horse"); nil != err {
		t.Fatal(err)
	} else {
		req := upnp.AddURIToQueueIn{
			EnqueuedURI: objs[0].Res(),
		}
		if result, err := s.AddURIToQueue(0 /*instanceId*/, &req); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("Added %d new track(s)", result.NumTracksAdded)
		}
	}
}

func TestRemoveRange(t *testing.T) {
	// This test removes the first two tracks from the queue.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	updateId, err := s.RemoveTrackRangeFromQueue(0, /*instanceId*/
		0, /*updateId*/
		1, /*startingIndex*/
		2 /*numberofTracks*/)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Logf("New UpdateId = %d", updateId)
	}
}

func TestReorder(t *testing.T) {
	// This test moves the first two tracks of Sgt. Pepper to the end.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	err := s.ReorderTracksInQueue(0, /*instanceId*/
		1,  /*startingIndex*/
		2,  /*numberOfTracks*/
		14, /*insertBefore*/
		0 /*updateId*/)
	if nil != err {
		t.Fatal(err)
	}
}

func TestGetCurrentTransportActions(t *testing.T) {
	// This test prints the currently available transport actions
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if actions, err := s.GetCurrentTransportActions(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	} else {
		for _, action := range actions {
			t.Logf(action)
		}
	}
}

func TestNamedQueue(t *testing.T) {
	// This test creates a new named queue and then deletes it
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	obj, err := s.SaveQueue(0 /*instanceId*/, "borodin" /*title*/, "" /*objectId*/)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", obj)
		if err = s.DestroyObject(obj); nil != err {
			t.Fatal(err)
		} else {
			t.Logf("queue removed")
		}
	}
}

func TestPositionInfo(t *testing.T) {
	// This test fetches and prints the current position info
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if info, err := s.GetPositionInfo(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%d", info.Track)
		t.Logf("%s", info.TrackDuration)
		t.Logf("%s", info.TrackMetaData)
		t.Logf("%s", info.TrackURI)
		t.Logf("%s", info.RelTime)
		t.Logf("%s", info.AbsTime)
		t.Logf("%d", info.RelCount)
		t.Logf("%d", info.AbsCount)
	}
}

func TestMediaInfo(t *testing.T) {
	// This test fetches and prints the current media info
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if info, err := s.GetMediaInfo(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", info.NrTracks)
		t.Logf("%#v", info.MediaDuration)
		t.Logf("%#v", info.CurrentURI)
		t.Logf("%#v", info.CurrentURIMetaData)
		t.Logf("%#v", info.NextURI)
		t.Logf("%#v", info.NextURIMetaData)
		t.Logf("%#v", info.PlayMedium)
		t.Logf("%#v", info.RecordMedium)
		t.Logf("%#v", info.WriteStatus)
	}
}

func TestDeviceCapabilities(t *testing.T) {
	// This test fetches and prints the capabilities of the current device
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if caps, err := s.GetDeviceCapabilities(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", caps.PlayMedia)
		t.Logf("%#v", caps.RecMedia)
		t.Logf("%#v", caps.RecQualityModes)
	}
}

func TestNavigation(t *testing.T) {
	// Test of the Next and Previous methods
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if err := s.Next(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	}
	if err := s.Previous(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	}
}

func TestPlayPauseStop(t *testing.T) {
	// Test of the Play, Pause, and Stop methods
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if info, err := s.GetTransportInfo(0); nil != err {
		t.Fatal(err)
	} else {
		if upnp.State_PAUSED_PLAYBACK == info.CurrentTransportState {
			s.Stop(0 /*instanceId*/)
		} else if upnp.State_STOPPED == info.CurrentTransportState {
			s.Play(0 /*instanceId*/, upnp.PlaySpeed_1)
		} else if upnp.State_PLAYING == info.CurrentTransportState {
			s.Pause(0 /*instanceId*/)
		}
	}
}

func TestAddMultiple(t *testing.T) {
	// This test adds Reptile from NIN and Reset from Big Boi to the
	// queue using AddMultipleURIsToQueue.  Note that the containing
	// album must be known to find the track.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)

	var tracks_uri string
	if objs, err := s.GetMetadata("A:TRACKS"); nil != err {
		t.Fatal(err)
	} else {
		tracks_uri = objs[0].Res()
	}

	var uris []string
	if objs, err := s.GetTrackFromAlbum("The Downward Spiral", "Reptile"); nil != err {
		t.Fatal(err)
	} else {
		for _, obj := range objs {
			uris = append(uris, obj.Res())
		}
	}

	if objs, err := s.GetTrackFromAlbum("Speakerboxxx", "Reset"); nil != err {
		t.Fatal(err)
	} else {
		for _, obj := range objs {
			uris = append(uris, obj.Res())
		}
	}

	enqueued_uris := strings.Join(uris, " ")

	req := &upnp.AddMultipleURIsToQueueIn{
		UpdateID:                        0,
		NumberOfURIs:                    uint32(len(uris)),
		EnqueuedURIs:                    enqueued_uris,
		EnqueuedURIsMetaData:            didl.EmptyDocuments(len(uris)),
		ContainerURI:                    tracks_uri,
		ContainerMetaData:               didl.EmptyDocument(),
		DesiredFirstTrackNumberEnqueued: 0,
		EnqueueAsNext:                   false,
	}
	if resp, err := s.AddMultipleURIsToQueue(0 /*instanceId*/, req); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", resp)
	}
}

func TestSeek(t *testing.T) {
	// A test that sets playback to 30 seconds from the beginning
	// of the fifth track in the queue.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if err := s.Seek(0 /*instanceId*/, upnp.SeekMode_TRACK_NR, "6"); nil != err {
		t.Fatal(err)
	}

	if err := s.Seek(0 /*instanceId*/, upnp.SeekMode_REL_TIME, "0:0:30"); nil != err {
		t.Fatal(err)
	}
}

func TestUseQueue(t *testing.T) {
	// Select the current queue (Q:0) (e.g., if we were already
	// listening to the radio)
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if info, err := s.GetMediaInfo(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	} else if sonos.ObjectID_Queue_AVT_Instance_0 != info.CurrentURI {
		if data, err := s.GetMetadata(sonos.ObjectID_Queue_AVT_Instance_0); nil != err {
			t.Fatal(err)
		} else {
			if err := s.SetAVTransportURI(0 /*instanceId*/, data[0].Res(), "" /*metadata*/); nil != err {
				t.Fatal(err)
			}
		}
	}
}

func TestTransportSettings(t *testing.T) {
	// Get the current playback mode and recording quality
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if data, err := s.GetTransportSettings(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", data)
	}
}

func TestSetPlayMode(t *testing.T) {
	// Changes the playback mode to shuffle, tries to set it to some
	// garbage value, then reverts it to normal.
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if err := s.SetPlayMode(0 /*instanceId*/, upnp.PlayMode_SHUFFLE); nil != err {
		t.Fatal(err)
	} else if err := s.SetPlayMode(0 /*instanceId*/, "SOMETHING_INVALID"); nil != err {
		s.SetPlayMode(0 /*instanceId*/, upnp.PlayMode_NORMAL)
	}
}

func TestSections(t *testing.T) {
	// Advance to the next section, then return to the previous section
	s := getTestSonos(sonos.SVC_CONTENT_DIRECTORY | sonos.SVC_AV_TRANSPORT)
	if err := s.NextSection(0); nil != err {
		t.Error(err)
	}
	if err := s.PreviousSection(0); nil != err {
		t.Error(err)
	}
}

func TestCrossfadeMode(t *testing.T) {
	// Tests setting and geting crossfade mode
	s := getTestSonos(sonos.SVC_AV_TRANSPORT)
	var mode, pmode bool
	var err error
	if pmode, err = s.GetCrossfadeMode(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	}
	mode = !pmode
	if err = s.SetCrossfadeMode(0 /*instanceId*/, mode); nil != err {
		t.Fatal(err)
	}
	if pmode, err = s.GetCrossfadeMode(0 /*instanceId*/); nil != err {
		t.Fatal(err)
	}
	if mode != pmode {
		t.Error("GetCrossfadeMode failed to set crossfade mode")
	}
	if err = s.SetCrossfadeMode(0 /*instanceId*/, false); nil != err {
		t.Fatal(err)
	}
}

func TestGetZoneInfo(t *testing.T) {
	// Tests getting basic information about the Sonos appliance
	s := getTestSonos(sonos.SVC_DEVICE_PROPERTIES)
	if info, err := s.GetZoneInfo(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", info)
	}
}

func TestGetZoneGroup(t *testing.T) {
	// Tests getting basic information about the Sonos appliance
	s := getTestSonos(sonos.SVC_ZONE_GROUP_TOPOLOGY)
	if info, err := s.GetZoneGroupAttributes(); nil != err {
		t.Fatal(err)
	} else {
		t.Logf("%#v", info)
	}
}

func handleEvent_TestEvent(reactor upnp.Reactor, c chan bool) {
	for {
		select {
		case evt := <-reactor.Channel():
			switch evt.Type() {
			case upnp.AlarmClock_EventType:
				b := evt.(upnp.AlarmClockEvent)
				log.Printf("%#v", b)
			case upnp.AVTransport_EventType:
				b := evt.(upnp.AVTransportEvent)
				log.Printf("%#v", b)
			case upnp.ConnectionManager_EventType:
				b := evt.(upnp.ConnectionManagerEvent)
				log.Printf("%#v", b)
			case upnp.ContentDirectory_EventType:
				b := evt.(upnp.ContentDirectoryEvent)
				log.Printf("%#v", b)
			case upnp.DeviceProperties_EventType:
				b := evt.(upnp.DevicePropertiesEvent)
				log.Printf("%#v", b)
			case upnp.GroupManagement_EventType:
				b := evt.(upnp.GroupManagementEvent)
				log.Printf("%#v", b)
			case upnp.MusicServices_EventType:
				b := evt.(upnp.MusicServicesEvent)
				log.Printf("%#v", b)
			case upnp.RenderingControl_EventType:
				b := evt.(upnp.RenderingControlEvent)
				log.Printf("%#v", b)
			case upnp.SystemProperties_EventType:
				b := evt.(upnp.SystemPropertiesEvent)
				log.Printf("%#v", b)
			case upnp.ZoneGroupTopology_EventType:
				b := evt.(upnp.ZoneGroupTopologyEvent)
				log.Printf("%#v", b)
			}
		}
	}
}

func handleEvent_TestEventBrief(reactor upnp.Reactor, c chan bool) {
	for {
		select {
		case evt := <-reactor.Channel():
			switch evt.Type() {
			case upnp.AlarmClock_EventType:
				log.Printf("ALARM_CLOCK")
			case upnp.AVTransport_EventType:
				log.Printf("AV_TRANSPORT")
			case upnp.ConnectionManager_EventType:
				log.Printf("CONNECTION_MANAGER")
			case upnp.ContentDirectory_EventType:
				log.Printf("CONTENT_DIRECTORY")
			case upnp.DeviceProperties_EventType:
				log.Printf("DEVICE_PROPERTIES")
			case upnp.GroupManagement_EventType:
				log.Printf("GROUP_MANAGEMENT")
			case upnp.MusicServices_EventType:
				log.Printf("MUSIC_SERVICES")
			case upnp.RenderingControl_EventType:
				log.Printf("RENDERING_CONTROL")
			case upnp.SystemProperties_EventType:
				log.Printf("SYSTEM_PROPERTIES")
			case upnp.ZoneGroupTopology_EventType:
				log.Printf("ZONE_GROUP_TOPOLOGY")
			}
		}
	}
}

func TestEvent(t *testing.T) {
	// Startup and listen to events
	log.SetFlags(log.Ltime | log.Lshortfile)
	c := config.MakeConfig(TEST_CONFIG)
	c.Init()
	if dev := c.Lookup(TEST_SONOS); nil != dev {
		exit_chan := make(chan bool)
		reactor := sonos.MakeReactor(TEST_NETWORK, TEST_EVENTING_PORT)
		go handleEvent_TestEventBrief(reactor, exit_chan)
		testSonos = sonos.Connect(dev, reactor, sonos.SVC_ALL)
		<-exit_chan
	} else {
		log.Fatal("Could not create test instance")
	}
}

func TestGetZoneGroupState(t *testing.T) {
	s := getTestSonos(sonos.SVC_ZONE_GROUP_TOPOLOGY)
	if state, err := s.GetZoneGroupState(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("%#v", state)
	}
}

func TestGrace(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	sonos.Coverage(r)
}

func TestListPresets(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if presets, err := r.ListPresets(0); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("PRESET: %v", presets)
	}
}

func TestIdArray(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if token, array, err := r.IdArray(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("--> %v %v", token, array)
	}
}

func TestTracksMax(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if tracksMax, err := r.TracksMax(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("%v", tracksMax)
	}
}

func TestPowerState(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if state, err := r.GetPowerState(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Power state is (%v)", state)
	}
	if state, err := r.SetPowerState("Off"); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Power state is (%v)", state)
	}
}

func TestDisplayLanguage(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if lang, iso, err := r.GetCurrentDisplayLanguage(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Language = %v", lang)
		log.Printf("ISO: %v", iso)
	}
}

func TestDisplayLanguageList(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if lang, iso, err := r.GetDisplayLanguages(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Language = %v", lang)
		log.Printf("ISO: %v", iso)
	}
}

func TestGetNumberOfPresets(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if n, err := r.GetNumberOfPresets(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Number of presets: %v", n)
	}
}

func TestRecivaGetDateTime(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if n, err := r.GetDateTime(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Date-time: %v", n)
	}
}

func TestRecivaGetTimeZone(t *testing.T) {
	r := getTestReciva(sonos.SVC_ALL)
	if n, err := r.GetTimeZone(); nil != err {
		log.Fatal(err)
	} else {
		log.Printf("Timezone: %v", n)
	}
}

////////////////////////////////////////////////////////
// Issue #4
////////////////////////////////////////////////////////
func read_events(c chan upnp.Event) {
	for {
		select {
		case <-c:
		}
	}
}

func TestIssue_4(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.Printf("Discovery: Starting")
	mgr, err := sonos.Discover("eth0", "13104")
	if nil != err {
		panic(err)
	}
	log.Printf("Discovery: Done; Reactor: Starting")
	reactor := sonos.MakeReactor("eth0", "13106")
	go read_events(reactor.Channel()) ///// <------------
	log.Printf("Reactor: Running; Query: Starting")
	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey(sonos.MUSIC_SERVICES): -1,
	}
	res := mgr.QueryServices(qry)
	log.Printf("Query: Done; Connect: Starting")
	if dev_list, has := res[sonos.MUSIC_SERVICES]; has {
		for _, dev := range dev_list {
			if sonos.SONOS == dev.Product() {
				if _, err := upnp.Describe(dev.Location()); nil != err {
					panic(err)
				} else {
					sonos.Connect(dev, reactor, sonos.SVC_ALL)
					log.Printf("Connect: Done")
					break
				}
			}
		}
	}
}
