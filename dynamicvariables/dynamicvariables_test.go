package dynamicvariables_test

import (
	"dynamic-variables-server/dynamicvariables"
	"gopkg.in/mgo.v2"
	"testing"
)

var gameID = "gid"

type Variable struct {
	GameID       string `bson:"game_id"`
	RuleIdx      int    `bson:"rule_idx"`
	RandomMax    string `bson:"randomMax"`
	WhammyChance string `bson:"whammyChance"`
}

func TestVarsFromFeatures_equals(t *testing.T) {
	mgoDBHost := "localhost"
	mgoDBName := "dynamicvariables_test"

	for _, tc := range testCases {
		err := populateDB(mgoDBHost, mgoDBName, tc)
		if err != nil {
			t.Fatal(err)
		}

		dvFactory, err := dynamicvariables.NewDynoVarFactory(mgoDBHost, mgoDBName)
		if err != nil {
			t.Fatal(err)
		}
		dvSource := dvFactory.NewDynoVarSource()

		varsResult, err := dvSource.VarsFromFeatures(tc.query, gameID)
		if err != nil {
			t.Fatal(err)
		}

		if varsResult["whammyChance"] != tc.expected["whammyChance"] {
			t.Error("Expected 'whammyChance' to be %q, but it was %v", tc.expected["whammyChance"], varsResult["whammyChance"])
		}
		if varsResult["randomMax"] != tc.expected["randomMax"] {
			t.Error("Expected 'randomMax' to be %q, but it was %v", tc.expected["randomMax"], varsResult["randomMax"])
		}
	}
}

type testCase struct {
	gameRuleData dynamicvariables.GameRuleData
	features     [][]dynamicvariables.Feature
	variables    []Variable
	query        map[string]interface{}
	expected     map[string]interface{}
}

var testCases = []testCase{
	{
		gameRuleData: dynamicvariables.GameRuleData{
			NumRules: 6,
		},
		features: [][]dynamicvariables.Feature{
			{
				dynamicvariables.Feature{Val: "CA", Mod: "="},
				dynamicvariables.Feature{Val: "iOS", Mod: "="},
				dynamicvariables.Feature{Val: "M", Mod: "="},
			},
			{
				dynamicvariables.Feature{Val: "US", Mod: "="},
				dynamicvariables.Feature{Val: "iOS", Mod: "="},
				dynamicvariables.Feature{Val: "M", Mod: "="},
			},
			{
				dynamicvariables.Feature{Val: "JP", Mod: "="},
				dynamicvariables.Feature{Val: "Android", Mod: "="},
				dynamicvariables.Feature{Val: "any", Mod: "="},
			},
			{
				dynamicvariables.Feature{Val: "CA", Mod: "="},
				dynamicvariables.Feature{Val: "Android", Mod: "="},
				dynamicvariables.Feature{Val: "F", Mod: "="},
			},
			{
				dynamicvariables.Feature{Val: "JP", Mod: "="},
				dynamicvariables.Feature{Val: "iOS", Mod: "="},
				dynamicvariables.Feature{Val: "any", Mod: "="},
			},
			{
				dynamicvariables.Feature{Val: "JP", Mod: "="},
				dynamicvariables.Feature{Val: "iOS", Mod: "="},
				dynamicvariables.Feature{Val: "F", Mod: "="},
			},
			{
				dynamicvariables.Feature{Val: "any", Mod: "="},
				dynamicvariables.Feature{Val: "any", Mod: "="},
				dynamicvariables.Feature{Val: "any", Mod: "="},
			},
		},
		variables: []Variable{
			Variable{RandomMax: "0", WhammyChance: "0"},
			Variable{RandomMax: "1", WhammyChance: "1"},
			Variable{RandomMax: "2", WhammyChance: "2"},
			Variable{RandomMax: "3", WhammyChance: "3"},
			Variable{RandomMax: "4", WhammyChance: "4"},
			Variable{RandomMax: "5", WhammyChance: "5"},
			Variable{RandomMax: "6", WhammyChance: "6"},
		},
		query: map[string]interface{}{
			"Country": "JP",
			"Device":  "iOS",
			"Gender":  "F",
		},
		expected: map[string]interface{}{
			"randomMax":    "4",
			"whammyChance": "4",
		},
	},
}

func populateDB(mgoDBHost, mgoDBName string, tc testCase) error {
	tc.gameRuleData.GameID = gameID
	tc.gameRuleData.FeatureTypes = []string{"Country", "Device", "Gender"}
	tc.gameRuleData.VariableTypes = []string{"randomMax", "whammyChance"}

	mgoSess, err := mgo.Dial(mgoDBHost)
	if err != nil {
		return err
	}
	mgoDB := mgoSess.DB(mgoDBName)
	err = mgoDB.DropDatabase()
	if err != nil {
		return err
	}

	err = mgoDB.C("game_rule_data").Insert(tc.gameRuleData)
	if err != nil {
		return err
	}

	featuresColl := mgoDB.C("features")
	variablesColl := mgoDB.C("variables")
	for ruleIdx, _ := range tc.features {
		variable := tc.variables[ruleIdx]
		variable.GameID = gameID
		variable.RuleIdx = ruleIdx
		err := variablesColl.Insert(variable)
		if err != nil {
			return err
		}

		for featIdx, feat := range tc.features[ruleIdx] {
			feat.GameID = gameID
			feat.RuleIdx = ruleIdx
			feat.Type = tc.gameRuleData.FeatureTypes[featIdx]
			err := featuresColl.Insert(feat)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
