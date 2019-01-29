package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/spf13/afero"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/mrombout/vika"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type indexViewModel struct {
	Issues []vika.Issue
}

type createViewModel struct {
}

type readViewModel struct {
	Issue vika.Issue
}

type updateViewModel struct {
	Issue vika.Issue
}

type deleteViewModel struct {
	Issue vika.Issue
}

var box packr.Box
var repository vika.IssuesRepository

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "index.html")
	viewModel := indexViewModel{}

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
	viewModel := createViewModel{}
	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	issue := vika.Issue{
		ID:          vika.ID(r.FormValue("id")),
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
	viewModel := readViewModel{}

	vars := mux.Vars(r)
	id := vika.ID(vars["id"])

	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func commentPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vika.ID(vars["id"])
	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	comment := vika.Comment{
		Author:  r.FormValue("author"),
		Message: r.FormValue("message"),
	}
	issue.Comments = append(issue.Comments, comment)

	repository.SaveIssue(&issue)

	http.Redirect(w, r, "/read/"+string(issue.ID)+"/", 303)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "update.html")
	viewModel := updateViewModel{}

	vars := mux.Vars(r)

	id := vika.ID(vars["id"])
	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vika.ID(vars["id"])

	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(issue)
	}

	issue.ID = vika.ID(r.FormValue("id"))
	issue.Title = r.FormValue("title")
	issue.Description = r.FormValue("description")

	repository.SaveIssue(&issue)

	http.Redirect(w, r, "/read/"+string(issue.ID)+"/", 303)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := loadTemplate(&box, "delete.html")
	viewModel := deleteViewModel{}

	vars := mux.Vars(r)
	id := vika.ID(vars["id"])
	issue, err := repository.GetIssue(id)
	if err != nil {
		log.Fatal(err)
	}

	viewModel.Issue = issue

	tmpl.ExecuteTemplate(w, "layout", viewModel)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vika.ID(vars["id"])

	repository.DeleteIssue(id)

	http.Redirect(w, r, "/", 303)
}

func main() {
	box = packr.NewBox("../../templates")
	repository = vika.FilesystemYamlIssuesRepository{
		Fs: vika.AferoFilesystem{
			Fs: afero.NewOsFs(),
		},
	}

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
