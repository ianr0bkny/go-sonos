#!/usr/bin/wish

package require Tk
package require json
package require http

namespace eval state {
	variable volume 0 mute 0
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

proc on_mute_button {} {
	puts "set mute state $state::mute"
}

wm title . "UPnP Console"
pack [ttk::frame .top_level]
pack [ttk::notebook .top_level.notebook]
pack [ttk::frame .top_level.notebook.control]
pack [ttk::frame .top_level.notebook.content]
.top_level.notebook add .top_level.notebook.control -text "Control"
.top_level.notebook add .top_level.notebook.content -text "Content"

set vf .top_level.notebook.control.volume
pack [ttk::labelframe $vf -text "Volume"]
pack [ttk::label $vf.volume_label -text "Volume"]
pack [scale $vf.volume_scale -command on_volume_scale -from 0 -to 100 -variable state::volume -showvalue 1 -orient horizontal]
pack [ttk::label $vf.mute_label -text "Mute"]
pack [ttk::checkbutton $vf.mute -command on_mute_button -variable state::mute]
get_volume

