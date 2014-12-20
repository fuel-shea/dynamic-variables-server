package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"go-expt/gdv"
	"net/http"
)

var (
	db   *mgo.Database
	coll *mgo.Collection
)

func main() {
	gdv.RunCSVReader()

	sess, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer sess.Close()
	coll = sess.DB("game_vars").C("game_vars")

	populateColl(coll)

	r := mux.NewRouter()

	r.Path("/features").
		HandlerFunc(FeaturesHandler).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Name("setFeatures")

	http.Handle("/", r)
	http.ListenAndServe(":3031", nil)
}

func FeaturesHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData interface{}
	decodeErr := json.NewDecoder(r.Body).Decode(&jsonData)
	if decodeErr != nil {
		sendError(w, "Invalid JSON")
		return
	}

	gameVars, gvErr := getGameVars(jsonData)
	if gvErr != nil {
		sendError(w, "Error parsing CSV")
		return
	}

	sendSuccess(w, gameVars)
	return
}

func getGameVars(data interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"max_score":     "1000",
		"whammy_chance": "5",
	}
	return res, nil
}

func sendSuccess(w http.ResponseWriter, data map[string]interface{}) {
	succObj := SuccRespObj{
		Result: data,
	}
	succObj.Init()

	json.NewEncoder(w).Encode(succObj)
}

func sendError(w http.ResponseWriter, msg string) {
	errObj := ErrRespObj{
		ErrorObj: ErrDetailsObj{
			Message:   msg,
			ErrorCode: "ERROR",
		},
	}
	errObj.Init()

	json.NewEncoder(w).Encode(errObj)
}

// DB

func populateColl(coll *mgo.Collection) {
	rows := []Row{
		Row{
			F: Feature{
				Type: "country",
				Val:  "CAN",
				Expr: "=",
			},
			GV: GameVar{
				Type: "whammy_chance",
				Val:  "5",
			},
		},
	}
	for _, r := range rows {
		if err := coll.Insert(r); err != nil {
			panic(err)
		}
	}
}

// structs

type Row struct {
	F  Feature `bson:"feature"`
	GV GameVar `bson:"game_var"`
}

type Feature struct {
	Type string `bson:"type"`
	Val  string `bson:"val"`
	Expr string `bson:"expr"`
}

type GameVar struct {
	Type string `bson:"type"`
	Val  string `bson:"val"`
}

type SuccRespObj struct {
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result"`
}

func (sro *SuccRespObj) Init() {
	sro.Success = true
}

type ErrRespObj struct {
	Success  bool          `json:"success"`
	ErrorObj ErrDetailsObj `json:"error"`
}

func (ero *ErrRespObj) Init() {
	ero.Success = false
}

type ErrDetailsObj struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
}
