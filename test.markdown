---
layout: default
title: go-sonos
---

<!--
# Documentation
******

## Discovery [go-sonos/ssdp]
******
* MakeManager() ssdp.Manager

### ssdp.Manager interface
******
* Discover(ifiname, port string, subscribe bool) error
	* device
		* The network device to query for UPnP devices via SSDP [e.g. 'eth0'];
	* port
		* A free port to use to listen for responses;
	* subscribe [unimplemented]
		* Listen to asynchronous updates after the initial query is complete;
* QueryServices(query ServiceQueryTerms) ServiceMap
	* query
		* Query terms, consisting of pairs of service keys and minimim required versions;
* Devices() DeviceMap
* Close() error
-->

### Device
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

Structs implementing ssdp.Device describe a UPnP device on the network,
which can provide any number of services.

### Manager
******

{% highlight go %}
type Manager interface {
	Discover(ifiname, port string, subscribe bool) error
	QueryServices(query ServiceQueryTerms) ServiceMap
	Devices() DeviceMap
	Close() error
}
{% endhighlight %}

Structs implementing ssdp.Manager provide access to the Simple Service
Discovery Protocol (SSDP), which allows control points to discover UPnP
devices on the network.

#### Example
******
{% highlight go %}
mgr := ssdp.MakeManager()
{% endhighlight %}

### ServiceKey
******

{% highlight go %}
type ServiceKey string
{% endhighlight %}

Typo protection for an SSDP service key

#### Example
******
{% highlight go %}
ssdp.ServiceKey("schemas-upnp-org-MusicServices")
{% endhighlight %}

### ServiceQueryTerms
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


<!--
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
-->

