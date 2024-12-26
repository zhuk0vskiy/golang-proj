
psql -U studios -d studios -c "\i /sql/scripts/create.sql"
psql -U studios -d studios -c "\i /sql/scripts/fill.sql"

http://172.28.1.1:8082/api/v1/login