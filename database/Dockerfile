FROM postgres:10.3

# database config 
# 1.sql first execute
COPY up.sql /docker-entrypoint-initdb.d/1.sql

CMD [ "postgres" ]
