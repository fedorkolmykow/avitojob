#Тестовое задание Авито в юнит Job

Запуск:

`
docker-compose -f docker-compose.yml up --build -d
`

Изменение баланаса пользователя. Отрицательное значение параметра change снимает средства.

`
curl -d '{"change":200,"comment":"My First","source":"Sberbank"}' -H "Content-Type: application/json" -X PATCH http://localhost:9000/users/1/balance
`

Перевод средств с баланса одного пользователя на баланс другого.

`
curl -d '{"change":200,"comment":"My First","target_id":2}' -H "Content-Type: application/json" -X PATCH http://localhost:9000/users/1/balance/transfer
`

Получение баланса пользователя.

`
curl http://localhost:9000/users/1/balance
`

Получение списка транзакций пользователя. change_sort - сортировка по изменению баланса, change_time - сортировка по времени совершения транзакции.
Page - номер страницы, per_page - количество транзакций на странице.

`
curl -d '{"page":1,"per_page":3,"change_sort":false,"time_sort":true}' -H "Content-Type: application/json" -X POST http://localhost:9000/users/1/transactions
`