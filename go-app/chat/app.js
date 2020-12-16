$( document ).ready(function() {
	$("#chat").css("display", "none");
	$("#username").focus();

	$("#username").keyup(function(event) {
		if (event.keyCode === 13) {
			$("#enter").click();
		}
	});
	$("#enter").click(function() {
		var username = $("#username").val().trim();

		if (usernameValid(username)) {
			$("#messages").append("Logged in as " + username + "<br>");
			window.ws = new WebSocket("ws://" + window.location.host + "/chatclient");
			window.ws.onmessage = function(event) {
				appendMessage(JSON.parse(event.data))
			}

			waitForConnection(window.ws, function() {
				window.username = username;
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

	$("#message").val("");
	window.ws.send(JSON.stringify(msg));
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