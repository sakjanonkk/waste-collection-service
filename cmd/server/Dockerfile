# syntax=docker/dockerfile:1
FROM golang AS builder
ARG timezone=Asia/Bangkok
ENV TZ $timezone
WORKDIR /app

# install tini
RUN apt-get update && \
    apt-get -y install tini && \
    apt-get -y clean
COPY ./waste-collection-service/go.* ./
RUN go mod download && go mod verify

# let's build project
COPY ./waste-collection-service/. .
RUN CGO_ENABLED=0 go build -v \
    -installsuffix 'static' \
    -ldflags="-X 'main.version=$(git rev-parse --short HEAD)' -X 'main.build=$(date --iso-8601=seconds)'" \
    -o dist/server ./cmd/server

# pack PRD image
FROM gcr.io/distroless/base:nonroot
LABEL maintainer="Chatchanan Panyaprasirtkit <changnoi2547@gmail.com>"

ARG timezone=Asia/Bangkok

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

# Create app dir
WORKDIR /app
COPY --from=builder /usr/bin/tini-static /usr/bin/tini
COPY --from=builder /app/dist/server /app/server
COPY ./waste-collection-service/configs /app/configs
COPY ./waste-collection-service/internal/assets /app/internal/assets
EXPOSE 8080 8443

# default run entrypoint
ENTRYPOINT ["tini", "--", "/app/server"]
# CMD ["--env=prd"]
CMD ["--env=prd"]
