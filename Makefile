# Copyright © 2020 Jonathan Cope jcope@redhat.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TAG = localhost/kni-install

.PHONY: all build
all: build

build:
	docker build -t "$(TAG)" -f build/Dockerfile .


# clean destroys the intermediate build images which consume a couple Gb of space in total.  It does not
# delete the app image or the base build images.
clean:
	docker rmi $(shell docker images -aq --filter "dangling=true" --filter "label=KNI_BUILDER=true") 2> /dev/null


# clean-all destroys all image artifacts created by the build target (app, intermediate, and base images).
clean-all: clean
	-docker rmi $(TAG) 2>/dev/null
	-docker rmi $(shell docker images -aq registry.redhat.io/ubi8/ubi-minimal:8.2) 2> /dev/null
	-docker rmi $(shell docker images -aq registry.redhat.io/ubi8/go-toolset:1.13.4) 2> /dev/null

