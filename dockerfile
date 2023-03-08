ARG PORT=8080
ARG BUILD_MODE=build
FROM node:alpine AS client_build
WORKDIR /client/
COPY ./client ./
RUN yarn install && yarn build

FROM golang:alpine AS server_build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN go build -o ./bin/blog-backend

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=client_build /client/build ./build/
COPY --from=server_build /go/src/app/bin/blog-backend .
EXPOSE $PORT
ENTRYPOINT ["./blog-backend"]