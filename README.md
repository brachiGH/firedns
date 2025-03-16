## FireDNS
### The new firewall for the modern Internet.

FireDNS protects you from all kinds of security threats, blocks ads and trackers on websites and in apps and provides a safe and supervised Internet for kids — on all devices and on all networks. While allowing you to see what's happening on your devices with in-depth analytics and real-time logs.

## Setup and run FireDNS localy

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