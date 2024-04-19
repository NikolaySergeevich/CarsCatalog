curl -vvv -H 'Content-Type: application/json' -XPOST 'http://localhost:8111/api/v1/cars' -d '{"regNums": ["X123XX150"]}'

curl -vv -XGET 'http://localhost:8111/api/v1/cars'