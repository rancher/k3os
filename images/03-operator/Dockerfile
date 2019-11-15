### BASE ###
ARG REPO
ARG TAG
FROM ${REPO}/k3os-progs:${TAG} as progs

FROM scratch
COPY --from=progs /output/k3os-operator /sbin/
ENTRYPOINT ["k3os-operator"]
CMD ["--help"]
