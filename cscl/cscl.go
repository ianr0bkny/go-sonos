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

// cscl := (c)ontrol (s)onos from the (c)ommand (l)ine

package main

import (
	"flag"
	"fmt"
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	"log"
	"os"
	"path"
)

var CONFIG *config.Config

func initConfig(dir string) {
	if "" == dir {
		dir = path.Join(os.Getenv("HOME"), ".go-sonos")
	}
	CONFIG = config.MakeConfig(dir)
	CONFIG.Init()
}

func cleanup() {
	CONFIG.Save()
}

func discover(args *Args) {
	if found, err := sonos.Discover(*args.discoveryDevice, fmt.Sprintf("%d", *args.discoveryPort)); nil != err {
		panic(err)
	} else {
		for _, s := range found {
			log.Printf("%#v", s)
		}
	}
}

type Args struct {
	alias           *string
	help, usage     *bool
	configDir       *string
	discoveryDevice *string
	discoveryPort   *int
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s [args] command:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	args := Args{}
	args.alias = flag.String("S", "", "device alias name")
	args.configDir = flag.String("C", "", "go-sonos configuration directory")
	args.discoveryDevice = flag.String("D", "eth0", "discovery device")
	args.discoveryPort = flag.Int("P", 13104, "discovery response port")
	args.help = flag.Bool("help", false, "show the usage message")
	args.usage = flag.Bool("usage", false, "show the usage message")
	flag.Usage = Usage
	flag.Parse()

	if 0 == len(flag.Args()) {
		flag.Usage()
		os.Exit(0)
	}

	initConfig(*args.configDir)

	for _, cmd := range flag.Args() {
		switch cmd {
		case "discover":
			discover(&args)
		default:
			flag.Usage()
		}
		break
	}

	cleanup()
}
