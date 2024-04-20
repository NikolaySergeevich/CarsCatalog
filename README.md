Примеры запросов к нашему API:

- curl -vvv -H 'Content-Type: application/json' -XPOST 'http://localhost:8111/api/v1/cars' -d '{"regNums": ["X123XX150"]}'

- curl -X GET "https://api.example.com/cars?mark=Toyota&color=red&limit=10&page=1"



Пример запроса во внешний API:

- curl -X GET "https://api.example.com/info?regNum=X123XX150"