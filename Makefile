swagger:
	swagger generate spec -o ./swagger.yaml --scan-models
run:
	go run main.go
build:
	go build main.go