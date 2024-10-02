package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type logError struct {
	Message string
	Time    time.Time
}

func (e *logError) Error() string {
	return fmt.Sprintf("[%v] %s", e.Time, e.Message)
}

func GetArcotel() (string, error) {
	resp, err := http.Get("https://arcotel.ru/studentam/raspisanie-i-grafiki/raspisanie-zanyatiy-studentov-ochnoy-i-vecherney-form-obucheniya?group=54651")

	if err != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: %v", err),
			Time:    time.Now(),
		}
	}

	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка при получении страницы: %v", err),
			Time:    time.Now(),
		}
	}
	return string(content), nil
}
