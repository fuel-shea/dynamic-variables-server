package responder

import (
	"encoding/json"
	"net/http"
)

type ErrorType struct {
	ErrorCode string `json:"errorcode"`
	ErrorMsg  string `json:"error"`
}

var ErrTypes = map[string]ErrorType{
	"general_error": ErrorType{
		ErrorCode: "GENERAL_ERROR",
		ErrorMsg:  "general error",
	},
	"empty_result": ErrorType{
		ErrorCode: "EMPTY_RESULT",
		ErrorMsg:  "empty result",
	},
	"authorization_failed": ErrorType{
		ErrorCode: "AUTHORIZATION_FAILED",
		ErrorMsg:  "authorization failed",
	},
	"parameters_missing": ErrorType{
		ErrorCode: "PARAMETERS_MISSING",
		ErrorMsg:  "parameters missing",
	},
	"invalid_request": ErrorType{
		ErrorCode: "INVALID_REQUEST",
		ErrorMsg:  "invalid request",
	},
}

func SendSuccess(w http.ResponseWriter, data map[string]interface{}) {
	succObj := SuccRespObj{Result: data}
	succObj.Init()
	json.NewEncoder(w).Encode(succObj)
}

func SendError(w http.ResponseWriter, errType ErrorType) {
	errObj := ErrRespObj{ErrorType: errType}
	errObj.Init()
	json.NewEncoder(w).Encode(errObj)
}

type SuccRespObj struct {
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result"`
}

func (sro *SuccRespObj) Init() {
	sro.Success = true
}

type ErrRespObj struct {
	Success   bool      `json:"success"`
	ErrorType ErrorType `json:"error"`
}

func (ero *ErrRespObj) Init() {
	ero.Success = false
}
