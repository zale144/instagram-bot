build:
	docker-compose build
	
up:
	docker-compose up
	
clean:
	@echo "Remove all non running containers"
	-docker rm `docker ps -q -f status=exited`
	@echo "Delete all untagged/dangling (<none>) images"
	-docker rmi `docker images -q -f dangling=true`
	
dep:
	cd api && dep ensure
	cd web && dep ensure
	cd htmlToimage && dep ensure
	cd sessions && dep ensure

proto_api:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. api/proto/api.proto

proto_htmltoimage:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. htmlToimage/proto/htmlToimage.proto
	
proto_sessions:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. sessions/proto/session.proto
	
proto_web:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. web/proto/web.proto
	
proto: proto_api proto_htmltoimage proto_sessions proto_web


