#!/bin/sh

if [ "$STORAGE_TYPE" = "memory" ]; then
    echo "Running in-memory storage"
    go run /GraphQLOzon/cmd/main.go
else
    echo "Running with PostgreSQL"
    pg_ctl start -D /usr/local/pgsql/data -l logfile
    sleep 5 # подождем пока PostgreSQL запустится
    go run /GraphQLOzon/cmd/main.go
fi
