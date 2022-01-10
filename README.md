# Incolore 🎨

## Stack

- Go http api
- [boltdb](https://github.com/etcd-io/bbolt) database
- [nanoid](https://github.com/matoous/go-nanoid) for id generation

Go HTTP and boltdb should scale well if you have money.

## How

- POST multipart/form-data f=file /

## Configuration

- `INCOLORE_DB` (default=bolt://data.bolt) db path
- `INCOLORE_HOSTNAME` (default=localhost:5377) hostname
- `INCOLORE_ID_LENGTH` (default=12) nanoid length (see [collision calculator](https://zelark.github.io/nano-id-cc/))
- `INCOLORE_ID_ALPHABET` (default=0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ) nanoid alphabet)
- `INCOLORE_PORT` (default=5376)
- `INCOLORE_DIRECTORY` (default=upload)
- `INCOLORE_MAX_SIZE` (default=10000000)

## Docker

```
docker pull soyuka/incolore
docker run -d --name incolore -p 5377:5376 -e INCOLORE_DB=bolt:///bolt/incolore.bolt -e INCOLORE_DIRECTORY="/bolt/incolore/uploads" -e INCOLORE_HOSTNAME=https://incolo.re -v "/home/soyuka/incolore:/bolt" soyuka/incolore
```
