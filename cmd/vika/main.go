package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Comment struct {
	Author  string
	Message string
}

type Label string

type Issue struct {
	Id          string
	Title       string
	Description string
	Author      string
	Milestone   string
	Comments    []Comment
	Labels      []Label
}

type IndexViewModel struct {
	Issues []Issue
}

type CreateViewModel struct {
}

type ReadViewModel struct {
	Issue Issue
}

type UpdateViewModel struct {
	Issue Issue
}

type DeleteViewModel struct {
	Issue Issue
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/index.html"))

	viewModel := IndexViewModel{}

	files, err := ioutil.ReadDir("./issues")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile("./issues/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		issue := Issue{
			Id: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
		}
		err = yaml.Unmarshal([]byte(data), &issue)
		viewModel.Issues = append(viewModel.Issues, issue)
	}

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/create.html"))
	viewModel := CreateViewModel{}
	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	issue := Issue{
		Id:          r.FormValue("id"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	d, err := yaml.Marshal(&issue)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("./issues/"+issue.Id+".yml", d, 0644)

	http.Redirect(w, r, "/read/"+issue.Id+"/", 303)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/read.html"))

	viewModel := ReadViewModel{}

	vars := mux.Vars(r)

	fileName := vars["id"] + ".yml"
	filePath := "./issues/" + fileName
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	issue := Issue{
		Id: strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	}
	err = yaml.Unmarshal([]byte(data), &issue)
	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func commentPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fileName := vars["id"] + ".yml"
	filePath := "./issues/" + fileName
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	issue := Issue{
		Id: strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	}
	err = yaml.Unmarshal([]byte(data), &issue)

	comment := Comment{
		Message: r.FormValue("message"),
	}
	issue.Comments = append(issue.Comments, comment)

	d, err := yaml.Marshal(&issue)
	ioutil.WriteFile(filePath, d, 0644)

	http.Redirect(w, r, "/read/"+issue.Id+"/", 303)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/update.html"))

	viewModel := UpdateViewModel{}

	vars := mux.Vars(r)

	fileName := vars["id"] + ".yml"
	filePath := "./issues/" + fileName
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	issue := Issue{
		Id: strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	}
	err = yaml.Unmarshal([]byte(data), &issue)
	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	issue := Issue{
		Id:          r.FormValue("id"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	d, err := yaml.Marshal(&issue)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("./issues/"+issue.Id+".yml", d, 0644)

	http.Redirect(w, r, "/read/"+issue.Id+"/", 303)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/delete.html"))

	viewModel := DeleteViewModel{}

	vars := mux.Vars(r)

	fileName := vars["id"] + ".yml"
	filePath := "./issues/" + fileName
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	issue := Issue{
		Id: strings.TrimSuffix(fileName, filepath.Ext(fileName)),
	}
	err = yaml.Unmarshal([]byte(data), &issue)
	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fileName := vars["id"] + ".yml"
	filePath := "./issues/" + fileName

	err := os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", 303)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/create/", createHandler).Methods("GET")
	r.HandleFunc("/create/", createPostHandler).Methods("POST")
	r.HandleFunc("/read/{id}/", readHandler).Methods("GET")
	r.HandleFunc("/read/{id}/comment/", commentPostHandler).Methods("POST")
	r.HandleFunc("/update/{id}/", updateHandler).Methods("GET")
	r.HandleFunc("/update/{id}/", updatePostHandler).Methods("POST")
	r.HandleFunc("/delete/{id}/", deleteHandler).Methods("GET")
	r.HandleFunc("/delete/{id}/", deletePostHandler).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}