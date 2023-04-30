#!/bin/sh

DBSTRING="host=$POSTGRES_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=$POSTGRES_SSL"

#sleep 3
goose postgres "$DBSTRING" up