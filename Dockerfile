FROM ubuntu:24.10

RUN apt-get update && \
    apt-get install -y build-essential git cmake make \
                       zlib1g-dev libevent-dev m4 \
                       libelf-dev llvm \
                       clang libc6-dev-i386 libpcap-dev \
		       curl tar sudo linux-tools-6.11.0-19-generic \
		       linux-headers-6.11.0-19-generic pkg-config

# Link asm/byteorder.h into eBPF
RUN ln -s /usr/include/x86_64-linux-gnu/asm/ /usr/include/asm

# installing go
RUN mkdir -p "/usr/local"
RUN curl -L -o /usr/local/go.tar.gz https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf /usr/local/go.tar.gz && rm /usr/local/go.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"

# installing golangci-lint
WORKDIR /usr/local/go
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.0.1

RUN mkdir /src
WORKDIR /src

COPY . .

RUN git clone --branch v1.3.0 https://github.com/libbpf/libbpf.git
WORKDIR /src/libbpf/src
RUN make
RUN make install

WORKDIR /src
RUN go mod download

WORKDIR /src/monitor
RUN go generate

WORKDIR /src
RUN make cleanmake

EXPOSE 2053

USER root
