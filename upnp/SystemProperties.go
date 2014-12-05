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

package upnp

import (
	"encoding/xml"
	_ "log"
)

var (
	SystemProperties_EventType = registerEventType("SystemProperties")
)

type SystemPropertiesState struct {
}

type SystemPropertiesEvent struct {
	SystemPropertiesState
	Svc *Service
}

func (this SystemPropertiesEvent) Service() *Service {
	return this.Svc
}

func (this SystemPropertiesEvent) Type() int {
	return SystemProperties_EventType
}

type SystemProperties struct {
	SystemPropertiesState
	Svc *Service
}

func (this *SystemProperties) BeginSet(svc *Service, channel chan Event) {
}

type systemPropertiesUpdate_XML struct {
	XMLName xml.Name `xml:"SystemPropertiesState"`
	Value   string   `xml:",innerxml"`
}

func (this *SystemProperties) HandleProperty(svc *Service, value string, channel chan Event) error {
	update := systemPropertiesUpdate_XML{
		Value: value,
	}
	if bytes, err := xml.Marshal(update); nil != err {
		return err
	} else {
		xml.Unmarshal(bytes, &this.SystemPropertiesState)
	}
	return nil
}

func (this *SystemProperties) EndSet(svc *Service, channel chan Event) {
	evt := SystemPropertiesEvent{SystemPropertiesState: this.SystemPropertiesState, Svc: svc}
	channel <- evt
}

func (this *SystemProperties) SetString(variableName, stringValue string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"VariableName", variableName},
		{"StringValue", stringValue},
	}
	response := this.Svc.Call("SetString", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) SetStringX(variableName, stringValue string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"VariableName", variableName},
		{"StringValue", stringValue},
	}
	response := this.Svc.Call("SetStringX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) GetString(variableName string) (stringValue string, err error) {
	type Response struct {
		XMLName     xml.Name
		StringValue string
		ErrorResponse
	}
	args := []Arg{
		{"VariableName", variableName},
	}
	response := this.Svc.Call("GetString", args)
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
		ErrorResponse
	}
	args := []Arg{
		{"VariableName", variableName},
	}
	response := this.Svc.Call("GetStringX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	stringValue = doc.StringValue
	err = doc.Error()
	return
}

func (this *SystemProperties) Remove(variableName string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"VariableName", variableName},
	}
	response := this.Svc.Call("Remove", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) GetWebCode(accountType uint32) (webCode string, err error) {
	type Response struct {
		XMLName xml.Name
		WebCode string
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
	}
	response := this.Svc.Call("GetWebCode", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	webCode = doc.WebCode
	err = doc.Error()
	return
}

func (this *SystemProperties) ProvisionTrialAccount(accountType uint32) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
	}
	response := this.Svc.Call("ProvisionTrialAccount", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) ProvisionCredentialedTrialAccountX(accountType uint32, accountId, accountPassword string) (isExpired bool,
	err error) {
	type Response struct {
		XMLName   xml.Name
		IsExpired bool
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountPassword", accountPassword},
	}
	response := this.Svc.Call("ProvisionCredentialedTrialAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	isExpired = doc.IsExpired
	err = doc.Error()
	return
}

func (this *SystemProperties) MigrateTrialAccountX(accountType uint32, accountId, accountPassword string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountPassword", accountPassword},
	}
	response := this.Svc.Call("MigrateTrialAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) AddAccountX(accountType uint32, accountId, accountPassword string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountPassword", accountPassword},
	}
	response := this.Svc.Call("AddAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) AddAccountWithCredentialsX(accountType uint32, accountToken, accountKey string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountToken", accountToken},
		{"AccountKey", accountKey},
	}
	response := this.Svc.Call("AddAccountWithCredentialsX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) RemoveAccount(accountType uint32, accountId string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
	}
	response := this.Svc.Call("RemoveAccount", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) EditAccountPasswordX(accountType uint32, accountId, newAccountPassword string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"NewAccountPassword", newAccountPassword},
	}
	response := this.Svc.Call("EditAccountPasswordX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) EditAccountMd(accountType uint32, accountId, accountMd string) (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountID", accountId},
		{"AccountMD", accountMd},
	}
	response := this.Svc.Call("EditAccountMd", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) DoPostUpdateTasks() (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	response := this.Svc.CallVa("DoPostUpdateTasks")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) ResetThirdPartyCredentials() (err error) {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	response := this.Svc.CallVa("ResetThirdPartyCredentials")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	err = doc.Error()
	return
}

func (this *SystemProperties) RemoveX(variableName string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"VariableName", variableName},
	}
	response := this.Svc.Call("RemoveX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *SystemProperties) EnableRDM(rdmValue bool) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"RDMValue", rdmValue},
	}
	response := this.Svc.Call("EnableRDM", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *SystemProperties) GetRDM() (rdmValue bool, err error) {
	type Response struct {
		XMLName  xml.Name
		RDMValue bool
		ErrorResponse
	}
	response := this.Svc.CallVa("GetRDM")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.RDMValue, doc.Error()
}

func (this *SystemProperties) ApplyRDMDefaultSettings() error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	response := this.Svc.CallVa("ApplyRDMDefaultSettings")
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

func (this *SystemProperties) RefreshAccountCredentialsX(accountType uint32, accountToken, accountKey string) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountToken,", accountToken},
		{"AccountKey,", accountKey},
	}
	response := this.Svc.Call("RefreshAccountCredentialsX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

type A_ARG_TYPE_AccountType uint32
type A_ARG_TYPE_AccountCredential string
type A_ARG_TYPE_OAuthDeviceID string
type A_ARG_TYPE_AccountUDN string

func (this *SystemProperties) AddOAuthAccountX(accountType A_ARG_TYPE_AccountType, accountToken, accountKey A_ARG_TYPE_AccountCredential,
	oauthDeviceID A_ARG_TYPE_OAuthDeviceID) (accountUDN A_ARG_TYPE_AccountUDN, err error) {
	type Response struct {
		XMLName    xml.Name
		AccountUDN A_ARG_TYPE_AccountUDN
		ErrorResponse
	}
	args := []Arg{
		{"AccountType", accountType},
		{"AccountToken,", accountToken},
		{"AccountKey,", accountKey},
		{"OAuthDeviceID,", oauthDeviceID},
	}
	response := this.Svc.Call("AddOAuthAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	accountUDN = doc.AccountUDN
	err = doc.Error()
	return
}

type A_ARG_TYPE_AccountNickname string

func (this *SystemProperties) SetAccountNicknameX(accountUDN A_ARG_TYPE_AccountUDN, accountNickname A_ARG_TYPE_AccountNickname) error {
	type Response struct {
		XMLName xml.Name
		ErrorResponse
	}
	args := []Arg{
		{"AccountUDN,", accountUDN},
		{"AccountNickname,", accountNickname},
	}
	response := this.Svc.Call("SetAccountNicknameX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	return doc.Error()
}

type A_ARG_TYPE_AccountID string
type A_ARG_TYPE_AccountPassword string

func (this *SystemProperties) ReplaceAccountX(accountUDN A_ARG_TYPE_AccountUDN, newAccountID A_ARG_TYPE_AccountID,
	newAccountPassword A_ARG_TYPE_AccountPassword) (newAccountUDN A_ARG_TYPE_AccountUDN, err error) {
	type Response struct {
		XMLName       xml.Name
		NewAccountUDN A_ARG_TYPE_AccountUDN
		ErrorResponse
	}
	args := []Arg{
		{"AccountUDN", accountUDN},
		{"NewAccountID,", newAccountID},
		{"NewAccountPassword,", newAccountPassword},
	}
	response := this.Svc.Call("ReplaceAccountX", args)
	doc := Response{}
	xml.Unmarshal([]byte(response), &doc)
	newAccountUDN = doc.NewAccountUDN
	err = doc.Error()
	return
}
