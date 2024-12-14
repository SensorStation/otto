const clientId = 'mqttjs_' + Math.random().toString(16).substr(2, 8);
const host = "ws://" + window.location.hostname + ":8080";
console.log(window.location.hostname);

const options = {
    keepalive: 60,
    clientId: clientId,
    protocolId: 'MQTT',
    protocolVersion: 4,
    clean: true,
    reconnectPeriod: 1000,
    connectTimeout: 30 * 1000,
    will: {
        topic: 'WillMsg',
        payload: 'Connection Closed abnormally..!',
        qos: 0,
        retain: false
    },
}
console.log('Connecting mqtt client');
const client = mqtt.connect(host, options);
client.on('error', (err) => {
    console.log('Connection error: ', err);
    client.end();
})
client.on('reconnect', () => {
    console.log('Reconnecting...');
})

client.on('connect', () => {
    console.log(`Client connected: ${clientId}`);
    // Subscribe
    console.log("subscribing to ss/c/station/relay");
    client.subscribe('ss/c/station/relay', { qos: 0 });
})
// Unsubscribe
/* client.unsubscribe('tt', () => {
 *   console.log('Unsubscribed');
 * })
 */
// Publish
/* console.log("sending to tt")
 * client.publish('ss/c/station/relay', 'on', { qos: 0, retain: false })
 * // Receive
 * client.on('message', (topic, message, packet) => {
 *   console.log(`Received Message: ${message.toString()} On topic: ${topic}`)
 * })
 */

var val = "on"
console.log(val)
var toggle = function() {
    console.log("val: ", val);     
    client.publish('ss/c/station/relay', val, { qos: 0, retain: false })

    if (val == "on") {
        val = "off"
    } else {
        val = "on"
    }
}
var id = window.setInterval(toggle, 2000);
