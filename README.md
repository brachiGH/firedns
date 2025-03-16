## FireDNS
### The new firewall for the modern Internet.

FireDNS protects you from all kinds of security threats, blocks ads and trackers on websites and in apps and provides a safe and supervised Internet for kids — on all devices and on all networks. While allowing you to see what's happening on your devices with in-depth analytics and real-time logs.

## Run FireDNS localy

```bash
docker-compose up --build -d
```

> Change default passwords in the docker-compose.yml:
>
> 	MONGO_INITDB_ROOT_USERNAME: root
>
>	MONGO_INITDB_ROOT_PASSWORD: root
> 

And setup the .env file
```txt
ifname=eth0
MONGO_DB_URI=mongodb://root:root@mongo:27017
```