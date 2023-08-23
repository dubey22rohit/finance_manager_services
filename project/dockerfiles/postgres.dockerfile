FROM postgres

COPY ./scripts/postgres/init.sql /docker-entrypoint-initdb.d/