TAG="localhost/kni-install"

all: build

.PHONY: build
build:
	docker build --no-cache -t "$(TAG)" ./build/
