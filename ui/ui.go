package ui

import (
	"encoding/json"
	"fmt"
	"strconv"
	"net"
	"net/http"
	"time"
	"io/ioutil"
	"log"

	"github.com/crazyyi/goweb/model"
)

type Config struct {
	Assets http.FileSystem
}

func Start(cfg Config, m *model.Model, listener net.Listener) {
	server := &http.Server{
		ReadTimeout:	60 * time.Second,
		WriteTimeout: 60 * time.Second,
		MaxHeaderBytes: 1 << 16}

	http.Handle("/", indexHandler(m))
	http.Handle("/people", peopleHandler(m))
	http.Handle("/create", createNewRecordHandler(m))
	http.Handle("/delete", deleteRowHandler(m))
	http.Handle("/update", updateRecordHandler(m))
	http.Handle("/js/", http.FileServer(cfg.Assets))

	go server.Serve(listener)
}

const (
	// cdnReact           = "https://unpkg.com/react@16/umd/react.production.min.js"
	// cdnReactDom        = "https://unpkg.com/react-dom@16/umd/react-dom.production.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.26.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.18.0/axios.min.js"
)

const indexHTML = `
	<!DOCTYPE HTML>
	<html>
		<head>
			<meta charset="utf-8">
			<title>Simple Go Web App</title>
		</head>
		<body>
			<div id='root'></div>
			<script crossorigin src="` + cdnBabelStandalone + `"></script>
			<script crossorigin src="` + cdnAxios + `"></script>
			<script crossorigin src="/js/dist/bundle.js" data-presets="es2015,react" ></script>
		</body>
	</html>
`

func indexHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, indexHTML)
	})
}

func peopleHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		people, err := m.People()
		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}

		js, err := json.Marshal(people)
		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(js))
	})
}

func createNewRecordHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		data := make(map[string]string)
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}

		var firstname = data["firstname"]
		var lastname = data["lastname"]
		log.Printf("%v", data)
		fmt.Printf(firstname + ", " + lastname + "\n")

		id, err := m.CreateRecord(&model.Person{Id: 10, First: string(firstname), Last: string(lastname)})
		if err != nil {
			http.Error(w, "Error creating new record", http.StatusBadRequest)
			panic(err)
		}

		fmt.Fprintf(w, string(id))

	})
}

func deleteRowHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id = r.URL.Query().Get("id")
		v, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}

		fmt.Printf("Pressed delete, v = %v\n", v)
		count, err := m.DeleteRow(v)

		if err != nil {
			http.Error(w, "This is an error", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(count))
	})
}

func updateRecordHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		data := make(map[string]interface{})
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Error parsing json object", http.StatusBadRequest)
			return
		}

		var id = data["id"].(float64)
		var firstname = data["firstname"].(string)
		var lastname = data["lastname"].(string)

		log.Println(firstname + ", " + lastname)

		if err != nil {
			http.Error(w, "Error parsing int", http.StatusBadRequest)
			return
		}

		count, err := m.UpdateRow(int64(id), firstname, lastname)

		if err != nil {
			http.Error(w, "Error updating row", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, string(count))
	})
}