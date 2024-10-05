package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type logError struct {
	Message string
	Time    time.Time
}

func (e *logError) Error() string {
	return fmt.Sprintf("[%v] %s", e.Time, e.Message)
}

var (
	count    = 1
	schedule string
)

func GetSchedule(selGroup string, day string) (string, error) {
	groups, errp := GetGroups()
	if errp != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка при парсинге групп: \n%v", errp),
			Time:    time.Now(),
		}
	}

	if time.Now().Weekday().String() == "Sunday" {
		var errs error
		schedule, errs = GetSchSunday(groups, selGroup, day)
		if errs != nil {
			return "", &logError{
				Message: fmt.Sprintf("Ошибка при парсинге групп: \n%v", errs),
				Time:    time.Now(),
			}
		}
	} else {
		var errs error
		schedule, errs = GetSch(groups, selGroup, day)
		if errs != nil {
			return "", &logError{
				Message: fmt.Sprintf("Ошибка при парсинге групп: \n%v", errs),
				Time:    time.Now(),
			}
		}
	}

	return schedule, nil
}

func GetGroups() (map[string]string, error) {
	resp, errg := http.Get("https://arcotel.ru/studentam/raspisanie-i-grafiki/raspisanie-zanyatiy-studentov-ochnoy-i-vecherney-form-obucheniya")

	if errg != nil {
		return nil, &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: %v", errg),
			Time:    time.Now(),
		}
	}
	defer resp.Body.Close()

	doc, errq := goquery.NewDocumentFromReader(resp.Body)
	if errq != nil {
		return nil, &logError{
			Message: fmt.Sprintf("Ошибка при получении страницы: \n%v", errq),
			Time:    time.Now(),
		}
	}

	groups := make(map[string]string)

	doc.Find(".vt256").Each(func(i int, s *goquery.Selection) {
		dataI, _ := s.Attr("data-i")
		dataNm := s.Text()
		dataNm = strings.TrimSpace(dataNm)
		groups[dataNm] = dataI
	})

	return groups, nil
}

func GetSch(groups map[string]string, selGroup string, day string) (string, error) {
	v, ok := groups[selGroup]
	if !ok {
		return "", &logError{
			Message: fmt.Sprintf("%s группы не существует.", selGroup),
			Time:    time.Now(),
		}
	}

	resp, errg := http.Get(fmt.Sprintf("https://arcotel.ru/studentam/raspisanie-i-grafiki/raspisanie-zanyatiy-studentov-ochnoy-i-vecherney-form-obucheniya?group=%s", v))

	if errg != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: \n%v", errg),
			Time:    time.Now(),
		}
	}
	defer resp.Body.Close()

	doc, errq := goquery.NewDocumentFromReader(resp.Body)
	if errq != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка при получении страницы: \n%v", errq),
			Time:    time.Now(),
		}
	}

	schedule := [5]string{}

	doc.Find(".vt237").Each(func(i int, s *goquery.Selection) {
		if v, ok := s.Attr("data-i"); ok && v == day {
			text := s.Text()
			text = strings.TrimSpace(text)
			text = strings.Trim(text, "\n")
			schedule[0] = text
		}
	})

	doc.Find(fmt.Sprintf(".vt239.rasp-day.rasp-day%s", day)).Each(func(i int, s *goquery.Selection) {
		text := s.Text() + " "
		text = strings.TrimSpace(text)
		text = strings.Trim(text, "\n")
		schedule[i] = fmt.Sprintf("-------------------------\n%v.\n%s\n-------------------------", count, text)
		count++
	})

	count = 1
	var text string

	for i := 0; i < 5; i++ {
		text += schedule[i]
	}

	return text, nil
}

func GetSchSunday(groups map[string]string, selGroup string, day string) (string, error) {
	href, err := getHref(groups, selGroup)
	if err != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: \n%v", err),
			Time:    time.Now(),
		}
	}

	resp, errg := http.Get(fmt.Sprintf("https://arcotel.ru/studentam/raspisanie-i-grafiki/raspisanie-zanyatiy-studentov-ochnoy-i-vecherney-form-obucheniya%s", href))

	if errg != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: \n%v", errg),
			Time:    time.Now(),
		}
	}

	defer resp.Body.Close()

	doc, errq := goquery.NewDocumentFromReader(resp.Body)
	if errq != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка при получении страницы: \n%v", errq),
			Time:    time.Now(),
		}
	}

	schedule := [5]string{}

	doc.Find(".vt237").Each(func(i int, s *goquery.Selection) {
		if v, ok := s.Attr("data-i"); ok && v == day {
			text := s.Text()
			text = strings.TrimSpace(text)
			text = strings.Trim(text, "\n")
			schedule[0] = text
		}
	})

	doc.Find(fmt.Sprintf(".vt239.rasp-day.rasp-day%s", day)).Each(func(i int, s *goquery.Selection) {
		text := s.Text() + " "
		text = strings.TrimSpace(text)
		text = strings.Trim(text, "\n")
		schedule[i] = fmt.Sprintf("-------------------------\n%v:\n%s\n-------------------------", count, text)
		count++
	})

	count = 1

	var text string

	for i := 0; i < 5; i++ {
		text += schedule[i]
	}

	return text, nil
}

func getHref(groups map[string]string, selGroup string) (string, error) {
	v, ok := groups[selGroup]
	if !ok {
		return "", &logError{
			Message: fmt.Sprintf("%s группы не существует.", selGroup),
			Time:    time.Now(),
		}
	}

	resp, errg := http.Get(fmt.Sprintf("https://arcotel.ru/studentam/raspisanie-i-grafiki/raspisanie-zanyatiy-studentov-ochnoy-i-vecherney-form-obucheniya?group=%s", v))

	if errg != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: \n%v", errg),
			Time:    time.Now(),
		}
	}

	defer resp.Body.Close()

	doc, errq := goquery.NewDocumentFromReader(resp.Body)
	if errq != nil {
		return "", &logError{
			Message: fmt.Sprintf("Ошибка при получении страницы: \n%v", errq),
			Time:    time.Now(),
		}
	}

	var href string

	doc.Find(".vt233.vt235").Each(func(i int, s *goquery.Selection) {
		href, _ = s.Attr("href")
	})

	return href, nil
}
