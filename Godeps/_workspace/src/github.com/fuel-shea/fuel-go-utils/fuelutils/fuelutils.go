package fuelutils

func InterfaceArr2StringArr(ifArray []interface{}) []string {
	nElems := len(ifArray)
	strs := make([]string, nElems)
	for i := 0; i < nElems; i++ {
		strs[i] = ifArray[i].(string)
	}
	return strs
}
