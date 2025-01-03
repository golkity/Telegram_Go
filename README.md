# Телеграмм бот для анализа акций

![version](https://shields.microej.com/github/go-mod/go-version/golkity/Calc?style=for-the-badge)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

![img_title](intro-bot.png)

>[!IMPORTANT]
> Данный телеграм-бот позволяет управлять своими инвестициями, а также прогназировать ротс и падания цены портфеля.

>[!TIP]
> ### Функционал:
> - **Добавление акций**: Укажите навзание акции, цену и процент, который будет начисляться каждый месяц.
> - **Просмотр портфеля**: Показывает кол-во акций в портфеле, общую стоимость и какие акции есть в нем.
> - **Продажа акций**: Укажите название акции, которую хотите продать, и цену покупки/продажи.
> - **Анализ портфеля**: Выводит прогноз портфеля с графиком роста акции.

## Установка зависимостей
```shell
go mod tidy
```

## Запуск

```shell
go run ./cmd/main.go
```

>[!WARNING]
> Чтобы бот работал, вам надо создать в папке config config.json с такой структурой:
> ```json
>  {
>   "token" : "ТОКЕН"
>  }
> ```