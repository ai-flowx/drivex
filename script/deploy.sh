#!/bin/bash

docker run -v $(pwd)/config/config.yaml:/app/config.yaml -p 4000:4000 craftslab/drivex:latest --config /app/config.yaml --detailed_debug
