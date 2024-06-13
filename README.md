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

### Общее

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
    "email": ["Неподдерживаемые символы", "Неверный формат почты],
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
    "email": ["Неподдерживаемые символы", "Неверный формат почты],
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
    "msg": "created",
    "image": "/clothes/2a78df8a-0277-4c72-a2d9-43fb8fef1d2c",
  },
  ...
]</pre></td>
  </tr>
</table>

### GET /clothes/:id

### PUT /clothes/:id

### DELETE /clothes/:id

### GET /user/:id/clothes                                

### GET /types

### GET /subtypes

### GET /styles

### GET /tags
 
### GET /tags/favourite

### GET /photos
 
### GET /photos/:id

### POST /photos

### DELETE /photos/:id

### POST /try-on

### POST /try-on/outfit

### POST /try-on/post
 
### GET /try-on

### GET /try-on/:id

### DELETE /try-on/:id

### PATCH /try-on/:id/rate

### GET /outfits/purposes

### GET /outfits/gen

### POST /outfits

### GET /outfits

### GET /user/:id/outfits

### GET /outfits/:id

### DELETE /outfits/:id

### PUT /outfits/:id

### GET /posts

### GET /posts/recommended

### GET /users/:id/posts

### GET /posts/:id/comments

### POST /posts/:id/comments

### POST /posts/:id/rate

### POST /comments/:id/rate

### PUT /comments/:id

### DELETE /comments/:id

### GET /posts/liked

### GET /posts/subs

### POST /users/:id/sub

### DELETE /users/:id/sub
