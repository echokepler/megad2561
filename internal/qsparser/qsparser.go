// Package qsparser (querystring parser) нужен для двустороннего биндинга
// структуры и параметров.
package qsparser

import (
	"github.com/echokepler/megad2561/core"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const tagName = "qs"

var settingRe = regexp.MustCompile(`(?m)setting`)
var setterRe = regexp.MustCompile(`(?m)setter`)
var requiredRe = regexp.MustCompile(`(?m)required`)

type MarshalOptions struct {
	OnlySetters  bool
	OnlySettings bool
}

func Marshal(s interface{}, opts MarshalOptions) core.ServiceValues {
	values := core.ServiceValues{}
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = reflect.Indirect(v)
	}

	for i := 0; i < t.NumField(); i++ {
		value := ""
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		tagArr := strings.Split(tag, ",")
		tagValue := tagArr[len(tagArr)-1]
		isSettingMode := true
		isSetterMode := false
		isRequiredMode := false

		if tagValue == "skip" {
			continue
		}

		if field.Type.Kind() == reflect.Ptr {
			iField := reflect.Indirect(v.Field(i))

			if iField.CanInterface() {
				deepValues := Marshal(reflect.Indirect(v.Field(i)).Interface(), opts)

				for key := range deepValues {
					values.Add(key, deepValues.Get(key))
				}
			}

			//fmt.Println(reflect.Indirect(v.Field(i)).Interface(), tagValue)

			continue
		}

		if field.Type.Kind() == reflect.Struct {
			deepValues := Marshal(v.Field(i).Interface(), opts)

			for key := range deepValues {
				values.Add(key, deepValues.Get(key))
			}
		}

		if len(tag) == 0 {
			continue
		}

		if len(tagArr) > 1 {
			isSetterMode = setterRe.MatchString(tag)
			isSettingMode = settingRe.MatchString(tag)
			isRequiredMode = requiredRe.MatchString(tag)
		}

		if !isRequiredMode {
			if opts.OnlySettings && !isSettingMode {
				continue
			}

			if opts.OnlySetters && !isSetterMode {
				continue
			}
		}

		valField := reflect.ValueOf(v.Field(i).Interface())

		switch field.Type.Kind() {
		case reflect.String:
			value = valField.String()

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(valField.Uint(), 10)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(valField.Int(), 10)

		case reflect.Bool:
			value = strconv.FormatBool(valField.Bool())
		}

		values.Add(tagValue, value)
	}

	return values
}

func UnMarshal(values core.ServiceValues, s interface{}) error {
	v := reflect.ValueOf(s)
	in := reflect.Indirect(v)
	t := in.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		qsValue := values.Get(tag)

		if len(qsValue) == 0 {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			in.Field(i).SetString(qsValue)

		case reflect.Uint, reflect.Uint8:
			val, err := strconv.ParseUint(qsValue, 10, 8)
			if err != nil {
				return err
			}
			in.Field(i).SetUint(val)

		case reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(qsValue, 10, 64)
			if err != nil {
				return err
			}
			in.Field(i).SetUint(val)

		case reflect.Int, reflect.Int8:
			val, err := strconv.ParseInt(qsValue, 10, 8)
			if err != nil {
				return err
			}
			in.Field(i).SetInt(val)

		case reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(qsValue, 10, 64)
			if err != nil {
				return err
			}
			in.Field(i).SetInt(val)

		case reflect.Bool:
			val, err := strconv.ParseBool(qsValue)
			if err != nil {
				return err
			}
			in.Field(i).SetBool(val)
		}

	}

	return nil
}
