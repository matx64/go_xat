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
        if(username == window.username){
            messagesList.innerHTML += `<div class="d-inline-flex flex-column text-white rounded bg-primary p-3 my-2">
                <div class="text-break">${text}</div> 
                <div class="fs-6 fw-light text-end">${sentAt}</div>
            </div>`
        }else{
            messagesList.innerHTML += `<div class="d-inline-flex flex-column text-white rounded bg-warning p-3 my-2">
                <div class="fw-bold">${username}</div> 
                <div class="text-break">${text}</div> 
                <div class="fs-6 fw-light text-end">${sentAt}</div>
            </div>`
        }
    },
    "join": (text, sentAt, username) => {
        roomUsers.add(username)
        usersList.innerHTML += `<div class="text-white" id="user-${username}"><strong>${username}</strong></div>`
        messagesList.innerHTML += `<div class="text-success text-center my-2"><strong>${text}</strong> ${sentAt}</div>`
    },
    "left": (text, sentAt, username) => {
        roomUsers.delete(username)
        document.getElementById("user-" + username).remove()
        messagesList.innerHTML += `<div class="text-danger text-center">${sentAt} <strong>${text}</strong></div>`
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