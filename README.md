# snippetbox

MySQL работает в Docker, Go-приложение запускается на компьютере.
Нужны Go и запущенный Docker Desktop с Linux containers; устанавливать MySQL
или MySQL CLI на Windows не требуется.

Из корня проекта запустите базу и дождитесь её готовности:

```powershell
docker compose up -d --wait mysql
```

При первом запуске создаются база `snippetbox` с таблицами `snippets`, `users`
и `sessions`, а также отдельная база `test_snippetbox` для существующих тестов.
Таблицы в тестовой базе создаются и удаляются самими тестами.

Запуск приложения с явным TCP-подключением к контейнеру:

```powershell
go run ./cmd/web -dsn "web:1123@tcp(127.0.0.1:3306)/snippetbox?parseTime=true"
```

Приложение использует HTTPS: для запуска также нужны существующие
`tls/cert.pem` и `tls/key.pem`. Адрес: https://localhost:4000.

Проверка проекта при работающем контейнере:

```powershell
go test ./...
```

MySQL доступен только с этого компьютера на `127.0.0.1:3306`.
Если порт занят другим сервером, измените левый порт в `compose.yaml`
и порт в DSN приложения; существующие тесты используют порт 3306.
Пароли в Compose предназначены только для локальной разработки.

Статус, логи и SQL-консоль внутри контейнера:

```powershell
docker compose ps
docker compose logs -f mysql
docker compose exec mysql mysql -u web -p snippetbox
```

Для SQL-консоли введите пароль `1123`.

Остановка:

```powershell
docker compose down
```

Данные сохраняются в Docker volume `mysql_data` между остановками и
пересозданиями контейнера. Скрипты из `docker/mysql/init` выполняются только
при инициализации пустого volume. Изменение SQL-скриптов или переменных
`MYSQL_*` не обновляет уже созданную базу. Команда `docker compose down -v`
удаляет volume вместе со всеми данными базы.
