# ImageStoreService
Image Store Service project helps create,retrieve and delete *.png,*.jpeg files into MySql DB.


## Running on local machine

### Prerequisites
- Install docker and docker-compose
- Install swagger
    - `brew tap go-swagger/go-swagger`
    - `brew install go-swagger`
- [`go`](https://golang.org/doc/install): built in (1.14 or later)
- Install Kafka
    - git clone https://github.com/wurstmeister/kafka-docker
    - cd kafka-docker
    - Edit docker-compose.yml - KAFKA_ADVERTISED_HOST_NAME: 127.0.0.1 ( Refer https://github.com/wurstmeister/kafka-docker#kafka-docker Readme.md for more details.)
    
1) Build the application  
    - `make build`
2) Run the application
    - `make Run`
3) Generate swagger.yaml
    - `make swagger`

