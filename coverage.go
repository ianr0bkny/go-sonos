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
	"fmt"
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
	"reflect"
)

type coverageData struct {
	total       int
	missing     []string
	implemented int
}

func (this *coverageData) pct() float32 {
	return 100. * float32(this.implemented) / float32(this.total)
}

func (this *coverageData) add(other *coverageData) {
	this.total += other.total
	this.missing = append(this.missing, other.missing...)
	this.implemented += other.implemented
}

func (this *coverageData) log(name string, missing bool) {
	log.Printf("%20s %8.2f%% %3d/%-3d", name, this.pct(), this.implemented, this.total)
	if missing {
		for _, action := range this.missing {
			log.Printf("%20s * %s", "", action)
		}
	}
}

func Coverage(s interface{}) {
	sv := reflect.Indirect(reflect.ValueOf(s))
	st := sv.Type()
	total_cd := coverageData{}
	for i := 0; i < st.NumField(); i++ {
		superclass := sv.Field(i)
		svc := superclass.FieldByName("Svc").Interface().(*upnp.Service)
		if nil == svc {
			cd := coverageData{}
			cd.log(fmt.Sprintf("Service %s not implemented", superclass.Type().Name()), true)
			total_cd.add(&cd)
			continue
		}
		actions := svc.Actions()
		superclass_type := reflect.PtrTo(superclass.Type())
		cd := coverageData{total: len(actions)}
		for _, action := range actions {
			if _, has := superclass_type.MethodByName(action); has {
				cd.implemented++
			} else {
				cd.missing = append(cd.missing, action)
			}
		}
		cd.log(superclass.Type().Name(), true)
		total_cd.add(&cd)
	}
	total_cd.log("TOTAL", false)
}
