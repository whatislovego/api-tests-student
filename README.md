# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```env
API_URL=[https://testboard.avito.com/api/v1](https://testboard.avito.com/api/v1)
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
DEBUG=false
```

8. Для последовательного запуска всех тестов (с учётом ограничения частоты одновременных запросов учебного стенда) выполните команду в корне проекта:

```bash
go test -p 1 -count=1 -v ./...
```

9. Отдельный пакет с тестами можно запустить командой:

```bash
go test -count=1 -v ./tests/scenarios/advertisementCreate/
go test -count=1 -v ./tests/scenarios/favorites/
go test -count=1 -v ./tests/scenarios/advertisementSearch/
go test -count=1 -v ./tests/scenarios/myAdvertisement/
```

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```