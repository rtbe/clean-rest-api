check-swagger-install:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger 
	
swagger: check-swagger-install
	swagger generate spec -o ./swagger.yaml --scan-models
	
docker-up:
	docker-compose -f "docker-compose.yml" up --detach --remove-orphans
	
docker-down:
	docker-compose -f 'docker-compose.yml' down --remove-orphans

docker-clean:
	docker system prune -f		

run: swagger docker-up

test-repository-order:
	go test ./repository/order -count=1

test-repository-order-item:
	go test ./repository/order_item -count=1

test-repository-product:
	go test ./repository/product -count=1

test-repository-user:
	go test ./repository/user -count=1
	
test-repository: test-repository-order test-repository-order-item test-repository-product test-repository-user

test-middleware:
	go test ./delivery/web/middlewares -count=1

staticcheck:
	staticcheck ./...	

# To run repository tests you should stop postgresql service: ```sudo systemctl stop postgresql```
test: test-middleware test-repository staticcheck