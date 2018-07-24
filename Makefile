# deploy to Kubernetes
deploy_api:
    kubectl apply -f deployments/api/deployment.yaml
    kubectl apply -f deployments/api/service.yaml

deploy_web:
    kubectl apply -f deployments/web/deployment.yaml
    kubectl apply -f deployments/web/service.yaml

deploy_htmltoimage:
    kubectl apply -f deployments/htmltoimage/deployment.yaml
    kubectl apply -f deployments/htmltoimage/service.yaml

deploy_facedetect:
    kubectl apply -f deployments/facedetect/deployment.yaml
    kubectl apply -f deployments/facedetect/service.yaml

deploy_sessions:
    kubectl apply -f deployments/sessions/deployment.yaml
    kubectl apply -f deployments/sessions/service.yaml

deploy_db:
    kubectl apply -f deployments/db/volume.yaml
    kubectl apply -f deployments/db/deployment.yaml
    kubectl apply -f deployments/db/service.yaml

deploy: deploy_db deploy_sessions deploy_facedetect deploy_htmltoimage deploy_web deploy_api

# build docker images
build_api:
    docker build -t livelance/api:v0.0.1 services/api

build_web:
    docker build -t livelance/web:v0.0.1 services/web

build_htmltoimage:
    docker build -t livelance/htmltoimage:v0.0.1 services/htmlToimage

build_facedetect:
    docker build -t livelance/facedetect:v0.0.1 services/facedetect

build_sessions:
    docker build -t livelance/sessions:v0.0.1 services/sessions

build: build_api build_web build_htmltoimage build_facedetect build_sessions


# compile proto bufers
proto_api:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/api/proto/api.proto

proto_htmltoimage:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/htmlToimage/proto/htmlToimage.proto

proto_sessions:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/sessions/proto/session.proto

proto_web:
	protoc --proto_path=$(GOPATH)/src:. --micro_out=. --go_out=. services/web/proto/web.proto

proto: proto_api proto_htmltoimage proto_sessions proto_web


# ensure dependencies
dep_api:
    cd services/api && dep ensure
dep_web:
    cd services/web && dep ensure
dep_htmltoimage:
    cd services/htmlToimage && dep ensure
dep_sessions:
    cd services/sessions && dep ensure

dep: dep_api dep_web dep_htmltoimage dep_sessions


# clean unused docker images and containers
clean:
	@echo "Remove all non running containers"
	-docker rm `docker ps -q -f status=exited`
	@echo "Delete all untagged/dangling (<none>) images"
	-docker rmi `docker images -q -f dangling=true`


# run docker-compose
up:
	docker-compose up