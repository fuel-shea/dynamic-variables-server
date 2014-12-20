package gdv

import (
	"bitbucket.org/kardianos/osext"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path"
)

type FeatureMap []map[string]map[string]interface{}

func ResultFromFeatures(featureMatch map[string]string, featureMap FeatureMap) map[string]interface{} {

	// go and find the first row that matches this criteria
	result := map[string]interface{}{}

	// run through the map for each critera and match against the incoming information
	for _, featureLine := range featureMap {
		matched := true
		for criteriaName, criteria := range featureLine["Criteria"] {
			// compare the criteria against the incoming
			matchValue := featureMatch[criteriaName]
			criteriaStringMap := criteria.(map[string]string)
			criteriaValue := criteriaStringMap["Value"]
			criteriaExp := criteriaStringMap["Exp"]

			if criteriaValue != "any" {
				if matchValue == "" {
					matched = false
					break
				} else {
					// use the expression to perform the test
					switch criteriaExp {
					case "<":
						if matchValue >= criteriaValue {
							matched = false
							break
						}
					case ">":
						if matchValue <= criteriaValue {
							matched = false
							break
						}
					case "=":
						if matchValue != criteriaValue {
							matched = false
							break
						}
					}

				}
			}
		}
		if matched {
			matchResult := featureLine["Result"]
			result = matchResult
			break
		}

	}

	return result
}

func ReadDynamicsMap(gameId string) (FeatureMap, error) {

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

	// sanity check, display to standard output

	featureMap := FeatureMap{}

	valueStartIndex := math.MaxUint32
	featureNames := []string{}
	variableNames := []string{}
	for rowIndex, row := range rawCSVdata {
		if rowIndex == 0 {
			for colIndex, value := range row {
				if value == "|" {
					valueStartIndex = colIndex
				}
				if colIndex > valueStartIndex {
					// this is a variable
					variableNames = append(variableNames, value)
				} else if colIndex < valueStartIndex {
					if colIndex%2 == 0 {
						featureNames = append(featureNames, value)
					}
				}
			}
		} else {
			featureLine := map[string]map[string]interface{}{}
			features := map[string]interface{}{}
			variables := map[string]interface{}{}

			for colIndex, value := range row {
				if colIndex > valueStartIndex {
					// this is a variable
					variableName := variableNames[colIndex-valueStartIndex-1]
					variables[variableName] = value
				} else if colIndex < valueStartIndex {
					featureName := featureNames[colIndex/2]
					if colIndex%2 == 0 {
						// this is a feature value
						if features[featureName] == nil {
							feature := map[string]string{}
							features[featureName] = feature
						}
						existingFeature := features[featureName].(map[string]string)
						existingFeature["Value"] = value
					} else {
						// this is a feature modifier
						if features[featureName] == nil {
							feature := map[string]string{}
							features[featureName] = feature
						}
						existingFeature := features[featureName].(map[string]string)
						existingFeature["Exp"] = value
					}
				}
			}

			featureLine["Criteria"] = features
			featureLine["Result"] = variables

			featureMap = append(featureMap, featureLine)

		}
	}

	return featureMap, nil
}
