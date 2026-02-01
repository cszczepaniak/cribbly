package config

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	DSN      string `env:"DSN"`
	SeedUser struct {
		Username string `env:"SEED_USERNAME"`
		Password string `env:"SEED_PASSWORD"`
	}
	Environment string `env:"RAILWAY_ENVIRONMENT_NAME"`
}

func Load(cfg *Config) error {
	return populate(reflect.ValueOf(cfg))
}

func populate(val reflect.Value) error {
	if val.Kind() != reflect.Pointer {
		return errors.New("value must be a pointer")
	}

	typ := val.Elem().Type()
	for i := range val.Elem().NumField() {
		field := val.Elem().Field(i)
		if field.Kind() == reflect.Struct {
			err := populate(field.Addr())
			if err != nil {
				return err
			}
			continue
		}

		envVarName := typ.Field(i).Tag.Get("env")
		if envVarName == "" {
			continue
		}

		envVal, ok := os.LookupEnv(envVarName)
		if !ok {
			continue
		}

		err := setFromEnv(field, envVal)
		if err != nil {
			return err
		}
	}

	return nil
}

func setFromEnv(val reflect.Value, envVal string) error {
	tu, ok := val.Interface().(encoding.TextUnmarshaler)
	if ok {
		return tu.UnmarshalText([]byte(envVal))
	}

	switch val.Kind() {
	case reflect.String:
		val.SetString(envVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(envVal, 10, 64)
		if err != nil {
			return err
		}
		val.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(envVal, 10, 64)
		if err != nil {
			return err
		}
		val.SetUint(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(envVal)
		if err != nil {
			return err
		}
		val.SetBool(b)
	default:
		return fmt.Errorf("unsupported config kind: %v", val.Kind())
	}

	return nil
}
