---
layout: default
title: go-sonos
---

# go-sonos
******

go-sonos is a [Go language](http://golang.org/ "The Go Programming Language") library for accessing [UPnP](http://en.wikipedia.org/wiki/Universal_Plug_and_Play "Universal Plug and Play - Wikipedia, the free encyclopedia") devices, including [Sonos](http://www.sonos.com "Sonos WIRELESS Hifi") wireless music devices.

## Discovery
******

### Example
******
{% highlight go %}
mgr := ssdp.MakeManager()
mgr.Discover("eth0", "13104", false)
qry := ssdp.ServiceQueryTerms{
	ssdp.ServiceKey("schemas-upnp-org-MusicServices"): -1,
}
result := mgr.QueryServices(qry)
if device_list, has := result["schemas-upnp-org-MusicServices"]; has {
	for _, device := range device_list {
		...
	}
}
mgr.Close()
{% endhighlight %}

### interface Device
******

{% highlight go %}
type Device interface {
	Product() string
	Name() string
	Location() Location
	UUID() UUID
	Service(key ServiceKey) (service Service, has bool)
}
{% endhighlight %}

Types implementing ssdp.Device describe a UPnP device on the network,
which can provide any number of services.

#### Methods
******
{% highlight go %}
Product() string
{% endhighlight %}
{% highlight go %}
Name() string
{% endhighlight %}
{% highlight go %}
Location() Location
{% endhighlight %}
{% highlight go %}
UUID() UUID
{% endhighlight %}
{% highlight go %}
Service(key ServiceKey) (service Service, has bool)
{% endhighlight %}

### type DeviceMap
******

{% highlight go %}
type DeviceMap map[UUID]Device
{% endhighlight %}

Indexes leaf devices by UUID

### type Location
******

{% highlight go %}
type Location string
{% endhighlight %}

Type protection for a device URI

### interface Manager
******

{% highlight go %}
type Manager interface {
	Discover(ifiname, port string, subscribe bool) error
	QueryServices(query ServiceQueryTerms) ServiceMap
	Devices() DeviceMap
	Close() error
}
{% endhighlight %}

Types implementing ssdp.Manager provide access to the Simple Service
Discovery Protocol (SSDP), which allows control points to discover UPnP
devices on the network.

#### Methods
******
{% highlight go %}
Discover(ifiname, port string, subscribe bool) error
{% endhighlight %}
{% highlight go %}
QueryServices(query ServiceQueryTerms) ServiceMap
{% endhighlight %}
{% highlight go %}
Devices() DeviceMap
{% endhighlight %}
{% highlight go %}
Close() error
{% endhighlight %}

#### Example
******
{% highlight go %}
mgr := ssdp.MakeManager()
{% endhighlight %}

### interface Service
******

{% highlight go %}
type Service interface {
	Version() int64
}
{% endhighlight %}

Types implementing ssd.Service provide an abstraction of a UPnP service.

#### Methods
******
{% highlight go %}
Version() int64
{% endhighlight %}

### type ServiceKey
******

{% highlight go %}
type ServiceKey string
{% endhighlight %}

Type protection for an SSDP service key

#### Example
******

{% highlight go %}
ssdp.ServiceKey("schemas-upnp-org-MusicServices")
{% endhighlight %}

### type ServiceQueryTerms
******

{% highlight go %}
type ServiceQueryTerms map[ServiceKey]int64
{% endhighlight %}

The ServiceQueryTerms type maps a service key to a minimum version,
and is used to query the ssdp.Manager database for matching services.

#### Example
******

{% highlight go %}
query := ssdp.ServiceQueryTerms{
	ssdp.ServiceKey("schemas-upnp-org-MusicServices"): -1,
}
result := mgr.QueryServices(qry)
{% endhighlight %}

### type UUID
******

{% highlight go %}
type UUID string
{% endhighlight %}

Type protection for a Universally Unique ID

