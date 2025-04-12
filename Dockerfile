ARG GO_IMAGE
ARG BASE_IMAGE
FROM $GO_IMAGE AS builder

WORKDIR /workspace/

ARG IMAGE_TAG
ARG GOPROXY

COPY . .

RUN make build

FROM $BASE_IMAGE

COPY migrations migrations
COPY --from=builder /workspace/bin/holiday /holiday

HEALTHCHECK CMD exit 0

# 65534 is the uid/gid of nobody, if you use scratch keep as id
USER nobody

# Run the binary.
ENTRYPOINT ["/holiday"]
