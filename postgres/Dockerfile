FROM postgres:alpine
ENV POSTGRES_PASSWORD="avitopass"
COPY ./entrypoint/db.sql /docker-entrypoint-initdb.d
