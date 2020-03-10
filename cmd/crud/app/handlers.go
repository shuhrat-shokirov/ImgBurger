package app

import (
	"context"
	"crud/pkg/crud/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const multipartMaxBytes = 10 * 1024 * 1024

func (receiver *server) handleBurgersList(isAdmin bool) func(http.ResponseWriter, *http.Request) {
	var (
		tpl *template.Template
		err error
	)
	if isAdmin {
		tpl, err = template.ParseFiles(
			filepath.Join(receiver.templatesPath, "admin", "burgers.gohtml"),
			filepath.Join(receiver.templatesPath, "base.gohtml"),
		)
	} else {
		tpl, err = template.ParseFiles(
			filepath.Join(receiver.templatesPath, "index.gohtml"),
			filepath.Join(receiver.templatesPath, "base.gohtml"),
		)
	}
	if err != nil {
		panic(err)
	}
	// -> go concurrency (paraller/thread-safe)
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := receiver.burgersSvc.BurgersList(request.Context())
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Burgers []models.Burger
		}{
			Title:   "McDonalds",
			Burgers: list,
		}

		err = tpl.Execute(writer, data)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (receiver *server) handleBurgersSave() func(responseWriter http.ResponseWriter, request *http.Request) {
	// POST
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(multipartMaxBytes)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		file, header, err := request.FormFile("image")
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer file.Close()

		contentType := header.Header.Get("Content-Type")

		fileName, err := receiver.filesSvc.Save(file, contentType)

		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// TODO: bad practice - save burger
		name := request.PostForm.Get("name")
		price := request.PostForm.Get("price")
		parsedPrice, err := strconv.Atoi(price)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}
		receiver.burgersSvc.Save(context.Background(), models.Burger{Name: name, Price: parsedPrice, FileName: fileName})
		http.Redirect(writer, request, "/admin/burgers", http.StatusFound)
		return
	}
}

func (receiver *server) handleBurgersRemove() func(responseWriter http.ResponseWriter, request *http.Request) {
	_, err := template.ParseFiles(filepath.Join(receiver.templatesPath, "index.gohtml"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		err := request.ParseForm()
		if err != nil {
			panic(err)
		}
		id := request.PostForm.Get("id")
		log.Print(id)
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400 - Bad Request
			return
		}
		receiver.burgersSvc.RemoveById(context.Background(),parsedId)

		http.Redirect(writer, request, "/admin/burgers", http.StatusFound)
		return
	}
}

func (receiver *server) handleSlow() func(responseWriter http.ResponseWriter, request *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		select {
		case <-request.Context().Done():
			return
		case <-time.After(time.Second * 10): // эмулируем долгий ответ
			_, _ = responseWriter.Write([]byte("Hello slow!"))
		}
	}
}

func (receiver *server) handleFavicon() func(http.ResponseWriter, *http.Request) {
	file, err := ioutil.ReadFile(filepath.Join(receiver.assetsPath, "favicon.ico"))
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(file)
		if err != nil {
			log.Print(err)
		}
	}
}
