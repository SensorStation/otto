#!/bin/sh

while [ true ]; do 
	echo "blasting a message ... "
	mosquitto_pub -t "ss/d/station-011/tempc" -m 77.4
	mosquitto_pub -t "ss/d/station-011/humidity" -m 15.8
        mosquitto_pub -t "ss/d/station-011/uptime" -m 234244

	mosquitto_pub -t "ss/d/station-012/tempf" -m 45.6
	mosquitto_pub -t "ss/d/station-012/humidity" -m 12.3
        mosquitto_pub -t "ss/d/station-012/uptime" -m 234234
        
	mosquitto_pub -t "ss/d/station-017/tempf" -m 45.16
	mosquitto_pub -t "ss/d/station-017/humidity" -m 112.23
        mosquitto_pub -t "ss/d/station-017/uptime" -m 1234455
        
	sleep 5
done

