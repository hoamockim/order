package util

import "reflect"

func Clone(obj interface{}) interface{} {
	clone := cloneS(obj)
	copyV(obj, clone)
	return clone
}

func cloneS(obj interface{}) interface{} {
	var newObj reflect.Value
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Struct:
		newObj = reflect.New(reflect.TypeOf(obj))
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		newObj = reflect.New(reflect.TypeOf(obj).Elem())
	default:
		newObj = reflect.ValueOf(obj)
	}
	return newObj.Interface()
}

func copyV(src interface{}, dst interface{}) {
	idrS := reflect.Indirect(reflect.ValueOf(src))
	idrD := reflect.Indirect(reflect.ValueOf(dst))
	elmS := idrS.Type()
	if elmS.Kind() == reflect.Struct {
		for i := 0; i < elmS.NumField(); i++ {
			fKind := elmS.Field(i).Type.Kind()
			field := idrS.Field(i)
			fieldName := elmS.Field(i).Name
			f := idrD.FieldByName(fieldName)
			switch fKind {
			case reflect.Array, reflect.Slice:
				fVal := reflect.ValueOf(field.Interface())
				nVal := copyAr(fVal, reflect.New(field.Type()).Elem(), fKind == reflect.Array)
				f.Set(nVal)
				break
			case reflect.Ptr, reflect.Chan, reflect.Interface:
				if field.IsNil() || field.IsZero() {
					break
				}
				itfType := reflect.Indirect(reflect.ValueOf(field.Interface())).Type()
				nVal := reflect.New(itfType)
				copyV(field.Interface(), nVal.Interface())
				f.Set(nVal)
				break
			case reflect.Struct:
				nVal := reflect.ValueOf(field.Interface())
				f.Set(nVal)
				break
			default:
				nVal := atomic(field)
				f.Set(reflect.ValueOf(nVal))
				break
			}
		}
	} else {
		nVal := atomic(idrS)
		idrD.Set(reflect.ValueOf(nVal))
	}
}

func copyAr(src reflect.Value, dst reflect.Value, isArr bool) reflect.Value {
	for i := 0; i < src.Len(); i++ {
		fIdx := src.Index(i)
		var nfVal interface{}
		if fIdx.Kind() == reflect.Struct || fIdx.Kind() == reflect.Ptr ||
			fIdx.Kind() == reflect.Interface || fIdx.Kind() == reflect.Chan {
			itfType := reflect.Indirect(reflect.ValueOf(fIdx.Interface())).Type()
			nObj := reflect.New(itfType)
			copyV(fIdx.Interface(), nObj.Interface())
			nfVal = nObj.Interface()
			if isArr {
				dst.Index(i).Set(reflect.Indirect(reflect.ValueOf(nfVal)))
			} else {
				dst = reflect.Append(dst, reflect.Indirect(reflect.ValueOf(nfVal)))
			}

		} else {
			nfVal = atomic(fIdx)
			if isArr {
				dst.Index(i).Set(reflect.ValueOf(nfVal))
			} else {
				dst = reflect.Append(dst, reflect.ValueOf(nfVal))
			}
		}
	}
	return dst
}

func atomic(src reflect.Value) (dst interface{}) {
	switch src.Kind() {
	case reflect.Int:
		dst = int(src.Int())
	case reflect.Int8:
		dst = int8(src.Int())
	case reflect.Int16:
		dst = int16(src.Int())
	case reflect.Int32:
		dst = int32(src.Int())
	case reflect.Int64:
		dst = src.Int()
	case reflect.Uint:
		dst = src.Uint()
	case reflect.Uint8:
		dst = uint8(src.Int())
	case reflect.Uint16:
		dst = uint16(src.Int())
	case reflect.Uint32:
		dst = uint32(src.Int())
	case reflect.Uint64:
		dst = uint64(src.Int())
	case reflect.String:
		dst = src.String()
	case reflect.Bool:
		dst = src.Bool()
	default:
		dst = src
	}
	return
}
