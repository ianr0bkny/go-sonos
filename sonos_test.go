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
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
	"testing"
)

const (
	TEST_CONFIG        = "/home/ianr/.go-sonos"
	TEST_DEVICE        = "kitchen"
	TEST_DISCOVER_PORT = "13104"
	TEST_EVENTING_PORT = "13105"
	TEST_NETWORK       = "eth0"
)

var testSonos *sonos.Sonos

func initTestSonos() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	c := config.MakeConfig(TEST_CONFIG)
	c.Init()
	if dev := c.Lookup(TEST_DEVICE); nil != dev {
		reactor := sonos.MakeReactor(TEST_NETWORK, TEST_EVENTING_PORT)
		testSonos = sonos.Connect(dev, reactor)
	} else {
		log.Fatal("Could not create test instance")
	}
}

func getTestSonos() *sonos.Sonos {
	if nil == testSonos {
		initTestSonos()
	}
	return testSonos
}

//
// AlarmClock
//
func TestAlarmClock(t *testing.T) {
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()

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
	s := getTestSonos()
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
		found := sonos.ConnectAny(mgr, reactor)
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
	s := getTestSonos()

	t.Logf("Root Level Children")
	t.Logf("-------------------")
	if result, err := s.GetRootLevelChildren(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-15s %s", container.ID, container.Title, container.Class)
		}
	}

	t.Logf("===================")
	t.Logf("Queues")
	t.Logf("-------------------")
	if result, err := s.ListQueues(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-15s %s", container.ID, container.Title, container.Class)
		}
	}

	t.Logf("===================")
	t.Logf("Saved Queues")
	t.Logf("-------------------")
	if result, err := s.ListSavedQueues(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-25s %s", container.ID, container.Title, container.Class)
		}
	}

	/*
	t.Logf("===================")
	t.Logf("Internet Radio")
	t.Logf("-------------------")
	if result, err := s.ListInternetRadio(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%3s %-25s %s", container.ID, container.Title, container.Class)
		}
	}
	*/

	t.Logf("===================")
	t.Logf("Attributes")
	t.Logf("-------------------")
	if result, err := s.ListAttributes(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%13s %-25s %s", container.ID, container.Title, container.Class)
		}
	}

	t.Logf("===================")
	t.Logf("Music Shares")
	t.Logf("-------------------")
	if result, err := s.ListMusicShares(); nil != err {
		t.Fatal(err)
	} else {
		for _, container := range result {
			t.Logf("%18s %-20s %s", container.ID, container.Title, container.Class)
		}
	}

	//s.GetQueueContents()
}

//			/*
//				// browse root-level metadata
//				s.Browse("0", "BrowseMetadata", "*", 0, 0, "")
//				// browse children of the root
//				s.Browse("0", "BrowseDirectChildren", "*", 0, 0, "")
//				// browse music shares
//				s.Browse("S:", "BrowseDirectChildren", "*", 0, 0, "")
//				// browse the //perseus/sonos share
//				s.Browse("S://perseus/sonos", "BrowseDirectChildren", "*", 0, 0, "")
//				// browse the //perseus/sonos/iTunes share
//				s.Browse("S://perseus/sonos/iTunes", "BrowseDirectChildren", "*", 0, 0, "")
//				// browse the //perseus/sonos/iTunes/Music share
//				s.Browse("S://perseus/sonos/iTunes/Music", "BrowseDirectChildren", "*", 0, 0, "")
//				// browse the //perseus/sonos/iTunes/Music/The Who share
//				s.Browse("S://perseus/sonos/iTunes/Music/The Who", "BrowseDirectChildren", "*", 0, 0, "")
//			*/
//			// browse the //perseus/sonos/iTunes/Music/The Who/Tommy share
//			x, _ := s.Browse("S://perseus/sonos/iTunes/Music/The Who/Tommy", "BrowseDirectChildren", "*", 0, 0, "")
//			log.Printf("%#v", x)
//
//			/*
//				// browse music attributes
//				s.Browse("A:", "BrowseDirectChildren", "*", 0, 0, "")
//				// return list of composers
//				s.Browse("A:COMPOSER", "BrowseDirectChildren", "dc:title", 0, 0, "")
//			*/
