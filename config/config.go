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

package config

import (
	"encoding/json"
	"github.com/ianr0bkny/go-sonos/ssdp"
	"io"
	"log"
	"os"
	"path"
)

type Config struct {
	dirname   string
	dir       *os.File
	Bookmarks Bookmarks
}

type Bookmark struct {
	Alias    string        `json:"alias,omitempty"`
	Location ssdp.Location `json:"location,omitempty"`
	UUID     ssdp.UUID     `json:"uuid"`
}

type Bookmarks map[string]Bookmark

func MakeConfig(dir string) *Config {
	return &Config{dir, nil, Bookmarks{}}
}

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

func (this *Config) AddBookmark(alias string, location ssdp.Location, uuid ssdp.UUID) {
	if alias != string(uuid) {
		this.Bookmarks[alias] = Bookmark{alias, location, uuid}
	} else {
		this.Bookmarks[alias] = Bookmark{"", location, uuid}
	}
}

func (this *Config) AddAlias(uuid ssdp.UUID, alias string) {
	old := this.Bookmarks[string(uuid)]
	this.AddBookmark(alias, ssdp.Location(""), old.UUID)
}

func (this *Config) ClearAliases() {
	for key, rec := range this.Bookmarks {
		if 0 < len(rec.Alias) {
			delete(this.Bookmarks, key)
		}
	}
}

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
