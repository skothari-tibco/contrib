package syncaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support"
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
	var ref string
	if a.settings.Ref[0] == '#' {
		ref, _ = support.GetAliasRef("action", a.settings.Ref[1:])
	} else {
		ref, _ = support.GetAliasRef("action", a.settings.Ref)
	}

	factory := action.GetFactory(ref)

	var act action.Action
	settingsURI := make(map[string]interface{})

	settingsURI["catalystMlURI"] = a.settings.ResURI //a.settings.ResURI

	act, _ = factory.New(&action.Config{Settings: settingsURI})

	if syncAct, ok := act.(action.SyncAction); ok {
		result, err := syncAct.Run(context.Background(), input.ToMap())
		if err != nil {
			fmt.Println("Error...", err)
		}

		out.Output = result

		ctx.SetOutputObject(out)

		return true, nil
	}

	return true, errors.New("Not a Sync Action")
}
