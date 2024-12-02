# ubuntu image
FROM ubuntu:20.04

# Set environment variables to prevent prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Update package manager and install prerequisites 
# install make for building local makefile
RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    wget \
    make \
    software-properties-common \
    build-essential \
    git \
    && apt-get clean

# sudo \

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash - \
    && apt-get install -y nodejs

# Verify Node.js installation
RUN node -v && npm -v

# Install Go
ENV GO_VERSION=1.22.0
# RUN wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go

# Verify Go installation
RUN go version

# copy local code
WORKDIR /corteza
COPY . .

# Go to server/webapp directory and run make with VERSION
RUN cd server/webapp && make VERSION=2023.9.8

# Go to server directory and run make watch with sudo
# RUN cd server && sudo make watch
# RUN cd server && make watch

# Copy the entrypoint script into the container
COPY entrypoint.sh /entrypoint.sh

# Make the script executable
RUN chmod +x /entrypoint.sh

# Set the entrypoint to run the script
ENTRYPOINT ["/entrypoint.sh"]