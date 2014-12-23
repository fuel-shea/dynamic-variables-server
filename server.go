package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-expt/dynovars"
	"go-expt/responder"
	"net/http"
)

func main() {
	gameID := "53e256d96170706e28063201"
	rules, err := dynovars.BuildDynoVars(gameID)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Path("/features").
		HandlerFunc(FeaturesHandler(rules)).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Name("setFeatures")

	http.Handle("/", r)
	http.ListenAndServe(":3030", nil)
}

func FeaturesHandler(rules dynovars.RuleSet) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jsonData map[string]interface{}
		decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
		if decodeErr != nil {
			responder.SendError(w, responder.ErrTypes["invalid_request"])
			return
		}

		params := buildParams(jsonData, rules)
		dVars := dynovars.VarsFromFeatures(params, rules)
		responder.SendSuccess(w, dVars)
		return
	}
}

func buildParams(jsonData map[string]interface{}, rules dynovars.RuleSet) map[string]string {
	params := map[string]string{}
	for _, feat := range rules.FeatureNames {
		if val, ok := jsonData[feat]; ok {
			params[feat] = val.(string)
		}
	}

	return params
}
