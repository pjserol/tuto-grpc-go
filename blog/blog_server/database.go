package main

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func insertBlog(ctx context.Context, item blogItem) (blogItem, error) {
	res, err := collection.InsertOne(ctx, item)
	if err != nil {
		return blogItem{}, err
	}

	var ok bool
	item.ID, ok = res.InsertedID.(primitive.ObjectID)
	if !ok {
		return blogItem{}, errors.New("Cannot convert ObjectID")
	}

	return item, nil
}

func findBlogByID(ctx context.Context, id primitive.ObjectID) (blogItem, error) {
	data := &blogItem{}
	filter := bson.M{"_id": id}

	res := collection.FindOne(ctx, filter)

	if err := res.Decode(data); err != nil {
		return blogItem{}, err
	}

	return *data, nil
}

func updateBlog(ctx context.Context, id primitive.ObjectID, item blogItem) (blogItem, error) {
	filter := bson.M{"_id": id}

	// use UpdateOne to update only the field that we want
	res, err := collection.ReplaceOne(ctx, filter, item)
	if err != nil {
		return blogItem{}, err
	}

	if res.ModifiedCount == 0 {
		return blogItem{}, fmt.Errorf("Cannot find blog id: %s", id)
	}

	item.ID = id
	return item, nil
}

func deleteBlog(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("Cannot find blog id: %s", id)
	}

	return nil
}

func listBlog(ctx context.Context) (*mongo.Cursor, error) {
	cur, err := collection.Find(ctx, primitive.D{{}})
	if err != nil {
		return nil, err
	}

	return cur, nil
}
