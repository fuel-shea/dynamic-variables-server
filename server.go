package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"strconv"
)

type CohortValues struct {
	MaxScore int
	MinScore int
	Random   int
}

func GetCohortValues(w rest.ResponseWriter, r *rest.Request) {
	vals := [3]CohortValues{
		CohortValues{1000, 10, 5},
		CohortValues{1100, 0, 6},
		CohortValues{900, 20, 5},
	}
	cid, err := strconv.Atoi(r.URL.Query()["cid"][0])
	if err != nil {
		log.Fatal(err)
		w.WriteJson("Error parsing cid")
		return
	}
	if cid >= len(vals) {
		w.WriteJson("cid not found")
		return
	}
	val := vals[cid]
	w.WriteJson(val)
}

func main() {
	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}

	err := handler.SetRoutes(
		&rest.Route{"GET", "/cohortvalues", GetCohortValues},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", &handler))
}
