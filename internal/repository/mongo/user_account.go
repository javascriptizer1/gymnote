package mongodb

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"gymnote/internal/entity"
	"gymnote/internal/errs"
)

const (
	userAccountsCollection = "user_accounts"
	linkCodesCollection    = "link_codes"
	linkCodeLength         = 6
	linkCodeExpiryMinutes  = 15
)

func (m *mongodb) GetOrCreateAccount(ctx context.Context, platform entity.Platform, platformUserID string, username string) (*entity.UserAccount, error) {
	collection := m.db.Collection(userAccountsCollection)

	var accountRow UserAccountRow
	err := collection.FindOne(ctx, bson.M{
		"platform":         platform,
		"platform_user_id": platformUserID,
	}).Decode(&accountRow)

	if err == nil {
		return accountRow.ToEntity(), nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	account := entity.NewUserAccount(
		entity.WithUserAccountInitSpec(entity.UserAccountInitSpecification{
			Platform:       platform,
			PlatformUserID: platformUserID,
			Username:       username,
		}),
	)

	accountRow = UserAccountRow{
		ID:             account.ID(),
		InternalUserID: account.InternalUserID(),
		Platform:       account.Platform(),
		PlatformUserID: account.PlatformUserID(),
		Username:       account.Username(),
		LinkedAt:       account.LinkedAt(),
		CreatedAt:      account.CreatedAt(),
	}

	if _, err = collection.InsertOne(ctx, accountRow); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (m *mongodb) GetAccountByPlatform(ctx context.Context, platform entity.Platform, platformUserID string) (*entity.UserAccount, error) {
	collection := m.db.Collection(userAccountsCollection)

	var accountRow UserAccountRow
	err := collection.FindOne(ctx, bson.M{
		"platform":         platform,
		"platform_user_id": platformUserID,
	}).Decode(&accountRow)

	if err == mongo.ErrNoDocuments {
		return nil, errs.ErrAccountNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	return accountRow.ToEntity(), nil
}

func (m *mongodb) GetAccountsByInternalUserID(ctx context.Context, internalUserID string) ([]entity.UserAccount, error) {
	collection := m.db.Collection(userAccountsCollection)

	cursor, err := collection.Find(ctx, bson.M{
		"internal_user_id": internalUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find accounts: %w", err)
	}
	defer cursor.Close(ctx)

	var accountRows []UserAccountRow
	if err := cursor.All(ctx, &accountRows); err != nil {
		return nil, fmt.Errorf("failed to decode accounts: %w", err)
	}

	accounts := make([]entity.UserAccount, len(accountRows))
	for i, row := range accountRows {
		accounts[i] = *row.ToEntity()
	}

	return accounts, nil
}

func (m *mongodb) CreateAccount(ctx context.Context, account *entity.UserAccount) error {
	collection := m.db.Collection(userAccountsCollection)

	accountRow := UserAccountRow{
		ID:             account.ID(),
		InternalUserID: account.InternalUserID(),
		Platform:       account.Platform(),
		PlatformUserID: account.PlatformUserID(),
		Username:       account.Username(),
		LinkedAt:       account.LinkedAt(),
		CreatedAt:      account.CreatedAt(),
	}

	if _, err := collection.InsertOne(ctx, accountRow); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (m *mongodb) UpdateAccountInternalUserID(ctx context.Context, accountID string, internalUserID string) error {
	collection := m.db.Collection(userAccountsCollection)
	if _, err := collection.UpdateOne(ctx,
		bson.M{"id": accountID},
		bson.M{"$set": bson.M{"internal_user_id": internalUserID}},
	); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func (m *mongodb) CreateLinkCode(ctx context.Context, internalUserID string, platform entity.Platform) (*entity.LinkCode, error) {
	collection := m.db.Collection(linkCodesCollection)

	code, err := generateLinkCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	linkCode := entity.NewLinkCode(
		entity.WithLinkCodeInitSpec(entity.LinkCodeInitSpecification{
			Code:           code,
			InternalUserID: internalUserID,
			Platform:       platform,
			ExpiresAt:      time.Now().Add(linkCodeExpiryMinutes * time.Minute),
		}),
	)

	linkCodeRow := LinkCodeRow{
		ID:             linkCode.ID(),
		Code:           linkCode.Code(),
		InternalUserID: linkCode.InternalUserID(),
		Platform:       linkCode.Platform(),
		ExpiresAt:      linkCode.ExpiresAt(),
		CreatedAt:      linkCode.CreatedAt(),
		Used:           linkCode.Used(),
	}

	if _, err = collection.InsertOne(ctx, linkCodeRow); err != nil {
		return nil, fmt.Errorf("failed to create link code: %w", err)
	}

	return linkCode, nil
}

func (m *mongodb) GetLinkCode(ctx context.Context, code string) (*entity.LinkCode, error) {
	collection := m.db.Collection(linkCodesCollection)

	var linkCodeRow LinkCodeRow
	err := collection.FindOne(ctx, bson.M{
		"code": code,
		"used": false,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&linkCodeRow)

	if err == mongo.ErrNoDocuments {
		return nil, errs.ErrLinkCodeNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find link code: %w", err)
	}

	return linkCodeRow.ToEntity(), nil
}

func (m *mongodb) MarkLinkCodeAsUsed(ctx context.Context, codeID string) error {
	collection := m.db.Collection(linkCodesCollection)

	_, err := collection.UpdateOne(ctx,
		bson.M{"id": codeID},
		bson.M{"$set": bson.M{"used": true}},
	)
	if err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	return nil
}

func generateLinkCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, linkCodeLength)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}
