#!/usr/bin/wish

package require Tk
package require json
package require http

namespace eval state {
	variable volume
}

set URL "http://localhost:8080"
set CONTROL "$URL/control"

proc on_set_volume {token} {
	upvar $token state
	puts $state(code)
}

proc on_volume_scale {value} {
	http::geturl "$::CONTROL?method=set-volume&value=$state::volume" -command on_set_volume
}

proc on_get_volume {token} {
	upvar $token state
	set volume [dict get [json::json2dict $state(body)] Value]
	puts "got volume $volume"
	set state::volume $volume
}

proc get_volume {} {
	http::geturl "$::CONTROL?method=get-volume" -command on_get_volume
}

wm title . "UPnP Console"
pack [frame .top_level]
pack [label .top_level.volume_label -text "Volume"]
pack [scale .top_level.volume_scale -command on_volume_scale -from 0 -to 100 -showvalue 1 -variable state::volume]
get_volume

