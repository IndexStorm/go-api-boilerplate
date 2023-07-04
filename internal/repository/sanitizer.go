package repository

import (
	"indexstorm/go-api-boilerplate/pkg/nanoid"
	"reflect"
)

func (r *postgresRepo) sanitizeCreateModel(model interface{}) {
	value := reflect.Indirect(reflect.ValueOf(model))
	if field := value.FieldByName("ID"); field.IsValid() {
		field.SetUint(0)
	}
	if field := value.FieldByName("PublicID"); field.IsValid() {
		field.SetString(nanoid.RandomShortID())
	}
	if field := value.FieldByName("DateMixin"); field.IsValid() {
		field.SetZero()
	}
	if field := value.FieldByName("SoftDeleteMixin"); field.IsValid() {
		field.SetZero()
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		println(value.Type().Field(i).Name)
		if k := field.Kind(); k == reflect.Pointer {
			if k := field.Elem().Kind(); k == reflect.Slice {
				if field.Elem().Len() == 0 {
					field.SetZero()
				}
			} else if sp, ok := field.Interface().(*string); ok && sp != nil {
				if *sp == "" {
					field.SetZero()
				}
			}
		} else if k == reflect.Slice {
			if field.Len() == 0 {
				field.SetZero()
			}
		}
	}
}
