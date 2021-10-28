FROM mysql/mysql-server:latest
COPY my.cnf /etc/my.cnf
COPY schema.sql schema.sql


ENV MYSQL_DATABASE=todo \
    MYSQL_ROOT_PASSWORD=password

ADD schema.sql /docker-entrypoint-initdb.d
