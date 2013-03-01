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


/* Kepp track of the number of interface updates */
var updateCount = 0;

/* Position in the playback queue [1...n] */
var currentTrack = 0;

/*
 * A helper function to write a diagnostic error message to the screen.
 */
function onError(msg) {
	$("#result").empty().append(msg);
}

/*
 * Clear a displayed diagnostic error message.
 */
function clearError() {
	$("#result").empty();
}

/*
 * Called in response to control::get-volume; adjusted the volume slider
 * to reflect the current volume level.
 */
function onVolume(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else {
		$("#volume-slider").slider("value", data.Value);
	}
}

/*
 * Format a duration in seconds into a string similar to the output of
 * time(1), with hours, minutes, and seconds broken out.
 */
function formatDuration(d) {
	seconds = d % 60;
	d /= 60;
	minutes = d % 60;
	d /= 60;
	hours = d;
	return Math.floor(hours) + "h" + Math.floor(minutes) + "m" + seconds + "s";
}

/*
 * Called when playback moves to a new track.  This method should change
 * the shading in the queue display of the track currently playing.
 */
function setCurrentTrack(track) {
	if (track != currentTrack) {
		// TODO
	}
	currentTrack = track - 1;
}

/*
 * Callback associated with control::get-position-info.  This method
 * populated the position information table and update the progress.
 */
function onPositionInfo(data, queueSize) {
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

/*
 * Callback associated with control::remove-track-from-queue.  This method
 * just updates the currently displayed error message.  Consider merging
 * into a generic error callback.
 */
function onRemoveTrack(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else {
		clearError();
	}
}

/*
 * Called to remove a track from the playback queue.  Note that @num
 * starts at 1 for the first track in the queue.
 */
function removeTrack(num) {
	$.post("/control", {method: "remove-track-from-queue", track: num}, onRemoveTrack, "json");
}

/*
 * Callback associated with control::seek(TRACK_NR).  This method just
 * updates the currently displayed error message.  Consider merging into
 * a generic error callback.
 */
function onPlayTrack(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else {
		clearError();
	}
}

/*
 * Called to advance playback to the indicated track.  Note that @num
 * starts at 1 for the first track in the queue.
 */
function playTrack(num) {
	$.post("/control", {method: "seek", unit: "TRACK_NR", target: num}, onPlayTrack, "json");
}

/*
 * Escape single quotes in argument strings.
 */
function jsEscape(s) {
	s = s.replace(/'/g, "\\'");
	return s;
}

/*
 * Write an empty in the playback queue as a row in a table.  Note that
 * @num starts at 1 for the first track in the queue.
 */
function writeTrackRow(track, num) {
	$("#current-queue>tbody").append(
		  "<tr>"
		+ "<td>" + num + "</td>"
		+ "<td><a href=\"javascript:playTrack(" + num + ")\">Jump</a></td>"
		+ "<td><a href=\"javascript:removeTrack(" + num + ")\">Dele</a></td>"
		+ "<td>" + track.Creator + "</td>"
		+ "<td>" + track.Album + "</td>"
		+ "<td>" + track.Title + "</td>"
		+ "</tr>");
}

/*
 * Callback associated with get-queue-contents.  This metod rewrites
 * the contents of the #current-queue table.
 */
function onQueue(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if ("Value" in data) {
		var queueSize = data.Value.length;
		$("#current-queue>tbody").empty();
		for (i = currentTrack; i < data.Value.length; i++) {
			writeTrackRow(data.Value[i], i + 1, queueSize);
		}
		for (i = 0; i < currentTrack; i++) {
			writeTrackRow(data.Value[i], i + 1, queueSize);
		}
	}
}

/*
 * Toggles the Play/Pause button to Play.
 */
function playButtonPlay() {
	options = {
		label: "Play",
		icons: {
			primary: "ui-icon-play"
		}
	}
	$("#control-panel>#play").button("option", options);
}

/*
 * Toggles the Play/Pause button to Pause.
 */
function playButtonPause() {
	options = {
		label: "Pause",
		icons: {
			primary: "ui-icon-pause"
		}
	}
	$("#control-panel>#play").button("option", options);
}

/*
 * The callback associated with get-transport-info.  This method updates
 * the appearance of the Play/Pause button using the current playback state.
 */
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

/*
 * Callback associated with the control::set-volume request. This method
 * just updates the currently displayed error message.  Consider merging
 * into a generic error callback.
 */
function onSetVolume(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else {
		clearError();
	}
}

/*
 * Callback associated with an update to the value of the slider widget.
 */
function onVolumeSlider(event, ui) {
	$.post("/control", {method: "set-volume", value: ui.value}, onSetVolume, "json");
}

/*
 * Adapter for the browse::get-direct-children request to use to browse
 * the music library.  Results are sent to @callback.
 */
function getDirectChildren(root, callback) {
	$.post("/browse", {method: "get-direct-children", root: root}, callback, "json");
}

/*
 * Genre
 */

/*
 * The callback associated with the browse::get-all-genres request.
 * This method populates the Genre list in the navigation table.
 */
function onGenreList(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if("Value" in data) {
		clearError();
		$("#genre-table>tbody").empty();
		for (i in data.Value) {
			genre = data.Value[i];
			writeGenreRow(genre);
		}
	}
}

/*
 * Write a row to the list of available genres.  Selecting a genre in
 * this list returns a sublist of the artists associated with that genre.
 */
function writeGenreRow(genre) {
	$("#genre-table>tbody").append(
		  "<tr>"
		+ "<td><a href=\"javascript:getDirectChildren(\'"
		+ jsEscape(genre.ID)
		+ "\', onGetGenreArtists)\">"
		+ genre.Title +
		"</a></td>"
		+ "</tr>");
}

/*
 * Callback associated with getting a list of artists associated with
 * a given genre.  This method populates the Artist list in the Genre
 * navigation table.
 */
function onGetGenreArtists(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if("Value" in data) {
		clearError();
		$("#artist-table>tbody").empty();
		$("#album-table>tbody").empty();
		$("#track-table>tbody").empty();
		for (i in data.Value) {
			artist = data.Value[i];
			writeGenreArtistRow(artist);
		}
	}
}

/*
 * Write a row to the list of artists associated with the selected genre.
 * Selecting an artist in this list returns a list of albums associated
 * with the selected genre and artist.
 */
function writeGenreArtistRow(artist) {
	$("#artist-table>tbody").append(
		  "<tr>"
		+ "<td><a href=\"javascript:getDirectChildren(\'"
		+ jsEscape(artist.ID)
	       	+ "\', onGetGenreArtist)\">"
		+ artist.Title + "</a></td>"
		+ "</tr>");
}

/*
 * Callback associated with getting a list of albums associated with a
 * given genre and artist.  This method populates the Album list in the
 * Genre navigation table.
 */
function onGetGenreArtist(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if("Value" in data) {
		clearError();
		$("#album-table>tbody").empty();
		$("#track-table>tbody").empty();
		for (i in data.Value) {
			album = data.Value[i];
			writeGenreAlbumRow(album);
		}
	}
}

/*
 * Write a row to the list of albums associated with the selected genre
 * and artist.  Selecting an album in the list returns a list of tracks
 * from that album associated with the selected genre and artist.
 */
function writeGenreAlbumRow(album) {
	$("#album-table>tbody").append(
		  "<tr>"
		+ "<td><a href=\"javascript:getDirectChildren(\'"
		+ jsEscape(album.ID)
	       	+ "\', onGetGenreAlbum)\">"
		+ album.Title + "</a></td>"
		+ "</tr>");
}

/*
 * Callback associated with getting a list of tracks associated with a
 * given genre, artist, and album combination.  This method populates the
 * Track list in the Genre navigation table.
 */
function onGetGenreAlbum(data) {
	if ("Error" in data) {
		onError(data.Error);
	} else if("Value" in data) {
		clearError();
		$("#track-table>tbody").empty();
		for (i in data.Value) {
			track = data.Value[i];
			writeGenreTrackRow(track);
		}
	}
}

/*
 * Write a row to the list of traks associated with the selected genre,
 * artist, and album combination.
 */
function writeGenreTrackRow(track) {
	$("#track-table>tbody").append(
		  "<tr>"
		+ "<td>" + track.Title + "</td>"
		+ "</tr>");
}

/*
 * Get the current playback state.
 */
function getState() {
	$.post("/control", {method: "get-volume"}, onVolume, "json");
	$.post("/control", {method: "get-position-info"}, onPositionInfo, "json");
	$.post("/control", {method: "get-transport-info"}, onTransportInfo, "json");
}

/*
 * For each trip through the event loop load the current playback state.
 * Geting the current queue is more expensive, both from a messaging and
 * a UI standpoint, so try to do that less often.
 */
function eventLoop() {
	getState()
	if (++updateCount % 5 == 0) {
		$.post("/browse", {method: "get-queue-contents"}, onQueue, "json");
	}
}

/*
 * Called when the page loads.  Get the current playback state, queue
 * contents, and initialize the navigation pane.
 */
function initUi() {
	getState()
	$.post("/browse", {method: "get-queue-contents"}, onQueue, "json");
	$.post("/browse", {method: "get-all-genres"}, onGenreList, "json");
}

/* END!*/

