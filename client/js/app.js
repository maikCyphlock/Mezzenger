var ws = new WebSocket("ws://localhost:8080/ws");



/**
 * Shortcut for document.querySelector
 * @param {string} selector CSS selector
 * @return {Element} The first element matching the selector
 */

const $ =/** @type {(selector: string) => Element} */ (document.querySelector.bind(document));



const sendBtn = $('#send-button');
const messageInput = $('#message-input');

ws.onopen = function (event) {
    console.log("Connection established!");
    ws.send("Hello from the client!");
};

sendBtn.addEventListener('click', () => {
    const message = messageInput.value;
    ws.send(message);
})

ws.onmessage = function (event) {
    console.log("Received message:", event.data);
    $('#message').innerHTML += `<p>${event.data}</p>`;

};

ws.onerror = function (event) {
    console.error("WebSocket error:", event.error);
};

ws.onclose = function (event) {
    console.log("Connection closed.");
};