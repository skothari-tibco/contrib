package action

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	ResURI string `md:"resURI"`
}

type Input struct {
	FPSInput map[string]interface{} `md:"fpsinput"`
}

type Output struct {
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fpsinput": i.FPSInput,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error

	i.FPSInput, err = coerce.ToObject(values["fpsinput"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	//var err error

	return nil
}
