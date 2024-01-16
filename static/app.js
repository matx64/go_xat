const connectForm = document.getElementById("connect-form");
const messageForm = document.getElementById("message-form");
const quitForm = document.getElementById("quit-form");
const usernameInput = document.getElementById("username-input");
const roomInput = document.getElementById("room-input");
const messageInput = document.getElementById("message-input");
const conversationDiv = document.getElementById("conversation-div");
const messagesList = document.getElementById("messages-list");
const usersList = document.getElementById("users-list");
const roomDisplay = document.getElementById("room-display");
const chatScroll = document.getElementById("chat-scroll");

let roomUsers = new Set();

const handleMessage = {
  normal: (text, sentAt, username) => {
    if (username == window.username) {
      messagesList.innerHTML += `<div class="m-2"><div class="d-flex flex-column text-white rounded p-3 float-end msg-sent">
                <div class="text-break">${text}</div> 
                <div class="fs-6 fw-light text-end">${sentAt}</div>
            </div></div>`;
    } else {
      messagesList.innerHTML += `<div class="m-2"><div class="d-flex flex-column text-white rounded p-3 float-start msg-received">
                <div class="fw-bold">${username}</div> 
                <div class="text-break">${text}</div> 
                <div class="fs-6 fw-light text-end">${sentAt}</div>
            </div></div>`;
    }
  },
  join: (text, sentAt, username) => {
    roomUsers.add(username);
    usersList.innerHTML += `<div class="text-white m-2 fs-6 fw-bold text-break" id="user-${username}">${username}</div>`;
    messagesList.innerHTML += `<div class="text-success text-center my-2"><strong>${text}</strong> ${sentAt}</div>`;
  },
  left: (text, sentAt, username) => {
    roomUsers.delete(username);
    document.getElementById("user-" + username).remove();
    messagesList.innerHTML += `<div class="text-danger text-center my-2"><strong>${text}</strong> ${sentAt}</div>`;
  },
};

connectForm.addEventListener("submit", function (e) {
  e.preventDefault();

  window.username = usernameInput.value;
  window.roomId = roomInput.value;

  setupWebSocket();
  setupMessageForm();

  roomDisplay.innerHTML = "Room " + roomInput.value;
  connectForm.classList.add("d-none");
  conversationDiv.classList.remove("d-none");
  messageInput.focus();
});

function setupWebSocket() {
  window.ws = new WebSocket(
    "ws://" +
      window.location.host +
      "/connect?roomId=" +
      window.roomId +
      "&username=" +
      window.username
  );

  window.ws.addEventListener("message", function (e) {
    const data = JSON.parse(e.data);
    data.sentAt = new Date(data.sentAt * 1000).toLocaleTimeString();

    handleMessage[data.kind](data.text, data.sentAt, data.username);
    chatScroll.scrollIntoView({ behavior: "smooth", block: "end" });
  });
}

function setupMessageForm() {
  messageForm.addEventListener("submit", function (e) {
    e.preventDefault();

    window.ws.send(
      JSON.stringify({
        roomId: window.roomId,
        username: window.username,
        text: messageInput.value,
        kind: "normal",
        sentAt: Math.floor(Date.now() / 1000),
      })
    );

    messageInput.value = "";
  });
}
