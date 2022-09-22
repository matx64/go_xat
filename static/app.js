const connectForm = document.getElementById("connect-form");
const messageForm = document.getElementById("message-form");
const quitForm = document.getElementById("quit-form");
const usernameInput = document.getElementById("username-input")
const roomInput = document.getElementById("room-input")
const messageInput = document.getElementById("message-input")
const conversationDiv = document.getElementById("conversation-div")
const messagesList = document.getElementById("messages-list")
const usersList = document.getElementById("users-list")
const roomDisplay = document.getElementById("room-display")

let roomUsers = new Set()

const handleMessage = {
    "normal": (text, sentAt, username) => {
        messagesList.innerHTML += `<li class="list-group-item">${sentAt} <strong>${username}</strong>: ${text}</li>`
    },
    "join": (text, sentAt, username) => {
        roomUsers.add(username)
        usersList.innerHTML += `<li class="list-group-item" id="user-${username}"><strong>${username}</strong></li>`
        messagesList.innerHTML += `<li class="list-group-item text-success">${sentAt} <strong>${text}</strong></li>`
    },
    "left": (text, sentAt, username) => {
        roomUsers.delete(username)
        document.getElementById("user-" + username).remove()
        messagesList.innerHTML += `<li class="list-group-item text-danger">${sentAt} <strong>${text}</strong></li>`
    }
}

connectForm.addEventListener("submit", function (e) {
    e.preventDefault()

    window.username = usernameInput.value
    window.roomId = roomInput.value

    setupWebSocket()
    setupMessageForm()

    roomDisplay.innerHTML = "Room " + roomInput.value
    connectForm.classList.add("d-none")
    conversationDiv.classList.remove("d-none")
    messageForm.classList.remove("d-none")
    quitForm.classList.remove("d-none")
})

function setupWebSocket() {
    window.ws = new WebSocket("ws://" + window.location.host + "/connect?roomId=" + window.roomId + "&username=" + window.username);

    window.ws.addEventListener("message", function (e) {
        const data = JSON.parse(e.data)
        data.sentAt = new Date(data.sentAt * 1000).toLocaleTimeString()

        handleMessage[data.type](data.text, data.sentAt, data.username)
    })
}

function setupMessageForm() {
    messageForm.addEventListener("submit", function (e) {
        e.preventDefault()

        window.ws.send(
            JSON.stringify({
                roomId: window.roomId,
                username: window.username,
                text: messageInput.value,
                type: "normal",
                sentAt: Math.floor(Date.now() / 1000)
            })
        )

        messageInput.value = ""
    })
}