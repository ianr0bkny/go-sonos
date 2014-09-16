---
layout: default
title: go-sonos
---

# Documentation
******

## Discovery
******
* MakeManager() ssdp.Manager

### ssdp.Manager interface
******
* Discover(ifiname, port string, subscribe bool) error
	* device
		* The network device to query for SSDP network services [e.g. 'eth0'];
	* port
		* A free port to use to listen for responses;
	* subscribe [unimplemented]
		* Listen to asynchronous updates after the initial query is complete;
* QueryServices(query ServiceQueryTerms) ServiceMap
	* query
		* Query terms, consisting of pairs of service keys and minimim required versions;
* Devices() DeviceMap
* Close() error

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

