package fuelutils

import (
	"errors"
	"fmt"
)

func InterfaceArr2StringArr(ifArray []interface{}) ([]string, error) {
	nElems := len(ifArray)
	strs := make([]string, nElems)
	for i := 0; i < nElems; i++ {
		var isString bool
		strs[i], isString = ifArray[i].(string)
		if !isString {
			return []string{}, errors.New(fmt.Sprintf("Tried to convert %v to string", ifArray[i]))
		}
	}
	return strs, nil
}
