# athens

# local run
    Need GO 1.22.0 version and Docker in your system

    Run docker-compose file for initial setup of dynamo, redis and localstack setup.
    Run the DB migration script using 'go run cmd/migration/main.go dynamo up' 
    Run the application using 'go run cmd/Web/main.go' 
# local run with docker
    Run docker-compose file for initial setup of dynamo, redis and localstack setup.
    Run the DB migration script using 'go run cmd/migration/main.go dynamo up' 
    Build docker image with "docker build -t athens ." from project root directory.
    go to docker directory "cd docker"
    run application stack with "docker-compose up -d" (uncomment athens service from docker-compose)
