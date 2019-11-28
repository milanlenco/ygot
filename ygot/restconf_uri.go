package ygot

import (
	"fmt"
	"strings"
)

type pathElem struct {
	name    string
	module  string
	keyVals []interface{}
}

// GetRestconfURI returns RESTCONF data resource identifier for the given
// container/list element.
// For more information, see RFC8040, section 3.5.3.
func GetRestconfURI(obj ValidatedGoStructExtended) (string, error) {
	schemaPath, module, _ := obj.GetSchemaPath()

	// construct a slice of objects representing the branch in the YANG tree
	var branch []ValidatedGoStructExtended
	for obj != nil {
		if _, module, _ := obj.GetSchemaPath(); module == "" {
			// fake root
			break
		}
		branch = append([]ValidatedGoStructExtended{obj}, branch...)
		obj = obj.GetParent()
	}

	// create a slice of entries, one for each element in the path
	var pathElems []pathElem
	var topModule string
	for i, elem := range strings.Split(schemaPath, "/") {
		if i == 1 {
			topModule = elem // top module
		}
		if i < 2 {
			// skip leading forward slash and the top module name
			continue
		}
		pathElems = append(pathElems, pathElem{name: elem})
	}

	j := 0
	var path string
	prevModule := topModule
	for _, obj := range branch {
		// set module for element and all the preceding that were compressed
		path, module, _ = obj.GetSchemaPath()
		i := len(strings.Split(path, "/"))-3
		for ; j < i; j++ {
			pathElems[j].module = prevModule
		}
		pathElems[i].module = module
		prevModule = module

		// collect key values if this is a list entry
		if multiKeyList, ok := obj.(MultiKeyHelperGoStruct); ok {
			keyMap, err := multiKeyList.ΛListKeyMap()
			if err != nil {
				return "", err
			}
			for _, key := range multiKeyList.ΛOrderedListKeys() {
				keyVal := keyMap[key]
				pathElems[i].keyVals = append(pathElems[i].keyVals, keyVal)
			}
		} else if singleKeyList, ok := obj.(KeyHelperGoStruct); ok {
			keyMap, err := singleKeyList.ΛListKeyMap()
			if err != nil {
				return "", err
			}
			for _, keyVal := range keyMap {
				// single iteration
				pathElems[i].keyVals = append(pathElems[i].keyVals, keyVal)
			}
		}
	}

	// finally construct the RESTCONF URI
	prevModule = ""
	var uri strings.Builder
	for _, elem := range pathElems {
		uri.WriteRune('/')
		if elem.module != prevModule {
			uri.WriteString(elem.module)
			uri.WriteRune(':')
			prevModule = elem.module
		}
		uri.WriteString(elem.name)
		if len(elem.keyVals) > 0 {
			uri.WriteRune('=')
			for i, keyVal := range elem.keyVals {
				if i > 0 {
					uri.WriteRune(',')
				}
				uri.WriteString(fmt.Sprintf("%v", keyVal))
			}
		}
	}

	return uri.String(), nil
}
