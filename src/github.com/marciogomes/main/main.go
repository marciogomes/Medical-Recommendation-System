/* Server Mux */

package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strconv"
	"strings"

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

	results := kg.QuerySymptom(store, kg.QueryIRIs(store, p.Sintomas))
	// ... -> variadic funcion. para appendar duas []string
	results = append(results, kg.QueryRiskFactor(store, kg.QueryIRIs(store, p.RiskFactors))...)

	resultsName := kg.QueryNames(store, results)

  	/* fazer um metodo para filtragem */
	var ocorrencias map[string]int
	var resultsSet map[string]bool
	ocorrencias = make(map[string]int)
	resultsSet = make(map[string]bool)

	for i := range resultsName {
		ocorrencias[resultsName[i]]++
	}

	p.Diagnosticos = nil; // melhorar isso, ao inves de string usar map[string]float

	for i := range resultsName {
		if resultsSet[resultsName[i]] == false {
			p.Diagnosticos = append(p.Diagnosticos, resultsName[i] + "-" + strconv.FormatFloat(float64(ocorrencias[resultsName[i]]) / float64(len(resultsName)) * 100.0, 'f', 2, 64))
			resultsSet[resultsName[i]] = true
		}
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
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.ListenAndServe(":8080", nil)
}
