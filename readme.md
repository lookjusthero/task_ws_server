задание, разместить на гитхаб, время 2-4 часа
написать на Golang сервер который является websocket сервером
при запросе на сервер, должен возвращаться ответ в виде JSON - рандомное число *big.Int
сервер не должен позволяет с одного IP адреса устанавливать более 1 соединения
сервер должен гарантировать что число всегда будет уникальное, появившись один раз - больше не должно появляться
описать в readme как запустить и использовать сервер
unit тесты и не docker - не обязательно, если в коде будут

```
git clone https://github.com/lookjusthero/task_ws_server

cd task_ws_server

go mod tidy

go run main.go
```

Сервер запуститься на :8080 порту

localhost:8080/ws - прием вебсокет подключений

localhost:8080/ - список команд для консоли браузера, для создания подключений и спама сообщениями

* Ручка localhost:8080/ws держит только одно соединение для каждого ip адреса
* В ответ на в полученое сообщение localhost:8080/ws возвращается рандомное число из диапазона [1, 2**130)
** ответ является строкой, не находясь в объекте т.к. это удовлетворяет JSON формату

