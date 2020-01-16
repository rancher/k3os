FROM scratch
COPY build/k3os/system/ /k3os/system/
ENV PATH /k3os/system/k3os/current:/k3os/system/k3s/current:${PATH}
ENTRYPOINT ["k3os"]
CMD ["help"]
