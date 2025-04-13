ARG ALPINE=alpine:3.21.3

FROM $ALPINE

RUN apk --no-cache --no-progress add tzdata ca-certificates

COPY calendar /

USER 65534

ENTRYPOINT [ "/calendar" ]
