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
	"github.com/ianr0bkny/go-sonos/ssdp"
	"log"
)

// This code identifies UPnP devices on the netork that support the
// MusicServices API.
func main() {
	log.Print("go-sonos example discovery\n")

	mgr := ssdp.MakeManager()

	// Discover()
	//  eth0 := Network device to query for UPnP devices
	// 11209 := Free local port for discovery replies
	// false := Do not subscribe for asynchronous updates
	mgr.Discover("eth0", "11209", false)

	// SericeQueryTerms
	// A map of service keys to minimum required version
	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey("schemas-upnp-org-MusicServices"): -1,
	}

	// Look for the service keys in qry in the database of discovered devices
	result := mgr.QueryServices(qry)
	if dev_list, has := result["schemas-upnp-org-MusicServices"]; has {
		for _, dev := range dev_list {
			log.Printf("%s %s %s %s %s\n", dev.Product(), dev.ProductVersion(), dev.Name(), dev.Location(), dev.UUID())
		}
	}
	mgr.Close()
}
