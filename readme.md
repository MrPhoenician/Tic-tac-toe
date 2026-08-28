## FAQ for game

### 1. Регистрация
```http
POST http://localhost:8080/signup
Content-Type: application/json

{
  "login": "player1",
  "password": "password"
}
```

### 2. Авторизация и получение JWT

```http
POST http://localhost:8080/signin
Content-Type: application/json

{
  "login": "player1",
  "password": "password123"
}
```

Для всех защищённых endpoint'ов нужно передавать заголовок:
```http
Authorization: Bearer ACCESS_TOKEN
```

### 3. Обновить access token, если сессия истекла
Этот endpoint доступен без access token, нужен только refresh token.

```http
POST http://localhost:8080/auth/access
Content-Type: application/json

{
  "refreshToken": "REFRESH_TOKEN"
}
```

### 4. Обновить refresh token
```http
POST http://localhost:8080/auth/refresh
Content-Type: application/json

{
  "refreshToken": "REFRESH_TOKEN"
}
```

### 5. Получить текущего пользователя
```http
GET http://localhost:8080/me
Authorization: Bearer ACCESS_TOKEN
```

### 6. Создать игру с компьютером
```http
POST http://localhost:8080/games
Authorization: Bearer ACCESS_TOKEN
Content-Type: application/json

{
  "vs_computer": true
}
```

### 7. Создать игру с другим игроком
Первый игрок создаёт игру:

```http
POST http://localhost:8080/games
Authorization: Bearer ACCESS_TOKEN_PLAYER1
Content-Type: application/json

{
  "vs_computer": false
}
```

Второй игрок присоединяется к игре:

```http
POST http://localhost:8080/game/{game_id}/join
Authorization: Bearer ACCESS_TOKEN_PLAYER2
```

### 8. Получить список доступных игр
```http
GET http://localhost:8080/games
Authorization: Bearer ACCESS_TOKEN
```

### 9. Получить игру по id
```http
GET http://localhost:8080/game/{game_id}
Authorization: Bearer ACCESS_TOKEN
```

### 10. Сделать ход
```http
POST http://localhost:8080/game/{game_id}
Authorization: Bearer ACCESS_TOKEN
Content-Type: application/json

{
  "id": "{game_id}",
  "board": [
    [1, 0, 0],
    [0, 0, 0],
    [0, 0, 0]
  ]
}
```

Обозначения на доске:
- `1` — ход X
- `-1` — ход O
- `0` — пустая клетка

### 11. История завершённых игр пользователя
```http
GET http://localhost:8080/games/completed
Authorization: Bearer ACCESS_TOKEN
```

Возвращаются завершённые игры пользователя:
- победы этого пользователя;
- ничьи.

### 12. Таблица лидеров
```http
GET http://localhost:8080/leaderboard?n=10
Authorization: Bearer ACCESS_TOKEN
```

Где:
- `n` — количество лучших игроков.

В ответе возвращаются:
- UUID пользователя;
- логин;
- соотношение побед.

### 13. Быстрый сценарий игры
1. Зарегистрировать пользователя через `/signup`.
2. Выполнить вход через `/signin`.
3. Сохранить `accessToken` и `refreshToken`.
4. Для защищённых endpoint'ов передавать заголовок:
   ```http
   Authorization: Bearer ACCESS_TOKEN
   ```
5. Если access token истёк — вызвать `/auth/access`.
6. Если нужно перевыпустить refresh token — вызвать `/auth/refresh`.

## Поднять PostgreSQL и развернуть БД в один клик
```bash
bash scripts/setup_postgres.sh
```

Скрипт автоматически:
- скачает и запустит контейнер `postgres`, если его ещё нет;
- запустит уже существующий контейнер, если он был остановлен;
- дождётся готовности PostgreSQL;
- создаст базу `tic_tac_toe`, если её ещё нет;
- применит `schema.sql`.

Можно переопределить параметры через переменные окружения:

```bash
CONTAINER_NAME=my-postgres POSTGRES_PASSWORD=secret POSTGRES_PORT=5433 bash scripts/setup_postgres.sh
```

## Ручная настройка базы
### Выкачать образ PostgreSQL и установить пароль `1`
```bash
docker run -p5432:5432 --name some-postgres -e POSTGRES_PASSWORD=1 -d postgres
```

### Подключиться к PostgreSQL как суперпользователь
```bash
psql -h localhost -U postgres -d postgres
```

### Создать новую базу данных
```sql
CREATE DATABASE tic_tac_toe;
```

### Добавить таблицы в базу
```bash
psql -h localhost -U postgres -d tic_tac_toe < schema.sql
```

Пояснение:
- `-U postgres` — пользователь базы данных;
- `-d tic_tac_toe` — база, в которую выполняются команды;
- `< schema.sql` — выполнение SQL из файла.

### Удалить старые таблицы при необходимости
```sql
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;
```