# Use the official PostgreSQL image as the base image
FROM postgres:latest

# Set environment variables for PostgreSQL
ARG POSTGRES_PORT 
ARG POSTGRES_USERNAME 
ARG POSTGRES_PASSWORD 
ARG POSTGRES_NAME 

ENV POSTGRES_PORT=$POSTGRES_PORT
ENV POSTGRES_USERNAME=$POSTGRES_USERNAME
ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV POSTGRES_NAME=$POSTGRES_NAME

# Copy SQL scripts to initialize the database
COPY init.sql /docker-entrypoint-initdb.d/

# Grant necessary permissions to the SQL script
RUN chmod +r /docker-entrypoint-initdb.d/init.sql

EXPOSE $POSTGRES_PORT
