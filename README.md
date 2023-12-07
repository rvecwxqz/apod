# Сервис астрологических прогнозов для молодежи
![meme](https://sun9-61.userapi.com/impg/Ik6Uu0Bs72AykhG4EuL9Z3Q_memk_VcGLc6bfw/VGVK9PDxMv0.jpg?size=768x402&quality=96&sign=1ebfbb3e7beb4f3916cace9c52593747&type=album)

Сервис использует асинхронный воркер для получения метаданных и изображения [APOD](https://apod.nasa.gov/apod/astropix.html).

# Технологии
Для хранения метаданных используется PostgreSQL, в качестве бинарного хранилища для картинок используется MinIO. Для конфигурирования использованы переменные окружения.
# HTTP Server
Доступны две HTTP ручки:
- POST /get_all - Получение всех записей;
- POST /get - Получение записи за выбранный день.
Тело для get:
```json
{
    "date": "2023-12-07"
}
```
# Prerequisites 
- [Docker](https://www.docker.com/)
# Установка
- Склонируйте репозиторий;
- Используйте make run.
