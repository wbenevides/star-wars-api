all: build up

# Build the container
builder:
	cd ./deployments && docker-compose build

# Build and run the container
up:
	cd ./deployments && docker-compose up

# Down and remove container
stop: 
	cd ./deployments && docker-compose down

# Run all tests: 
test:
	go test ./internal/app/starwars/dao ./internal/app/starwars/resources

