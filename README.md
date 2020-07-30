# ImageStoreService
Image Store Service project helps create,retrieve and delete *.png,*.jpeg files into MySql DB.

## Prerequisites
- Install docker and docker-compose
- Install swagger
    - `brew tap go-swagger/go-swagger`
    - `brew install go-swagger`
- [`go`](https://golang.org/doc/install): built in (1.14 or later)
- Prequisites for Kafka
    - git clone https://github.com/wurstmeister/kafka-docker 
    - cd ~/WORKDIR/kafka-docker
    - Edit docker-compose.yml - KAFKA_ADVERTISED_HOST_NAME: 127.0.0.1 ( Refer https://github.com/wurstmeister/kafka-docker#kafka-docker Readme.md for more details.)
- Run Mysql server
    - ``` sudo docker pull mysql/mysql-server:latest
          sudo docker run --name test-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest]```

## Running on local machine

### Basic commands
1) Build the application  
    - `make build`
2) Run the application
    - `make Run`
3) Generate swagger.yaml
    - `make swagger`

### Steps to run the Application
1) Do the following
   - Change the KafkaURL in imageStoreService/api/core/utilities/kafka_utilities.go --> 127.0.0.1:9092
   - Change the mysql in  imageStoreService/api/core/utilities/utilities.go --> DB_HOST = "tcp(0.0.0.0:3306)" in 
2) Run the application in one terminal
3) Run Kafka broker and zookeeper in 2nd terminal using docker-compose (dir : ~/WORKDIR/kafka-docker) 
   - docker-compose -f docker-compose-single-broker.yml up
4) Run Kafka consumer in 3rd terminal (Once kafka container is up)
   ```sudo docker exec -t kafka-docker_kafka_1 kafka-console-consumer.sh --bootstrap-server :9092  --group jacek-japila-pl --topic NewTopic1 ```
5) Trigger any create/delete api to check kafka notification in the consumer terminal. Api definitions are given below.


### Running Dockerized Application

1) docker build -f Dockerfile -t imagestoreservice3:latest 
2) docker run --name imagestoreserver3 -d  imagestoreservice3:latest -p 8193:8193
3) Create network 
    - docker network create mynetwork
    - docker network connect mynetwork kafka-docker_kafka_1
    - docker network connect mynetwork imagestoreserver3
    - docker network connect mynetwork test-mysql
4) Run Kafka broker and Kafka consumer same as step 3 and step 4 in Steps to run the Application
5) Trigger any create/delete api to check kafka notification in the consumer terminal. Api definitions are given below.

## API Definitions

-  Create Album  : `curl -v -X POST "http://localhost:8193/createAlbum/Family" `
-  Create Image  : `curl -v -X POST -F "image=@/<location>"  "http://localhost:8193/createImage?imageCaption=&albumName=Family"`
    - Both imageCaption and albumName are manadatory parameters.
-  Get Image     :  `curl -v http://localhost:8193/getImage?imageCaption=&albumName=` 
    - Both imageCaption and albumName are manadatory parameters.
    - This will create an image to your working directory to check if the image is getting formed correctly.
-	Delete Album :  `curl -v -X DELETE http://localhost:8193/deleteAlbum/{albumName}` 
    - albumName is a manadator parameter. It will delete the Album and the underlying images with it.
-	Delete Image :   `curl -v -X DELETE http://localhost:8193/deleteImage?imageCaption=&albumName="`
    - AlbumName is not a mandatory parameter. But if used, Image will be deleted for imageCaption and albumName combination.
-   List All Images : `curl -v http://localhost:8193/getAllImages`
    - List all the images in the db. This will create a zip file to your local machine in the working directory so you can check if the images are formed correctly.
