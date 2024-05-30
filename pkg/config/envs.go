package config

import (
	"fmt"
	"golang/pkg/config/helpers"
	"os"
	"reflect"
)

type EnvTypes struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_HOST     string
	POSTGRES_PORT     int
	MIGRATION_URL     string
}

var Envs EnvTypes

type castingOptions struct {
	DefaultValue interface{}
	IsRequired   bool
}

func handleEnv(key string, handler func(string) (interface{}, bool), options castingOptions) {
	envVar, found := os.LookupEnv(key)

	if !found {
		if options.DefaultValue != nil {
			envVar = options.DefaultValue.(string)
		} else if options.IsRequired {
			panic(fmt.Sprintf("missing required envVar: %s", key))
		} else {
			return
		}
	}

	value, isValidEnvVar := handler(envVar)

	if !isValidEnvVar {
		panic(fmt.Sprintf("invalid value for envVar: %s", key))
	}

	v := reflect.ValueOf(&Envs).Elem()
	field := v.FieldByName(key)
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(value))
	} else {
		panic(fmt.Sprintf("unable to set value for envVar: %s", key))
	}
}

func InitEnvVars() {
	envVarHandlers := map[string]struct {
		handler func(string) (interface{}, bool)
		options castingOptions
	}{
		"POSTGRES_PASSWORD": {helpers.IsString, castingOptions{IsRequired: true}},
		"POSTGRES_USER":     {helpers.IsString, castingOptions{IsRequired: true}},
		"POSTGRES_DB":       {helpers.IsString, castingOptions{IsRequired: true}},
		"POSTGRES_HOST":     {helpers.IsString, castingOptions{IsRequired: true, DefaultValue: "localhost"}},
		"POSTGRES_PORT":     {helpers.IsInt, castingOptions{IsRequired: true}},
		"MIGRATION_URL":     {helpers.IsString, castingOptions{IsRequired: true, DefaultValue: "file://../pkg/migrations"}},
	}

	for key, handlerOpts := range envVarHandlers {
		handleEnv(key, handlerOpts.handler, handlerOpts.options)
	}

	fmt.Printf("%+v\n", Envs)
}
