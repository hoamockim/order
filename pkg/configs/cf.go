package configs

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

var app appConfig

func devEnv() {
	os.Setenv("MONGOURL", "mongodb://localhost:27017")
	os.Setenv("MONGO_DATABASE", "tipee_cart")

}

func InitConfig() {
	app = appConfig{Mongo: Mongo{}}
	bindingFromEnv(&app)
}

func bindingFromEnv(entity interface{}) {
	idr := reflect.Indirect(reflect.ValueOf(entity))
	tElm := idr.Type()
	if tElm.Kind() != reflect.Struct && tElm.Kind() != reflect.Ptr {
		panic("Error type config")
	}

	for i := 0; i < tElm.NumField(); i++ {
		fKind := tElm.Field(i).Type.Kind()
		if fKind == reflect.Ptr {
			bindingFromEnv(idr.Field(i).Interface())
			continue
		}

		if fKind == reflect.Struct {
			fVal := reflect.New(idr.Field(i).Type()).Interface()
			bindingFromEnv(fVal)
			f := idr.FieldByName(tElm.Field(i).Name)
			f.Set(reflect.Indirect(reflect.ValueOf(fVal)))
		}
		if !isAtomic(fKind) {
			continue
		}
		bindingField(tElm, idr, i)
	}
}

func bindingField(tElm reflect.Type, idr reflect.Value, idx int) {
	if envKey, ok := tElm.Field(idx).Tag.Lookup("env"); ok && strings.TrimSpace(envKey) != "" {
		envVal := os.Getenv(envKey)
		fieldName := tElm.Field(idx).Name
		f := idr.FieldByName(fieldName)
		setValueForField(tElm.Field(idx).Type.Kind(), f, envVal)
	}
}

func setValueForField(fKind reflect.Kind, f reflect.Value, val string) {
	switch fKind {
	case reflect.Int:
		iVal, _ := strconv.Atoi(val)
		f.Set(reflect.ValueOf(iVal))
		break
	case reflect.Uint:
		iVal, _ := strconv.Atoi(val)
		f.Set(reflect.ValueOf(absInt(iVal)))
		break
	case reflect.String:
		f.Set(reflect.ValueOf(val))
		break
	case reflect.Bool:
		bVal, _ := strconv.ParseBool(val)
		f.Set(reflect.ValueOf(bVal))
		break
	default:
		break
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isAtomic(kind reflect.Kind) bool {
	return (kind == reflect.Int) ||
		(kind == reflect.Int8) ||
		(kind == reflect.Int16) ||
		(kind == reflect.Int32) ||
		(kind == reflect.Int64) ||
		(kind == reflect.Uint) ||
		(kind == reflect.Uint8) ||
		(kind == reflect.Uint16) ||
		(kind == reflect.Uint32) ||
		(kind == reflect.Uint64) ||
		(kind == reflect.String) ||
		(kind == reflect.Bool)
}
