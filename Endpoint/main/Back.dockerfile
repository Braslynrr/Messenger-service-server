# First stage: build the Golang binary
FROM golang:1.20 AS builder
WORKDIR /messenger
COPY ./../../ .
WORKDIR /messenger/Endpoint/main
RUN go mod download
WORKDIR /messenger
RUN go build -o Endpoint/main/messengerservice ./Endpoint/main/

# Second stage: build the React app
FROM node:latest AS react-builder
WORKDIR /messenger/ServerFiles/messenger-ui
COPY ./../../ServerFiles/messenger-ui .
RUN npm install
RUN npm run build

# Third stage: create the final image
FROM golang:1.20
WORKDIR /messenger
COPY --from=builder /messenger/Endpoint/main/messengerservice ./Endpoint/main/
COPY --from=react-builder /messenger/ServerFiles/messenger-ui/build ./ServerFiles/messenger-ui/build
COPY  ./../../ServerFiles/countrycodes ./ServerFiles/countrycodes
EXPOSE 8080
WORKDIR /messenger/Endpoint/main
CMD ["./messengerservice"]