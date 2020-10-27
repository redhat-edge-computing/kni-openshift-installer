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

# Build kni-openshift-installer
FROM registry.redhat.io/ubi8/go-toolset:1.13.4 AS KNI_OPENSHIFT_INSTALL

LABEL KNI_BUILDER=true

ENV \
GOPATH="/opt/app-root" \
APP_FOLDER="/opt/app-root/src/kni-openshift-installer"

RUN mkdir -p "$APP_FOLDER"

WORKDIR "$APP_FOLDER"

COPY --chown=1001:root ./ .

RUN ./hack/build.sh

# Build knictl
FROM registry.redhat.io/ubi8/go-toolset:1.13.4 AS KNICTL

LABEL KNI_BUILDER=true

ENV \
GOPATH="/opt/app-root" \
PROJECT_PATH="/opt/app-root/src/gerrit.akraino.org/kni/installer"

RUN mkdir -p "$PROJECT_PATH"

WORKDIR "$PROJECT_PATH"

RUN git clone "https://gerrit.akraino.org/r/kni/installer" . && make build

# Build container image
FROM registry.redhat.io/ubi8/ubi-minimal:8.2

ARG PLUGIN_PATH="/usr/local/plugins/kustomize/plugin/kni.akraino.org/v1alpha1/siteconfig/"

RUN mkdir -p "$PLUGIN_PATH"

COPY --from=KNICTL /opt/app-root/src/gerrit.akraino.org/kni/installer/bin/knictl /usr/local/bin/knictl

COPY --from=KNICTL /opt/app-root/src/gerrit.akraino.org/kni/installer/plugins/kustomize/plugin/kni.akraino.org/v1alpha1/siteconfig/SiteConfig "$PLUGIN_PATH"

COPY --from=KNI_OPENSHIFT_INSTALL /opt/app-root/src/kni-openshift-installer/bin/kni-install /usr/bin/

RUN chmod +x /usr/bin/kni-install

RUN microdnf install git tar

ENTRYPOINT ["kni-install"]
