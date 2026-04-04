package repository

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BeforeCreateHook[ID comparable] interface {
	BeforeCreate() error
}

type BeforeUpdateHook[ID comparable] interface {
	BeforeUpdate() error
}

func BeforeCreate[ID comparable](entity any) {
	if hook, ok := entity.(BeforeCreateHook[ID]); ok {
		if err := hook.BeforeCreate(); err != nil {
			panic(err)
		}
	}
	setFieldValue(entity, "CreatedAt", time.Now())
	setFieldValue(entity, "UpdatedAt", time.Now())
}

func BeforeUpdate[ID comparable](entity any) {
	if hook, ok := entity.(BeforeUpdateHook[ID]); ok {
		if err := hook.BeforeUpdate(); err != nil {
			panic(err)
		}
	}
	setFieldValue(entity, "UpdatedAt", time.Now())
}

func newId[ID comparable]() ID {
	var idValue ID
	idType := reflect.TypeFor[ID]()
	if idType.Kind() == reflect.Pointer {
		idType = idType.Elem()
	}

	switch idType.Kind() {
	case reflect.String:
		reflect.ValueOf(&idValue).Elem().SetString("")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		reflect.ValueOf(&idValue).Elem().SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.ValueOf(&idValue).Elem().SetUint(0)
	default:
		if idType == reflect.TypeFor[primitive.ObjectID]() {
			reflect.ValueOf(&idValue).Elem().Set(reflect.ValueOf(primitive.NewObjectID()))
		}
	}

	return idValue
}

func getID(entity any) any {
	entityValue := getValue(entity)
	return entityValue.FieldByName("ID").Interface()
}

func setID(entity any, id any) {
	setFieldValue(entity, "ID", id)
}

func getValue(entity any) reflect.Value {
	entityType := reflect.TypeOf(entity)
	entityValue := reflect.ValueOf(entity)
	if entityType.Kind() == reflect.Pointer {
		entityType = entityType.Elem()
		entityValue = entityValue.Elem()
	}
	return entityValue
}

func getFieldValue[T any](entity any, field string) T {
	entityValue := getValue(entity)

	fieldValue := entityValue.FieldByName(field)

	return fieldValue.Interface().(T)
}

func setFieldValue[T any](entity any, field string, value T) {
	entityValue := getValue(entity)
	entityField := entityValue.FieldByName(field)
	if !entityField.IsValid() {
		return
	}

	val := reflect.ValueOf(value)

	// If the field is a pointer and the value is not a pointer
	if entityField.Kind() == reflect.Ptr && val.Kind() != reflect.Ptr {
		// And types match (value matches the type the pointer points to)
		if entityField.Type().Elem() == val.Type() {
			ptr := reflect.New(val.Type())
			ptr.Elem().Set(val)
			entityField.Set(ptr)
			return
		}
	}

	// If the field is NOT a pointer and the value IS a pointer
	if entityField.Kind() != reflect.Ptr && val.Kind() == reflect.Ptr {
		// And types match (field matches the type the pointer points to)
		if entityField.Type() == val.Type().Elem() {
			if !val.IsNil() {
				entityField.Set(val.Elem())
			}
			return
		}
	}

	entityField.Set(val)
}
