#Copyright © 2020 Jonathan Cope jcope@redhat.com
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

# Build knictl
FROM registry.redhat.io/ubi8/go-toolset:1.13.4 AS KNICTL

ENV GOPATH="/opt/app-root"
ENV PROJECT_PATH="$GOPATH/src/gerrit.akraino.org/kni/installer"

RUN \
mkdir -p "$PROJECT_PATH" && \
cd "$PROJECT_PATH" && \
git clone "https://gerrit.akraino.org/r/kni/installer" . && \
make build

# Build kni-openshift-installer
FROM registry.redhat.io/ubi8/go-toolset:1.13.4 AS KNI_OPENSHIFT_INSTALL

ARG GIT_BRANCH="master"

ENV GOPATH="/opt/app-root"
ENV APP_FOLDER="$GOPATH/src/kni-openshift-installer"

RUN \
mkdir -p "$APP_FOLDER" && \
cd "$APP_FOLDER" && \
git clone "https://github.com/copejon/kni-openshift-installer" . && \
git checkout "$GIT_BRANCH" && \
./hack/build.sh

# Build container image
FROM registry.redhat.io/ubi8/ubi-minimal:8.2

ENV PLUGIN_FOLDER="/usr/bin/plugins/kustomize/plugin/kni.akraino.org/v1alpha1/siteconfig"

RUN mkdir -p "$PLUGIN_FOLDER"

COPY --from=KNICTL /opt/app-root/src/gerrit.akraino.org/kni/installer/knictl /usr/bin/knictl
COPY --from=KNICTL /opt/app-root/src/gerrit.akraino.org/kni/installer/plugins/kustomize/plugin/kni.akraino.org/v1alpha1/siteconfig/SiteConfig "$PLUGIN_FOLDER/"
COPY --from=KNI_OPENSHIFT_INSTALL /opt/app-root/src/kni-openshift-installer/bin/kni-install /usr/bin/

RUN \
echo "Installing knictl system dependencies" && \
microdnf install git tar

ENTRYPOINT ["kni-install"]