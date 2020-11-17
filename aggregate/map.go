package main

func trimMap(orig map[string]int, length int) map[string]int {
	if len(orig) < length {
		return orig
	}

	var orderedKeys []string
	var orderedValues []int

	for k1, v1 := range orig {
		if len(orderedValues) == 0 || orderedValues[(len(orderedValues)-1)] > v1 {
			orderedValues = append(orderedValues, v1)
			orderedKeys = append(orderedKeys, k1)
		} else {
			for index, v2 := range orderedValues {
				if v1 > v2 {
					orderedValues = intIndexInsert(orderedValues, index, v1)
					orderedKeys = strIndexInsert(orderedKeys, index, k1)
					break
				}
			}
		}
	}

	for _, k := range orderedKeys[length:] {
		delete(orig, k)
	}

	return orig
}

func strIndexInsert(slice []string, index int, value string) []string {
	slice = append(slice, "")
	copy(slice[(index+1):], slice[index:])
	slice[index] = value
	return slice
}

func intIndexInsert(slice []int, index int, value int) []int {
	slice = append(slice, 0)
	copy(slice[(index+1):], slice[index:])
	slice[index] = value
	return slice
}

func copyMap(originalMap map[string]int) map[string]int {
	newMap := make(map[string]int)

	for key, value := range originalMap {
		newMap[key] = value
	}

	return newMap
}

func mapDiff(orig map[string]int, new map[string]int) bool {
	if len(orig) != len(new) {
		return true
	}

	for k, v := range orig {
		if val, ok := new[k]; ok {
			if val != v {
				return true
			}
		} else {
			return true
		}
	}

	return false
}
