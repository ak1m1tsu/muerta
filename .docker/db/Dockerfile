FROM postgres:alpine
COPY db/init.sql /docker-entrypoint-initdb.d/
COPY db/migrations/*.up.sql /docker-entrypoint-initdb.d/