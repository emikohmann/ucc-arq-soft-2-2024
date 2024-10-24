### Run with 

`docker compose up`

<!--
docker pull mysql:latest
docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=users-api -e MYSQL_PASSWORD=root --name mysql-container mysql:latest
docker exec -t mysql-container bash
mysql -uroot -p
create database `users-api`;

docker pull memcached:latest
docker run -d -p 11211:11211 --name memcached-container memcached:latest

docker pull mongo:4
docker run -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=root --name mongo-container mongo:4
docker exec -t mongo-container bash
mongo --username root --password root --authenticationDatabase admin

docker pull rabbitmq:4-management
docker run -d -p 5671:5671 -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=root --name rabbit-container rabbitmq:4-management

docker pull solr:latest
docker run -d -p 8983:8983 --name solr-container -v $(pwd)/search-api/solr-config:/opt/solr/server/solr/hotels solr:latest solr-create -c hotels
-->