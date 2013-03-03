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

package config_test

import (
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	"log"
	"os"
	"testing"
)

const (
	alias     = "kitchen"
	configdir = "dot_go-sonos"
	location  = "http://192.168.1.44:1400/xml/device_description.xml"
	uuid      = "RINCON_000E58741A8401400"
)

func TestConfig(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	if err := os.RemoveAll(configdir); nil != err {
		panic(err)
	}

	c := config.MakeConfig(configdir)
	c.Init()
	c.AddBookmark(uuid, sonos.SONOS, location, uuid)
	c.AddAlias(uuid, alias)
	c.Save()
	c = nil

	d := config.MakeConfig(configdir)
	d.Init()

	bookmark := d.Bookmarks[uuid]
	if location != bookmark.Location {
		panic("failed")
	}

	if dev := d.Lookup(alias); nil == dev {
		panic("failed")
	}

	os.RemoveAll(configdir)
}
