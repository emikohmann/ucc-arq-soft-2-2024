```
docker pull mongo:4.4.6
docker run -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=root --name mongo-container mongo:4.4.6

docker pull rabbitmq:4-management
docker run -d -p 5671:5671 -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password --name rabbit-container rabbitmq:4-management
```