# Start from a Golang image with the specified version
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /Micro

# Copy the source code into the container
COPY ./../../ .

# Install any dependencies
WORKDIR /Micro/Micro/MicroBlob

RUN go mod download

WORKDIR /Micro

# Build the executable binary
RUN go build -o Micro/MicroBlob/MicroBlob ./Micro/MicroBlob/

# Set the startup command for the container
CMD ["/Micro/Micro/MicroBlob/MicroBlob"]