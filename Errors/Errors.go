package Errors

import "errors"

var (
	ErrLoadConfig = errors.New("Ошибка в загрузке конфигурации!")
	ErrCreateBot  = errors.New("Ошибка в создании бота!")
)
