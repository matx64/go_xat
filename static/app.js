const connectForm = document.getElementById("connect-form");
const messageForm = document.getElementById("message-form");
const quitForm = document.getElementById("quit-form");
const usernameInput = document.getElementById("username-input")
const messageInput = document.getElementById("message-input")
const messagesList = document.getElementById("messages-list")

const preDefinedTexts = {
    "standard": (text, username) => `<li class="list-group-item"><strong>${username}</strong>: ${text}</li>`,
    "join": (text, _username) => `<li class="list-group-item text-success"><strong>${text}</strong></li>`,
    "left": (text, _username) => `<li class="list-group-item text-danger"><strong>${text}</strong></li>`,
}

connectForm.addEventListener("submit", function (e) {
    e.preventDefault()

    window.username = usernameInput.value

    setupWebSocket()
    setupMessageForm()

    connectForm.classList.add("d-none")
    messagesList.classList.remove("d-none")
    messageForm.classList.remove("d-none")
    quitForm.classList.remove("d-none")
})

function setupWebSocket() {
    window.ws = new WebSocket("ws://" + window.location.host + "/connect?username=" + window.username);

    window.ws.addEventListener("message", function (e) {
        const data = JSON.parse(e.data)
        const message = preDefinedTexts[data.type](data.text, data.username)

        messagesList.innerHTML += message
    })
}

function setupMessageForm() {
    messageForm.addEventListener("submit", function (e) {
        e.preventDefault()

        window.ws.send(
            JSON.stringify({
                username: window.username,
                text: messageInput.value,
                type: "standard"
            })
        )

        messageInput.value = ""
    })
}