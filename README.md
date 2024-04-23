**Примеры запросов к нашему API:**

- curl -vvv -H 'Content-Type: application/json' -XPOST 'http://localhost:8111/api/v1/cars' -d '{"regNums": ["X123XX150"]}'

- curl -vvv -X GET "http://localhost:8111/api/v1/cars?mark=Toyota&color=red&limit=10&page=1"

- curl -vv -XDELETE 'http://localhost:8111/api/v1/cars/123e4567-e89b-12d3-a456-426614174001'

- curl -vvv  -X PUT 'http://localhost:8111/api/v1/cars/123e4567-e89b-12d3-a456-426614174001' -d '{"mark":"NewMark","model":"NewModel","color":"NewColor", "owner":"NewName"}'



**Пример запроса во внешний API:**

- curl -X GET "https://api.example.com/info?regNum=X123XX150"

Когда на наш сервис будет поступать POST запрос с перечнем номеров машин, то наш сервис будет обращаться в стороннему сервису и получать у него всю дополнительную информацию по этим машинам. 