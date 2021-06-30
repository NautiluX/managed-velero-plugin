FROM registry.ci.openshift.org/openshift/release:golang-1.14

WORKDIR /src
COPY . .
RUN go build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
LABEL io.openshift.managed.name="managed-velero-plugin-status-patch" \
      io.openshift.managed.description="Velero sidekick to apply status after creation"

COPY --from=0 /src/managed-velero-plugin /bin/managed-velero-plugin

ENTRYPOINT ["bin/managed-velero-plugin"]
