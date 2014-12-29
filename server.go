package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-expt/dynovars"
	"go-expt/responder"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Path("/features").
		HandlerFunc(FeaturesHandler).
		Methods("POST").
		Name("setFeatures")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func FeaturesHandler(w http.ResponseWriter, r *http.Request) {
	dvSrc, err := dynovars.NewDynoVarSource()
	if err != nil {
		responder.SendError(w, responder.ErrTypes["empty_result"])
		return
	}

	var jsonData map[string]interface{}
	decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
	if decodeErr != nil {
		responder.SendError(w, responder.ErrTypes["invalid_request"])
		return
	}

	dVars, err := dvSrc.VarsFromFeatures(jsonData, "gid1")
	if err != nil {
		fmt.Println(err)
		responder.SendError(w, responder.ErrTypes["general_error"])
		return
	}

	responder.SendSuccess(w, dVars)
	return
}
