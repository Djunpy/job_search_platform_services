package data_processing

import (
	"fmt"
	"strings"
)

func RemoveSAtEnd(arr []string) []string {
	for i, str := range arr {
		str = strings.TrimSpace(str)
		str = strings.ToLower(str)
		if strings.HasSuffix(str, "s") {
			arr[i] = str[:len(str)-1]
		}
	}
	return arr
}

func ExtractGroups(groupsInterface interface{}) ([]string, error) {
	var groups []string

	if g, ok := groupsInterface.([]interface{}); ok {
		for _, group := range g {
			if groupStr, ok := group.(string); ok {
				groups = append(groups, groupStr)
			} else {
				return nil, fmt.Errorf("group is not a string")
			}
		}
		return groups, nil
	} else {
		return nil, fmt.Errorf("groupsInterface is not a []interface{}")
	}
}
