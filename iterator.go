package structx

import (
	"reflect"
	"strconv"
	"strings"
)

type Field struct {
	Val reflect.Value
	Tag reflect.StructTag
}

type FieldFilter func(path string, field Field) (err error)

func IterateObject(obj interface{}, filters ...FieldFilter) (values map[string]Field, err error) {
	val := reflect.ValueOf(obj)

	if !val.CanInterface() {
		return
	}

	v := reflect.ValueOf(val.Interface())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	rootPath := ""
	if v.Type().Kind() == reflect.Struct {
		rootPath = v.Type().Name()
	}

	tagMap := map[string]reflect.StructTag{}
	valMap := map[string]Field{}
	err = nested(rootPath, valMap, tagMap, val, filters)

	tagMap = nil

	if err != nil {
		return
	}

	values = valMap
	return
}

func nested(path string, valuesMap map[string]Field, tagMap map[string]reflect.StructTag, val reflect.Value, filters []FieldFilter) (err error) {

	if !val.CanInterface() {
		return
	}

	v := reflect.ValueOf(val.Interface())
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vType := v.Type()

	switch v.Kind() {
	case reflect.Struct:

		vTypeTag, _ := tagMap[path]

		if omitnestedField(vTypeTag) {
			err = applyFilters(path, valuesMap, tagMap, val, filters)
			if err != nil {
				return
			}

			break
		}

		for i := 0; i < v.NumField(); i++ {

			fieldVal := v.Field(i)
			fieldTyp := vType.Field(i)
			nextPath := path + "." + vType.Field(i).Name

			if skipField(fieldTyp.Tag) {
				continue
			}

			tagMap[nextPath] = fieldTyp.Tag

			err = nested(nextPath, valuesMap, tagMap, fieldVal, filters)
			if err != nil {
				return
			}
		}

	case reflect.Map:
		mapElem := val.Type()
		switch val.Type().Kind() {
		case reflect.Ptr, reflect.Array, reflect.Map,
			reflect.Slice, reflect.Chan:
			mapElem = val.Type().Elem()
			if mapElem.Kind() == reflect.Ptr {
				mapElem = mapElem.Elem()
			}
		}

		if mapElem.Kind() == reflect.Struct ||
			(mapElem.Kind() == reflect.Slice &&
				mapElem.Elem().Kind() == reflect.Struct) {
			for _, k := range val.MapKeys() {
				err = nested(path+"["+k.String()+"]", valuesMap, tagMap, val.MapIndex(k), filters)
				if err != nil {
					return
				}
			}
			break
		}

	case reflect.Slice, reflect.Array:

		if val.Type().Kind() == reflect.Interface {
			break
		}

		// do not iterate of non struct types, just pass the value. Ie: []int,
		// []string, co... We only iterate further if it's a struct.
		// i.e []foo or []*foo
		if val.Type().Elem().Kind() != reflect.Struct &&
			!(val.Type().Elem().Kind() == reflect.Ptr &&
				val.Type().Elem().Elem().Kind() == reflect.Struct) {
			break
		}

		for x := 0; x < val.Len(); x++ {
			err = nested(path+"["+strconv.Itoa(x)+"]", valuesMap, tagMap, val.Index(x), filters)
			if err != nil {
				return
			}
		}

	default:
		applyFilters(path, valuesMap, tagMap, val, filters)
	}

	return
}

func applyFilters(path string, valuesMap map[string]Field, tagMap map[string]reflect.StructTag, val reflect.Value, filters []FieldFilter) (err error) {
	for i := 0; i < len(filters); i++ {
		if val.IsValid() {

			tag, _ := tagMap[path]

			field := Field{
				Val: val,
				Tag: tag,
			}

			err = filters[i](path, field)
			if err != nil {
				return
			}

			valuesMap[path] = field
		}
	}

	return
}

func skipField(tag reflect.StructTag) bool {
	tagVal, exist := tag.Lookup("structs")
	if !exist {
		return false
	}

	tagVals := strings.Split(tagVal, ",")

	return hasTagVal(tagVals, "-")
}

func omitnestedField(tag reflect.StructTag) bool {
	tagVal, exist := tag.Lookup("structs")
	if !exist {
		return false
	}

	tagVals := strings.Split(tagVal, ",")

	return hasTagVal(tagVals, "omitnested")
}

func hasTagVal(tagVals []string, findVal string) bool {
	for _, val := range tagVals {
		if val == findVal {
			return true
		}
	}

	return false
}
