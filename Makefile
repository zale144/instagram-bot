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
	cd services/api && dep ensure
	cd services/web && dep ensure
	cd services/htmlToimage && dep ensure
	cd services/sessions && dep ensure

proto_api:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/api/proto/api.proto

proto_htmltoimage:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/htmlToimage/proto/htmlToimage.proto
	
proto_sessions:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/sessions/proto/session.proto
	
proto_web:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/web/proto/web.proto
	
proto: proto_api proto_htmltoimage proto_sessions proto_web


