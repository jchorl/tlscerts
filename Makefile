UID=$(shell id -u)
GID=$(shell id -g)

dev-serve:
	docker run -it --rm \
		-v $(PWD):/tlscerts:ro \
		-w /tlscerts \
		-p 8080:8080 \
		python:3.9 \
		python -m http.server 8080

prettier:
	docker run -it --rm \
		-v $(PWD):/tlscerts \
		-u $(UID):$(GID) \
		node:15.8 \
		sh -c "mkdir ~/.npm-global && npm config set prefix '~/.npm-global' && npm install -g prettier && ~/.npm-global/bin/prettier --write /tlscerts/index.html"

wasm:
	docker run -it --rm \
		-v $(PWD):/tlscerts \
		-w /tlscerts \
		-e GOOS=js \
		-e GOARCH=wasm \
		-e GOCACHE=/tmp/.cache \
		-u $(UID):$(GID) \
		jchorl/golang \
		go build -o main.wasm

test:
	docker run -it --rm \
		-v $(PWD):/tlscerts:ro \
		-w /tlscerts \
		-e GOCACHE=/tmp/.cache \
		-u $(UID):$(GID) \
		jchorl/golang \
		go test .

compose:
	docker run -it --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):$(PWD):ro \
		-w $(PWD)/examples \
		--env=LANGUAGE=python \
		docker/compose:1.27.4 \
		sh

integration-build:
	docker run -it --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):$(PWD):ro \
		-w $(PWD)/examples \
		--env=LANGUAGE=python \
		docker/compose:1.27.4 \
		build

integration:
	docker run -it --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):$(PWD):ro \
		-w $(PWD)/examples \
		--env=LANGUAGE=python \
		docker/compose:1.27.4 \
		up client
