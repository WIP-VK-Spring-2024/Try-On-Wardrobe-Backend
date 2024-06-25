# try_on-wardrobe-backend

Для запуска надо:
- Иметь установленными docker compose, make, go 1.22 (версия go принципиальна)
- ```docker network create shared-api-network```
- Создать .env файл аналогичный .env.example
- Иметь поднятый инстанс RabbitMQ на хосте, указанном в .env (подразумевается деплой в паре с [МЛ-сервером](https://github.com/WIP-VK-Spring-2024/Try-On-Wardrobe-ML))
- ```make docker```

## API Documentation

Все запросы, за исключением /hearbeat, /login и /register требуют токена в хэдере X-Session-ID.

Все POST/PUT/PATCH запросы требуют ```Content-Type=application/json```, за исключением тех, где происходит загрузка картинок. Там требуется ```multipart/form```

### GET /heartbeat

Возвращает текущее состояние сервера, проверяет подключения к Postgres, Centrifugo, Redis

#### Ответы
<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">{}</pre></td>
  </tr>
  <tr>
    <td>503</td>
    <td><pre lang="json">
{
  "db": "error message",
  "centrifugo": "error message",
  "redis": "error message"
}</pre></td>
  </tr>
</table>
                        
### GET /users

Возвращает список пользователей

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                                                                                 |
| ----- | ---------- | :----------: | ---------------------------------------------------------------------------------------- |
| name  | string     |      Да      | Регистронезависимая строка для поиска пользователей                                      |
| limit | number     |      Да      | Максимальное кол-во пользователей в ответе                                               |
| since | string     |      Да      | Имя пользователя, начиная с которого искать (невключительно). Используется для пагинации |

#### Ответы
<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Nikita",
    "email": "nikita@mail.ru",
    "avatar": "/avatars/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "gender": "male",
    "privacy": "private"
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /users/subbed

Возвращает список подписок текущего пользователя

#### Ответы
<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Nikita",
    "email": "nikita@mail.ru",
    "avatar": "/avatars/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "gender": "male",
    "privacy": "private"
  },
  ...
]</pre></td>
  </tr>
</table>

### POST /users

Регистрирует пользователя

#### Параметры

| Имя      | Тип данных         | Опциональный | Описание            |
| -------- | ------------------ | :----------: | ------------------- |
| name     | string             |     Нет      | Логин пользователя  |
| email    | string             |     Нет      | Почта пользователя  |
| password | string             |     Нет      | Пароль пользователя |
| gender   | 'male' \| 'female' |     Нет      | Пол пользователя    |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "token": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_name": "Nikita",
  "email": "nikita@mail.ru",
  "gender": "male",
  "privacy": "private"
}</pre></td>
  </tr>
  <tr>
    <td>409</td>
    <td><pre lang="json">
{
  "msg": "Такой пользователь уже существует"
}</pre></td>
  <tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Неправильный формат запроса",
  "errors": {
    "name": ["Неподдерживаемые символы"],
    "email": ["Неподдерживаемые символы", "Неверный формат почты]
  }
}</pre></td>
  </tr>
</table>

### PUT /users/:id

Обновляет данные пользователя

#### Параметры

| Имя     | Тип данных            | Опциональный | Описание              |
| ------- | --------------------- | :----------: | --------------------- |
| gender  | 'male' \| 'female'    |      Да      | Пол пользователя      |
| privacy | 'private' \| 'public' |      Да      | Приватность аккаунта  |
| img     | image                 |      Да      | Аватарка пользователя |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "avatar": "avatars/new-avatar-path"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Нельзя изменять параметры других пользователей"
}</pre></td>
  </tr>
</table>

### POST /login

Вход в аккаунт

#### Параметры

| Имя      | Тип данных | Опциональный | Описание |
| -------- | ---------- | :----------: | -------- |
| name     | string     |     Нет      | Логин    |
| password | string     |     Нет      | Пароль   |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "token": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_name": "Nikita",
  "email": "nikita@mail.ru",
  "gender": "male",
  "privacy": "private"
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Такого пользователя не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Неправильный логин или пароль"
}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Неправильный формат запроса",
  "errors": {
    "name": ["Неподдерживаемые символы"],
    "email": ["Неподдерживаемые символы", "Неверный формат почты]
  }
}</pre></td>
</tr>
</table>

### POST /renew

Выдает новый JWT-токен

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "token": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_name": "Nikita",
  "email": "nikita@mail.ru",
  "gender": "male",
  "privacy": "private"
}</pre></td>
  </tr>
  <tr>
    <td>401</td>
    <td><pre lang="json">
{
  "msg": "Токен истёк или отсутствует в запросе"
}</pre></td>
  </tr>
</table>

### GET /clothes

Список вещей текущего пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Футболка в зал",
    "tryonable": "true",
    "style_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "type_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "subtype_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "image": "/clothes/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "seasons": ["summer", "winter", "spring", "autumn"],
    "tags": ["Для спорта", "Лёгкое"]
  },
  ...
]</pre></td>
  </tr>
</table>

### POST /clothes

Загрузка новой вещи

#### Параметры

| Имя | Тип данных | Опциональный | Описание        |
| --- | ---------- | :----------: | --------------- |
| img | image      |     Нет      | Фотография вещи |

#### Ответы

При получении ответа 200 путь к обрезанному фото и определённые категории, стиль, возвращаются через канал Centrifugo, указанный в config.json (по умолчанию - "processing#\<user_id\>")

Пример ответа из Centrifugo:
```json
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "msg": "processed",
  "image": "cut/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "tryonable": "true",
  "classification": {
    "type": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "subtype": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "style": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "seasons": ["autumn", "winter"],
    "tags": ["Тёплое", "С мехом"]
  }
}
```

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "msg": "created",
  "image": "/clothes/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}
</pre></td>
  </tr>
</table>

### GET /clothes/:id

Получение вещи

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "created_at": "2024-04-13T12:00:07.458144Z",
  "updated_at": "2024-04-13T12:00:07.458144Z",
  "name": "Футболка в зал",
  "tryonable": "true",
  "style_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "type_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "subtype_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "image": "/clothes/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "seasons": ["summer", "winter", "spring", "autumn"],
  "tags": ["Для спорта", "Лёгкое"]
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
</table>

### PUT /clothes/:id

Обновление полей вещи

#### Параметры

| Имя     | Тип данных | Опциональный | Описание          |
| ------- | ---------- | :----------: | ----------------- |
| img     | image      |      Да      | Фотография вещи   |
| name    | string     |      Да      | Название вещи     |
| tags    | string[]   |      Да      | Теги вещи         |
| seasons | string[]   |      Да      | Времена года вещи |
| style   | uuid       |      Да      | Стиль вещи        |
| type    | uuid       |      Да      | Категория вещи    |
| subtype | uuid       |      Да      | Подкатегория вещи |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Изменять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### DELETE /clothes/:id

Удаление вещи

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Удалять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>                          

### GET /types

Получение списка категорий

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Верх",
    "tryonable": "true",
    "subtypes": [
      {
        "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
        "created_at": "2024-04-13T12:00:07.458144Z",
        "updated_at": "2024-04-13T12:00:07.458144Z",
        "name": "Футболка",
        "type_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
      },
      ...
    ]
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /subtypes

Получение списка подкатегорий

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Футболка",
    "type_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /styles

Получение списка стилей

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Спортивный"
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /tags

Получение списка тегов

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                            |
| ----- | ---------- | :----------: | ----------------------------------- |
| limit | int        |      Да      | Максимальное кол-во тегов  в ответе |
| from  | int        |      Да      | Offset для пагинации                |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "В зал"
  },
  ...
]</pre></td>
  </tr>
</table>
 
### GET /tags/favourite

Получение списка наиболее используемых тегов пользователя

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                            |
| ----- | ---------- | :----------: | ----------------------------------- |
| limit | int        |      Да      | Максимальное кол-во тегов  в ответе |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "В зал"
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /photos

Получение списка загруженных фото пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
  },
  ...
]</pre></td>
  </tr>
</table>
 
### GET /photos/:id

Получение фото пользователя по id

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "created_at": "2024-04-13T12:00:07.458144Z",
  "updated_at": "2024-04-13T12:00:07.458144Z",
  "image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
</table>

### POST /photos

Загрузка фото пользователя

#### Параметры

| Имя | Тип данных | Опциональный | Описание                |
| --- | ---------- | :----------: | ----------------------- |
| img | image      |     Нет      | Фотография пользователя |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}</pre></td>
  </tr>
</table>

### DELETE /photos/:id

Удаление фото пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Удалять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### POST /try-on

Примерка одной или нескольких вещей

#### Параметры

| Имя           | Тип данных | Опциональный | Описание                      |
| ------------- | ---------- | :----------: | ----------------------------- |
| user_image_id | uuid       |     Нет      | ID фотографии пользователя    |
| clothes_id    | uuid[]     |     Нет      | Массив ID одежды для примерки |

#### Ответы

При получении ответа 200 результат примерки возвращается через канал Centrifugo, указанный в config.json (по умолчанию - "try-on#\<user_id\>")

Пример ответа из Centrifugo:
```json
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_image_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "clothes_id": ["2a78df8a-0277-4c72-a2d9-43fb8fef1d2c", ...],
  "image": "/try-on/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}
``` 

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Невозможно примерить запрашиваемую одежду"
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Указанной фотографии/вещей не существует"
}</pre></td>
  </tr>
  <tr>
    <td>503</td>
    <td><pre lang="json">
{
  "msg": "Сервис примерки недоступен"
}</pre></td>
  </tr>
</table>

### POST /try-on/outfit

Примерка образа

#### Параметры

| Имя           | Тип данных | Опциональный | Описание                   |
| ------------- | ---------- | :----------: | -------------------------- |
| user_image_id | uuid       |     Нет      | ID фотографии пользователя |
| outfit_id     | uuid       |     Нет      | ID образа для примерки     |

#### Ответы

При получении ответа 200 результат примерки возвращается через канал Centrifugo, указанный в config.json (по умолчанию - "try-on#\<user_id\>")

Пример ответа из Centrifugo:
```json
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_image_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "clothes_id": ["2a78df8a-0277-4c72-a2d9-43fb8fef1d2c", ...],
  "image": "/try-on/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}
```

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Невозможно примерить запрашиваемый образ"
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Указанной фотографии/образа не существует"
}</pre></td>
  </tr>
  <tr>
    <td>503</td>
    <td><pre lang="json">
{
  "msg": "Сервис примерки недоступен"
}</pre></td>
  </tr>
</table>

### POST /try-on/post

Примерка образа из поста

#### Параметры

| Имя           | Тип данных | Опциональный | Описание                   |
| ------------- | ---------- | :----------: | -------------------------- |
| user_image_id | uuid       |     Нет      | ID фотографии пользователя |
| post_id       | uuid       |     Нет      | ID **образа** для примерки |

#### Ответы

При получении ответа 200 результат примерки возвращается через канал Centrifugo, указанный в config.json (по умолчанию - "try-on#\<user_id\>")

Пример ответа из Centrifugo:
```json
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_image_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "clothes_id": ["2a78df8a-0277-4c72-a2d9-43fb8fef1d2c", ...],
  "image": "/try-on/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}
```

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Невозможно примерить запрашиваемый образ"
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Указанной фотографии/образа не существует"
}</pre></td>
  </tr>
  <tr>
    <td>503</td>
    <td><pre lang="json">
{
  "msg": "Сервис примерки недоступен"
}</pre></td>
  </tr>
</table>
 
### GET /try-on

Получение списка результатов примерок текущего пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "clothes_id": ["2a78df8a-0277-4c72-a2d9-43fb8fef1d2c", ...],
    "image": "/try-on/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /try-on/:id

Получение результата примерки по ID

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "user_image_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "clothes_id": ["2a78df8a-0277-4c72-a2d9-43fb8fef1d2c", ...],
  "image": "/try-on/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
</table>

### DELETE /try-on/:id

Удаление результата примерки

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Удалять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### PATCH /try-on/:id/rate

Оценка результата примерки (планировалось использовать оценки для сбора статистики качества модели)

#### Параметры

| Имя    | Тип данных   | Опциональный | Описание                   |
| ------ | ------------ | :----------: | -------------------------- |
| rating | -1 \| 0 \| 1 |     Нет      | Оценка результата примерки |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Изменять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### GET /outfits/purposes

Получение списка назначений образов. Используются в генерации образов

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Для активного отдыха"
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /outfits/gen

Получение сгенерированных образов

#### Параметры

| Имя         | Тип данных               | Опциональный | Описание                                                         |
| ----------- | ------------------------ | :----------: | ---------------------------------------------------------------- |
| amount      | int                      |      Да      | Максимальное число образов в ответе                              |
| use_weather | bool                     |      Да      | Учитывать ли погоду в генерации                                  |
| prompt      | string                   |      Да      | Описание желаемого образа                                        |
| purposes    | string[]                 |      Да      | Список назначений образа                                         |
| pos         | {lat: float, lon: float} |      Да      | Координаты пользователя (будет использован IP, если не переданы) |

#### Ответы


При получении ответа 200 результат примерки возвращается через канал Centrifugo, указанный в config.json (по умолчанию - "outfit-gen#\<user_id\>")

Пример ответа из Centrifugo:
```json
{
  "outfits": [
    {
      "clothes": [
        {
          "clothes_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
        },
        ...
      ]
    },
    ...
  ]
}
```

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Не хватает одежды для генерации"
}</pre></td>
  </tr>
  <tr>
    <td>503</td>
    <td><pre lang="json">
{
  "msg": "Сервис генерации образов недоступен"
}</pre></td>
  </tr>
</table>

### POST /outfits

Создание образа

#### Параметры

| Имя        | Тип данных         | Опциональный | Описание           |
| ---------- | ------------------ | :----------: | ------------------ |
| transforms | map[uuid]Transform |     Нет      | Одежда образа      |
| img        | image              |     Нет      | Изображение образа |

```json
Transform = {
  "x": 42,
  "y": 42,
  "width": 42,
  "height": 42,
  "angle": 42,
  "scale": 42,
  "zindex": 0
}
```

Transform не используется на бэкенде, служит для отрисовки одежды в редакторе


#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "created_at": "2024-04-13T12:00:07.458144Z",
  "updated_at": "2024-04-13T12:00:07.458144Z",
}</pre></td>
  </tr>
</table>

### GET /outfits

Получение образов текущего пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "name": "Образ для активного отдыха",
    "image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "tags": ["В поход", "Для спорта"],
    "privacy": "private",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "transforms": {
      "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c": {
        "x": 42,
        "y": 42,
        "width": 42,
        "height": 42,
        "angle": 42,
        "scale": 42,
        "zindex": 0
      },
      ...
    }
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /outfits/:id

Получение образа по ID

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "created_at": "2024-04-13T12:00:07.458144Z",
  "updated_at": "2024-04-13T12:00:07.458144Z",
  "name": "Образ для активного отдыха",
  "image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "tags": ["В поход", "Для спорта"],
  "privacy": "private",
  "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  "transforms": {
    "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c": {
      "x": 42,
      "y": 42,
      "width": 42,
      "height": 42,
      "angle": 42,
      "scale": 42,
      "zindex": 0
    },
    ...
  }
}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
</table>

### DELETE /outfits/:id

Удаление образа

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Удалять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### PUT /outfits/:id

Обновление образа

#### Параметры

| Имя        | Тип данных         | Опциональный | Описание           |
| ---------- | ------------------ | :----------: | ------------------ |
| transforms | map[uuid]Transform |      Да      | Одежда образа      |
| img        | image              |      Да      | Изображение образа |
| tags       | string[]           |      Да      | Теги образа        |
| name       | string             |      Да      | Название образа    |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Изменять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### GET /posts

Получение новых постов. Каждый пост - это образ пользователя + примерка (при наличии). Любой образ попадает в ленту, если он публичный.

#### Параметры

| Имя     | Тип данных | Опциональный | Описание                                   |
| ------- | ---------- | :----------: | ------------------------------------------ |
| limit   | integer    |      Да      | Макс. кол-во постов в ответе               |
| since   | string     |      Да      | Дата, начиная с которой возвращаются посты |
| genders | string[]   |      Да      | Пол пользователей, чьи посты возвращаем    |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image": "avatar/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_name": "Nikita",
    "is_subbed": true,
    "rating": 42,
    "user_rating": 1, // оценка текущего пользователя. 0, если не поставил лайк, 1 - если поставил
    "tryonable": true
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /posts/recommended

Рекомендованные посты. При первом запросе всегда возвращается пустой массив, т.к. отправляется запрос в рекомендательную систему. Далее ответы будут возвращать посты до тех пор, пока полученные от рекомендательной системы посты не кончатся. Следующий после этого запрос также вернёт пустой массив.

| Имя           | Тип данных | Опциональный | Описание                                                      |
| ------------- | ---------- | :----------: | ------------------------------------------------------------- |
| limit         | integer    |      Да      | Макс. кол-во постов в ответе                                  |
| sample_amount | integer    |      Да      | Кол-во постов, которое запрашиваем у рекомендательной системы |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image": "avatar/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_name": "Nikita",
    "is_subbed": true,
    "rating": 42,
    "user_rating": 1, // оценка текущего пользователя. 0, если не поставил лайк, 1 - если поставил
    "tryonable": true
  },
  ...
]</pre></td>
  </tr>
  <tr>
    <td>503</td>
    <td><pre lang="json">
{
  "msg": "Рекомендательная система недоступна"
}</pre></td>
  </tr>
</table>

### GET /users/:id/posts

Получение постов пользователя

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                                   |
| ----- | ---------- | :----------: | ------------------------------------------ |
| limit | integer    |      Да      | Макс. кол-во постов в ответе               |
| since | string     |      Да      | Дата, начиная с которой возвращаются посты |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image": "avatar/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_name": "Nikita",
    "is_subbed": true,
    "rating": 42,
    "user_rating": 1, // оценка текущего пользователя. 0, если не поставил лайк, 1 - если поставил
    "tryonable": true
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /posts/liked

Получение понравившихся постов

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                                   |
| ----- | ---------- | :----------: | ------------------------------------------ |
| limit | integer    |      Да      | Макс. кол-во постов в ответе               |
| since | string     |      Да      | Дата, начиная с которой возвращаются посты |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image": "avatar/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_name": "Nikita",
    "is_subbed": true,
    "rating": 42,
    "user_rating": 1, // оценка текущего пользователя. 0, если не поставил лайк, 1 - если поставил
    "tryonable": true
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /posts/subs

Получение постов от подписок

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                                   |
| ----- | ---------- | :----------: | ------------------------------------------ |
| limit | integer    |      Да      | Макс. кол-во постов в ответе               |
| since | string     |      Да      | Дата, начиная с которой возвращаются посты |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image": "avatar/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_name": "Nikita",
    "is_subbed": true,
    "rating": 42,
    "user_rating": 1, // оценка текущего пользователя. 0, если не поставил лайк, 1 - если поставил
    "tryonable": true
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /posts/:id/comments

Получение комментариев поста

#### Параметры

| Имя   | Тип данных | Опциональный | Описание                                      |
| ----- | ---------- | :----------: | --------------------------------------------- |
| limit | integer    |      Да      | Макс. кол-во комментов в ответе               |
| since | string     |      Да      | Дата, начиная с которой возвращаются комменты |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
[
  {
    "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "created_at": "2024-04-13T12:00:07.458144Z",
    "updated_at": "2024-04-13T12:00:07.458144Z",
    "outfit_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "outfit_image": "outfits/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "try_on_image": "photos/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_image": "avatar/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
    "user_name": "Nikita",
    "body": "Отличный образ!",
    "level": 0, // для реализации комментов с ответами
    "parent_id": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c", // для реализации комментов с ответами
    "rating": 42,
    "user_rating": 1 // оценка текущего пользователя. 0, если не поставил лайк, 1 - если поставил
  },
  ...
]</pre></td>
  </tr>
</table>

### POST /posts/:id/comments

Создание комментария к посту

#### Параметры

| Имя       | Тип данных | Опциональный | Описание                                                         |
| --------- | ---------- | :----------: | ---------------------------------------------------------------- |
| body      | string     |     Нет      | Тело комментария                                                 |
| parent_id | uuid       |      Да      | ID комментария, на который хотим ответить **(не тестировалось)** |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{
  "uuid": "2a78df8a-0277-4c72-a2d9-43fb8fef1d2c"
}</pre></td>
  </tr>
    <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Такого поста не существует"
}</pre></td>
  </tr>
</table>

### POST /posts/:id/rate

Оценка поста

#### Параметры

| Имя    | Тип данных   | Опциональный | Описание     |
| ------ | ------------ | :----------: | ------------ |
| rating | -1 \| 0 \| 1 |     Нет      | Оценка поста |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
    <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Такого поста не существует"
}</pre></td>
  </tr>
</table>

### POST /comments/:id/rate

Оценка комментария

#### Параметры

| Имя    | Тип данных   | Опциональный | Описание           |
| ------ | ------------ | :----------: | ------------------ |
| rating | -1 \| 0 \| 1 |     Нет      | Оценка комментария |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
    <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Такого комментария не существует"
}</pre></td>
  </tr>
</table>

### PUT /comments/:id

Изменение комментария

#### Параметры

| Имя  | Тип данных | Опциональный | Описание               |
| ---- | ---------- | :----------: | ---------------------- |
| body | string     |     Нет      | Новое тело комментария |

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
    <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Такого комментария не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Изменять этот ресурс может только владелец"
}</pre></td>
  </tr>
</table>

### DELETE /comments/:id

Удаление комментария

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Запрашиваемого ресурса не существует"
}</pre></td>
  </tr>
  <tr>
    <td>403</td>
    <td><pre lang="json">
{
  "msg": "Удалять этот ресурс может только его владелец"
}</pre></td>
  </tr>
</table>

### POST /users/:id/sub

Подписка на пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>404</td>
    <td><pre lang="json">
{
  "msg": "Такого пользователя не существует"
}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Нельзя подписываться на себя"
}</pre></td>
  </tr>
  <tr>
    <td>409</td>
    <td><pre lang="json">
{
  "msg": "Вы уже подписаны на этого пользователя"
}</pre></td>
  </tr>
</table>

### DELETE /users/:id/sub

Отписка от пользователя

#### Ответы

<table>
  <tr>
    <th>Код</th>
    <th>Пример</th>
  </tr>
  <tr>
    <td>200</td>
    <td><pre lang="json">
{}</pre></td>
  </tr>
  <tr>
    <td>400</td>
    <td><pre lang="json">
{
  "msg": "Нельзя отписаться от себя"
}</pre></td>
  </tr>
</table>
