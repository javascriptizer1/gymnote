package mongodb

import (
	"time"

	"gymnote/internal/entity"
)

type UserAccountOption func(o *UserAccountRow)
type LinkCodeOption func(o *LinkCodeRow)

type UserAccountRow struct {
	ID             string          `bson:"id"`
	InternalUserID string          `bson:"internal_user_id"`
	Platform       entity.Platform `bson:"platform"`
	PlatformUserID string          `bson:"platform_user_id"`
	Username       string          `bson:"username"`
	LinkedAt       time.Time       `bson:"linked_at"`
	CreatedAt      time.Time       `bson:"created_at"`
}

func (ua *UserAccountRow) ToEntity() *entity.UserAccount {
	return entity.NewUserAccount(
		entity.WithUserAccountRestoreSpec(entity.UserAccountRestoreSpecification{
			ID:             ua.ID,
			InternalUserID: ua.InternalUserID,
			Platform:       ua.Platform,
			PlatformUserID: ua.PlatformUserID,
			Username:       ua.Username,
			LinkedAt:       ua.LinkedAt,
			CreatedAt:      ua.CreatedAt,
		}),
	)
}

func NewUserAccountRow(opts ...UserAccountOption) *UserAccountRow {
	account := &UserAccountRow{}

	for _, opt := range opts {
		opt(account)
	}

	return account
}

type UserAccountRowRestoreSpecification struct {
	ID             string
	InternalUserID string
	Platform       entity.Platform
	PlatformUserID string
	Username       string
	LinkedAt       time.Time
	CreatedAt      time.Time
}

func WithUserAccountRowRestoreSpec(spec UserAccountRowRestoreSpecification) UserAccountOption {
	return func(ua *UserAccountRow) {
		ua.ID = spec.ID
		ua.InternalUserID = spec.InternalUserID
		ua.Platform = spec.Platform
		ua.PlatformUserID = spec.PlatformUserID
		ua.Username = spec.Username
		ua.LinkedAt = spec.LinkedAt
		ua.CreatedAt = spec.CreatedAt
	}
}

type LinkCodeRow struct {
	ID             string          `bson:"id"`
	Code           string          `bson:"code"`
	InternalUserID string          `bson:"internal_user_id"`
	Platform       entity.Platform `bson:"platform"`
	ExpiresAt      time.Time       `bson:"expires_at"`
	CreatedAt      time.Time       `bson:"created_at"`
	Used           bool            `bson:"used"`
}

func (lc *LinkCodeRow) ToEntity() *entity.LinkCode {
	return entity.NewLinkCode(
		entity.WithLinkCodeRestoreSpec(entity.LinkCodeRestoreSpecification{
			ID:             lc.ID,
			Code:           lc.Code,
			InternalUserID: lc.InternalUserID,
			Platform:       lc.Platform,
			ExpiresAt:      lc.ExpiresAt,
			CreatedAt:      lc.CreatedAt,
			Used:           lc.Used,
		}),
	)
}

func NewLinkCodeRow(opts ...LinkCodeOption) *LinkCodeRow {
	linkCode := &LinkCodeRow{}

	for _, opt := range opts {
		opt(linkCode)
	}

	return linkCode
}

type LinkCodeRowRestoreSpecification struct {
	ID             string
	Code           string
	InternalUserID string
	Platform       entity.Platform
	ExpiresAt      time.Time
	CreatedAt      time.Time
	Used           bool
}

func WithLinkCodeRowRestoreSpec(spec LinkCodeRowRestoreSpecification) LinkCodeOption {
	return func(lc *LinkCodeRow) {
		lc.ID = spec.ID
		lc.Code = spec.Code
		lc.InternalUserID = spec.InternalUserID
		lc.Platform = spec.Platform
		lc.ExpiresAt = spec.ExpiresAt
		lc.CreatedAt = spec.CreatedAt
		lc.Used = spec.Used
	}
}
