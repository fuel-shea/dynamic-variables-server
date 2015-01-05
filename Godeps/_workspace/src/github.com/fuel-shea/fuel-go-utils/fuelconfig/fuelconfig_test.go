package fuelconfig_test

import (
	"bitbucket.org/kardianos/osext"
	"fmt"
	"github.com/fuel-shea/fuel-go-utils/fuelconfig"
	"os"
	"path"
	"testing"
)

type testCase struct {
	env     string
	appName string
	dbHost  string
	dbName  string
}

var testCases = []testCase{
	testCase{"TESTCONFIG", "testconfigapp", "testinghost", "testingname"},
}

func TestCreateConfig(t *testing.T) {
	for _, tc := range testCases {
		runTestCase(t, tc)
	}
}

func runTestCase(t *testing.T, tc testCase) {
	// change the APP_ENV to something that doesn't clash with "TEST" (currently in use)
	origEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", tc.env)
	defer os.Setenv("APP_ENV", origEnv)

	// determine where to put the config file
	configFilename := tc.appName + "." + tc.env + ".config.json"
	execDir, err := osext.ExecutableFolder()
	if err != nil {
		t.Fatal(err)
	}
	configPath := path.Join(execDir, configFilename)
	defer os.Remove(configPath)

	// write the JSON to the config file
	configFile, err := os.Create(configPath)
	if err != nil {
		t.Fatal(err)
	}
	configFile.WriteString(fmt.Sprintf(`{"DBHost": %q, "DBName": %q}`, tc.dbHost, tc.dbName))
	defer configFile.Close()

	// run the function to test
	conf, err := fuelconfig.CreateConfig(tc.appName)
	if err != nil {
		t.Error(err)
	}
	if conf.DBHost != tc.dbHost {
		t.Error("Expected conf.DBHost to be '", tc.dbHost, "', but got", conf.DBHost)
	}
	if conf.DBName != tc.dbName {
		t.Error("Expected conf.DBName to be '", tc.dbName, "', but got", conf.DBHost)
	}
}
