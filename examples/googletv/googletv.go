//
// go-sonos
// ========
//
// Copyright (c) 2012-2015, Ian T. Richards <ianr@panix.com>
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
package main

import (
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/ssdp"
	"log"
	"strings"
)

const (
	instanceId = 0
)

// This code locates a GoogleTV device on the network
func main() {
	log.Print("go-sonos example discovery\n")

	mgr := ssdp.MakeManager()
	mgr.Discover("eth0", "11209", false)
	dev_map := mgr.Devices()
	for _, dev := range dev_map {
		if "NSZ-GS7" == dev.Product() {
			log.Printf("%s %s %s %s %s\n", dev.Product(), dev.ProductVersion(), dev.Name(), dev.Location(), dev.UUID())
			keys := dev.Services()
			for _, key := range keys {
				log.Printf("\t%s\n", key)
			}

			s := sonos.Connect(dev, nil, sonos.SVC_CONNECTION_MANAGER|sonos.SVC_RENDERING_CONTROL|sonos.SVC_AV_TRANSPORT)

			if source, sink, err := s.GetProtocolInfo(); nil != err {
				panic(err)
			} else {
				log.Printf("Source: %v", source)
				for _, sink := range strings.Split(sink, ",") {
					log.Printf("Sink: %v", sink)
				}
			}

			if connection_ids, err := s.GetCurrentConnectionIDs(); nil != err {
				panic(err)
			} else {
				log.Printf("ConnectionIDs: %v", connection_ids)
			}

			if presets, err := s.ListPresets(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("Preset: %v", presets)
			}

			if media_info, err := s.GetMediaInfo(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("%v", media_info)
			}

			if transport_info, err := s.GetTransportInfo(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("%v", transport_info)
			}

			if position_info, err := s.GetPositionInfo(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("%v", position_info)
			}

			if device_capabilities, err := s.GetDeviceCapabilities(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("%v", device_capabilities)
			}

			if transport_settings, err := s.GetTransportSettings(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("%v", transport_settings)
			}

			if actions, err := s.GetCurrentTransportActions(instanceId); nil != err {
				panic(err)
			} else {
				log.Printf("%v", actions)
			}

			/*TODO*/
			/*
				if err := s.SetAVTransportURI(instanceId, uri, metadata); nil != err {
					panic(err)
				}
			*/
		}
	}
	mgr.Close()
}
