package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Mycv struct {
	Name      string
	Age       string
	Education string
	Role      string
	Location  string
	Phone     string
	Email     string
	Github    string
	Skills    string
	Projects  string
	Languages string
}

var templates = make(map[string]*template.Template, 4)
var responceCVdata Mycv

func loadTemplates() {

	/*
		files := []string{
			"./tmpl/base.html",
			"./tmpl/partials/about.html",
			"./tmpl/partials/education.html",
			"./tmpl/partials/experience.html",
			"./tmpl/partials/keywords.html",
			"./tmpl/partials/languages.html",
			"./tmpl/partials/mission.html",
			"./tmpl/partials/projects.html",
			"./tmpl/partials/skills.html",
		}

		templates, err := template.ParseFiles(files...)
		if err != nil {
			log.Print(err.Error())
			return
		}*/

	templateNames := [4]string{"form", "main", "about", "skills"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("tmpl/layout.html", "tmpl/parts/"+name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

// aboutHandler handles about data from cv
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	templates["about"].Execute(w, responceCVdata)
}

// skillsHandler handles skills data from cv
func skillsHandler(w http.ResponseWriter, r *http.Request) {
	templates["skills"].Execute(w, responceCVdata)
}

// mainHandler handles main page after create cv
func mainHandler(w http.ResponseWriter, r *http.Request) {
	templates["main"].Execute(w, responceCVdata) //templates["main"].Execute(w, responceCVdata.Name)
}

type formData struct {
	*Mycv
	Errors []string
}

// formHandler handles "/form" URL
func cvFormHandler(w http.ResponseWriter, request *http.Request) {

	switch request.Method {

	case http.MethodGet:
		templates["form"].Execute(w, formData{
			Mycv: &Mycv{
				Name:      "Егор",
				Age:       "36",
				Education: "высшее",
				Role:      "специалист",
				Location:  "Москва",
				Phone:     "+79169998877",
				Email:     "ieo@rambler.ru",
				Github:    "github.com/impr0ver",
				Skills:    "templates on golang",
				Projects:  "cv on golang",
				Languages: "C and golang",
			}, Errors: []string{}, //&Mycv{} for empty form
		})

	case http.MethodPost:
		request.ParseForm()

		responceCVdata = Mycv{
			Name:      request.FormValue("name"),
			Age:       request.FormValue("age"),
			Education: request.FormValue("education"),
			Role:      request.FormValue("role"),
			Location:  request.FormValue("location"),
			Phone:     request.FormValue("phone"),
			Email:     request.FormValue("email"),
			Github:    request.FormValue("github"),
			Skills:    request.FormValue("skills"),
			Projects:  request.FormValue("projects"),
			Languages: request.FormValue("languages"),
		}

		err := checkFormErrors(&responceCVdata)
		if err != nil {
			templates["form"].Execute(w, formData{
				Mycv: &responceCVdata, Errors: err,
			})
		} else {
			templates["main"].Execute(w, responceCVdata)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func checkFormErrors(respData *Mycv) []string {
	errors := []string{}

	if respData.Name == "" {
		errors = append(errors, "Please enter your name")
	}
	if respData.Age == "" {
		errors = append(errors, "Please enter your age")
	}
	if respData.Education == "" {
		errors = append(errors, "Please enter your education")
	}
	if respData.Role == "" {
		errors = append(errors, "Please enter your role")
	}
	if respData.Location == "" {
		errors = append(errors, "Please enter your location")
	}
	if respData.Phone == "" {
		errors = append(errors, "Please enter your phone number")
	}
	if respData.Email == "" {
		errors = append(errors, "Please enter your email")
	}
	if respData.Github == "" {
		errors = append(errors, "Please enter your github repository")
	}
	if respData.Skills == "" {
		errors = append(errors, "Please enter your skills")
	}
	if respData.Projects == "" {
		errors = append(errors, "Please enter your projects")
	}
	if respData.Languages == "" {
		errors = append(errors, "Please enter your programming languages")
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func main() {
	loadTemplates()

	http.HandleFunc("/", mainHandler /*helloHandler*/)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/skills", skillsHandler)
	http.HandleFunc("/form", cvFormHandler)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
