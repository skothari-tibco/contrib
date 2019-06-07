package action

import (
	"context"
	"errors"
	"path"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
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
	out := &Output{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}
	if input.Input == nil {
		return true, errors.New("Input not here")

	}

	factory := action.GetFactory(a.settings.Ref)

	var act action.Action
	settingsURI := make(map[string]interface{})

	switch path.Base(a.settings.Ref) {
	case "fps":
		settingsURI["catalystMlURI"] = a.settings.ResURI
	case "flow":
		settingsURI["flowURI"] = a.settings.ResURI
	}

	act, _ = factory.New(&action.Config{Settings: settingsURI})

	if syncAct, ok := act.(action.SyncAction); ok {
		result, _ := syncAct.Run(context.Background(), input.Input)

		out.Output = result

		ctx.SetOutputObject(out)

		return true, nil
	}

	return true, errors.New("Not a Sync Action")
}
