# Sensor Station

Sensor Station gathers data via MQTT from any number of _publishers_,
which are typically battery powered, wireless sensors spread around a
certain location.

# Overview

## MQTT is the key

### MQTT Broker 

- Run MQTT broker, e.g. mosquitto

- Base topic "ss/<id>/data/<data-type>"

Example: ```ss/00:95:fb:3f:34:95/data/tempc 25.00```

### Web Sockets

We sockets or HTTP/2 will be used to send data to and from the IOTe
device (otto) in our case.

### Subscribe to Topics

- announce/station  - announces stations that control or collect
- announce/hub      - announces hubs, typ

- data/tempc/       - data can have option /index at the end
- data/humidity

- control/relay/idx - control can have option /index at the end

## REST API

- GET   /api/config 
- PUT   /api/config     data => { config: id, ... }

- GET /api/data
- GET /api/stations

## Station Manager 

- Collection of stations
- Stations can age out 

### Stations

- ID (name, IP and mac address)
- Capabilities
  - sensors
  - relay

## Data

Data can be optimized and we expect we will want to optimize different
data for all kinds of reasons and we won't preclude that from
happening, we'll give applications the flexibility to handle data
elements as they see fit (can optimize).

We will take an memory expensive approach, every data point can be
handled on it's own. The data structure will be:

    struct Data
        Source ID
        Type
        Timestamp
        Value

# User Interface

This project haws 


# Build

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
