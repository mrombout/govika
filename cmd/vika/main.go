package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/mrombout/govika"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type IndexViewModel struct {
	Issues []govika.Issue
}

type CreateViewModel struct {
}

type ReadViewModel struct {
	Issue govika.Issue
}

type UpdateViewModel struct {
	Issue govika.Issue
}

type DeleteViewModel struct {
	Issue govika.Issue
}

var box packr.Box
var repository govika.IssuesRepository

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "index.html")
	viewModel := IndexViewModel{}

	issues, err := repository.GetIssues()
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issues = issues

	err = tmpl.ExecuteTemplate(w, "layout", viewModel)
	if err != nil {
		log.Fatal(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "create.html")
	viewModel := CreateViewModel{}
	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	issue := govika.Issue{
		ID:          govika.ID(r.FormValue("id")),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	err := repository.SaveIssue(&issue)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/read/"+string(issue.ID)+"/", 303)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "read.html")
	viewModel := ReadViewModel{}

	vars := mux.Vars(r)
	id := govika.ID(vars["id"])

	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func commentPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := govika.ID(vars["id"])
	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	comment := govika.Comment{
		Message: r.FormValue("message"),
	}
	issue.Comments = append(issue.Comments, comment)

	repository.SaveIssue(&issue)

	http.Redirect(w, r, "/read/"+string(issue.ID)+"/", 303)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "update.html")
	viewModel := UpdateViewModel{}

	vars := mux.Vars(r)

	id := govika.ID(vars["id"])
	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := govika.ID(vars["id"])

	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(issue)
	}

	issue.ID = govika.ID(r.FormValue("id"))
	issue.Title = r.FormValue("title")
	issue.Description = r.FormValue("description")

	repository.SaveIssue(&issue)

	http.Redirect(w, r, "/read/"+string(issue.ID)+"/", 303)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "delete.html")
	viewModel := DeleteViewModel{}

	vars := mux.Vars(r)
	id := govika.ID(vars["id"])
	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := govika.ID(vars["id"])

	repository.DeleteIssue(id)

	http.Redirect(w, r, "/", 303)
}

func main() {
	box = packr.NewBox("../../templates")
	repository = govika.FilesystemIssuesRepository{}

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

func loadTemplate(box *packr.Box, file string) *template.Template {
	layoutHTML, err := box.FindString("layout.html")
	if err != nil {
		log.Fatal(err)
	}

	layoutTmpl := template.New("layout")
	layoutTmpl, err = layoutTmpl.Parse(layoutHTML)
	layoutTmpl.Funcs(template.FuncMap{
		"markdown": markdown,
	})
	if err != nil {
		log.Fatal(err)
	}

	html, err := box.FindString(file)
	if err != nil {
		log.Fatal(err)
	}
	layoutTmpl.Parse(html)

	return layoutTmpl
}

func markdown(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
