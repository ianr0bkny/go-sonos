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
// A go-language implementation of the Sonos UPnP API.
//
package sonos

import (
	"github.com/ianr0bkny/go-sonos/ssdp"
	"github.com/ianr0bkny/go-sonos/upnp"
	_ "log"
)

const MUSIC_SERVICES = "schemas-upnp-org-MusicServices"
const SONOS = "Sonos"

type Sonos struct {
	upnp.AlarmClock
	upnp.AVTransport
	upnp.ConnectionManager
	upnp.ContentDirectory
	upnp.DeviceProperties
	upnp.GroupManagement
	upnp.MusicServices
	upnp.RenderingControl
	upnp.SystemProperties
	upnp.ZoneGroupTopology
}

const (
	SVC_ALARM_CLOCK         = 1
	SVC_AV_TRANSPORT        = SVC_ALARM_CLOCK << 1
	SVC_CONNECTION_MANAGER  = SVC_AV_TRANSPORT << 1
	SVC_CONTENT_DIRECTORY   = SVC_CONNECTION_MANAGER << 1
	SVC_DEVICE_PROPERTIES   = SVC_CONTENT_DIRECTORY << 1
	SVC_GROUP_MANAGEMENT    = SVC_DEVICE_PROPERTIES << 1
	SVC_MUSIC_SERVICES      = SVC_GROUP_MANAGEMENT << 1
	SVC_RENDERING_CONTROL   = SVC_MUSIC_SERVICES << 1
	SVC_SYSTEM_PROPERTIES   = SVC_RENDERING_CONTROL << 1
	SVC_ZONE_GROUP_TOPOLOGY = SVC_SYSTEM_PROPERTIES << 1
	//
	SVC_ALL = SVC_ALARM_CLOCK |
		SVC_AV_TRANSPORT |
		SVC_CONNECTION_MANAGER |
		SVC_CONTENT_DIRECTORY |
		SVC_DEVICE_PROPERTIES |
		SVC_GROUP_MANAGEMENT |
		SVC_MUSIC_SERVICES |
		SVC_RENDERING_CONTROL |
		SVC_SYSTEM_PROPERTIES |
		SVC_ZONE_GROUP_TOPOLOGY
)

func sonosCheckServiceFlags(svc_type string, flags int) bool {
	switch svc_type {
	case "AlarmClock":
		return flags&SVC_ALARM_CLOCK > 0
	case "AVTransport":
		return flags&SVC_AV_TRANSPORT > 0
	case "ConnectionManager":
		return flags&SVC_CONNECTION_MANAGER > 0
	case "ContentDirectory":
		return flags&SVC_CONTENT_DIRECTORY > 0
	case "DeviceProperties":
		return flags&SVC_DEVICE_PROPERTIES > 0
	case "GroupManagement":
		return flags&SVC_GROUP_MANAGEMENT > 0
	case "MusicServices":
		return flags&SVC_MUSIC_SERVICES > 0
	case "RenderingControl":
		return flags&SVC_RENDERING_CONTROL > 0
	case "SystemProperties":
		return flags&SVC_SYSTEM_PROPERTIES > 0
	case "ZoneGroupTopology":
		return flags&SVC_ZONE_GROUP_TOPOLOGY > 0
	}
	return false
}

func MakeSonos(svc_map upnp.ServiceMap, reactor upnp.Reactor, flags int) (sonos *Sonos) {
	sonos = &Sonos{}
	for svc_type, svc_list := range svc_map {
		if !sonosCheckServiceFlags(svc_type, flags) {
			continue
		}
		switch svc_type {
		case "AlarmClock":
			for _, svc := range svc_list {
				sonos.AlarmClock.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.AlarmClock)
				}
				break
			}
		case "AVTransport":
			for _, svc := range svc_list {
				sonos.AVTransport.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.AVTransport)
				}
				break
			}
		case "ConnectionManager":
			for _, svc := range svc_list {
				sonos.ConnectionManager.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.ConnectionManager)
				}
				break
			}
		case "ContentDirectory":
			for _, svc := range svc_list {
				sonos.ContentDirectory.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.ContentDirectory)
				}
				break
			}
		case "DeviceProperties":
			for _, svc := range svc_list {
				sonos.DeviceProperties.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.DeviceProperties)
				}
				break
			}
		case "GroupManagement":
			for _, svc := range svc_list {
				sonos.GroupManagement.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.GroupManagement)
				}
				break
			}
		case "MusicServices":
			for _, svc := range svc_list {
				sonos.MusicServices.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.MusicServices)
				}
				break
			}
		case "RenderingControl":
			for _, svc := range svc_list {
				sonos.RenderingControl.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.RenderingControl)
				}
				break
			}
		case "SystemProperties":
			for _, svc := range svc_list {
				sonos.SystemProperties.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.SystemProperties)
				}
				break
			}
		case "ZoneGroupTopology":
			for _, svc := range svc_list {
				sonos.ZoneGroupTopology.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &sonos.ZoneGroupTopology)
				}
				break
			}
		}
	}
	return
}

func ConnectAny(mgr ssdp.Manager, reactor upnp.Reactor, flags int) (sonos []*Sonos) {
	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey(MUSIC_SERVICES): -1,
	}
	res := mgr.QueryServices(qry)
	if dev_list, has := res[MUSIC_SERVICES]; has {
		for _, dev := range dev_list {
			if SONOS == dev.Product() {
				if svc_map, err := upnp.Describe(dev.Location()); nil != err {
					panic(err)
				} else {
					sonos = append(sonos, MakeSonos(svc_map, reactor, flags))
				}
				break
			}
		}
	}
	return
}

func Connect(dev ssdp.Device, reactor upnp.Reactor, flags int) (sonos *Sonos) {
	if svc_map, err := upnp.Describe(dev.Location()); nil != err {
		panic(err)
	} else {
		sonos = MakeSonos(svc_map, reactor, flags)
	}
	return
}

func MakeReactor(ifiname, port string) upnp.Reactor {
	reactor := upnp.MakeReactor()
	reactor.Init(ifiname, port)
	return reactor
}

func Discover(ifiname, port string) (mgr ssdp.Manager, err error) {
	mgr = ssdp.MakeManager()
	mgr.Discover(ifiname, port, false)
	return
}
