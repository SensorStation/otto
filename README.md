# Sensor Station

Sensor Station gathers data via MQTT from any number of _publishers_,
which are typically battery powered, wireless sensors spread around a
certain location.

## Build

1. Install Go 
2. go get ./...
3. cd ss; go build 

That should leave the execuable 'sensors' in the 'sensors' directory as so:

> ./station/sensors/sensors

## Deploy

1. Install and run an MQTT broker on the sensors host
(e.g. mosquitto).

2. Start the _sensors_ program ensuring the sensor station has
connected to a wifi network.

3. Put batteries in sensors and let the network build itself.

## Testing

### Fake Websocket Data

```bash
% ./ss -fake-ws
```

```bash
% ./ss -help
```

This will open the following URL for the fake websocket data:

> http://localhost:8011/ws

Replace localhost with a hostname or IP if needed. Have the websocket
connect to the URL and start spitting out fake data formatted like
this:

```json
{"year":2020,"month":12,"day":10,"hour":20,"minute":48,"second":8,"action":"setTime"}
{"K":"tempf","V":88}
{"K":"soil","V":0.49}
{"K":"light","V":0.62}
{"K":"humid","V":0.12}
```

## Adding Automated Builds

Automated builds will be using Github events.
