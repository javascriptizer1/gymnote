package repository

import (
	"context"

	"gymnote/internal/entity"
)

type UserAccountRepository interface {
	GetOrCreateAccount(ctx context.Context, platform entity.Platform, platformUserID string, username string) (*entity.UserAccount, error)
	GetAccountByPlatform(ctx context.Context, platform entity.Platform, platformUserID string) (*entity.UserAccount, error)
	GetAccountsByInternalUserID(ctx context.Context, internalUserID string) ([]entity.UserAccount, error)
	CreateAccount(ctx context.Context, account *entity.UserAccount) error
	UpdateAccountInternalUserID(ctx context.Context, accountID string, internalUserID string) error
	CreateLinkCode(ctx context.Context, internalUserID string, platform entity.Platform) (*entity.LinkCode, error)
	GetLinkCode(ctx context.Context, code string) (*entity.LinkCode, error)
	MarkLinkCodeAsUsed(ctx context.Context, codeID string) error
}
