package repository

import (
	"context"
	"errors"
	"reflect"
)

type CRUDRepository[T any, ID comparable] interface {
	/// Create one entity of type T
	Create(ctx context.Context, entity *T, opts ...OneOpts) error
	/// Update one entity of type T
	Update(ctx context.Context, entity *T, opts ...OneOpts) error
	/// Update one entity of type T
	UpdateOne(ctx context.Context, filter any, update any, opts ...OneOpts) error
	/// Update many entities of type T
	UpdateMany(ctx context.Context, filter any, update any, opts ...ManyOpts) error
	/// Save (Insert or Update) one entity of type T
	Save(ctx context.Context, entity *T, opts ...OneOpts) error
	/// Delete one entity of type T
	Delete(ctx context.Context, id ID, opts ...OneOpts) error
	/// Delete one entity of type T
	DeleteOne(ctx context.Context, filter any, opts ...OneOpts) error
	/// Delete many entities of type T
	DeleteMany(ctx context.Context, filter any, opts ...ManyOpts) error
	/// Find one entity of type T by ID
	FindById(ctx context.Context, id ID, opts ...OneOpts) (*T, error)
	/// Count entities of type T
	Count(ctx context.Context, filter any, opts ...ManyOpts) (int64, error)
	/// Check if entity of type T exists
	Exists(ctx context.Context, filter any, opts ...OneOpts) (bool, error)
	/// Find one entity of type T
	FindOne(ctx context.Context, filter any, opts ...ManyOpts) (*T, error)
	/// Find many entities of type T
	FindMany(ctx context.Context, filter any, opts ...ManyOpts) (*[]T, error)
	/// Find page of entities of type T
	FindPage(ctx context.Context, filter any, opts ...ManyOpts) (*Page[T], error)
}

type SortDirection string

const (
	SortDirection_Unspecified SortDirection = ""
	SortDirection_Ascending   SortDirection = "asc"
	SortDirection_Descending  SortDirection = "desc"
)

type SortBy struct {
	Field     string
	Direction SortDirection
}

type ManyOpts struct {
	Select []string
	Expand []string
	SortBy []SortBy
	Limit  uint32
	Offset uint32
}

type Page[T any] struct {
	Items      []T
	Total      int64
	Page       int32
	Size       int32
	TotalPages int32
}

type OneOpts struct {
	Select []string
	Expand []string
}

func GetEntityID[T any, ID comparable](entity *T) (*ID, error) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("entity is not a struct")
	}

	id := val.FieldByName("ID")
	if !id.IsValid() {
		return nil, errors.New("entity does not have ID field")
	}

	return new(id.Interface().(ID)), nil
}
