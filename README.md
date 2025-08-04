# 29.07

При тестировании использовал инструмент Postman.

## Запрос на создание задачи:
```shell
метод:POST
URI:http://localhost:8080/api/v1/tasks
тело запроса:"1"
тело ответа:
{
"1": {
"Link": null,
"Status": "created",
"ArchiveUrl": "",
"ErrorLoad": {}
}
}
код ответа:201

При создании задачи сервис возвращает в ответ все задачи которые находятся в работе.
При попытке создать более 3 задач, сервер вернет сообщение:
"the server is currently busy"
код ответа:400
```

## Запрос на добавление ссылки в задачу:
```shell
метод:POST
URI:http://localhost:8080/api/v1/links/1
тело запроса:"https://www.orimi.com/pdf-test.pdf"
тело ответа:
{
"Link": [
"https://www.orimi.com/pdf-test.pdf"
],
"Status": "created",
"ArchiveUrl": "",
"ErrorLoad": {}
}
код ответа:201

После добавления в задачу 3 ссылки, ее статус меняется.
"Status": "loading"
```

## Запрос на проверку статуса и запуск архивирования.
```shell
метод:GET
URI:http://localhost:8080/api/v1/status/1
тело ответа:
{
"Link":[
"https://www.orimi.com/pdf-test.pdf",
"https://upload.wikimedia.org/wikipedia/commons/3/3f/Fronalpstock_big.jpg",
"https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
],
"Status":"completed",
"ArchiveUrl":"http://localhost:8080/archives/1.zip",
"ErrorLoad":{"https://www.orimi.com/pdf-test.pdf":"unexpected content type: text/html"}
}
код ответа:200
```
## Запрос на получение архива
```shell
метод:GET
URI:http://localhost:8080/archives/1.zip
код ответа:200
```
