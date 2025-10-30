<h1 align="center">CommentTree — древовидные комментарии с навигацией и поиском</h1>

> - Представляет собой базовую систему комментариев 
> - Каждый комментарий может иметь родительский комментарий и дочерние
---

```
git clone https://github.com/child6yo/wbtech-l3-comment-tree

docker compose up 
```
- практически вся система конфигурируема через .env
- UI будет доступен по адресу localhost:80

## API

### POST /comments — создание комментария (с указанием родительского);

#### Request
```
curl -X POST 'localhost:8080/comments' \
--header 'Content-Type: application/json' \
--data '{
    "content": "1234",
    "id": 1 - 
}'
```

- content - содержание комментария
- id - ID родительского комментария, опционально

#### Response 
*200 OK*
{
    "id": 2 
}

id - ID созданного комментария

*400 Bad Request/500 Internal Server Error*
```
    "error": "some error"
```

### GET /comments?parent={id} — получение комментария и всех вложенных;

#### Query Param

- parent - ID родительского комментария. Опциональный параметр. Если передать parent=0, вернется тоже самое, что и если параметр не передавать вовсе, т.е. все комментарии.

#### Request

```
curl -X GET 'localhost:8080/comments?parent={id}'
```

#### Response

*200 OK*
```
[
    {
        "id": 1,
        "content": "1",
        "answer_at": 0,
        "children": [
            {
                "id": 2,
                "content": "12",
                "answer_at": 1,
                "children": [
                    {
                        "id": 4,
                        "content": "123",
                        "answer_at": 2,
                        "children": null
                    }
                ]
            },
            {
                "id": 3,
                "content": "122",
                "answer_at": 1,
                "children": null
            }
        ]
    }
    ...
]
```

*500 Internal Server Error*
```
    "error": "some error"
```

### DELETE /comments/{id} — удаление комментария и всех вложенных под ним.

#### Request

```
curl -X DELETE 'localhost:8080/comments/{id}'
```

#### Response

*200 OK*

*500 Internal Server Error*
```
    "error": "some error"
```
