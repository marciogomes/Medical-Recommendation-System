/* Server Mux */

package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strings"
	"strconv"

	"github.com/marciogomes/kg"
	"github.com/cayleygraph/cayley"
)

var store *cayley.Handle

var TodosSintomas []string

type Paciente struct {
	Name, DateBirth, CidadeAtual string
	Sintomas, RiskFactors, Diagnosticos []string
}

var p Paciente

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := r.URL.Path[len("/home/"):]
	t, err := template.ParseFiles("tmpl/" + tmpl)
	if err != nil {
		log.Fatalln(err)
	}
	t.Execute(w, nil)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := r.URL.Path[len("/edit/"):]
		t, err := template.ParseFiles("tmpl/" + tmpl)
		if err != nil {
			log.Fatalln(err)
		}
		t.Execute(w, TodosSintomas)
	} else {
		action := r.URL.Path[len("/edit/"):]
		if action == "register" {

			p.Name = r.FormValue("name")
			p.DateBirth = r.FormValue("dateBirth");
			p.CidadeAtual = r.FormValue("city");
			value := r.FormValue("fuma")
			if value == "Fumante" {
				log.Println("Novo habito inserido: " + value)
				p.RiskFactors = append(p.RiskFactors, "Fumo")
			}

			log.Println("Paciente: " + p.Name)
			log.Println("Data Nascimento: " + p.DateBirth)
			log.Println("Cidade Atual: " + p.CidadeAtual)

			http.Redirect(w, r, "/edit/register.html", http.StatusFound)

		} else if action == "symptom" {

			value := r.FormValue("body")
			log.Println("Novo sintoma inserido: " + value)
			p.Sintomas = append(p.Sintomas, value)

			test := convertToList(p.Sintomas)
			w.Write([]byte(test))

		} else if action == "signals" {

			fmt.Println(r.FormValue("temperature"))
			fmt.Println(r.FormValue("pressure"))

		}
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := r.URL.Path[len("/results/"):]
	t, err := template.ParseFiles("tmpl/" + tmpl)
	if err != nil {
		log.Fatalln(err)
	}

	resultsMap := kg.QuerySymptom(store, kg.QueryIRIs(store, p.Sintomas))
	var resultsKeys []string

	for k := range resultsMap {
		resultsKeys = append(resultsKeys, k)
	}

	resultsNames := kg.QueryNames(store, resultsKeys)

	p.Diagnosticos = nil

	// lembrando que resultsName e resultsKeys estao na mesma ordem
	for i := range resultsNames {
		p.Diagnosticos = append(p.Diagnosticos, resultsNames[i] + "-" + strconv.Itoa(resultsMap[resultsKeys[i]]))
	}

	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	req := r.URL.Path[len("/view/"):]
	busca := strings.Split(req, "-")

	t, err := template.ParseFiles("tmpl/view.html")
	if err != nil {
		log.Fatalln(err)
	}
	var result kg.Doenca
	result = kg.QueryDoenca(store, kg.QueryIRI(store, busca[0]))
	t.Execute(w, result)
}

func convertToList(terms []string) string{
	var res string
	res = "<h2>Sintomas inseridos</h2>\n<ul class=\"list-group\">\n"
	for i := range terms {
		res = res + "<li class=\"list-group-item\">" + terms[i] + "</li>\n"
	}
	res = res + "</ul>";
	return res
}

func main() {

	store = kg.Init()

	TodosSintomas = kg.QueryAll(store)

	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/results/", resultsHandler)
	http.HandleFunc("/view/", viewHandler)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img/"))))
	http.ListenAndServe(":8080", nil)
}
