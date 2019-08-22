package simpleconfig

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const DefaultTag = "d"

func uniq(input []string) []string {
	u := []string{}
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func extractFields(appConfig interface{}) []string {
	fieldNames := []string{}
	t := reflect.TypeOf(appConfig)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			subStrings := extractFields(reflect.ValueOf(appConfig).Field(i).Interface())
			prefix := t.Field(i).Name + "_"
			for _, s := range subStrings {
				fieldNames = append(fieldNames, strings.ToUpper(prefix+s))
			}
		} else {
			fieldNames = append(fieldNames, strings.ToUpper(t.Field(i).Name))
		}
	}
	return uniq(fieldNames)
}

func populateDefaultValue(key string, appConfig interface{}) error {
	if reflect.TypeOf(appConfig).Kind() != reflect.Ptr {
		return errors.New("Need a pointer value")
	}

	keyTree := strings.Split(key, "_")
	v := reflect.ValueOf(appConfig).Elem().FieldByNameFunc(
		func(s string) bool {
			return (strings.ToUpper(s) == keyTree[0])
		})
	oldv := reflect.ValueOf(appConfig).Elem()
	if len(keyTree) > 1 {
		for _, k := range keyTree[1:] {
			oldv = v
			v = v.FieldByNameFunc(func(s string) bool {
				return (strings.ToUpper(s) == k)
			})
		}
	}

	def, _ := oldv.Type().FieldByNameFunc(func(s string) bool {
		return (strings.ToUpper(s) == keyTree[len(keyTree)-1])
	})

	defaultValue := ""
	localDefaultValue, defaultFound := def.Tag.Lookup(DefaultTag)
	if !defaultFound {
		return nil
	}

	defaultValue = localDefaultValue

	/* Add a new type here! */
	if v.Type().Kind() == reflect.Int {
		writeValue, err := strconv.Atoi(defaultValue)
		if err == nil {
			v.Set(reflect.ValueOf(writeValue))
		}
		return err
	} else if v.Type().Kind() == reflect.Bool {
		writeValue, err := strconv.ParseBool(defaultValue)
		if err == nil {
			v.Set(reflect.ValueOf(writeValue))
		}
		return err
	}

	v.Set(reflect.ValueOf(defaultValue))
	return nil

}

func populateValue(key string, appConfig interface{}, value string) error {
	if reflect.TypeOf(appConfig).Kind() != reflect.Ptr {
		return errors.New("Need a pointer value")
	}

	keyTree := strings.Split(key, "_")
	v := reflect.ValueOf(appConfig).Elem().FieldByNameFunc(
		func(s string) bool {
			return (strings.ToUpper(s) == keyTree[0])
		})

	if len(keyTree) > 1 {
		for _, k := range keyTree[1:] {
			v = v.FieldByNameFunc(func(s string) bool {
				return (strings.ToUpper(s) == k)
			})
		}
	}

	/* Add a new type here! */
	if v.Type().Kind() == reflect.Int {
		writeValue, err := strconv.Atoi(value)
		if err == nil {
			v.Set(reflect.ValueOf(writeValue))
		}
		return err
	} else if v.Type().Kind() == reflect.Bool {
		writeValue, err := strconv.ParseBool(value)
		if err == nil {
			v.Set(reflect.ValueOf(writeValue))
		}
		return err
	}

	v.Set(reflect.ValueOf(value))
	return nil
}

func LoadConfig(cfg interface{}) {
	fldList := extractFields(reflect.ValueOf(cfg).Elem().Interface())

	fldMap := map[string]string{}

	for _, key := range fldList {
		if value, ok := os.LookupEnv(key); ok {
			fldMap[key] = value
		}
	}

	for _, key := range fldList {
		envValue, ok := fldMap[key]
		if ok {
			populateValue(key, cfg, envValue)
		} else {
			populateDefaultValue(key, cfg)
		}
	}
}
