#!/bin/sh

while :; do curl -H "X-Real-IP: $((1+(RANDOM%254))).$((1+(RANDOM%254))).$((1+(RANDOM%254))).$((1+(RANDOM%254)))" localhost:8081; done
