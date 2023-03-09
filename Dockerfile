FROM node:alpine AS client_build
ARG PORT
ENV REACT_APP_API_URL=api/
WORKDIR /client/
COPY ./client ./
RUN yarn install && yarn build

FROM golang:alpine AS server_build
ARG PORT
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
RUN go mod tidy
RUN go build -o ./bin/blog-backend

FROM alpine:latest
ARG PORT
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=client_build /client/build ./build/
COPY --from=server_build /go/src/app/bin/blog-backend .
EXPOSE ${PORT}
ENTRYPOINT ["./blog-backend"]