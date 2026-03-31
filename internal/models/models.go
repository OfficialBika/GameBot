package models

import "time"

type User struct {
	UserID             int64      `bson:"userId"`
	Username           string     `bson:"username,omitempty"`
	FirstName          string     `bson:"firstName,omitempty"`
	LastName           string     `bson:"lastName,omitempty"`
	FullName           string     `bson:"fullName,omitempty"`
	Balance            int64      `bson:"balance"`
	IsVIP              bool       `bson:"isVip"`
	StartedBot         bool       `bson:"startedBot"`
	StartBonusClaimed  bool       `bson:"startBonusClaimed"`
	LastDailyClaimDate *string    `bson:"lastDailyClaimDate,omitempty"`
	CreatedAt          time.Time  `bson:"createdAt"`
	UpdatedAt          time.Time  `bson:"updatedAt"`
}

type Group struct {
	GroupID        int64     `bson:"groupId"`
	Title          string    `bson:"title"`
	Username       string    `bson:"username,omitempty"`
	ApprovalStatus string    `bson:"approvalStatus"`
	ApprovedBy     *int64    `bson:"approvedBy,omitempty"`
	BotIsAdmin     bool      `bson:"botIsAdmin"`
	InviteLink     string    `bson:"inviteLink,omitempty"`
	CreatedAt      time.Time `bson:"createdAt"`
	UpdatedAt      time.Time `bson:"updatedAt"`
}

type Treasury struct {
	Key              string    `bson:"key"`
	OwnerUserID      int64     `bson:"ownerUserId"`
	TotalSupply      int64     `bson:"totalSupply"`
	OwnerBalance     int64     `bson:"ownerBalance"`
	MaintenanceMode  bool      `bson:"maintenanceMode"`
	VIPWinRate       int       `bson:"vipWinRate"`
	ShopEnabled      bool      `bson:"shopEnabled"`
	BroadcastRunning bool      `bson:"broadcastRunning"`
	BroadcastRunID   any       `bson:"broadcastRunId"`
	SlotRTP          float64   `bson:"slotRtp"`
	CreatedAt        time.Time `bson:"createdAt"`
	UpdatedAt        time.Time `bson:"updatedAt"`
}
