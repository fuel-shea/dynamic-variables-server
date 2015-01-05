package fuelresponder_test

import (
	"encoding/json"
	"github.com/fuel-shea/fuel-go-utils/fuelresponder"
	"net/http/httptest"
	"testing"
)

type testCase struct{}

var testCases = []testCase{}

func TestSendSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{"somekey": "someval"}
	fuelresponder.SendSuccess(w, data)

	respBytes := []byte(w.Body.String())
	var respJSON interface{}
	err := json.Unmarshal(respBytes, &respJSON)
	if err != nil {
		t.Fatal(err)
	}

	respJSONMap, isMap := respJSON.(map[string]interface{})
	if !isMap {
		t.Fatal("Response body is not a map...")
	}

	successVal, successFound := respJSONMap["success"]
	if !successFound {
		t.Error("Response has no 'success' field")
	}
	if successVal != true {
		t.Error("'success' should be true, found", successVal)
	}

	resultVal, resultFound := respJSONMap["result"]
	if !resultFound {
		t.Error("Response has no 'result' field")
	}
	resultMap, isMap := resultVal.(map[string]interface{})
	if !isMap {
		t.Error("'result' field is not a map")
	}

	if resultMap["somekey"] != "someval" {
		t.Error("'Expected 'somekey' to be 'someval', got", resultMap["somekey"])
	}
}

func TestSendError(t *testing.T) {
	errorSent := fuelresponder.ErrTypes["general_error"]
	w := httptest.NewRecorder()
	fuelresponder.SendError(w, errorSent)

	respBytes := []byte(w.Body.String())
	var respJSON interface{}
	err := json.Unmarshal(respBytes, &respJSON)
	if err != nil {
		t.Fatal(err)
	}

	respJSONMap, isMap := respJSON.(map[string]interface{})
	if !isMap {
		t.Fatal("Response body is not a map...")
	}

	successVal, successFound := respJSONMap["success"]
	if !successFound {
		t.Error("Response has no 'success' field")
	}
	if successVal != false {
		t.Error("'success' should be false, found", successVal)
	}

	errVal, errFound := respJSONMap["error"]
	if !errFound {
		t.Error("Response has no 'error' field")
	}
	errMap, isMap := errVal.(map[string]interface{})
	if !isMap {
		t.Error("'error' field is not a map")
	}

	if errMap["errorcode"] != errorSent.ErrorCode {
		t.Error("'Expected 'error.errorcode' to be ", errorSent.ErrorCode, ", got", errMap["errorcode"])
	}
	if errMap["error"] != errorSent.ErrorMsg {
		t.Error("'Expected 'error.error' to be ", errorSent.ErrorMsg, ", got", errMap["error"])
	}
}
