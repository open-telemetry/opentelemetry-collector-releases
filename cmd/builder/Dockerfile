FROM golang:alpine

WORKDIR /build

# Copy default builder config and ocb binary into the container
COPY builder-config.yaml /build/builder-config.yaml
COPY --chmod=755 ocb /build/ocb

ENTRYPOINT [ "./ocb"]
# default arguments to use provided builder config file from documentation;
# replace if using a different file location or file name
CMD [ "--config", "builder-config.yaml" ]