/* Knowledge Graph */

package kg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sort"
	"strconv"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/quad/cquads"
)

type Doenca struct {
	Name, Description, Code, Image string
	Sintomas, RiskFactors, Drugs []string
}

func Init() *cayley.Handle {
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open("data/quads.nq")
	if(err != nil) {
		log.Fatalln(err)
	}

	decoder := cquads.NewDecoder(file)
	for {
		aux, err := decoder.Unmarshal()
		if err != nil {
			break
		}
		store.AddQuad(aux)
	}

	return store
}

func QuerySymptom(store *cayley.Handle, sintomas []string) map[string]int {
	var diagnosticos map[string]int
	diagnosticos = make(map[string]int)
  	// fazendo o caminho
	// faço Raw porque ja passo uma IRI
	for i := range sintomas {

		// percorre os pesos
		for j := 1; j <= 5; j++ {

			p := cayley.StartPath(store, quad.Raw(sintomas[i])).In(quad.IRI("http://health-lifesci.schema.org/MedicalSignOrSymptom:" + strconv.Itoa(j)))

			err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
				nativeValue := quad.NativeOf(value)
				diagnosticos[fmt.Sprint(nativeValue)] += j;
				})
			if err != nil {
				log.Fatalln(err)
			}
		}
		
	}
	return diagnosticos
}

func QueryRiskFactor(store *cayley.Handle, riskFactors []string) map[string]int {
	var diagnosticos map[string]int
	diagnosticos = make(map[string]int)
  	// fazendo o caminho
	// faço Raw porque ja passo uma IRI
	for i := range riskFactors {

		// percorre os pesos
		for j := 1; j <= 5; j++ {

			p := cayley.StartPath(store, quad.Raw(riskFactors[i])).In(quad.IRI("http://health-lifesci.schema.org/MedicalRiskFactor:" + strconv.Itoa(j)))

			err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
				nativeValue := quad.NativeOf(value)
				diagnosticos[fmt.Sprint(nativeValue)] += j;
				})
			if err != nil {
				log.Fatalln(err)
			}
		}
		
	}
	return diagnosticos
}


// query para obter informações de uma doença
func QueryDoenca(store *cayley.Handle, diagnostico string) Doenca {
	var d Doenca
  	// fazendo o caminho

	p := cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/name"))

	err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		d.Name = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	p = cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/description"))

	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		d.Description = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	p = cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/code"))

	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		d.Code = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	d.Code = QueryCodeValue(store, d.Code)

	// Obtem os sintomas

	for i := 1; i <= 5; i++ {

		p = cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/MedicalSignOrSymptom:" + strconv.Itoa(i)))

		err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			nativeValue := quad.NativeOf(value)
			d.Sintomas = append(d.Sintomas, fmt.Sprint(nativeValue))
			})

		if err != nil {
			log.Fatalln(err)
		}
	}

	d.Sintomas = QueryNames(store, d.Sintomas)

	for i := 1; i <= 5; i++ {

		p = cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/MedicalRiskFactor:" + strconv.Itoa(i)))

		err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			nativeValue := quad.NativeOf(value)
			d.RiskFactors = append(d.RiskFactors, fmt.Sprint(nativeValue))
			})

		if err != nil {
			log.Fatalln(err)
		}
	}
	
	d.RiskFactors = QueryNames(store, d.RiskFactors)

	p = cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/Drug"))

	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		d.Drugs = append(d.Drugs, fmt.Sprint(nativeValue))
		})

	if err != nil {
		log.Fatalln(err)
	}

	d.Drugs = QueryNames(store, d.Drugs)

	p = cayley.StartPath(store, quad.Raw(diagnostico)).Out(quad.IRI("http://health-lifesci.schema.org/image"))

	err = p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		d.Image = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	// remove os '<' e ">"
	d.Image = strings.Replace(d.Image, "<", "", -1)
	d.Image = strings.Replace(d.Image, ">", "", -1)

	return d
}

// query para obter o nome do objeto
func QueryName(store *cayley.Handle, entidade string) string {
	var res string
  	// fazendo o caminho
	// eu faco raw porque ja estou passando uma IRI
	p := cayley.StartPath(store, quad.Raw(entidade)).Out(quad.IRI("http://health-lifesci.schema.org/name"))

	err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		res = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	return res
}

// obrigado por nao ter polimorfismo =)
func QueryNames(store *cayley.Handle, entidades []string) []string {
	var res []string
  	// fazendo o caminho

	for i := range entidades {
		p := cayley.StartPath(store, quad.Raw(entidades[i])).Out(quad.IRI("http://health-lifesci.schema.org/name"))

		err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			nativeValue := quad.NativeOf(value)
			res = append(res, fmt.Sprint(nativeValue))
			})
		if err != nil {
			log.Fatalln(err)
		}
	}
	return res
}

// query para obter o IRI do objeto
func QueryIRI(store *cayley.Handle, entidade string) string {
	var res string
  	// fazendo o caminho
  	// passo string porque é o tipo do nome, oras
	p := cayley.StartPath(store, quad.String(entidade)).In(quad.IRI("http://health-lifesci.schema.org/name"))

	err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		res = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	return res
}

// obrigado por nao ter polimorfismo =)
func QueryIRIs(store *cayley.Handle, entidades []string) []string {
	var res []string
  	// fazendo o caminho

	for i := range entidades {
		p := cayley.StartPath(store, quad.String(entidades[i])).In(quad.IRI("http://health-lifesci.schema.org/name"))

		err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			nativeValue := quad.NativeOf(value)
			res = append(res, fmt.Sprint(nativeValue))
			})
		if err != nil {
			log.Fatalln(err)
		}
	}
	return res
}

// query para obter o codigo CID10 do objeto
func QueryCodeValue(store *cayley.Handle, entidade string) string {
	var res string
  	// fazendo o caminho
	// eu faco raw porque ja estou passando uma IRI
	p := cayley.StartPath(store, quad.Raw(entidade)).Out(quad.IRI("http://health-lifesci.schema.org/codeValue"))

	err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value)
		res = fmt.Sprint(nativeValue)
		})

	if err != nil {
		log.Fatalln(err)
	}

	return res
}

// retorna todos os sintomas do nq
func QueryAll(store *cayley.Handle) []string {
	var sintomas []string
	var sintomasSet map[string]bool
	sintomasSet = make(map[string]bool)

  	// fazendo o caminho
	// parte de todos
	// percorre os pesos
	for i := 1; i <= 5; i++ {
		p := cayley.StartPath(store).Out(quad.IRI("http://health-lifesci.schema.org/MedicalSignOrSymptom:" + strconv.Itoa(i)))

		err := p.Iterate(nil).EachValue(nil, func(value quad.Value) {
			nativeValue := quad.NativeOf(value)
			v := fmt.Sprint(nativeValue)
			if sintomasSet[v] == false {
				sintomasSet[v] = true
				sintomas = append(sintomas, v)
			}
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
	sintomas = QueryNames(store, sintomas)
	sort.Strings(sintomas)
	
	return sintomas
}