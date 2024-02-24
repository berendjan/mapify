function connectWebSocket() {
    var socket = new WebSocket('ws://127.0.0.1:7878');
    socket.onmessage = function (event) {
        var message = event.data;
        console.log('WebSocket message:', message);
        if (message.trim() === 'refresh') {
            window.location.reload();
        }
    };
    socket.onopen = function () {
        console.log('WebSocket connection established');
    };
    socket.onerror = function (error) {
        console.error('WebSocket error:', error);
    };
    socket.onclose = function () {
        console.log('WebSocket connection closed. Attempting to reconnect...');
        setTimeout(connectWebSocket, 1000);
    };
}
connectWebSocket();
