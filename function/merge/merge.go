package merge

import (
	"errors"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	_ = function.Register(&fnMergeMap{})
}

type fnMergeMap struct {
}

func (fnMergeMap) Name() string {
	return "mergeMaps"
}

func (fnMergeMap) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeArray}, false
}

func (fnMergeMap) Eval(params ...interface{}) (interface{}, error) {

	if len(params) == 0 {
		return nil, errors.New("Array of map is empty")
	}

	result := make(map[string][]interface{})

	temp, err := coerce.ToArray(params[0])

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(temp); i++ {

		if temp[i] != nil {

			for key, val := range temp[i].(map[string]interface{}) {
				result[key] = append(result[key], val)
			}
		}

	}
	return result, nil
}
