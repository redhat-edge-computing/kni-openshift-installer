# Build knictl
FROM registry.redhat.io/ubi8/go-toolset:1.13.4 AS KNICTL

ENV GOPATH="/opt/app-root"
ENV PROJECT_PATH="$GOPATH/src/gerrit.akraino.org/kni/installer"

RUN mkdir -p "$PROJECT_PATH"

WORKDIR "$PROJECT_PATH"

RUN git clone "https://gerrit.akraino.org/r/kni/installer" .
RUN make build

# Build kni-openshift-installer
FROM registry.redhat.io/ubi8/go-toolset:1.13.4 AS KNI_OPENSHIFT_INSTALL

ARG GIT_BRANCH="master"

ENV GOPATH="/opt/app-root"
ENV APP_FOLDER="$GOPATH/src/kni-openshift-installer"

RUN mkdir -p "$APP_FOLDER"

WORKDIR "$APP_FOLDER"

RUN git clone "https://github.com/copejon/kni-openshift-installer" .
RUN git checkout "$GIT_BRANCH"

ENTRYPOINT ./hack/build.sh

# Build container image
FROM registry.redhat.io/ubi8/ubi-minimal:8.2

ENV PLUGIN_FOLDER="/usr/bin/plugins/kustomize/plugin/kni.akraino.org/v1alpha1/siteconfig"

RUN mkdir -p "$PLUGIN_FOLDER"
RUN microdnf install git tar

COPY --from=KNICTL /opt/app-root/src/gerrit.akraino.org/kni/installer/knictl /usr/bin/knictl
COPY --from=KNICTL /opt/app-root/src/gerrit.akraino.org/kni/installer/plugins/kustomize/plugin/kni.akraino.org/v1alpha1/siteconfig/SiteConfig "$PLUGIN_FOLDER/"
COPY --from=KNI_OPENSHIFT_INSTALL /opt/app-root/src/kni-openshift-installer/bin/kni-install /usr/bin/

ENTRYPOINT ["kni-install"]