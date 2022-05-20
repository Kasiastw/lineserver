package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

type (
	line struct {
		text string
		the	 bool
	}
	FilesStore struct {
		sync.Mutex
		store map[int]interface{}
	}
)

var (

	mockStore = map[int]interface{}{
		0: line{text: "the", the: true},
		1: line{text: "quick brown", the: false},
		2: line{text: "fox jumps over the", the: true},
		3: line{text: "lazy dog", the: false},
	}
	////store with word without words with 'the'
	mockStoreWithThe = map[int]interface{}{
		0: line{text: "", the: false},
		1: line{text: "quick brown", the: false},
		2: line{text: "fox jumps over", the: false},
		3: line{text: "lazy dog", the: false},
	}
)

func (fs FilesStore) getLine(ctx echo.Context) error {
	lineNumber := ctx.Param("lines")
	if lineNumber == "" {
		return errors.New("There is no line number")
	}
	intId, _ := strconv.Atoi(lineNumber)
	if _, ok := fs.store[intId]; ok {
		return ctx.String(http.StatusOK, fs.store[intId].(line).text)
	}
	return ctx.JSON(http.StatusNotFound, "the required line does not exist")
}

func (fs FilesStore) getTheLines(ctx echo.Context) error {

	var tempList []string

	for _, v := range fs.store {
		if v.(line).the == true {
			tempList = append(tempList, v.(line).text)
		}
	}
	if len(tempList)==0 {
		fmt.Println("len(tempList)", len(tempList))
		return ctx.JSON(http.StatusNotFound, "there are no words with the")
	}
	return ctx.JSON(http.StatusOK, tempList)
}

func TestGetLines(t *testing.T) {

	t.Run("test when the number of line is out of range", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/files/lines")
		c.SetParamNames("lines")
		c.SetParamValues("5")
		h := &FilesStore{sync.Mutex{} , mockStore }
		if assert.NoError(t, h.getLine(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})

	t.Run("test when the number of line exists", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/files/lines")
		c.SetParamNames("lines")
		c.SetParamValues("1")
		h := &FilesStore{sync.Mutex{}, mockStore }
		if assert.NoError(t, h.getLine(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestGetTheLines(t *testing.T) {

	t.Run("test when the file consists of the words with the", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/files/the")
		h := &FilesStore{sync.Mutex{}, mockStore }
		if assert.NoError(t, h.getTheLines(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("test the file doesn't consist of the words with the", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/files/name/the")
		h := &FilesStore{sync.Mutex{}, mockStoreWithThe}
		if assert.NoError(t, h.getTheLines(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}

