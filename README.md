# Тестовое задание Sycret.ru

Генератор документов. Необходимо заменить в шаблоне документа тэги TEXT на реальные данные. 

# Содержание

1. [Запуск](#Запуск)
2. [Юнит-тесты](#Юнит-тесты)
3. [API](#API)
4. [Реализация](#Реализация)

# Запуск

```
make run
```

or 

```
go run cmd/app/main.go
```

# Юнит-тесты

```
make test
```

or 

```
go test ./... -cover
```

# API

> 1) Тело запроса/ответа - в формате JSON.
> 2) В случае ошибки возвращается необходимый HTTP код.



## GET /api/doc

- Параметры тела запроса:
    - URLTemplate - ссылка на ворд документ,
    - RecordID - номер записи.
- Тело ответа:
    - URLWord - ссылка на сгенерированный ворд документ.

## Пример

Запрос:

```
curl -X GET localhost:8080/api/doc \
-H "Content-Type: application/json" \
-d '{
    "URLTemplate": "https://sycret.ru/service/apigendoc/forma_025u.doc",
    "RecordID": "20"
}'
```

Ответ:

```
{
    "URLWord": "http://res.cloudinary.com/miragost/raw/upload/v1654802194/sycret/2022-06-09%2022-16-33.doc"
}
```

Еще пример запроса

```
http://localhost:8080/api/doc
{
    "URLTemplate": "https://sycret.ru/service/apigendoc/forma_025u.doc",
    "RecordID": "20"
}
```

# Реализация

- REST API
- Конфигурация приложения - библиотека [viper](https://github.com/spf13/viper).
- Юнит-тестирование библиотека [testify](https://github.com/stretchr/testify)
- Сервис для хранения сгенерированных документов - [cloudinary](https://cloudinary.com/)
