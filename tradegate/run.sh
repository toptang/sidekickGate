#!/bin/bash

nohup ./tradegate \
	-port=5010 \
	-server_okex="" \
	-data_service="http://127.0.0.1:6001" \
	-redis="172.19.16.15:7000" \
	&
