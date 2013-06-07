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

package sonos

import (
	"github.com/ianr0bkny/go-sonos/ssdp"
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
)

const RECIVA_RADIO = "reciva-com-RecivaRadio"
const RADIO = "Radio"

type Reciva struct {
	upnp.AVTransport
	upnp.ConnectionManager
	upnp.RenderingControl
}

func MakeReciva(svc_map upnp.ServiceMap, reactor upnp.Reactor, flags int) (reciva *Reciva) {
	reciva = &Reciva{}
	for svc_type, svc_list := range svc_map {
		switch svc_type {
		case "AVTransport":
			for _, svc := range svc_list {
				reciva.AVTransport.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &reciva.AVTransport)
				}
				break
			}
		case "ConnectionManager":
			for _, svc := range svc_list {
				reciva.ConnectionManager.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &reciva.ConnectionManager)
				}
				break
			}
		case "RenderingControl":
			for _, svc := range svc_list {
				reciva.RenderingControl.Svc = svc
				svc.Describe()
				if nil != reactor {
					reactor.Subscribe(svc, &reciva.RenderingControl)
				}
				break
			}
		default:
			log.Print(svc_type)
		}
	}
	return
}

func ConnectReciva(mgr ssdp.Manager, reactor upnp.Reactor, flags int) (reciva []*Reciva) {
	qry := ssdp.ServiceQueryTerms{
		ssdp.ServiceKey(RECIVA_RADIO): -1,
	}
	res := mgr.QueryServices(qry)
	if dev_list, has := res[RECIVA_RADIO]; has {
		for _, dev := range dev_list {
			if RADIO == dev.Product() {
				if svc_map, err := upnp.Describe(dev.Location()); nil != err {
					panic(err)
				} else {
					reciva = append(reciva, MakeReciva(svc_map, reactor, flags))
				}
				break
			}
		}
	}
	return
}
