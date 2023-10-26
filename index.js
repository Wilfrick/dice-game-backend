const ws = new WebSocket("ws://localhost:12345/echo");
const sendTextBox = document.querySelector(".controls #sendTextBox")
const sendTextBtn = document.querySelector(".controls #sendTextBtn")
ws.onopen = (event) => {
    ws.send("Hello server")
    log("Send Hello server")
    ws.send("Hi there")
    log("Sent Hi there")
}

ws.onmessage = (event) => {
    console.log(event.data);
    log("Received " + event.data)
}

function log(message) {
    let new_message = document.createElement("li");
    new_message.textContent = message;
    document.querySelector("body #log ol").appendChild(new_message)
}

sendTextBtn.addEventListener("click", () => { log(`Clicked btn with text ${sendTextBox.value}`); ws.send(sendTextBox.value) })