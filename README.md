## ЗАПУСК ПРИЛОЖЕНИЯ
docker build -t my-go-server .

docker run --rm -p 8080:8080 my-go-server

Пинг: сurl -i localhost:8080/api/ping

Далее как описано в задании, например, добавление тендера

curl -i -X POST --data '{"name": "БДоставка", "description": "Доставить товары из Казани в Москву", "serviceType": "Delivery", "creatorUsername": "user4"}' localhost:8080/api/tenders/new

и т.д.

Сервер на старте ждет переменных среды POSTGRES_*


------------------------------------------------------
-------------------------------------------------------

## СТАРЫЙ КОММЕНТАРИЙ
## Структура проекта
В данном проекте находится типой пример для сборки приложения в докере из находящящегося в проекте Dockerfile. Пример на Gradle используется исключительно в качестве примера, вы можете переписать проект как вам хочется, главное, что бы Dockerfile находился в корне проекта и приложение отвечало по порту 8080. Других требований нет.

## Задание
В папке "задание" размещена задча.

## Сбор и развертывание приложения
Приложение должно отвечать по порту `8080` (жестко задано в настройках деплоя). После деплоя оно будет доступно по адресу: `https://<имя_проекта>-<уникальный_идентификатор_группы_группы>.avito2024.codenrock.com`

Пример: Для кода из репозитория `/avito2024/cnrprod-team-27437/task1` сформируются домен

```
task1-5447.avito2024.codenrock.com
```

**Для удобства домен указывается в логе сборки**

Логи сборки проекта находятся на вкладке **CI/CD** -> **Jobs**.

Ссылка на собранный проект находится на вкладке **Deployments** -> **Environment**. Вы можете сразу открыть URL по кнопке "Open".

## Доступ к сервисам

### Kubernetes
На вашу команду выделен kubernetes namespace. Для подключения к нему используйте утилиту `kubectl` и `*.kube.config` файл, который вам выдадут организаторы.

Состояние namespace, работающие pods и логи приложений можно посмотреть по адресу [https://dashboard.avito2024.codenrock.com/](https://dashboard.avito2024.codenrock.com/). Для открытия дашборда необходимо выбрать авторизацию через Kubeconfig и указать путь до выданного вам `*.kube.config` файла

