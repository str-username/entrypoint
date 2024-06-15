# Entrypoint
The challenge was to implement a simple API:
- there is some set of backends to which the client connects via UDP
- There are clients who, having defined their region, want to know where the backends are located in order to connect 
to them.

Backends do not support balancing in front of them. And the client wants to be sure of a good "latency" to the backend 
(by searching some list of IPs and pinging them). Backends have a "discovery service" that somehow knows about them, 
and at a certain interval sends information about all backends in a particular region in JSON. This implied a small 
write load, but - a large read load.

# Solution
It was decided to make a simple API interface with several endpoints:
- GET api/v1/backends - get backend addresses by region
- POST api/v1/backends - record new backends

The discovery service registers new backends via POST a query /api/v1/backends:
```shell
{
    "region": string,
    "protocol": string,
    "maintenance": bool,
    "allowedVersions": string,
    "serverParameters": {
        "tickRate": string,
        "tickRateValue": {
            "min": int,
            "max": int,
            "default": int
        }
    },
    "serverAddresses": []stirng
}
```
The client, defining his region receives all necessary information by making a request in the 
api/v1/backends?region=europe:
```shell
{
    "maintenance": false,
    "allowedVersions": "1.0.0-rc",
    "serverParameters": {
        "tickRate": "static",
        "tickRateValue": {
            "max": {
                "$numberInt": "120"
            },
            "default": {
                "$numberInt": "80"
            },
            "min": {
                "$numberInt": "60"
            }
        }
    },
    "serverAddresses": [
        "10.10.10.10",
        "10.10.10.11",
        "10.10.10.12",
        "10.10.10.13"
    ],
    "region": "europe",
    "protocol": "UDP"
}
```
The body is not structured (there was no need for that).