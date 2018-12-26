FROM golang:1.11.4-alpine3.8 AS build

RUN apk add --no-cache git
WORKDIR /src/project/
RUN git clone https://github.com/theshadow/rolld.git .
RUN mkdir -p build
RUN GO111MODULE=on go mod vendor
RUN GO111MODULE=on go env
RUN pwd
RUN GO111MODULE=on go build -o /bin/rolld
RUN cd client && GO111MODULE=on go mod vendor && GO111MODULE=on go build -o /bin/rolld-cli

FROM scratch
COPY --from=build /bin/rolld /bin/rolld
COPY --from=build /bin/rolld-cli /bin/rolld-cli
ENTRYPOINT ["/bin/rolld"]
CMD ["start", "--address=:50051"]