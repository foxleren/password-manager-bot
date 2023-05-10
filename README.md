# password-manager-bot

## Тестовое задание на направление "Разработчик системы кластеризации Tarantool". Телеграм бот для хранения паролей. Выполнил Мохов Сергей Александрович.

## [Ссылка на бота](https://t.me/password_manager_by_foxleren_bot)

![image](https://github.com/foxleren/OS_HW_4/assets/64990498/b49f9517-ce12-4bbf-b761-cc0149dc9f4b)

## Стек:
* Golang
* PostgreSQL

### Реализованы все пункты ТЗ.

## Взаимодействие с ботом

### Команды
/start - начало работы

/set - добавить сервис в библиотеку

/get - получить логин и пароль для сервиса

/del - удалить сервис из библиотеки

/subscribe - подписаться на бота для работы с сервисами

/check_subscription - проверить статус подписки

/unsubscribe - отписаться от бота

### Автоматическое удаление сообщений с паролями пользователя
Файл configs/config.yml содержит следующие параметры
```
bot:
  messageTTLInMinutes: 0 // интервал в минутах, в течение которого сообщение с паролем не подлежит удалению

  messageTTLInHours: 10 // интервал в часах, в течение которого сообщение с паролем не подлежит удалению

  outdatedMessagesSleepIntervalInHours: 10 // интервал, раз в который выполняется автоматическое удаление всех сообщений с паролями пользователей
```

## Запуск приложения

Для локального запуска приложения необходимо сохранить файл ```.env``` в корень репозитория:

.env
```
DB_PASSWORD=<password>
BOT_TOKEN=<token>
```

configs/config.yml
```
port: "8000"

bot:
  messageTTLInMinutes: 0
  messageTTLInHours: 10
  outdatedMessagesSleepIntervalInHours: 10

db:
  host: "db"
  port: "5432"
  username: "postgres"
  dbname: "postgres"
  sslmode: "disable"
```

## Запуск в контейнере
```
> make migrate-up
> docker-compose up --build
```