TAG = "localhost/kni-install"
BRANCH ?= "master"

all: build

.PHONY: build
build:
	docker build --build-arg "GIT_BRANCH=$(BRANCH)" --no-cache -t "$(TAG)" ./build/

