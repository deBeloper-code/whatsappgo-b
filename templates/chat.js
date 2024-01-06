document.addEventListener("DOMContentLoaded", function () {
  var socket: any = null;
  var msgBox = document.querySelector("#chatbox textarea");
  var messages: any = document.querySelector("#messages");

  document
    .querySelector("#chatbox")
    .addEventListener("submit", function (event) {
      event.preventDefault();
      if (!msgBox.value) return false;
      if (!socket) {
        alert("Error: There is no socket connection.");
        return false;
      }
      socket.send(msgBox.value);
      msgBox.value = "";
      return false;
    });

  if (!window.WebSocket) {
    alert("Error: Your browser does not support web sockets.");
  } else {
    socket = new WebSocket("ws://localhost:8080/room");
    socket.onclose = function () {
      alert("Connection has been closed.");
    };
    socket.onmessage = function (event) {
      var newMessage = document.createElement("li");
      newMessage.textContent = event.data;
      messages.appendChild(newMessage);
    };
  }
});
