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
		var msg = {
			username: window.username,
			message: $("#message").val().trim()
		};
		sendMessage(msg);
		$("#message").focus();
	});
});

function connect() {
	window.ws = new WebSocket("ws://" + window.location.host + "/chatclient/" + window.room);
	window.ws.onmessage = function(event) {
		appendMessage(JSON.parse(event.data))
	}
	window.ws.onclose = function() {
		connect();
	};

	waitForConnection(window.ws, function() {
		window.history.pushState({}, null, window.location.origin + "/chat/" + window.room);
		$("#landing").css("display", "none");
		$("#chat").css("display", "block");
		$("#message").focus();
		var msg = {
			username: window.username,
			message: window.username + " just joined the chat! Say hello."
		}
		sendMessage(msg);
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

function sendMessage(msg) {
	if (!messageValid(msg)) return;

	window.ws.send(JSON.stringify(msg));
	$("#message").val("");
}

function roomValid(room) {
	if (room.length == 0) return false;

	return true;
}


function appendMessage(msg) {
	$("#messages").append(msg.username + ": " + msg.message + "<br>");
}

function usernameValid(username) {
	if (!username) return false;
	if (username.length > 16) return false;

	return true;
}

function messageValid(msg) {
	if (!usernameValid(msg.username)) return false;
	if (!msg.message) return false;
	if (msg.length > 256) return false;

	return true;
}