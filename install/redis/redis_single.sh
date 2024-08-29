#!/bin/sh

# redis single
docker pull redis

# redis execute
docker run -d --name chat-redis -p 6379:6379 redis