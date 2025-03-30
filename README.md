## FireDNS
### The new firewall for the modern Internet.

FireDNS protects you from all kinds of security threats, blocks ads and trackers on websites and in apps and provides a safe and supervised Internet for kids — on all devices and on all networks. While allowing you to see what's happening on your devices with in-depth analytics and real-time logs.

## Setup and run FireDNS localy

### Using Docker

First start with updating the docker-compose.yml and change default passwords

```yml
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: root


    environment:
      ifname: eth0
      MONGO_DB_URI: mongodb://root:root@mongo:27017

```
then run 

```bash
docker-compose up --build -d
```

### On my machine


First, ensure you have the required dependencies installed:

```bash
apt-get install -y build-essential git cmake make \
                       zlib1g-dev libevent-dev m4 \
                       libelf-dev llvm \
                       clang libc6-dev-i386 libpcap-dev \
		       curl tar sudo linux-tools-$(uname -r) \
		       linux-headers-$(uname -r)
```

and install GO from [go.dev](https://go.dev/).

Next, create a `.env`

```text
ifname=eth0
MONGO_DB_URI=mongodb://root:root@localhost:27017
APP_ENV=development
CertFile=The absolute path to public key
KeyFile=The absolute path to private key
```

Note: CertFile and KeyFile are optional. If they are not provided, only the HTTP server will start.

build:

```bash
make cleabmake
```

run:

```bash
docker-compose up mongo -d
sudo ./main
```

### Testing

```bash
sudo -E env MONGO_DB_URI="mongodb://localhost:27017" APP_ENV="production" ifname="eth0" go test ./test -v
```