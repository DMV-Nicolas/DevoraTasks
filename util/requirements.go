package util

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var avalaibleRequirements = []string{
	"required", "email", "min", "max",
}

type requirement struct {
	fieldName  string
	fieldValue string
	fieldType  string
	reqName    string
	reqValue   string
}

func (r requirement) Verify() error {
	var err error
	switch r.reqName {
	case "required":
		err = r.Required()
	case "email":
		err = r.Email()
	case "min":
		err = r.Min()
	case "max":
		err = r.Max()
	}
	return err
}

func (r requirement) Required() (err error) {
	switch r.fieldType {
	case "string":
		valueWithoutSpaces := strings.TrimSpace(r.fieldValue)
		if valueWithoutSpaces == "" {
			err = fmt.Errorf("Requirements error: '%s' is required", r.fieldName)
		}
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float64", "float32":
		if r.fieldValue == "0" {
			err = fmt.Errorf("Requirements error: '%s' is required", r.fieldName)
		}
	case "bool":
		if r.fieldValue == "false" {
			err = fmt.Errorf("Requirements error: '%s' is required", r.fieldName)
		}
	}

	return err
}

func (r requirement) Min() (err error) {
	minNum, err := strconv.Atoi(r.reqValue)
	if err != nil {
		err = fmt.Errorf("Min requirement cannot be '%s', should be an integer", r.reqValue)
		panic(err)
	}

	var num float64
	switch r.fieldType {
	case "string":
		num = float64(len(r.fieldValue))
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		x, _ := strconv.Atoi(r.fieldValue)
		num = float64(x)
	case "float64", "float32":
		num, _ = strconv.ParseFloat(r.fieldValue, 64)
	}

	if num < float64(minNum) {
		err = fmt.Errorf("Requirements error: '%s' is less than the minimum value (%d)", r.fieldValue, minNum)
	}

	return err
}

func (r requirement) Max() error {
	maxNum, err := strconv.Atoi(r.reqValue)
	if err != nil {
		err = fmt.Errorf("Max requirement cannot be '%s', should be an integer", r.reqValue)
		panic(err)
	}

	var num float64
	switch r.fieldType {
	case "string":
		num = float64(len(r.fieldValue))
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		x, _ := strconv.Atoi(r.fieldValue)
		num = float64(x)
	case "float64", "float32":
		num, _ = strconv.ParseFloat(r.fieldValue, 64)
	}

	if num > float64(maxNum) {
		err = fmt.Errorf("Requirements error: '%s' is greater than the maximum value (%d)", r.fieldValue, maxNum)
	}

	return err
}

func (r requirement) Email() (err error) {
	regexPattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	isEmail := regexPattern.MatchString(r.fieldValue)
	if !isEmail {
		err = fmt.Errorf("Requirements error: '%s' is not an Email", r.fieldValue)
	}
	return err
}

func GetRequirements(obj any) []requirement {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	numFields := t.NumField()

	requirements := make([]requirement, 0, numFields)
	for i := 0; i < numFields; i++ {
		fieldName := t.Field(i).Name
		fieldType := t.Field(i).Type.String()
		fieldValue := ""

		switch fieldType {
		case "string":
			fieldValue = v.Field(i).String()
		case "int", "int8", "int16", "int32", "int64":
			fieldValue = fmt.Sprint(v.Field(i).Int())
		case "uint", "uint8", "uint16", "uint32", "uint64":
			fieldValue = fmt.Sprint(v.Field(i).Uint())
		case "float64", "float32":
			fieldValue = fmt.Sprint(v.Field(i).Float())
		case "bool":
			fieldValue = fmt.Sprint(v.Field(i).Bool())
		}

		tags := strings.Split(t.Field(i).Tag.Get("requirements"), ";")
		if tags[0] == "" {
			continue
		}
		for _, tag := range tags {
			index := strings.Index(tag, "=")
			reqName := tag
			reqValue := "true"

			if index > 0 {
				reqName = tag[:index]
				if len(tag)-1 == index {
					err := fmt.Errorf("Requirement '%s' has no value", reqName)
					panic(err)
				}
				reqValue = tag[index+1:]
			} else if !slices.Contains(avalaibleRequirements, reqName) {
				err := fmt.Errorf("'%s' requirement is not avalaible", reqName)
				panic(err)
			}

			requirements = append(requirements, requirement{
				fieldName:  fieldName,
				fieldValue: fieldValue,
				fieldType:  fieldType,
				reqName:    reqName,
				reqValue:   reqValue,
			})
		}

	}

	return requirements
}

func VerifyRequirements(obj any) error {
	requirements := GetRequirements(obj)
	var errors string
	for _, r := range requirements {
		err := r.Verify()
		if err != nil {
			errors += fmt.Sprint(err.Error()) + ";"
		}
	}

	if errors == "" {
		return nil
	} else {
		return fmt.Errorf(errors)
	}
}
