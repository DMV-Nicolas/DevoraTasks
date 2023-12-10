package util

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRequirements(t *testing.T) {
	type test struct {
		name         string
		obj          any
		requirements []requirement
		isErr        bool
	}

	tests := []test{
		{
			name: "Simple user data",
			obj: struct {
				Username string `json:"username" requirements:"required"`
				Email    string `json:"email" requirements:"required;email"`
				Password string `json:"password" requirements:"required;min=8;max=30"`
			}{},
			requirements: []requirement{
				{
					fieldName:  "Username",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Email",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Email",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "email",
					reqValue:   "true",
				},
				{
					fieldName:  "Password",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Password",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "min",
					reqValue:   "8",
				},
				{
					fieldName:  "Password",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "max",
					reqValue:   "30",
				},
			},
			isErr: false,
		},
		{
			name: "All data types",
			obj: struct {
				String  string  `requirements:"required"`
				Bool    bool    `requirements:"required"`
				Int     int     `requirements:"required"`
				Int8    int8    `requirements:"required"`
				Int16   int16   `requirements:"required"`
				Int32   int32   `requirements:"required"`
				Int64   int64   `requirements:"required"`
				Uint    uint    `requirements:"required"`
				Uint8   uint8   `requirements:"required"`
				Uint16  uint16  `requirements:"required"`
				Uint32  uint32  `requirements:"required"`
				Uint64  uint64  `requirements:"required"`
				Float32 float32 `requirements:"required"`
				Float64 float64 `requirements:"required"`
			}{},
			requirements: []requirement{
				{
					fieldName:  "String",
					fieldValue: "",
					fieldType:  "string",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Bool",
					fieldValue: "false",
					fieldType:  "bool",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Int",
					fieldValue: "0",
					fieldType:  "int",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Int8",
					fieldValue: "0",
					fieldType:  "int8",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Int16",
					fieldValue: "0",
					fieldType:  "int16",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Int32",
					fieldValue: "0",
					fieldType:  "int32",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Int64",
					fieldValue: "0",
					fieldType:  "int64",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Uint",
					fieldValue: "0",
					fieldType:  "uint",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Uint8",
					fieldValue: "0",
					fieldType:  "uint8",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Uint16",
					fieldValue: "0",
					fieldType:  "uint16",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Uint32",
					fieldValue: "0",
					fieldType:  "uint32",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Uint64",
					fieldValue: "0",
					fieldType:  "uint64",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Float32",
					fieldValue: "0",
					fieldType:  "float32",
					reqName:    "required",
					reqValue:   "true",
				},
				{
					fieldName:  "Float64",
					fieldValue: "0",
					fieldType:  "float64",
					reqName:    "required",
					reqValue:   "true",
				},
			},
			isErr: false,
		},
	}

	for _, x := range tests {
		t.Run(x.name, func(t *testing.T) {
			requirements := GetRequirements(x.obj)
			require.True(t, reflect.DeepEqual(requirements, x.requirements))
		})
	}
}

func TestVerifyRequirements(t *testing.T) {
	type testRequest struct {
		Username   string  `requirements:"required"`
		Email      string  `requirements:"required;email"`
		Password   string  `requirements:"required;min=8;max=20"`
		Age        int     `requirements:"required;min=18;max=60"`
		Money      float64 `requirements:"required;min=100;max=1000"`
		IsCriminal bool    `requirements:"required"`
	}

	type test struct {
		name    string
		request testRequest
		isErr   bool
	}

	tests := []test{
		{
			name: "Verify OK",
			request: testRequest{
				Username:   "dmvnicolas",
				Email:      "dmvnicolas@gmail.com",
				Password:   "83nicomoreno19",
				Age:        18,
				Money:      525.5,
				IsCriminal: true,
			},
			isErr: false,
		},
		{
			name:    "Verify required",
			request: testRequest{},
			isErr:   true,
		},
		{
			name: "Verify Min",
			request: testRequest{
				Username:   "dmvnicolas",
				Email:      "dmvnicolas@gmail.com",
				Password:   "123456",
				Age:        13,
				Money:      50,
				IsCriminal: true,
			},
			isErr: true,
		},
		{
			name: "Verify Max",
			request: testRequest{
				Username:   "dmvnicolas",
				Email:      "dmvnicolas@gmail.com",
				Password:   "contraseñahypermegalarga",
				Age:        75,
				Money:      2500.2,
				IsCriminal: true,
			},
			isErr: true,
		},
	}

	for _, x := range tests {
		t.Run(x.name, func(t *testing.T) {
			err := VerifyRequirements(x.request)
			require.Equal(t, x.isErr, err != nil)
		})
	}

}
