# ARM v7 app sampler

## Deploy through Docker

```shell
docker-compose up -d
```

## Create record

```shell
curl -X POST http://localhost:18080/records \
  -H "Content-Type: application/json" \
  -d '{"first_name": "John", "age": 30}'
```

## Get records

```shell
curl http://localhost:18080/records
```