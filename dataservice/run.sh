#!/bin/bash

nohup java -jar dataservice-0.0.1-SNAPSHOT.jar \
--server.port=6001 \
--spring.datasource.url=jdbc:mysql://127.0.0.1:3306/sidekick?zeroDateTimeBehavior=convertToNull \
--spring.datasource.username=apple_archer \
--spring.datasource.password=sidekick123 \
&
