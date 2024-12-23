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

function On() {
    console.log("on")
    client.publish('ss/c/station/relay', "on", { qos: 0, retain: false })    
}

function Off() {
    console.log("off")
    client.publish('ss/c/station/relay', "off", { qos: 0, retain: false })    
}
