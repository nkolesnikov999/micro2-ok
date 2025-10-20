# micro2-OK

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/nkolesnikov999/03d42dc72a3bdf072f3b1d011eea206e/raw/coverage.json)

Этот репозиторий содержит проект из курса [Микросервисы, как в BigTech 2.0](https://olezhek28.courses/microservices) от [Олега Козырева](http://t.me/olezhek28go).

Для того чтобы вызывать команды из Taskfile, необходимо установить Taskfile CLI:

```bash
brew install go-task
```

## CI/CD

Проект использует GitHub Actions для непрерывной интеграции и доставки. Основные workflow:

- **CI** (`.github/workflows/ci.yml`) - проверяет код при каждом push и pull request
  - Линтинг кода
  - Проверка безопасности
  - Выполняется автоматическое извлечение версий из Taskfile.yml
