package runpod

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var CloudTypeSecure = "SECURE"
var CloudTypeCommunity = "COMMUNITY"
var CloudTypeAll = "ALL"

var DefaultCloudType = CloudTypeSecure

func (c *Client) mapToStruct(map_ map[string]any, struct_ any) (err error) {
	jsonString, err := json.Marshal(map_)
	if err != nil {
		return &ErrMapToStruct{Msg: err.Error()}
	}
	err = json.Unmarshal(jsonString, &struct_)
	if err != nil {
		return &ErrMapToStruct{Msg: err.Error()}
	}
	return nil
}

type ErrMapToStruct struct {
	Msg string
}

func (e *ErrMapToStruct) Error() string {
	return fmt.Sprintf("error converting map[string]any to type struct; %s", e.Msg)
}

func ErrorIs(err error, type_ any) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(type_)
}
