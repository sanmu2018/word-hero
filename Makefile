# note: call scripts from /scripts

image_name =word-hero
harbor_addr=${image_name}
tag        =v1.0
arch       =$(shell arch)

testarch:
ifeq (${arch}, x86_64)
	@echo "current build host is amd64 ..."
	$(eval arch=amd64)
else ifeq (${arch},aarch64)
	@echo "current build host is arm64 ..."
	$(eval arch=arm64)
else
	echo "cannot judge host arch:${arch}"
	exit -1
endif
	@echo "arch type:$(arch)"


get-deps:
	go mod tidy
	go mod download

vet-check-all: get-deps
	go vet ./...

gosec-check-all: get-deps
	gosec ./...

bin: get-deps
	go build -o ${image_name} cmd/${image_name}.go

docker:
	docker build -f build/Dockerfile . -t ${image_name}:v0.1.0

gen-swagger:
	swag init -g cmd/${image_name}.go -o api

builder:
	docker build -t ${image_name} -f build/Dockerfile .

push:
	docker tag ${image_name} ${harbor_addr}:${tag}
	docker push ${harbor_addr}:${tag}

dist: testarch
	docker build -t ${image_name} -f build/Dockerfile .
	docker tag ${image_name} ${harbor_addr}/${arch}:${tag}
	docker push ${harbor_addr}/${arch}:${tag}


dockertag:
	docker build --add-host gitlab.apulis.com.cn:120.79.42.225 --build-arg GITLAB_USER=${GITLAB_USER} --build-arg GITLAB_PWD=${GITLAB_PWD} -t ${image_name} -f build/Dockerfile .
    docker tag ${image_name} ${harbor_addr}:${tag}

dockerpush: dockertag
	docker push ${harbor_addr}:${tag}

localbuild: get-deps bin

dkpush: localpush
