package repository

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCRUDRepository[T any, ID comparable] struct {
	collection *mongo.Collection
}

func NewMongoCRUDRepository[T any, ID comparable](db *mongo.Database, collectionName string) CRUDRepository[T, ID] {
	return &mongoCRUDRepository[T, ID]{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoCRUDRepository[T, ID]) Create(ctx context.Context, entity *T, opts ...OneOpts) error {
	res, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	// If the entity doesn't have an ID set, populate it from the InsertOneResult
	if r.getID(entity) == nil || r.isZeroValue(r.getID(entity)) {
		r.setID(entity, res.InsertedID)
	}
	return nil
}

func (r *mongoCRUDRepository[T, ID]) Update(ctx context.Context, entity *T, opts ...OneOpts) error {
	id := r.getID(entity)
	filter := bson.M{"_id": id}
	// ReplaceOne is often cleaner than UpdateOne with $set for whole entities
	// as it avoids issues with updating the _id field if it's in the struct.
	_, err := r.collection.ReplaceOne(ctx, filter, entity)
	return err
}

func (r *mongoCRUDRepository[T, ID]) UpdateOne(ctx context.Context, filter any, update any, opts ...OneOpts) error {
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoCRUDRepository[T, ID]) UpdateMany(ctx context.Context, filter any, update any, opts ...ManyOpts) error {
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *mongoCRUDRepository[T, ID]) Save(ctx context.Context, entity *T, opts ...OneOpts) error {
	id := r.getID(entity)
	filter := bson.M{"_id": id}
	upsert := true
	_, err := r.collection.ReplaceOne(ctx, filter, entity, &options.ReplaceOptions{Upsert: &upsert})
	return err
}

func (r *mongoCRUDRepository[T, ID]) Delete(ctx context.Context, id ID, opts ...OneOpts) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoCRUDRepository[T, ID]) DeleteOne(ctx context.Context, filter any, opts ...OneOpts) error {
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *mongoCRUDRepository[T, ID]) DeleteMany(ctx context.Context, filter any, opts ...ManyOpts) error {
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *mongoCRUDRepository[T, ID]) FindById(ctx context.Context, id ID, opts ...OneOpts) (*T, error) {
	var entity T
	findOpts := options.FindOne()
	if len(opts) > 0 {
		opt := opts[0]
		selectFields := opt.Select
		if selectFields != nil && len(*selectFields) > 0 {
			projection := bson.M{}
			for _, f := range *selectFields {
				projection[f] = 1
			}
			findOpts.SetProjection(projection)
		}
	}
	err := r.collection.FindOne(ctx, bson.M{"_id": id}, findOpts).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *mongoCRUDRepository[T, ID]) Count(ctx context.Context, filter any, opts ...ManyOpts) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

func (r *mongoCRUDRepository[T, ID]) Exists(ctx context.Context, filter any, opts ...OneOpts) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	return count > 0, err
}

func (r *mongoCRUDRepository[T, ID]) FindOne(ctx context.Context, filter any, opts ...ManyOpts) (*T, error) {
	var entity T
	findOpts := options.FindOne()
	if len(opts) > 0 {
		opt := opts[0]
		selectFields := opt.Select
		if selectFields != nil && len(*selectFields) > 0 {
			projection := bson.M{}
			for _, f := range *selectFields {
				projection[f] = 1
			}
			findOpts.SetProjection(projection)
		}
		sortFields := opt.SortBy
		if sortFields != nil && len(*sortFields) > 0 {
			sort := bson.D{}
			for _, s := range *sortFields {
				dir := 1
				if s.Direction == SortDirection_Descending {
					dir = -1
				}
				sort = append(sort, bson.E{Key: s.Field, Value: dir})
			}
			findOpts.SetSort(sort)
		}
	}
	err := r.collection.FindOne(ctx, filter, findOpts).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *mongoCRUDRepository[T, ID]) FindMany(ctx context.Context, filter any, opts ...ManyOpts) (*[]T, error) {
	findOpts := options.Find()
	if len(opts) > 0 {
		o := opts[0]
		selectFields := o.Select
		if selectFields != nil && len(*selectFields) > 0 {
			projection := bson.M{}
			for _, f := range *selectFields {
				projection[f] = 1
			}
			findOpts.SetProjection(projection)
		}
		sortFields := o.SortBy
		if sortFields != nil && len(*sortFields) > 0 {
			sort := bson.D{}
			for _, s := range *sortFields {
				dir := 1
				if s.Direction == SortDirection_Descending {
					dir = -1
				}
				sort = append(sort, bson.E{Key: s.Field, Value: dir})
			}
			findOpts.SetSort(sort)
		}
		limit := o.Limit
		if limit != nil && *limit > 0 {
			findOpts.SetLimit(*limit)
		}
		offset := o.Offset
		if offset != nil && *offset > 0 {
			findOpts.SetSkip(*offset)
		}
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

func (r *mongoCRUDRepository[T, ID]) FindPage(ctx context.Context, filter any, opts ...ManyOpts) (*Page[T], error) {
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	items, err := r.FindMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var page, size int32
	if len(opts) > 0 {
		o := opts[0]
		limit := o.Limit
		if limit != nil && *limit > 0 {
			size = int32(*limit)
			page = int32(*o.Offset)/size + 1
		}
	}

	totalPages := int32(0)
	if size > 0 {
		totalPages = int32((total + int64(size) - 1) / int64(size))
	}

	return &Page[T]{
		Items:      *items,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}

// getID is a helper to extract the ID from the entity using reflection.
func (r *mongoCRUDRepository[T, ID]) getID(entity *T) any {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	// Try to find field with bson:"_id" tag first
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("bson")
		if tag == "_id" || tag == "_id,omitempty" {
			return val.Field(i).Interface()
		}
	}

	// Fallback to field named "ID"
	idField := val.FieldByName("ID")
	if idField.IsValid() {
		return idField.Interface()
	}

	return nil
}

// setID is a helper to set the ID of the entity using reflection.
func (r *mongoCRUDRepository[T, ID]) setID(entity *T, id any) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	// Try to find field with bson:"_id" tag first
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("bson")
		if tag == "_id" || tag == "_id,omitempty" {
			targetField := val.Field(i)
			if targetField.CanSet() {
				targetField.Set(reflect.ValueOf(id))
			}
			return
		}
	}

	// Fallback to field named "ID"
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		idField.Set(reflect.ValueOf(id))
	}
}

func (r *mongoCRUDRepository[T, ID]) isZeroValue(v any) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).IsZero()
}
