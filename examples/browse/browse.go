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
)

const (
	instanceId = 0
)

// This code locates a content directory device on the network
func main() {
	log.Print("go-sonos example discovery\n")

	mgr := ssdp.MakeManager()
	mgr.Discover("eth0", "11209", false)
	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey("schemas-upnp-org-ContentDirectory"): -1,
	}

	result := mgr.QueryServices(qry)
	if dev_list, has := result["schemas-upnp-org-ContentDirectory"]; has {
		for _, dev := range dev_list {
			log.Printf("%s %s %s %s %s\n", dev.Product(), dev.ProductVersion(), dev.Name(), dev.Location(), dev.UUID())
			s := sonos.Connect(dev, nil, sonos.SVC_CONTENT_DIRECTORY)

			//Method 1
			if tracks, err := s.GetAlbumTracks("The Beatles"); nil != err {
				panic(err)
			} else {
				for _, track := range tracks {
					if "Long, Long, Long" == track.Title() {
						log.Printf("%#v", track)
						log.Printf("%#v", track.Res())
						if objects, err := s.GetMetadata(track.ID()); nil != err {
							panic(err)
						} else {
							for _, object := range objects {
								log.Printf("--> %#v", object)
							}
						}
					}
				}
			}

			//Method 2
			if tracks, err := s.GetTrackFromAlbum("The Beatles", "Long, Long, Long"); nil != err {
				panic(err)
			} else {
				for _, track := range tracks {
					log.Printf("%#v", track)
					log.Printf("%#v", track.Res())
					if objects, err := s.GetMetadata(track.ID()); nil != err {
						panic(err)
					} else {
						for _, object := range objects {
							log.Printf("--> %#v", object)
						}
					}
				}
			}
			break
		}
	}
	mgr.Close()
}
