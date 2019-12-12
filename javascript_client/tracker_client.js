URL_MASTER = "ws://localhost:1323/ws/notAToken";

webSocket = new WebSocket(URL_MASTER, "tracker");

websocket.onopen = async function() {
  while ("always") {
    await sleep(30_000);
    exampleSocket.send("Keep connection alive");
  }
};

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
