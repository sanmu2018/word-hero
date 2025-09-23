# note: call scripts from /scripts

image_name =word-hero
harbor_addr=sanmu2018/${image_name}
tag        =v0.0.2
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


bin: get-deps
	go build -o ${image_name} cmd/${image_name}.go

tag:
	docker build -f build/Dockerfile . -t ${image_name}:${tag}
	docker tag ${image_name} ${harbor_addr}:${tag}

push: tag
	docker push ${harbor_addr}:${tag}

dist: testarch
	docker build -t ${image_name} -f build/Dockerfile .
	docker tag ${image_name} ${harbor_addr}/${arch}:${tag}
	docker push ${harbor_addr}/${arch}:${tag}

dockertag:
	docker build --build-arg GITLAB_USER=${GITLAB_USER} --build-arg GITLAB_PWD=${GITLAB_PWD} -t ${image_name} -f build/Dockerfile .
	docker tag ${image_name} ${harbor_addr}:${tag}

dockerpush: dockertag
	docker push ${harbor_addr}:${tag}

localbuild: get-deps bin

dkpush: localpush

run:
	docker compose up -d

stop:
	docker compose stop

start:
	docker compose start