package utils

//
func DelSlice(origin []interface{}, items ...interface{}) {

	for _, item := range items {
		for i := 0; i < len(origin); i++ {
			if origin[i] == item {
				origin = append(origin[:i], origin[i+1:]...)
				i-- // maintain the correct index
			}
		}
	}

}

//
func DelStrSlice(origin []string, items ...string) {

	for _, item := range items {
		for i := 0; i < len(origin); i++ {
			if origin[i] == item {
				origin = append(origin[:i], origin[i+1:]...)
				i-- // maintain the correct index
			}
		}
	}

}
