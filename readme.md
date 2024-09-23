# MCS (Messaging Campaign System)

# 1. client: 
this is a cmdline program used to create campaign, usage:

``` ./client [filepath] [Campaign_ID] [Message_template] [Scheduled_send_time](using RFC3339)```

kafka brokers and mysql dsn should be configed in the same directory using a json file.(or using etcd/zookeeper etc.)

client read recipients info from local(net) file, assemble with the input campaign params, send msgs into kafka. 

## (Why command line program?) 
1. command line program do not need frontend developer to participate, easy to develop. 
2. as the provided document does not specify who is going to use this program, suggest that user is internal staff, it's ok that kafka/mySQL address in config file is visible to them. But if not, we should not trust the user, we can integrate a http server in the server program to handle the create campaign request.

# 2. server: 
server start a endless server, keep receiving msg from kafka and persist the msg in mysql, monitoring campaign created in db, and send msgs when it's shscheduled time. 

server can be deployed distributed. If mutil instance deployed, should set instance_id and instance_count in env. kafka brokers and mysql dsn should be configed in the same directory using a json file.(or using etcd etc .)

## Build
Build client

```
cd client
go build
```

Build server

```
cd server
go build
```

## config(server && client)
```
{
    "dsn": "root:123456@/?charset=utf8&parseTime=True&loc=Local",
    "borkers": ["127.0.0.1"]
}
```