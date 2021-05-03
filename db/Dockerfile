FROM postgres:12.0

ENV TZ Asia/Tokyo

COPY init/setup.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
