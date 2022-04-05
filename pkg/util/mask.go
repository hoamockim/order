package util

import (
	"fmt"
	"reflect"
)

func Mask(obj interface{})  {
	idr := reflect.Indirect(reflect.ValueOf(obj))
	tElm := idr.Type()
	if tElm.Kind() == reflect.Struct {
		for i := 0; i < tElm.NumField(); i++ {
			fKind := tElm.Field(i).Type.Kind()
			if fKind == reflect.Ptr {
				Mask(idr.Field(i).Interface())
				continue
			}

			if fKind == reflect.Struct {
				fVal := Clone(idr.Field(i).Interface())
				Mask(fVal)
				f := idr.FieldByName(tElm.Field(i).Name)
				f.Set(reflect.Indirect(reflect.ValueOf(fVal)))
			}

			if isSensitive, ok := tElm.Field(i).Tag.Lookup("sensitive"); ok && isSensitive == "true" {
				if fKind != reflect.String {
					continue
				}

				input:=  idr.Field(i).String()
				maskData, _ := maskField(input)
				f := idr.FieldByName(tElm.Field(i).Name)
				f.Set(maskData)
			}
		}
	}
}

func maskFieldV2(f reflect.Value, val string) {
	var oStr string

	if len(val) <= 3 {
		for i := 0; i < len(val); i++ {
			oStr = fmt.Sprintf("%v*", oStr)
		}
	} else {
	}
	f.Set(reflect.ValueOf(oStr))
}

func maskField(input string) (reflect.Value, *string) {
	var output reflect.Value
	var oStr string

	if len(input) <= 3 {
		for i := 0; i < len(input); i++ {
			oStr = fmt.Sprintf("%v*", oStr)
		}
	} else {
		for i := 0; i < len(input); i++ {
			if i < len(input)-3 {
				oStr = fmt.Sprintf("%v*", oStr)
			} else {
				oStr = fmt.Sprintf("%v%v",oStr, string(input[i]))
			}
		}
	}

	output = reflect.ValueOf(oStr)
	return output, &oStr
}
