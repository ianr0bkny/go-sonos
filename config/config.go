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

//
// A module to support bookmarking discovered UPnP devices.
//
// This module is intended to solve the problem of addressing commands to
// a specific UPnP device, where each device is known by a user-unfriendly
// UUID, where network addresses are subject to change, ala DHCP, and where
// it is not desirable to repeat device discovery each time a command is
// to be sent.
//
// Since device discovery is slow, it should only be run once in a while
// when DHCP leases would have expired.  The discovery process writes the
// list of discovered devices to the list of bookmarks.
//
// Once, when the device is installed, an association is added to the
// list of bookmarks, making a memorable alias, i.e. 'kitchen', point to
// a particular UUID.
//
// Now whenever network addresses change disovery can be rerun, with
// aliases automatically pointing to the new address associated with the
// static UUID.
//
// Look at sonos_test.go for examples of how this class is used.
//
package config

import (
	"encoding/json"
	"github.com/ianr0bkny/go-sonos/ssdp"
	"io"
	"log"
	"os"
	"path"
)

//
// A container for the runtime configuration used by go-sonos application.
//
type Config struct {
	// The path to the configuration directory
	dirname string
	// A handle to the configuration directory
	dir *os.File
	// A set of discovered devices
	Bookmarks Bookmarks
}

//
// A strucutre that holds all of the fields required to build a UPnP
// device without first trying to discover it.
//
type Bookmark struct {
	// A memorable string standing in for a UUID
	Alias string `json:"alias,omitempty"`
	// The name of the device's product, e.g. 'Sonos'
	Product string `json:"product,omitempty"`
	// The product version (e.g. "28.1-83040 (BR100)")
	ProductVersion string `json:"productVersion,omitempty"`
	// The last know location of the device
	Location ssdp.Location `json:"location,omitempty"`
	// The device's UUID
	UUID ssdp.UUID `json:"uuid"`
}

//
// A map holding a set of bookmarks, where the key is alternately the
// UUID of the device, or an alias to a device.
//
type Bookmarks map[string]Bookmark

type configDevice struct {
	product        string
	productVersion string
	location       ssdp.Location
	uuid           ssdp.UUID
}

func (this *configDevice) Product() string {
	return this.product
}

func (this *configDevice) ProductVersion() string {
	return this.productVersion
}

func (this *configDevice) Name() string {
	panic("Not implemented")
}

func (this *configDevice) Location() ssdp.Location {
	return this.location
}

func (this *configDevice) UUID() ssdp.UUID {
	return this.uuid
}

func (this *configDevice) Service(key ssdp.ServiceKey) (service ssdp.Service, has bool) {
	return
}

func (this *configDevice) Services() (keys []ssdp.ServiceKey) {
	return
}

//
// Create a configuration object where @dir is the path to the
// configuration directory.  Note that Init() must be called in order to
// use the newly created object.
//
func MakeConfig(dir string) *Config {
	return &Config{dir, nil, Bookmarks{}}
}

//
// Initialize the configuration object by loading any existing
// configuration from disk.  This method creates the configuration directory,
// if needed.
//
func (this *Config) Init() {
	var err error
	if this.dir, err = os.Open(this.dirname); nil != err {
		if err = os.Mkdir(this.dirname, 0755); nil != err {
			log.Printf("Config: %s", err.(*os.PathError).Error())
			return
		} else if this.dir, err = os.Open(this.dirname); nil != err {
			log.Printf("Config: %s", err.(*os.PathError).Error())
			return
		}
	}
	if fi, err := this.dir.Stat(); nil != err {
		log.Printf("Config: %s", err.(*os.PathError).Error())
		return
	} else if !fi.IsDir() {
		log.Printf("Config: %s: Not a directory", this.dirname)
		return
	}
	this.loadFromDisk()
}

//
// Write the current configuration to disk.
//
func (this *Config) Save() {
	if nil == this.dir {
		return
	}
	this.saveBookmarks()
}

func (this *Config) saveBookmarks() {
	path := path.Join(this.dirname, "bookmarks")
	if fd, err := os.Create(path); nil != err {
		panic(err)
	} else if err := json.NewEncoder(fd).Encode(this.Bookmarks); err != nil {
		if io.EOF != err {
			panic(err)
		}
		fd.Close()
	}
}

//
// Add a bookmark to the bookmark set where @ident is either the uuid
// or the alias to add; @product is the product string, such as 'Sonos';
// @localtion is the network location of the resource; and @uuid is the
// device's UUID.  When adding a device @ident and @uuid should be the same;
// when adding an alias @ident and @uuid will be different.
//
func (this *Config) AddBookmark(ident, product, productVersion string, location ssdp.Location, uuid ssdp.UUID) {
	if ident != string(uuid) {
		this.Bookmarks[ident] = Bookmark{ident, product, productVersion, location, uuid}
	} else {
		this.Bookmarks[ident] = Bookmark{"", product, productVersion, location, uuid}
	}
}

//
// Add @alias as an alias for @uuid.
//
func (this *Config) AddAlias(uuid ssdp.UUID, alias string) {
	old := this.Bookmarks[string(uuid)]
	this.AddBookmark(alias, "", "", ssdp.Location(""), old.UUID)
}

//
// Remove all aliases from the set of bookmarks.
//
func (this *Config) ClearAliases() {
	for key, rec := range this.Bookmarks {
		if 0 < len(rec.Alias) {
			delete(this.Bookmarks, key)
		}
	}
}

//
// Remove any association of a device to the alias @alias.  If @alias
// is a UUID this method is a noop.
//
func (this *Config) ClearAlias(alias string) {
	if rec, has := this.Bookmarks[alias]; has {
		if 0 < len(rec.Alias) {
			delete(this.Bookmarks, alias)
		}
	}
}

func (this *Config) loadFromDisk() {
	for {
		if files, err := this.dir.Readdir(16); nil != err {
			if io.EOF != err {
				panic(err)
			} else {
				break
			}
		} else {
			for _, file := range files {
				this.maybeLoadFile(file)
			}
		}
	}
}

func (this *Config) maybeLoadFile(f os.FileInfo) {
	if "bookmarks" == f.Name() {
		this.maybeLoadBookmarks(f)
	}
}

func (this *Config) maybeLoadBookmarks(f os.FileInfo) {
	if f.IsDir() {
		log.Printf("%s/%s: Expected regular file", this.dirname, f.Name())
		return
	} else {
		path := path.Join(this.dirname, f.Name())
		if fd, err := os.Open(path); nil != err {
			panic(err)
		} else if err := json.NewDecoder(fd).Decode(&this.Bookmarks); err != nil {
			if io.EOF != err {
				panic(err)
			}
			fd.Close()
		}
	}
}

func (this *Config) lookupImpl(ident string, history map[string]bool) (dev ssdp.Device) {
	if _, has := history[ident]; !has {
		history[ident] = true
		if bookmark, has := this.Bookmarks[ident]; has {
			if 0 < len(bookmark.Alias) {
				dev = this.lookupImpl(string(bookmark.UUID), history)
			} else {
				dev = &configDevice{bookmark.Product, bookmark.ProductVersion, bookmark.Location, bookmark.UUID}
			}
		}
	}
	return
}

//
// Try to find the device associated with the UUID or alias
// @ident. Returns nil if there is no device associated with @ident.
//
func (this *Config) Lookup(ident string) ssdp.Device {
	return this.lookupImpl(ident, map[string]bool{})
}
