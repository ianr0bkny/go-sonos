---
layout: default
title: go-sonos
---

# Documentation

## Discovery

### Example

{% highlight go linenos %}
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

