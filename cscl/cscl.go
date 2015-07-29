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
// A client to demonstrate controlling Sonos from the command line.
//
// cscl := (c)ontrol (s)onos from the (c)ommand (l)ine
//
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/config"
	"github.com/ianr0bkny/go-sonos/ssdp"
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

func alias(flags *Args, args []string) (err error) {
	switch len(args) {
	case 0:
		for key, rec := range CONFIG.Bookmarks {
			if 0 < len(rec.Alias) {
				fmt.Printf("%s is an alias for %s\n", key, rec.UUID)
			}
		}
	case 1:
		key := args[0]
		if rec, has := CONFIG.Bookmarks[key]; has {
			if 0 < len(rec.Alias) {
				fmt.Printf("%s is an alias for %s\n", key, rec.UUID)
			} else {
				fmt.Printf("%s is not an alias\n", key)
			}
		} else {
			fmt.Printf("%s is not an alias\n", key)
		}
	case 2:
		CONFIG.AddAlias(ssdp.UUID(args[0]), args[1])
	default:
		err = errors.New("usage: alias [alias | {uuid alias}]")
	}
	return
}

func discover(flags *Args) {
	port := fmt.Sprintf("%d", *flags.discoveryPort)
	if mgr, err := sonos.Discover(*flags.discoveryDevice, port); nil != err {
		panic(err)
	} else {
		query := ssdp.ServiceQueryTerms{
			ssdp.ServiceKey(sonos.MUSIC_SERVICES): -1,
			ssdp.ServiceKey(sonos.RECIVA_RADIO):   -1,
		}
		res := mgr.QueryServices(query)
		if dev_list, has := res[sonos.MUSIC_SERVICES]; has {
			for _, dev := range dev_list {
				if sonos.SONOS == dev.Product() {
					fmt.Printf("%s %s\n", string(dev.UUID()), dev.Location())
					CONFIG.AddBookmark(string(dev.UUID()), dev.Product(), dev.ProductVersion(), dev.Location(), dev.UUID())
				}
			}
		}
		if dev_list, has := res[sonos.RECIVA_RADIO]; has {
			for _, dev := range dev_list {
				if sonos.RADIO == dev.Product() {
					fmt.Printf("%s %s\n", string(dev.UUID()), dev.Location())
					CONFIG.AddBookmark(string(dev.UUID()), dev.Product(), dev.ProductVersion(), dev.Location(), dev.UUID())
				}
			}
		}
	}
}

func devices(flags *Args, args []string) (err error) {
	port := fmt.Sprintf("%d", *flags.discoveryPort)
	if mgr, err := sonos.Discover(*flags.discoveryDevice, port); nil != err {
		panic(err)
	} else {
		dm := mgr.Devices()
		for uuid, dev := range dm {
			fmt.Printf("%s {\n\tProduct = %s\n\tName = %s\n\tLocation = %s\n}\n",
				uuid, dev.Product(), dev.Name(), dev.Location())
		}
	}
	return
}

func queue(flags *Args, args []string) (err error) {
	if 1 != len(args) {
		log.Fatal("usage: cscl queue alias")
	}
	if dev := CONFIG.Lookup(args[0]); nil != dev {
		s := sonos.Connect(dev, nil, sonos.SVC_CONTENT_DIRECTORY)
		if q, err := s.GetQueueContents(); nil != err {
			log.Fatalf("GetQueueContents: %#v", err)
		} else {
			for _, track := range q {
				log.Printf("%s\n", track.Title())
			}
		}
	} else {
		log.Fatal("unknown device")
	}
	return
}

func unalias(flags *Args, args []string) (err error) {
	switch len(args) {
	case 0:
		CONFIG.ClearAliases()
	case 1:
		CONFIG.ClearAlias(args[0])
	default:
		err = errors.New("usage: unalias [alias]")
	}
	return
}

type Args struct {
	alias           *string
	help, usage     *bool
	configDir       *string
	discoveryDevice *string
	discoveryPort   *int
}

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: cscl [-S <uuid/alias>] [-C <configdir=~/.go-sonos/>] [-D <discovery device=eth0>]\n")
	fmt.Fprintf(os.Stderr, "            [-P <discovery port=13104>]\n")
	fmt.Fprintf(os.Stderr, "            [--help|--usage]\n")
	fmt.Fprintf(os.Stderr, "            <command> [args ...]\n\n")
	fmt.Fprintf(os.Stderr, "The available commands are:\n")
	fmt.Fprintf(os.Stderr, "   alias      Add an alias binding\n")
	fmt.Fprintf(os.Stderr, "   devices    Report devices found during discovery\n")
	fmt.Fprintf(os.Stderr, "   discover   Start SSDP device discovery\n")
	fmt.Fprintf(os.Stderr, "   unalias    Remove an alias binding\n")
	fmt.Fprintf(os.Stderr, "\n")
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
		case "alias":
			alias(&args, flag.Args()[1:])
		case "devices":
			devices(&args, flag.Args()[1:])
		case "discover":
			discover(&args)
		case "queue":
			queue(&args, flag.Args()[1:])
		case "unalias":
			unalias(&args, flag.Args()[1:])
		default:
			flag.Usage()
		}
		break
	}

	cleanup()
}
