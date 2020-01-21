all: build up

# Build the container
build:
	docker-compose build

# Build and run the container
up:
	docker-compose up -d

# Down and remove container
stop: 
	docker-compose down

# Run all tests: 
test:
	go test ./dao ./resources

