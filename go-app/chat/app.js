goog.require('proto.ChatData');
goog.require('proto.ChatMessage');
goog.require('proto.ChatNewClient');

$( document ).ready(function() {
	$("#chat").css("display", "none");
	$("#username").focus();

	var parts = window.location.pathname.split('/chat/');
	var room = parts.pop() || parts.pop();
	if (room.length > 0) {
		$("#room").val(room);
	}

	$(".login-input").keyup(function(event) {
		if (event.keyCode === 13) {
			$("#enter").click();
		}
	});
	$("#enter").click(function() {
		var username = $("#username").val().trim();
		if (!usernameValid(username)) {
			$("#username").focus()
			$("#error").html("Username cannot be empty and must be less than 16 characters");
			return;
		}

		var room = $("#room").val().trim();
		if (!roomValid(room)) {
			$("#room").focus();
			$("#error").html("Room cannot be empty");
			return;
		}

		$("#messages").html("Joined " + room + " as " + username + "<br>");
		window.username = username;
		window.room = room;

		connect()
	});

	$("#message").keyup(function(event) {
		if (event.keyCode === 13) {
			$("#send").click();
		}
	});
	$("#send").click(function() {
		var msg = new proto.ChatMessage();
		msg.setUsername(window.username);
		msg.setMessage($("#message").val().trim());

		var data = new proto.ChatData();
		data.setChatMessage(msg);

		sendData(data);
		$("#message").focus();
	});
});

function connect() {
	window.ws = new WebSocket("ws://" + window.location.host + "/chatclient/" + window.room);
	window.ws.binaryType = 'arraybuffer';

	window.ws.onmessage = function(event) {
		handleData(proto.ChatData.deserializeBinary(event.data));
	}
	window.ws.onclose = function() {
		connect();
	};

	waitForConnection(window.ws, function() {
		window.history.pushState({}, null, window.location.origin + "/chat/" + window.room);
		$("#landing").css("display", "none");
		$("#chat").css("display", "block");
		$("#message").focus();

		var msg = new proto.ChatMessage();
		msg.setUsername(window.username);
		msg.setMessage(window.username + " just joined the chat! Say hello.");

		var data = new proto.ChatData();
		data.setChatMessage(msg);

		sendData(data);
	});
}

function waitForConnection(socket, callback){
    setTimeout(
        function () {
            if (socket.readyState === 1) {
                if (callback != null){
                    callback();
                }
            } else {
                waitForConnection(socket, callback);
            }

        }, 5);
}

function handleData(data) {
	if (!messageValid(data)) {
		console.log("invalid message received: %s", data.ToString())
	}
	if (data.hasChatMessage()) {
		appendMessage(data.getChatMessage());
	}
}

function sendData(data) {
	if (!messageValid(data)) return;

	window.ws.send(data.serializeBinary(), {binary : true});
	$("#message").val("");
}

function roomValid(room) {
	if (room.length == 0) return false;

	return true;
}


function appendMessage(msg) {
	$("#messages").append(msg.getUsername() + ": " + msg.getMessage() + "<br>");
}

function usernameValid(username) {
	if (!username) return false;
	if (username.length > 16) return false;

	return true;
}

function messageValid(data) {
	if (data.hasChatMessage()) {
		var msg = data.getChatMessage();

		if (!usernameValid(msg.getUsername())) return false;
		if (!msg.getMessage()) return false;
		if (msg.getMessage().length > 256) return false;

		return true;
	} 

	return false;
}