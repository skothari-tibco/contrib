package action

import (
	"errors"
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/fps/pipeline"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}

	//ctx.Logger().Debugf("flowURI: %+v", s.FlowURI)

	return act, nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Activity is an Activity that is used to log a message to the console
// inputs : {message, flowInfo}
// outputs: none
type Activity struct {
	settings *Settings
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}

	input.FPSInput = make(map[string]interface{})
	input.FPSInput["input"] = 4

	res := pipeline.GetManager().GetResource(a.settings.ResURI)
	if res != nil {
		def, ok := res.Object().(*pipeline.Definition)
		if !ok {
			return true, errors.New("unable to resolve fps: " + a.settings.ResURI)
		}
		inst := pipeline.NewInstance(def, "instId", ctx.Logger())

		output, err := inst.Run(input.FPSInput)

		if err != nil {
			return true, err
		}

		fmt.Println("Output..", output)
	}

	ctx.Logger().Info("Something")

	return true, nil
}
