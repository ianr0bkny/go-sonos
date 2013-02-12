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

var queueSize = 0;
var updateCount = 0;
var currentTrack = 0;

/*
$(function() {
	$.getScript("test.js");
});
*/

function onError(msg) {
	$("#result").empty().append(msg);
}

function onVolume(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else {
		$("#volume-slider").slider("value", data.Value);
	}
}

function formatDuration(d) {
	seconds = d % 60;
	d /= 60;
	minutes = d % 60;
	d /= 60;
	hours = d;
	return Math.floor(hours) + "h" + Math.floor(minutes) + "m" + seconds + "s";
}

function setCurrentTrack(track) {
	if (track != currentTrack) {
		//.remvoveClass()
	}
	currentTrack = track - 1;
}

function onPositionInfo(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if ("Value" in data) {
		obj = data.Value;
		$("#track").text(obj.Track + "/" + queueSize);
		$("#track-duration").text(formatDuration(obj.TrackDuration));
		$("#rel-time").text(formatDuration(obj.RelTime));
		$("#title").text(obj.Title);
		$("#album").text(obj.Album);
		$("#progress-bar").progressbar("value", 100 * (obj.RelTime / obj.TrackDuration));
		$(".progress-label").text(formatDuration(obj.TrackDuration - obj.RelTime));
		setCurrentTrack(obj.Track);
	}
}

function onRemoveTrack(data) {
	if ("Error" in data) {
		onError(data.Error);
	}
}

function removeTrack(num) {
	$.post("/control", {method: "remove-track", track: num}, onRemoveTrack, "json");
}

function onPlayTrack(data) {
	if ("Error" in data) {
		onError(data.Error);
	}
}

function playTrack(num) {
	$.post("/control", {method: "play-track", track: num}, onPlayTrack, "json");
}

function xmlUnescape(s) {
	s = s.replace(/%3a/g, ":");
	s = s.replace(/%2f/g, "/");
	s = s.replace(/%2520/g, " ");
	return s;
}

function writeTrackRow(track, num) {
	$("#current-queue>tbody").append(
		  "<tr>"
		+ "<td>" + num + "</td>"
		+ "<td><a href=\"javascript:playTrack(" + num + ")\">Jump</a></td>"
		+ "<td><a href=\"javascript:removeTrack(" + num + ")\">Dele</a></td>"
		+ "<td>" + track.Creator + "</td>"
		+ "<td>" + track.Album + "</td>"
		+ "<td>" + track.Title + "</td>"
		//+ "<td><img src=\"http://192.168.1.44:1400" + xmlUnescape(track.AlbumArtURI) + "\"/></td>"
		+ "</tr>");
}

function onQueue(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if ("Value" in data) {
		queueSize = data.Value.length;
		$("#current-queue>tbody").empty();
		for (i = currentTrack; i < data.Value.length; i++) {
			writeTrackRow(data.Value[i], i + 1);
		}
		for (i = 0; i < currentTrack; i++) {
			writeTrackRow(data.Value[i], i + 1);
		}
	}
}

function playButtonPlay() {
	options = {
		label: "Play",
		icons: {
			primary: "ui-icon-play"
		}
	}
	$("#control-panel>#play").button("option", options);
}

function playButtonPause() {
	options = {
		label: "Pause",
		icons: {
			primary: "ui-icon-pause"
		}
	}
	$("#control-panel>#play").button("option", options);
}

function onTransportInfo(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if("Value" in data) {
		state = data.Value.CurrentTransportState
		if ("STOPPED" == state || "PAUSED_PLAYBACK" == state) {
			playButtonPlay();
		} else if ("PLAYING" == state) {
			playButtonPause();
		}
	}
}

function writeGenreRow(genre) {
	$("#genre-table>tbody").append(
		  "<tr>"
		+ "<td>" + genre.Title + "</td>"
		+ "</tr>");
}

function onGenreList(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if("Value" in data) {
		$("#genre-table>tbody").empty();
		for (i in data.Value) {
			genre = data.Value[i];
			writeGenreRow(genre);
		}
	}
}

function eventLoop() {
	$.post("/control", {method: "get-volume"}, onVolume, "json");
	$.post("/control", {method: "get-position-info"}, onPositionInfo, "json");
	$.post("/control", {method: "get-transport-info"}, onTransportInfo, "json");
	if (++updateCount % 5 == 0) {
		$.post("/control", {method: "get-queue"}, onQueue, "json");
	}
}

function initUi() {
	$.post("/control", {method: "get-volume"}, onVolume, "json");
	$.post("/control", {method: "get-position-info"}, onPositionInfo, "json");
	$.post("/control", {method: "get-transport-info"}, onTransportInfo, "json");
	$.post("/control", {method: "get-queue"}, onQueue, "json");
	$.post("/control", {method: "list-genres"}, onGenreList, "json");
}

