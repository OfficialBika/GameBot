package store

import (
	"context"
	"time"

	"bikagame-go/internal/db"
	"bikagame-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureGroup(ctx context.Context, dbc *db.DB, groupID int64, title, username string, botIsAdmin bool) (*models.Group, error) {
	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"groupId":    groupID,
			"title":      title,
			"username":   username,
			"botIsAdmin": botIsAdmin,
			"updatedAt":  now,
		},
		"$setOnInsert": bson.M{
			"approvalStatus": "pending",
			"approvedBy":     nil,
			"createdAt":      now,
		},
	}

	_, err := dbc.Groups.UpdateOne(ctx, bson.M{"groupId": groupID}, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return GetGroup(ctx, dbc, groupID)
}

func ApproveGroup(ctx context.Context, dbc *db.DB, groupID int64, ownerID int64) error {
	now := time.Now()
	_, err := dbc.Groups.UpdateOne(ctx, bson.M{"groupId": groupID}, bson.M{
		"$set": bson.M{
			"approvalStatus": "approved",
			"approvedBy":     ownerID,
			"updatedAt":      now,
		},
	}, options.UpdateOne().SetUpsert(true))
	return err
}

func RejectGroup(ctx context.Context, dbc *db.DB, groupID int64, ownerID int64) error {
	now := time.Now()
	_, err := dbc.Groups.UpdateOne(ctx, bson.M{"groupId": groupID}, bson.M{
		"$set": bson.M{
			"approvalStatus": "rejected",
			"approvedBy":     ownerID,
			"updatedAt":      now,
		},
	}, options.UpdateOne().SetUpsert(true))
	return err
}

func GetGroup(ctx context.Context, dbc *db.DB, groupID int64) (*models.Group, error) {
	var group models.Group
	err := dbc.Groups.FindOne(ctx, bson.M{"groupId": groupID}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}
