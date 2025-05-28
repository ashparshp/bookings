#!/bin/bash
go build -o bookings cmd/web/*.go && ./bookings
./bookings -dbname=bookings -dbuser=ashparsh -cache=false -production=false