#!/bin/bash

if [ "$1" = "dev" ] || [ -z "$1" ]; then
    echo "Starting in development mode..."
    go run $(find cmd/web -name "*.go" -not -name "*_test.go") \
        -dbname=bookings \
        -dbuser=ashparsh \
        -dbpassword=postgres \
        -production=false \
        -cache=false \
        -mailhost=localhost \
        -mailport=1025 \
        -mailusername="" \
        -mailpassword="" \
        -mailencryption=none \
        -mailfrom=noreply@bookings.dev \
        -mailfromname="Bookings Dev"
fi

if [ "$1" = "prod" ]; then
    echo "Starting in production mode..."
    go run $(find cmd/web -name "*.go" -not -name "*_test.go") \
        -dbname=bookings_db_8szz \
        -dbuser=bookings_db_8szz_user \
        -dbpassword=20iNeVoaJgHVGEcnLHCxhbWNMlOlVvye \
        -dbhost=dpg-d0rhah15pdvs73e0csr0-a.singapore-postgres.render.com \
        -dbport=5432 \
        -dbssl=require \
        -production=true \
        -cache=true \
        -mailhost=smtp.gmail.com \
        -mailport=587 \
        -mailusername=ashparshp1@gmail.com \
        -mailpassword=nthugvfxhsqgtizw \
        -mailencryption=starttls \
        -mailfrom=noreply@bookings.com \
        -mailfromname="bookings"
fi

if [ "$1" = "build" ]; then
    echo "Building application..."
    go build -o bookings $(find cmd/web -name "*.go" -not -name "*_test.go")
    echo "Built successfully! Run with: ./bookings [flags]"
fi

# Add help option
if [ "$1" = "help" ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "Usage: ./run.sh [mode]"
    echo ""
    echo "Modes:"
    echo "  dev   - Development mode (default)"
    echo "  prod  - Production mode with Gmail SMTP"
    echo "  build - Build binary only"
    echo "  help  - Show this help"
    echo ""
    echo "Examples:"
    echo "  ./run.sh      # Run in development"
    echo "  ./run.sh dev  # Run in development"
    echo "  ./run.sh prod # Run in production"
    echo "  ./run.sh build # Build binary"
fi