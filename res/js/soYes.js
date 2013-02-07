/*
 * go-sonos
 * ========
 * 
 * Copyright (c) 2012, Ian T. Richards <ianr@panix.com>
 * All rights reserved.
 * 
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

$(function() {
	//$.getScript("test.js");
});

function onError(msg) {
	$("#result").empty().append(data.Error);
}

function onVolume(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else {
		$("#volume-slider").slider("value", data.Value);
	}
}

function onPositionInfo(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if ("Value" in data) {
		obj = data.Value;
		$("#track").text(obj.Track);
		$("#track-duration").text(obj.TrackDuration);
		$("#rel-time").text(obj.RelTime);
		$("#title").text(obj.Title);
		$("#album").text(obj.Album);
	}
}

function onQueue(data) {
	if ("Error" in data) {
		$("#result").empty().append(data.Error);
	} else if ("Value" in data) {
		$("#current-queue>tbody").empty();
		for (n in data.Value) {
			track = data.Value[n];
			$("#current-queue>tbody").append(
				  "<tr>"
				+ "<td>" + n + "</td>"
				+ "<td>" + track.Creator + "</td>"
				+ "<td>" + track.Album + "</td>"
				+ "<td>" + track.Title + "</td>"
				+ "</tr>");
		}
	}
}

function onTransportInfo(data) {
	if ("Error" in data) {
		$("#result").empty().append(data.Error);
	} else if("Value" in data) {
		state = data.Value.CurrentTransportState
		if ("STOPPED" == state || "PAUSED_PLAYBACK" == state) {
			options = {
				label: "Play",
				icons: {
					primary: "ui-icon-play"
				}
			}
			$("#control-panel>#play").button("option", options);
		} else if ("PLAYING" == state) {
			options = {
				label: "Pause",
				icons: {
					primary: "ui-icon-pause"
				}
			}
			$("#control-panel>#play").button("option", options);
		}
	}
}

function clearQueue() {
	$("#current-queue>tbody").empty();
}
