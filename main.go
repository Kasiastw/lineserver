package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"lineserver/helpers"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)
type (
	line struct {
		text string
		the	 bool
	}
 	FilesStore struct {
		sync.Mutex
		store map[int]interface{}
	})

func main() {

	address := flag.String("address", ":8080", "Server address")
	file := flag.String("file", "file.txt", "the file to be loaded")
	flag.Parse()

	store := FilesStore{store: make(map[int]interface{})}
	if fileStore, err := store.getStore(*file); err!=nil{
		log.Println(" couldn't get the file store")
	} else {
		store = fileStore
	}

	app := echo.New()

	app.GET("/files/:lines", store.getLine)
	app.GET("/files/the", store.getTheLines)

	log.Printf("Listening on %s...\n", *address)
	app.Logger.Fatal(app.Start(*address))
	//app.Start(":8080")
}

func (fs FilesStore) getStore(filePath string) (FilesStore, error) {
	fs.Lock()
	defer fs.Unlock()

	f, err := os.Open(filePath)
	if err != nil {
		log.Println("error opening file: err:", err)
		return FilesStore{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	fs.store = make(map[int]interface{})
	for scanner.Scan() {
		text := scanner.Text()
		fs.store[i] = line{text: text, the: helpers.GetRegex(text)}
		i++
	}
	log.Println(fs.store)
	if err := scanner.Err(); err != nil {
		log.Println("error reading file: err:", err)
	}
	return fs, nil
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

