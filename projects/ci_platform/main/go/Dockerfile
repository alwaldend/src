FROM node:18 as web
WORKDIR /src/web
COPY web .
RUN yarn install
RUN yarn build

FROM golang:1.19 as build
WORKDIR /src
COPY . .
COPY --from=web /src/web/dist web/dist
RUN CGO_ENABLED=0 go build -o app

FROM alpine:3.16 as final
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=build /src/app /usr/local/bin/ci_platform
COPY ci_platform.yml ci_platform.yml
EXPOSE 8080
CMD ["ci_platform", "server", "run"]
