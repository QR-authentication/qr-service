#!/bin/bash

goose -dir ./migrations postgres "user=$QR_SERVICE_POSTGRES_USER password=$QR_SERVICE_POSTGRES_PASSWORD dbname=$QR_SERVICE_POSTGRES_DB host=$QR_SERVICE_POSTGRES_HOST port=$QR_SERVICE_POSTGRES_PORT sslmode=disable" down
