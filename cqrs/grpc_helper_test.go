package cqrs_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-openapi/swag"
	st "github.com/golang/protobuf/ptypes/struct"
	"github.com/jedrp/go-core/cqrs"
)

func TestToStruct(t *testing.T) {
	// msg := map[string]interface{}{
	// 	"name":        2222,
	// 	"description": 1,
	// }
	jsonStr := `{
		"message": {
			"name": 2222,
			"description":1,
			"isSale": true
		},
		"referring": {
		  "correlationID": "string",
		  "requestID": "string"
		},
		"routing": {
		  "action": "insert.product",
		  "entityID": "string",
		  "topic": "product"
		}
	  }`

	var res SendEventRequest

	if err := swag.ReadJSON([]byte(jsonStr), &res); err != nil {
		t.Error("read json error")
	}
	var msg *st.Struct
	if v, ok := res.Message.(map[string]interface{}); ok {
		k := v["isSale"]
		fmt.Println(reflect.TypeOf(k))
		msg = cqrs.ToStruct(v)
	}
	v := msg.Fields["name"].GetNumberValue()
	if v != 2222 {
		t.Error("value not correct")
	}
}

type SendEventRequest struct {

	// message to send
	// Required: true
	Message interface{} `json:"message"`

	// referred request info
	Referring *Referring `json:"referring,omitempty"`

	// routing info
	// Required: true
	Routing *Routing `json:"routing"`
}
type Referring struct {

	// correlation ID
	CorrelationID string `json:"correlationID,omitempty"`

	// request ID
	RequestID string `json:"requestID,omitempty"`
}
type Routing struct {

	// Action specific
	// Required: true
	Action *string `json:"action"`

	// For which entity
	//
	// Guid value field
	// Required: true
	EntityID *string `json:"entityID"`

	// Which opic the event need to go to
	// Required: true
	Topic *string `json:"topic"`
}
