// This file contains the implementation of BlankSparkable.

package app

import (
	"blank-sparkable/util"
	"fmt"
	"github.com/Bitspark/go-bitnode/bitnode"
	"os"
	"reflect"
)

// Struct definition for BlankSparkable.

// BlankSparkable is the main sparkable.
// @@SPARKABLE_DESCRIPTION@@
type BlankSparkable struct {
	bitnode.System

	// @@SPARKABLE_FIELDS@@
}

// BlankSparkable methods.

// @@METHOD_STUBS@@

// @@HANDLER_STUBS@@

// Lifecycle callbacks.

// lifecycleCreate is called when the container has been created.
func (s *BlankSparkable) lifecycleCreate(vals ...bitnode.HubItem) error {
	// TODO: Add startup logic here which is called when the spark is created.
	return nil
}

// lifecycleLoad is called when the container has been started (after lifecycleCreate) or restarted.
func (s *BlankSparkable) lifecycleLoad() error {
	// TODO: Add startup logic here which is called after the spark has been created.

	s.LogInfo("BlankSparkable running...")
	s.SetStatus(bitnode.SystemStatusRunning)

	return nil
}

// lifecycleStop is called when the container is started.
func (s *BlankSparkable) lifecycleStop(timeout float64) error {
	// TODO: Add cleanup logic here which is called when the spark is stopped.

	s.SetStatus(bitnode.SystemStatusStopping)
	os.Exit(0)

	return nil
}

// DO NOT CHANGE THE FOLLOWING CODE UNLESS YOU KNOW WHAT YOU ARE DOING.

func (s *BlankSparkable) Update(values ...string) error {
	sv := reflect.ValueOf(*s)
	st := reflect.TypeOf(*s)
	if len(values) == 0 {
		for i := 0; i < st.NumField(); i++ {
			values = append(values, st.Field(i).Name)
		}
	}
	for _, value := range values {
		ft, ok := st.FieldByName(value)
		if !ok {
			return fmt.Errorf("field '%s' not found in %s", value, st.Name())
		}
		fv := sv.FieldByName(value)
		if !fv.IsValid() {
			return fmt.Errorf("field '%s' not found in %s", value, st.Name())
		}
		val, err := util.InterfaceFromValue(fv.Interface())
		if err != nil {
			return err
		}
		hubName := ft.Tag.Get("json")
		if err := s.GetHub(hubName).Set("", val); err != nil {
			return err
		}
	}
	return nil
}

// Init attaches the methods of the BlankSparkable to the respective handlers.
func (s *BlankSparkable) Init() error {
	// METHODS

	// @@METHOD_HANDLERS@@

	// VALUES

	// @@VALUE_HANDLERS@@

	// CHANNELS

	// @@CHANNEL_HANDLERS@@

	// LIFECYCLE EVENTS

	s.AddCallback(bitnode.LifecycleCreate, bitnode.NewNativeEvent(func(vals ...bitnode.HubItem) error {
		return s.lifecycleCreate(vals...)
	}))

	s.AddCallback(bitnode.LifecycleLoad, bitnode.NewNativeEvent(func(vals ...bitnode.HubItem) error {
		return s.lifecycleLoad()
	}))

	s.AddCallback(bitnode.LifecycleStop, bitnode.NewNativeEvent(func(vals ...bitnode.HubItem) error {
		timeout := vals[0].(float64)
		return s.lifecycleStop(timeout)
	}))

	return nil
}
