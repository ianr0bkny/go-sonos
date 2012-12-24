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

package sonos_test

import (
	"github.com/ianr0bkny/go-sonos"
	"log"
	"testing"
)

func TestDiscover(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)
	if mgr, err := sonos.Discover("eth0", "13104"); nil != err {
		panic(err)
	} else {
		reactor := sonos.MakeReactor("eth0", "13105")
		found := sonos.ConnectAny(mgr, reactor)
		for _, s := range found {
			id := s.GetHouseholdID()
			name, _ := s.GetZoneAttributes()
			caps := s.GetSearchCapabilities()

			s.SetPlayMode(0, sonos.REPEAT_ALL)
			s.GetCrossfadeMode(0)
			s.GetTransportSettings(0)
			s.GetDeviceCapabilities(0)
			s.GetPositionInfo(0)
			s.GetTransportInfo(0)
			s.GetMediaInfo(0)
			//s.Play(0, "1")
			//s.Stop(0)
			s.GetProtocolInfo()
			s.GetCurrentConnectionIDs()
			s.GetCurrentConnectionInfo(0)

			log.Printf("Found device %s `%s' '%s'", id, name, caps)

			/*
				// browse root-level metadata
				s.Browse("0", "BrowseMetadata", "*", 0, 0, "")
				// browse children of the root
				s.Browse("0", "BrowseDirectChildren", "*", 0, 0, "")
				// browse music shares
				s.Browse("S:", "BrowseDirectChildren", "*", 0, 0, "")
				// browse the //perseus/sonos share
				s.Browse("S://perseus/sonos", "BrowseDirectChildren", "*", 0, 0, "")
				// browse the //perseus/sonos/iTunes share
				s.Browse("S://perseus/sonos/iTunes", "BrowseDirectChildren", "*", 0, 0, "")
				// browse the //perseus/sonos/iTunes/Music share
				s.Browse("S://perseus/sonos/iTunes/Music", "BrowseDirectChildren", "*", 0, 0, "")
				// browse the //perseus/sonos/iTunes/Music/The Who share
				s.Browse("S://perseus/sonos/iTunes/Music/The Who", "BrowseDirectChildren", "*", 0, 0, "")
			*/
			// browse the //perseus/sonos/iTunes/Music/The Who/Tommy share
			s.Browse("S://perseus/sonos/iTunes/Music/The Who/Tommy", "BrowseDirectChildren", "*", 0, 0, "")

			/*
				// browse music attributes
				s.Browse("A:", "BrowseDirectChildren", "*", 0, 0, "")
				// return list of composers
				s.Browse("A:COMPOSER", "BrowseDirectChildren", "dc:title", 0, 0, "")
			*/
		}
	}
}
