// The Websocket connection
let conn;
// The room name, or blank for the global room
let roomName = "";

let msg = document.getElementById("msg");
let log = document.getElementById("log");

function appendLog(item) {
    let doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

window.onload = function () {
    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        let msgToSend = {
            type: "sendmessage",
            data: msg.value
        }
        conn.send(JSON.stringify(msgToSend));
        msg.value = "";
        return false;
    };
};

$(document).ready(function () {
    // Show room selection modal on startup
    $("#room-modal").modal("show");
});

// Make room text box active when the modal shows
// todo: else if websockets is not supported, alert the user to that instead
$('#room-modal').on('shown.bs.modal', function () {
    $('#room-name').trigger('focus')
})

function connectWebSocket() {
    conn = new WebSocket("wss://3fewjvg7q8.execute-api.ap-southeast-2.amazonaws.com/Prod");
    registerWebsocketHandlers()
}

// User enters a room
$("#room-form").submit(function (event) {
    event.preventDefault();
    connectWebSocket()
    roomName = $('#room-name').val()
    let item = document.createElement("div");
    item.innerHTML = "<b>Joined " + roomName + " room</b>";
    appendLog(item)
    $('#room-modal').modal('hide')
    $('#msg').focus()
});

// User enters the global room
$('#room-modal').on('hidden.bs.modal', function (event) {
    event.preventDefault();
    if (roomName === "") {
        connectWebSocket()
        let item = document.createElement("div");
        item.innerHTML = "<b>Joined global room</b>";
        appendLog(item)
        $('#msg').focus()
    }
})

function registerWebsocketHandlers() {
    // Register the user as soon as the socket connection is made
    conn.onopen = function (evt) {
        console.log("Registering user")
        let msgToSend = {
            type: "register",
            roomName: roomName
        }
        conn.send(JSON.stringify(msgToSend));
    }
    conn.onclose = function (evt) {
        let item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        appendLog(item);
    };
    conn.onmessage = function (evt) {
        let messages = evt.data.split('\n');
        for (let i = 0; i < messages.length; i++) {
            let item = document.createElement("div");
            item.innerText = messages[i];
            appendLog(item);
        }
    };
}