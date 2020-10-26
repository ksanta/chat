window.onload = function () {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        let msgToSend = {
            message: "sendmessage",
            data: msg.value
        }
        conn.send(JSON.stringify(msgToSend));
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("wss://3fewjvg7q8.execute-api.ap-southeast-2.amazonaws.com/Prod");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }

};

$(document).ready(function () {
    // Show room selection modal on startup
    $("#room-modal").modal("show");
});

// Make room text box active when the modal shows
$('#room-modal').on('shown.bs.modal', function () {
    $('#room-name').trigger('focus')
})

// todo: Entering global room... when user doesn't specify room

// todo: Entering "foobar" room... when user specifies room
