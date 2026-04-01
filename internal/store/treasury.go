package store

import (
	"context"
	"time"

	"bikagame-go/internal/config"
	"bikagame-go/internal/db"
	"bikagame-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureTreasury(ctx context.Context, dbc *db.DB, cfg config.Config) (*models.Treasury, error) {
	now := time.Now()
	update := bson.M{
		"$setOnInsert": bson.M{
			"key":              "treasury",
			"ownerUserId":      cfg.OwnerID,
			"totalSupply":      int64(0),
			"ownerBalance":     int64(0),
			"maintenanceMode":  false,
			"vipWinRate":       90,
			"shopEnabled":      true,
			"broadcastRunning": false,
			"broadcastRunId":   nil,
			"slotRtp":          0.90,
			"createdAt":        now,
		},
		"$set": bson.M{"updatedAt": now},
	}
	_, err := dbc.Config.UpdateOne(ctx, bson.M{"key": "treasury"}, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return GetTreasury(ctx, dbc)
}

func GetTreasury(ctx context.Context, dbc *db.DB) (*models.Treasury, error) {
	var t models.Treasury
	err := dbc.Config.FindOne(ctx, bson.M{"key": "treasury"}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func TreasuryPayToUser(ctx context.Context, dbc *db.DB, userID int64, amount int64) error {
	if amount <= 0 {
		return nil
	}
	res, err := dbc.Config.UpdateOne(ctx, bson.M{"key": "treasury", "ownerBalance": bson.M{"$gte": amount}}, bson.M{
		"$inc": bson.M{"ownerBalance": -amount},
		"$set": bson.M{"updatedAt": time.Now()},
	})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrTreasuryInsufficient
	}
	return AddBalance(ctx, dbc, userID, amount)
}

func UserPayToTreasury(ctx context.Context, dbc *db.DB, userID int64, amount int64) error {
	if amount <= 0 {
		return nil
	}
	ok, err := SubBalanceIfEnough(ctx, dbc, userID, amount)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserInsufficient
	}
	_, err = dbc.Config.UpdateOne(ctx, bson.M{"key": "treasury"}, bson.M{
		"$inc": bson.M{"ownerBalance": amount},
		"$set": bson.M{"updatedAt": time.Now()},
	})
	return err
}
