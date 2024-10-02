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

func GetSchedule(selGroup string, day int8) ([5]string, error) {
	groups, errp := GetGroups()
	if errp != nil {
		return [5]string{}, &logError{
			Message: fmt.Sprintf("Ошибка при парсинге групп: \n%v", errp),
			Time:    time.Now(),
		}
	}

	schedule, errs := GetSch(groups, selGroup, day)
	if errs != nil {
		return [5]string{}, &logError{
			Message: fmt.Sprintf("Ошибка при парсинге групп: \n%v", errs),
			Time:    time.Now(),
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

func GetSch(groups map[string]string, selGroup string, day int8) ([5]string, error) {
	v, ok := groups[selGroup]
	if !ok {
		return [5]string{}, &logError{
			Message: "Такой группы не существует.",
			Time:    time.Now(),
		}
	}

	resp, errg := http.Get(fmt.Sprintf("https://arcotel.ru/studentam/raspisanie-i-grafiki/raspisanie-zanyatiy-studentov-ochnoy-i-vecherney-form-obucheniya?group=%s", v))

	if errg != nil {
		return [5]string{}, &logError{
			Message: fmt.Sprintf("Ошибка подключения к arcotel.ru: \n%v", errg),
			Time:    time.Now(),
		}
	}
	defer resp.Body.Close()

	doc, errq := goquery.NewDocumentFromReader(resp.Body)
	if errq != nil {
		return [5]string{}, &logError{
			Message: fmt.Sprintf("Ошибка при получении страницы: \n%v", errq),
			Time:    time.Now(),
		}
	}

	schedule := [5]string{}

	doc.Find(".vt237").Each(func(i int, s *goquery.Selection) {
		if _, ok := s.Attr("data-i"); ok {
			text := s.Text()
			text = strings.TrimSpace(text)
			text = strings.Trim(text, "\n")
			schedule[i] = text
		}
	})

	doc.Find(".vt244").Each(func(i int, s *goquery.Selection) {
		s.Find(".vt239").Each(func(j int, st *goquery.Selection) {
			text := s.Text() + " " + st.Text()
			text = strings.TrimSpace(text)
			text = strings.Trim(text, "\n")
			schedule[i] = text[:100]
		})
	})

	for i := 0; i < 5; i++ {
		for j := 0; j < 7; j++ {
			fmt.Printf("%s\t", schedule[i][j])
		}
	}

	return schedule, nil
}
