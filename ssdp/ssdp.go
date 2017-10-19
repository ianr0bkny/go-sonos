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
// A client implementation of the SSDP protocol.
//
//	mgr := ssdp.MakeManager()
//	mgr.Discover("eth0", "13104", false)
//	qry := ssdp.ServiceQueryTerms{
//		ssdp.ServiceKey("schemas-upnp-org-MusicServices"): -1,
//	}
//	result := mgr.QueryServices(qry)
//	if device_list, has := result["schemas-upnp-org-MusicServices"]; has {
//		for _, device := range device_list {
//			...
//		}
//	}
//	mgr.Close()
//
package ssdp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ssdpHNAPRegex *regexp.Regexp
var ssdpOtherDeviceRegex *regexp.Regexp
var ssdpOtherDeviceUUIDRegex *regexp.Regexp
var ssdpOtherServiceRegex *regexp.Regexp
var ssdpOtherServiceUUIDRegex *regexp.Regexp
var ssdpResponseStartLineRegexp *regexp.Regexp
var ssdpServerStringRegexp *regexp.Regexp
var ssdpRootDeviceUUIDRegex *regexp.Regexp
var ssdpUPnPBareUUIDRegex *regexp.Regexp
var ssdpUPnPDeviceRegex *regexp.Regexp
var ssdpUPnPDeviceUUIDRegex *regexp.Regexp
var ssdpUPnPServiceRegex *regexp.Regexp
var ssdpUPnPServiceUUIDRegex *regexp.Regexp
var ssdpUPnPUIDRegex *regexp.Regexp

func init() {
	ssdpHNAPRegex = regexp.MustCompile("^hnap:(.+)$")
	ssdpOtherDeviceRegex = regexp.MustCompile("^urn:([^:]+):device:([^:]+)(:(.+))?$")
	ssdpOtherDeviceUUIDRegex = regexp.MustCompile("^uuid:([^:]+)::urn:([^:]+):device:([^:]+)(:(.+))?$")
	ssdpOtherServiceRegex = regexp.MustCompile("^urn:([^:]+):service:([^:]+)(:(.+))?$")
	ssdpOtherServiceUUIDRegex = regexp.MustCompile("^uuid:([^:]+)::urn:([^:]+):service:([^:]+)(:(.+))?$")
	ssdpResponseStartLineRegexp = regexp.MustCompile("^HTTP/([0-9\\.]+)$")
	ssdpServerStringRegexp = regexp.MustCompile("^([^ /,]+)(/([^ ,]+))?,?\\s+[Uu][Pp][Nn][Pp](/([^ ,]+))?,?\\s+([^/]+)(/(.+))?$")
	ssdpRootDeviceUUIDRegex = regexp.MustCompile("^uuid:([^:]+)::upnp:rootdevice$")
	ssdpUPnPBareUUIDRegex = regexp.MustCompile("^uuid:([^:]+)$")
	ssdpUPnPDeviceRegex = regexp.MustCompile("^urn:schemas-upnp-org:device:([^:]+)(:(.+))?$")
	ssdpUPnPDeviceUUIDRegex = regexp.MustCompile("uuid:([^:]+)::urn:schemas-upnp-org:device:([^:]+)(:(.+))?$")
	ssdpUPnPServiceRegex = regexp.MustCompile("^urn:schemas-upnp-org:service:([^:]+)(:(.+))?$")
	ssdpUPnPServiceUUIDRegex = regexp.MustCompile("uuid:([^:]+)::urn:schemas-upnp-org:service:([^:]+)(:(.+))?$")
	ssdpUPnPUIDRegex = regexp.MustCompile("^uuid:(.+)$")
}

const (
	ssdpBroadcastGroup   = "239.255.255.250:1900"
	ssdpBroadcastVersion = "udp"
)

// Type protection for a device URI
type Location string

// Type protection for a SSDP service key
type ServiceKey string

// Type protection for a Universally Unique ID
type UUID string

type ssdpServerDescription struct {
	os             string
	os_version     string
	upnp_version   string
	product        string
	productVersion string
}

type ssdpResourceBase struct {
	ssdpServerDescription
	location Location
	uuid     UUID
}

// An abstraction of an SSDP service
type Service interface {
	// The version of the service
	Version() int64
}

type ssdpService struct {
	ssdpResourceBase
	name    string
	version int64
	uri     string
}

func (this *ssdpService) Version() int64 {
	return this.version
}

type ssdpServiceSet map[ServiceKey]Service

// An abstraction of an SSDP leaf device
type Device interface {
	// The product name (e.g. "Sonos")
	Product() string
	// The product version (e.g. "28.1-83040 (BR100)")
	ProductVersion() string
	// The device name (e.h. "ZonePlayer")
	Name() string
	// A URI that can be queried for device capabilities
	Location() Location
	// The device's Universally Unique ID
	UUID() UUID
	// Search for a service in the device's capabilities; same
	// return semantics as querying a map
	Service(key ServiceKey) (service Service, has bool)
	// Return a list of services implemented by this device
	Services() []ServiceKey
}

type ssdpDevice struct {
	ssdpResourceBase
	name     string
	version  int64
	uri      string
	services ssdpServiceSet
}

func (this *ssdpDevice) Product() string {
	return this.product
}

func (this *ssdpDevice) ProductVersion() string {
	return this.productVersion
}

func (this *ssdpDevice) Name() string {
	return this.name
}

func (this *ssdpDevice) Location() Location {
	return this.location
}

func (this *ssdpDevice) UUID() UUID {
	return this.uuid
}

func (this *ssdpDevice) Service(key ServiceKey) (service Service, has bool) {
	service, has = this.services[key]
	return
}

func (this *ssdpDevice) Services() []ServiceKey {
	i := 0
	keys := make([]ServiceKey, len(this.services))
	for key, _ := range this.services {
		keys[i] = key
		i += 1
	}
	return keys
}

// Indexes leaf devices by UUID
type DeviceMap map[UUID]Device

type ssdpRootDevice struct {
	ssdpResourceBase
	Devices DeviceMap
}

type ssdpRootDeviceMap map[Location]*ssdpRootDevice

// Discovered devices sorted by supported service
type ServiceMap map[ServiceKey]DeviceMap

type ssdpResourceType int

const (
	ssdpTypessdpRootDevice ssdpResourceType = iota
	ssdpTypeService
	ssdpTypeDevice
	ssdpTypeHNAP
	ssdpTypeUUID
	ssdpTypeUnknown
)

type ssdpResource struct {
	ssdpServerDescription
	ssdptype ssdpResourceType
	uuid     UUID
	name     string
	version  int64
	uri      string
	location Location
	children []*ssdpResource
}

type ssdpMessageType int

const (
	ssdpSearch ssdpMessageType = iota
	ssdpNotify
	ssdpResponse
	ssdpInvalid
)

type ssdpRawMessage struct {
	msgtype ssdpMessageType
	httpver string
	header  map[string]string
}

type ssdpNotifyMessage struct {
	_01_nls            string
	cache_control      string
	host               string
	location           Location
	nts                string
	nt                 string
	opt                string
	server             string
	usn                string
	x_rincon_bootseq   string
	x_rincon_household string
	x_user_agent       string
}

type ssdpNotifyQueue chan *ssdpNotifyMessage

type ssdpResponseMessage struct {
	ssdpServerDescription
	_01_nls            string
	al                 string
	cache_control      string
	date               string
	ext                string
	location           Location
	opt                string
	st                 string
	usn                string
	x_rincon_bootseq   string
	x_rincon_household string
	x_rincon_variant   string
	x_rincon_wifimode  string
	x_user_agent       string
	content_length     string
	bootid_upnp_org    string
	configid_upnp_org  string
	household_smartspeaker_audio string
	x_av_server_info   string
}

type ssdpResponseQueue chan *ssdpResponseMessage

type ssdpConnection struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

// A map of service key to minimum required version
type ServiceQueryTerms map[ServiceKey]int64

// Encapsulates SSDP discovery, handles updates, and stores results
type Manager interface {
	// Initiates SSDP discovery, where ifiname names a network device
	// to query, port gives a free port on that network device to listen
	// for responses, and the subscribe flag (currrently unimplemented)
	// determines whether to listen to asynchronous updates after the
	// initial query is complete.
	Discover(ifiname, port string, subscribe bool) error
	// After discovery is complete searches for devices implementing
	// the services specified in query.
	QueryServices(query ServiceQueryTerms) ServiceMap
	// Return the list of devices that were found during discovery
	Devices() DeviceMap
	// Shuts down asynchronous subscriptions to device state
	Close() error
}

type ssdpDefaultManager struct {
	responseQueue ssdpResponseQueue
	notifyQueue   ssdpNotifyQueue
	unicast       ssdpConnection
	multicast     ssdpConnection
	rootDeviceMap ssdpRootDeviceMap
	deviceMap     DeviceMap
	serviceMap    ServiceMap
	readyChan     chan int
	closeChan     chan int
}

// Returns an empty manager ready for SSDP discovery
func MakeManager() Manager {
	mgr := &ssdpDefaultManager{}
	mgr.responseQueue = make(ssdpResponseQueue)
	mgr.notifyQueue = make(ssdpNotifyQueue)
	mgr.unicast = ssdpConnection{}
	mgr.multicast = ssdpConnection{}
	mgr.rootDeviceMap = make(ssdpRootDeviceMap)
	mgr.deviceMap = make(DeviceMap)
	mgr.serviceMap = make(ServiceMap)
	mgr.readyChan = make(chan int)
	mgr.closeChan = make(chan int)
	return mgr
}

func (this *ssdpDefaultManager) Discover(ifiname, port string, subscribe bool) (err error) {
	this.ssdpDiscoverImpl(ifiname, port, subscribe)
	return
}

func (this *ssdpDefaultManager) QueryServices(query ServiceQueryTerms) (results ServiceMap) {
	results = make(ServiceMap)
	for name, minver := range query {
		results[name] = make(DeviceMap)
		if dlist, has := this.serviceMap[name]; has {
			for uuid, _ := range dlist {
				de := this.deviceMap[uuid]
				if svc, has := de.Service(name); has {
					if minver <= svc.Version() {
						results[name][de.UUID()] = de
					}
				}
			}
		}
	}
	return
}

func (this *ssdpDefaultManager) Devices() DeviceMap {
	return this.deviceMap
}

func (this *ssdpDefaultManager) Close() (err error) {
	if nil != this.unicast.conn {
		this.unicast.conn.Close()
		<-this.closeChan
	}
	if nil != this.multicast.conn {
		this.multicast.conn.Close()
		<-this.closeChan
	}
	return
}

func ssdpNewDevice(res *ssdpResource) (de *ssdpDevice) {
	de = new(ssdpDevice)
	de.ssdpServerDescription = res.ssdpServerDescription
	de.uuid = res.uuid
	de.location = res.location
	de.name = res.name
	de.version = res.version
	de.uri = res.uri
	de.services = make(ssdpServiceSet)
	return
}

func ssdpNewssdpRootDevice(res *ssdpResource) (rd *ssdpRootDevice) {
	rd = new(ssdpRootDevice)
	rd.ssdpServerDescription = res.ssdpServerDescription
	rd.uuid = res.uuid
	rd.location = res.location
	rd.Devices = make(DeviceMap)
	return
}

func (rd *ssdpRootDevice) ssdpRequireEmbeddedDevice(de *ssdpDevice) {
	if _, has := rd.Devices[de.UUID()]; !has {
		rd.Devices[de.UUID()] = de
	}
}

func ssdpNewRawMessage() *ssdpRawMessage {
	return &ssdpRawMessage{ssdpInvalid, "1.1", make(map[string]string)}
}

func ssdpNewResource(uuid UUID) (res *ssdpResource) {
	res = new(ssdpResource)
	res.ssdptype = ssdpTypeUnknown
	res.uuid = uuid
	return
}

func ssdpNewService(res *ssdpResource) (sv *ssdpService) {
	sv = new(ssdpService)
	sv.ssdpServerDescription = res.ssdpServerDescription
	sv.uuid = res.uuid
	sv.location = res.location
	sv.name = res.name
	sv.version = res.version
	sv.uri = res.uri
	return
}

func (this *ssdpDefaultManager) ssdpParseStartLineFields(raw *ssdpRawMessage, fields []string) {
	if "M-SEARCH" == fields[0] {
		raw.msgtype = ssdpSearch
	} else if "NOTIFY" == fields[0] {
		raw.msgtype = ssdpNotify
	} else if m := ssdpResponseStartLineRegexp.FindStringSubmatch(fields[0]); 0 < len(m) {
		raw.msgtype = ssdpResponse
		raw.httpver = m[1]
	} else {
		panic(fmt.Sprintf("Invalid start line `%s'", fields[0]))
	}
}

func (this *ssdpDefaultManager) ssdpParseStartLine(raw *ssdpRawMessage, line []byte) {
	fields := strings.Fields(string(line))
	if 3 != len(fields) {
		panic("Invalid start line")
	} else {
		this.ssdpParseStartLineFields(raw, fields)
	}
}

func (this *ssdpDefaultManager) ssdpParseHeaderLine(raw *ssdpRawMessage, line []byte) {
	i := strings.Index(string(line), ":")
	if -1 == i {
		panic("Invalid header")
	}
	field := textproto.CanonicalMIMEHeaderKey(strings.TrimSpace(string(line[0:i])))
	value := strings.TrimSpace(string(line[i+1:]))
	if _, has := raw.header[field]; has {
		panic("Header field redefined")
	} else {
		raw.header[field] = value
	}
}

func (this *ssdpDefaultManager) ssdpParseInputLine(raw *ssdpRawMessage, line []byte, lineno int) {
	if 1 < lineno {
		this.ssdpParseHeaderLine(raw, line)
	} else {
		this.ssdpParseStartLine(raw, line)
	}
}

func (this *ssdpDefaultManager) ssdpParseInput(msg []byte) (raw *ssdpRawMessage) {
	raw = ssdpNewRawMessage()
	bin := bufio.NewReader(bytes.NewReader(msg))
	var line []byte
	lineno := 0
	for {
		if fragment, is_prefix, err := bin.ReadLine(); nil == err {
			line = append(line, fragment...)
			if !is_prefix {
				lineno += 1
				if 1 < lineno && 0 == len(line) {
					break
				}
				this.ssdpParseInputLine(raw, line, lineno)
				line = nil
			}
		} else if io.EOF == err {
			panic("Premature end of header")
		} else {
			panic(err)
		}
	}
	return
}

func (this *ssdpDefaultManager) ssdpHandleNotify(raw *ssdpRawMessage) *ssdpNotifyMessage {
	msg := new(ssdpNotifyMessage) /*asynchronous reply from multicast*/
	for key, value := range raw.header {
		switch key {
		case "Location":
			msg.location = Location(value)
		case "Server":
			msg.server = value
		case "Host":
			msg.host = value
		case "Usn":
			msg.usn = value
		case "Cache-Control":
			msg.cache_control = value
		case "X-User-Agent":
			msg.x_user_agent = value
		case "X-Rincon-Bootseq":
			msg.x_rincon_bootseq = value
		case "X-Rincon-Household":
			msg.x_rincon_household = value
		case "Nts":
			msg.nts = value
		case "Nt":
			msg.nt = value
		case "Opt":
			msg.opt = value
		case "01-Nls":
			msg._01_nls = value
		default:
			log.Printf("No support for field `%s' (value `%s')", key, value)
		}
	}
	return msg
}

func (this *ssdpDefaultManager) ssdpHandleResponse(raw *ssdpRawMessage) *ssdpResponseMessage {
	msg := new(ssdpResponseMessage) /*synchronous reply from unicast*/
	for key, value := range raw.header {
		switch key {
		case "Location":
			msg.location = Location(value)
		case "St":
			msg.st = value
		case "Server":
			m := ssdpServerStringRegexp.FindStringSubmatch(value)
			if 0 < len(m) {
				msg.os = m[1]
				msg.os_version = m[3]
				msg.upnp_version = m[5]
				msg.product = m[6]
				msg.productVersion = m[8]
			} else {
				log.Printf("Invalid server description `%s'", value)
			}
		case "Opt":
			msg.opt = value
		case "Usn":
			msg.usn = value
		case "Ext":
			msg.ext = value
		case "Date":
			msg.date = value
		case "Cache-Control":
			msg.cache_control = value
		case "X-User-Agent":
			msg.x_user_agent = value
		case "X-Rincon-Bootseq":
			msg.x_rincon_bootseq = value
		case "X-Rincon-Household":
			msg.x_rincon_household = value
		case "X-Rincon-Variant":
			msg.x_rincon_variant = value
		case "X-Rincon-Wifimode":
			msg.x_rincon_wifimode = value
		case "Al":
			msg.al = value
		case "01-Nls":
			msg._01_nls = value
		case "Content-Length":
			msg.content_length = value
		case "Bootid.upnp.org":
			msg.bootid_upnp_org = value
		case "Configid.upnp.org":
			msg.configid_upnp_org = value
		case "Household.smartspeaker.audio":
			msg.household_smartspeaker_audio = value
		case "X-Av-Server-Info":
			msg.x_av_server_info = value
		default:
			log.Printf("No support for field `%s' (value `%s')", key, value)
			log.Printf("%v", raw)
		}
	}
	return msg
}

func (this *ssdpDefaultManager) ssdpHandleMessage(raw *ssdpRawMessage) {
	switch raw.msgtype {
	case ssdpSearch: /*ignore*/
	case ssdpResponse:
		this.responseQueue <- this.ssdpHandleResponse(raw)
	case ssdpNotify:
		this.notifyQueue <- this.ssdpHandleNotify(raw)
	}
}

func (this *ssdpDefaultManager) ssdpDiscoverLoop(conn net.Conn) {
	this.readyChan <- 1
	msg := make([]byte, 65536) /*max size of a single UDP packet*/
	defer func() {
		recover()
		this.closeChan <- 1
	}()
	for {
		if n, err := conn.Read(msg); nil != err {
			panic(err)
		} else if raw := this.ssdpParseInput(msg[:n]); nil != raw {
			this.ssdpHandleMessage(raw)
		}
	}
}

func (this *ssdpDefaultManager) ssdpUnicastDiscoverImpl(ifi *net.Interface, port string) (err error) {
	addrs, err := ifi.Addrs()
	if nil != err {
		return
	} else if 0 == len(addrs) {
		err = errors.New(fmt.Sprintf("No addresses found for interface %s", ifi.Name))
		return
	}
	var lip net.IP
	for _, addr := range addrs {
		if nil != addr.(*net.IPNet).IP.DefaultMask() {
			lip = addr.(*net.IPNet).IP
			break;
		}
	}
	laddr, err := net.ResolveUDPAddr(ssdpBroadcastVersion, net.JoinHostPort(lip.String(), port))
	if nil != err {
		return
	}
	uc, err := net.ListenUDP(ssdpBroadcastVersion, laddr)
	if nil != err {
		return
	}
	this.unicast.addr = laddr
	this.unicast.conn = uc
	go this.ssdpDiscoverLoop(uc)
	<-this.readyChan
	return
}

func (this *ssdpDefaultManager) ssdpMulticastDiscoverImpl(ifi *net.Interface, subscribe bool) (err error) {
	maddr, err := net.ResolveUDPAddr(ssdpBroadcastVersion, ssdpBroadcastGroup)
	if nil != err {
		return
	}
	this.multicast.addr = maddr
	var mc *net.UDPConn
	if subscribe {
		mc, err = net.ListenMulticastUDP(ssdpBroadcastVersion, ifi, maddr)
		if nil != err {
			return
		}
		this.multicast.conn = mc
		go this.ssdpDiscoverLoop(mc)
		<-this.readyChan
	}
	return
}

func (this *ssdpDefaultManager) ssdpQueryMessage(timeout int) (msg *bytes.Buffer) {
	msg = new(bytes.Buffer)
	msg.WriteString("M-SEARCH * HTTP/1.1\r\n")
	msg.WriteString("HosT: 239.255.255.250:1900\r\n")
	msg.WriteString("MAN: \"ssdp:discover\"\r\n")
	msg.WriteString(fmt.Sprintf("MX: %d\r\n", timeout))
	msg.WriteString("ST: ssdp:all\r\n")
	msg.WriteString("USER-AGENT: unix/5.1 UPnP/1.1 crash/1.0\r\n")
	msg.WriteString("\r\n")
	return
}

func (this *ssdpDefaultManager) ssdpRequiressdpRootDevice(res *ssdpResource) (rd *ssdpRootDevice) {
	var has bool
	if rd, has = this.rootDeviceMap[res.location]; !has {
		rd = ssdpNewssdpRootDevice(res)
		this.rootDeviceMap[rd.location] = rd
	}
	return
}

func (this *ssdpDefaultManager) ssdpRequireDevice(res *ssdpResource) (de *ssdpDevice) {
	if raw, has := this.deviceMap[res.uuid]; !has {
		de = ssdpNewDevice(res)
		this.deviceMap[res.uuid] = de
	} else {
		de = raw.(*ssdpDevice)
	}
	return
}

func (this *ssdpDefaultManager) ssdpNotifyssdpRootDevice(res *ssdpResource) {
	this.ssdpRequiressdpRootDevice(res)
	this.ssdpRequireDevice(res)
}

func (this *ssdpDefaultManager) ssdpNotifyDevice(res *ssdpResource) {
	rd := this.ssdpRequiressdpRootDevice(res)
	de := this.ssdpRequireDevice(res)
	rd.ssdpRequireEmbeddedDevice(de)
}

func (this *ssdpDefaultManager) ssdpGetServiceKey(sv *ssdpService) (key ServiceKey) {
	uri := sv.uri
	if 0 >= len(uri) {
		uri = "schemas-upnp-org"
	}
	key = ServiceKey(fmt.Sprintf("%s-%s", uri, sv.name))
	return
}

func (this *ssdpDefaultManager) ssdpRequireService(de *ssdpDevice, sv *ssdpService) {
	key := this.ssdpGetServiceKey(sv)
	if _, has := de.services[key]; !has {
		de.services[key] = sv
		if _, has := this.serviceMap[key]; !has {
			this.serviceMap[key] = make(DeviceMap)
		}
		this.serviceMap[key][de.UUID()] = de
	}
}

func (this *ssdpDefaultManager) ssdpNotifyService(res *ssdpResource) {
	de := this.ssdpRequireDevice(res)
	sv := ssdpNewService(res)
	this.ssdpRequireService(de, sv)
}

func (this *ssdpDefaultManager) ssdpNotifyResource(res *ssdpResource) {
	switch res.ssdptype {
	case ssdpTypessdpRootDevice:
		this.ssdpNotifyssdpRootDevice(res)
	case ssdpTypeDevice:
		this.ssdpNotifyDevice(res)
	case ssdpTypeService:
		this.ssdpNotifyService(res)
	case ssdpTypeHNAP:
		/*TODO*/
	case ssdpTypeUUID:
		/*TODO ... not sure what to do*/
	default:
		log.Fatalf("Unhandled ssdptype %d", res.ssdptype)
	}
}

func (this *ssdpDefaultManager) ssdpIncludessdpRootDevice(ssdpsm *ssdpResponseMessage) {
	if n := ssdpRootDeviceUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.uuid = uuid
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypessdpRootDevice
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name for upnp:rootdevice: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeService(ssdpsm *ssdpResponseMessage) {
	if n := ssdpUPnPServiceUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypeService
		res.name = n[2]
		var err error
		if res.version, err = strconv.ParseInt(n[4], 10, 4); nil != err {
			log.Printf("Error in parsing service version `%s'", n[4])
		}
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name for UPnP service: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeDevice(ssdpsm *ssdpResponseMessage) {
	if n := ssdpUPnPDeviceUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypeDevice
		res.name = n[2]
		var err error
		if res.version, err = strconv.ParseInt(n[4], 10, 4); nil != err {
			log.Printf("Error in parsing device version `%s'", n[4])
		}
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name for UPnP device: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeUUID(ssdpsm *ssdpResponseMessage) {
	if n := ssdpUPnPBareUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypeUUID
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeHNAP(ssdpsm *ssdpResponseMessage, name string) {
	if n := ssdpUPnPBareUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypeHNAP
		res.name = name
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name for HNAP: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeOtherService(ssdpsm *ssdpResponseMessage) {
	if n := ssdpOtherServiceUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypeService
		res.uri = n[2]
		res.name = n[3]
		var err error
		if res.version, err = strconv.ParseInt(n[5], 10, 4); nil != err {
			log.Printf("Error in parsing service version `%s'", n[4])
		}
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name for third-party service: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeOtherDevice(ssdpsm *ssdpResponseMessage) {
	if n := ssdpOtherDeviceUUIDRegex.FindStringSubmatch(ssdpsm.usn); 0 < len(n) {
		uuid := UUID(n[1])
		res := ssdpNewResource(uuid)
		res.ssdpServerDescription = ssdpsm.ssdpServerDescription
		res.location = ssdpsm.location
		res.ssdptype = ssdpTypeDevice
		res.uri = n[2]
		res.name = n[3]
		var err error
		if res.version, err = strconv.ParseInt(n[5], 10, 4); nil != err {
			log.Printf("Error in parsing device version `%s'", n[4])
		}
		this.ssdpNotifyResource(res)
	} else {
		log.Printf("Invalid Unique Service Name for third-party device: `%s'", ssdpsm.usn)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeResponse(msg *ssdpResponseMessage) {
	if "upnp:rootdevice" == msg.st {
		this.ssdpIncludessdpRootDevice(msg)
	} else if m := ssdpUPnPServiceRegex.FindStringSubmatch(msg.st); 0 < len(m) {
		this.ssdpIncludeService(msg)
	} else if m := ssdpUPnPDeviceRegex.FindStringSubmatch(msg.st); 0 < len(m) {
		this.ssdpIncludeDevice(msg)
	} else if m := ssdpUPnPUIDRegex.FindStringSubmatch(msg.st); 0 < len(m) {
		this.ssdpIncludeUUID(msg)
	} else if m := ssdpHNAPRegex.FindStringSubmatch(msg.st); 0 < len(m) {
		this.ssdpIncludeHNAP(msg, m[1])
	} else if m := ssdpOtherServiceRegex.FindStringSubmatch(msg.st); 0 < len(m) {
		this.ssdpIncludeOtherService(msg)
	} else if m := ssdpOtherDeviceRegex.FindStringSubmatch(msg.st); 0 < len(m) {
		this.ssdpIncludeOtherDevice(msg)
	} else {
		log.Printf("Unsupported search term [ST] `%s'", msg.st)
	}
}

func (this *ssdpDefaultManager) ssdpIncludeNotification(msg *ssdpNotifyMessage) {
	/*TODO*/
}

func (this *ssdpDefaultManager) ssdpSendQuery(timeout int) (err error) {
	msg := this.ssdpQueryMessage(timeout)
	if _, err = this.unicast.conn.WriteTo(msg.Bytes(), this.multicast.addr); nil != err {
		return
	} else {
		timeout := time.NewTimer(time.Duration(timeout) * time.Second)
		done := false
		for !done {
			select {
			case m := <-this.responseQueue:
				this.ssdpIncludeResponse(m)
			case raw := <-this.notifyQueue:
				this.ssdpIncludeNotification(raw)
			case <-timeout.C:
				done = true
			}
		}
	}
	return
}

func (this *ssdpDefaultManager) ssdpQueryLoop() (err error) {
	for i := 0; i < 2; i++ {
		if err = this.ssdpSendQuery(3 /*timeout seconds*/); nil != err {
			return
		}
	}
	return
}

func (this *ssdpDefaultManager) ssdpDiscoverImpl(ifiname, port string, subscribe bool) {
	ifi, err := net.InterfaceByName(ifiname)
	if nil != err {
		panic(err)
	} else if err = this.ssdpUnicastDiscoverImpl(ifi, port); nil != err {
		panic(err)
	} else if err = this.ssdpMulticastDiscoverImpl(ifi, subscribe); nil != err {
		panic(err)
	} else if err = this.ssdpQueryLoop(); nil != err {
		panic(err)
	}
}
