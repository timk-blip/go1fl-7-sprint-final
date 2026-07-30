package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(len(cafeList[city]), 100)},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)
		strings.TrimSpace(response.Body.String())
		body := response.Body.String()
		bodyInt := len(strings.Split(body, ","))
		//не понимаю, как обработать пустую строку, потому что слайс равен 1
		//Учитывайте, что если вызвать Split() для пустой строки, то в слайсе будет один элемент — пустая строка.
		if body == "" {
			bodyInt = 0
		}
		assert.Equal(t, v.want, bodyInt)
		require.Equal(t, http.StatusOK, response.Code)
	}
}
func TestCafeSearch(t *testing.T) {
	searchVal := "/cafe?city=moscow&search="
	//city := "moscow"
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string
		wantCount int
	}{
		{"", 5},
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", searchVal+v.search, nil)
		handler.ServeHTTP(response, req)
		strings.TrimSpace(response.Body.String())
		body := response.Body.String()
		strings.ToLower(body)
		count := len(strings.Split(body, ","))
		if strings.Contains(body, v.search) {
			assert.Equal(t, v.wantCount, count, v.search)
		} else {
			//не понимаю, как обработать пустую строку, потому что слайс равен 1
			//Учитывайте, что если вызвать Split() для пустой строки, то в слайсе будет один элемент — пустая строка.
			if body == "" {
				count = 0
				assert.Equal(t, v.wantCount, count, v.search)
			}
		}
		require.Equal(t, http.StatusOK, response.Code)
	}
}
