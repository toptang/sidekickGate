#!/bin/bash

nohup ./new_quotegate \
	-host="172.19.16.15" \
	-port=5003 \
	-CORS_Origin="http://172.19.16.15" \
	-server_okex="ws://127.0.0.1:8080/ws" \
	-data_service="http://127.0.0.1:6001" \
	&
