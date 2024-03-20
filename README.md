# try_on-wardrobe-backend

Для запуска надо:
- Иметь установленными docker compose и go 1.22 (версия go принципиальна)
- docker network create shared-api-network
- Создать .env файл с переменными:
  - POSTGRES.PASSWORD
  - RABBIT.PASSWORD
- make docker
