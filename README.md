# redis-like
attempt to implement a Redis-like in-memory cache

## server API documentation
To get full API documenation for server package
just clone this repo and run 'godoc redis-like/server'
command in your favorite terminal

## Telnet-like API documentation
- SET key value
Set the string value of a key
- GET key
Get the string value of a key
- LSET key value
Set a list value of a key
  - Example:
    ```
    server> LSET lst "I'm a list" with "hello world inside"
    ```
- LGET key
Get the list value of a key
  - Example:
    ```
    server> LSET lst "I'm a list" with "hello world inside"
    server> LGET lst
    ["I'm a list" with "hello world inside"]
    ```
- LGETIT key index
Get a value by list index of a key
  - Example:
    ```
    server> LSET lst "I'm a list" with "hello world inside"
    server> LGETIT lst 2
    "hello world inside"
    ```
- LUPDATE key index value
Update a value in list index of a key
- HSET key value
Set the dict value of a key
  - Example:
    ```
    server> HSET dict a dict with "hello world"
    ```
- HGET key
Get the dict value of a key
  - Example:
    ```
    server> HSET dict a dict with "hello world"
    server> HGET dict
    map[a:dict with:"hello world"]
    ```
- HGETVAL outerKey innerKey
Get a value from a dict by innerKey of a outerKey
  - Example:
    ```
    server> HSET dict a dict with "hello world"
    server> HGETVAL dict with
    "hello world"
    ```
- HUPDATE outerKey innerKey value
Update a value of a innerKey of dict outerKey
Or create a new innerKey: value pair if innerKey
doesn't exists
- KEYS
Get all keys from current database
- SELECT dbID
Switch to dbID database
- TTL key
Get ttl of a key
- EXPIRE key seconds
Set a key's time to live in seconds
- EXPIREAT key timestamp
Set the expiration for a key as a UNIX timestamp
- PERSIST key
Remove the expiration from a key

## Deployment
- clone this repo
- cd to redis-like/memcache-server
- go build
- ./memcache-server

## How to connect to the server
You may user netcat, telnet or another simular solution
- nc SERVER_HOST SERVER_PORT
- Example:
    ```
    >itsyplenkov$ nc localhost 8000
    >localhost:8000[0]
    >localhost:8000[0] set str "hello world"
    >OK
    >localhost:8000[0]
    ```
