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
	dvSrc, err := dynovars.NewDynoVarSource()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Path("/features").
		HandlerFunc(FeaturesHandler(dvSrc)).
		Methods("POST").
		Name("setFeatures")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func FeaturesHandler(dvSrc dynovars.DynoVarSource) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData map[string]interface{}
		decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
		if decodeErr != nil {
			responder.SendError(w, responder.ErrTypes["invalid_request"])
			return
		}

		params := buildParams(jsonData, dvSrc)
		dVars, err := dvSrc.VarsFromFeatures(params, "gid1")
		if err != nil {
			fmt.Println(err)
			responder.SendError(w, responder.ErrTypes["general_error"])
			return
		}
		responder.SendSuccess(w, dVars)
		return
	}
}

func buildParams(jsonData map[string]interface{}, dvSrc dynovars.DynoVarSource) map[string]string {
	params := map[string]string{}
	for _, feat := range dvSrc.RuleSet.FeatureNames {
		if val, ok := jsonData[feat]; ok {
			params[feat] = val.(string)
		}
	}

	return params
}
