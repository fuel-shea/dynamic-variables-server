package gdv

import (
	"bitbucket.org/kardianos/osext"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path"
)

type RuleSet struct {
	FeatureNames  []string
	VariableNames []string
	Rules         []Rule
}

type Rule struct {
	Features  map[string]Feature
	Variables map[string]interface{}
}

func NewRule() Rule {
	return Rule{
		Features:  make(map[string]Feature),
		Variables: make(map[string]interface{}),
	}
}

type Feature struct {
	Value string
	Exp   string
}

func VarsFromFeatures(featureMatch map[string]string, ruleSet RuleSet) map[string]interface{} {

	// go and find the first row that matches this criteria
	result := map[string]interface{}{}

	// run through the map for each critera and match against the incoming information
	for _, rule := range ruleSet.Rules {
		matched := true
		for featureName, feature := range rule.Features {
			matchValue := featureMatch[featureName]

			if feature.Value != "any" {
				if matchValue == "" {
					matched = false
					break
				} else {
					// use the expression to perform the test
					switch feature.Exp {
					case "<":
						if matchValue >= feature.Value {
							matched = false
							break
						}
					case ">":
						if matchValue <= feature.Value {
							matched = false
							break
						}
					case "=":
						if matchValue != feature.Value {
							matched = false
							break
						}
					}
				}
			}
		}
		if matched {
			result = rule.Variables
			break
		}
	}

	return result
}

func BuildDynoVars(gameId string) (RuleSet, error) {
	ruleSet := RuleSet{}

	valueStartIndex := math.MaxUint32
	rawCSVData, err := readCSVByGameId(gameId)
	if err != nil {
		return ruleSet, err
	}
	for rowIndex, row := range rawCSVData {
		if rowIndex == 0 {
			for colIndex, value := range row {
				if value == "|" {
					valueStartIndex = colIndex
				}
				if colIndex > valueStartIndex {
					// this is a variable
					ruleSet.VariableNames = append(ruleSet.VariableNames, value)
				} else if colIndex < valueStartIndex {
					if colIndex%2 == 0 {
						ruleSet.FeatureNames = append(ruleSet.FeatureNames, value)
					}
				}
			}
		} else {
			rule := NewRule()
			for colIndex, value := range row {
				if colIndex > valueStartIndex {
					// this is a variable
					variableName := ruleSet.VariableNames[colIndex-valueStartIndex-1]
					rule.Variables[variableName] = value

				} else if colIndex < valueStartIndex {
					featureName := ruleSet.FeatureNames[colIndex/2]
					var feature Feature
					if _, wasFound := rule.Features[featureName]; wasFound {
						feature = rule.Features[featureName]
					} else {
						feature = Feature{}
					}

					if colIndex%2 == 0 {
						// this is a feature value
						feature.Value = value
					} else {
						// this is a feature modifier
						feature.Exp = value
					}
					rule.Features[featureName] = feature
				}
			}

			ruleSet.Rules = append(ruleSet.Rules, rule)
		}
	}

	return ruleSet, nil
}

func readCSVByGameId(gameId string) ([][]string, error) {
	execDir, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}
	csvPath := path.Join(execDir, gameId+".csv")

	csvfile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return rawCSVdata, nil
}
