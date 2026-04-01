package store

import (
	"context"
	"strings"
	"time"

	"bikagame-go/internal/db"
	"bikagame-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureUser(ctx context.Context, dbc *db.DB, userID int64, username, firstName, lastName, fullName string) (*models.User, error) {
	now := time.Now()
	username = strings.ToLower(username)

	update := bson.M{
		"$set": bson.M{
			"userId":    userID,
			"username":  username,
			"firstName": firstName,
			"lastName":  lastName,
			"fullName":  fullName,
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"balance":           int64(0),
			"isVip":             false,
			"startedBot":        false,
			"startBonusClaimed": false,
			"createdAt":         now,
		},
	}

	_, err := dbc.Users.UpdateOne(ctx, bson.M{"userId": userID}, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return nil, err
	}

	return GetUser(ctx, dbc, userID)
}

func GetUser(ctx context.Context, dbc *db.DB, userID int64) (*models.User, error) {
	var user models.User
	err := dbc.Users.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func AddBalance(ctx context.Context, dbc *db.DB, userID int64, amount int64) error {
	_, err := dbc.Users.UpdateOne(ctx, bson.M{"userId": userID}, bson.M{
		"$inc": bson.M{"balance": amount},
		"$set": bson.M{"updatedAt": time.Now()},
	}, options.UpdateOne().SetUpsert(true))
	return err
}

func SubBalanceIfEnough(ctx context.Context, dbc *db.DB, userID int64, amount int64) (bool, error) {
	res, err := dbc.Users.UpdateOne(ctx, bson.M{
		"userId":  userID,
		"balance": bson.M{"$gte": amount},
	}, bson.M{
		"$inc": bson.M{"balance": -amount},
		"$set": bson.M{"updatedAt": time.Now()},
	})
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}
func SetLastDailyClaimDate(ctx context.Context, dbc *db.DB, userID int64, dateKey string) error {
	_, err := dbc.Users.UpdateOne(ctx, bson.M{"userId": userID}, bson.M{
		"$set": bson.M{
			"lastDailyClaimDate": dateKey,
			"updatedAt":          time.Now(),
		},
	})
	return err
}

func TopUsersByBalance(ctx context.Context, dbc *db.DB, limit int64) ([]models.User, error) {
	cur, err := dbc.Users.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"balance": -1}).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]models.User, 0)
	for cur.Next(ctx) {
		var u models.User
		if err := cur.Decode(&u); err == nil {
			out = append(out, u)
		}
	}
	return out, nil
}
