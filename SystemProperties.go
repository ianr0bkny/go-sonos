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
	"encoding/xml"
	"github.com/ianr0bkny/go-sonos/upnp"
	_ "log"
)

type SystemProperties struct {
	Svc *upnp.Service
}

func (this *SystemProperties) SetString(variableName, stringValue string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"VariableName", variableName},
		{"StringValue", stringValue},
	}
	response := upnp.Call(this.Svc, "SetString", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) SetStringX(variableName, stringValue string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"VariableName", variableName},
		{"StringValue", stringValue},
	}
	response := upnp.Call(this.Svc, "SetStringX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) GetString(variableName string) (stringValue string, err error) {
	type Response struct {
		XMLName     xml.Name
		StringValue string
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"VariableName", variableName},
	}
	response := upnp.Call(this.Svc, "GetString", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	stringValue = doc.StringValue
	err = doc.Error()
	return
}

func (this *SystemProperties) GetStringX(variableName string) (stringValue string, err error) {
	type Response struct {
		XMLName     xml.Name
		StringValue string
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"VariableName", variableName},
	}
	response := upnp.Call(this.Svc, "GetStringX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	stringValue = doc.StringValue
	err = doc.Error()
	return
}

func (this *SystemProperties) Remove(variableName string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"VariableName", variableName},
	}
	response := upnp.Call(this.Svc, "Remove", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) GetWebCode(accountType uint32) (webCode string, err error) {
	type Response struct {
		XMLName xml.Name
		WebCode string
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
	}
	response := upnp.Call(this.Svc, "GetWebCode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	webCode = doc.WebCode
	err = doc.Error()
	return
}

func (this *SystemProperties) ProvisionTrialAccount(accountType uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
	}
	response := upnp.Call(this.Svc, "ProvisionTrialAccount", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) ProvisionCredentialedTrialAccountX(accountType uint32, accountId, accountPassword string) (isExpired bool,
	err error) {
	type Response struct {
		XMLName xml.Name
		IsExpired bool
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountPassword", accountPassword},
	}
	response := upnp.Call(this.Svc, "ProvisionCredentialedTrialAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	isExpired = doc.IsExpired
	err = doc.Error()
	return
}

func (this *SystemProperties) MigrateTrialAccountX(accountType uint32, accountId, accountPassword string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountPassword", accountPassword},
	}
	response := upnp.Call(this.Svc, "MigrateTrialAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) AddAccountX(accountType uint32, accountId, accountPassword string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountPassword", accountPassword},
	}
	response := upnp.Call(this.Svc, "AddAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) AddAccountWithCredentialsX(accountType uint32, accountToken, accountKey string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountToken", accountToken},
		{"AccountKey", accountKey},
	}
	response := upnp.Call(this.Svc, "AddAccountWithCredentialsX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) RemoveAccount(accountType uint32, accountId string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
	}
	response := upnp.Call(this.Svc, "RemoveAccount", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) EditAccountPasswordX(accountType uint32, accountId, newAccountPassword string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"NewAccountPassword", newAccountPassword},
	}
	response := upnp.Call(this.Svc, "EditAccountPasswordX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) EditAccountMd(accountType uint32, accountId, accountMd string) (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	args := []upnp.Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountMD", accountMd},
	}
	response := upnp.Call(this.Svc, "EditAccountMd", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) DoPostUpdateTasks() (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "DoPostUpdateTasks")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) ResetThirdPartyCredentials() (err error) {
	type Response struct {
		XMLName xml.Name
		upnp.ErrorResponse
	}
	response := upnp.CallVa(this.Svc, "ResetThirdPartyCredentials")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}
