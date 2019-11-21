UID=$(shell id -u)
GID=$(shell id -g)

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
