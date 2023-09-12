#!/bin/bash

# Set your desired database name
DB_NAME=$POSTGRES_NAME

# Check if the database already exists
if ! psql -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    createdb "$DB_NAME"
    echo "Database $DB_NAME created successfully"
else
    echo "Database $DB_NAME already exists"
fi
