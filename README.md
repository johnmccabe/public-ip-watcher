# public-ip-watcher

Keeps track of changes to your public IP by polling Cloudflares icanhazip endpoint.

Exposes two endpoints.

```
$ http localhost:8080/latest

{
    "Addr": "11.22.33.01",
    "Created": "2022-01-01 01:02:03"
}
```

```
$ http localhost:8080/history

[   
    {
        "Addr": "44.55.66.01",
        "Created": "2022-02-13 04:05:06"
    },
    {
        "Addr": "11.22.33.01",
        "Created": "2022-01-01 01:02:03"
    }
]
```

## Building

Build a local binary with;

```
make
```

And a container image with (update the container tag as needed in the Makefile);

```
make image
```

## Running

Local binary;

```
./ipwatcher
```

In docker;

```
docker run --rm --name ipwatcher -p 8080:8080 -v db:/db johnmccabe/ipwatcher:dev
```
