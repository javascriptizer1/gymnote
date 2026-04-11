package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	"gymnote/internal/entity"
	"gymnote/internal/errs"
	"gymnote/internal/repository"
)

type AccountService struct {
	db repository.DB
}

func NewAccountService(db repository.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (s *AccountService) GetOrCreateAccount(ctx context.Context, platform entity.Platform, platformUserID string, username string) (*entity.UserAccount, error) {
	return s.db.GetOrCreateAccount(ctx, platform, platformUserID, username)
}

func (s *AccountService) CreateLinkCode(ctx context.Context, platform entity.Platform, platformUserID string) (*entity.LinkCode, error) {
	account, err := s.db.GetOrCreateAccount(ctx, platform, platformUserID, "")
	if err != nil {
		log.Printf("Error getting or creating account: %v\n", err)

		return nil, err
	}

	return s.db.CreateLinkCode(ctx, account.InternalUserID(), platform)
}

func (s *AccountService) LinkAccountWithCode(ctx context.Context, code string, platform entity.Platform, platformUserID string, username string) error {
	linkCode, err := s.db.GetLinkCode(ctx, code)
	if err != nil {
		log.Printf("Error getting link code: %v\n", err)

		return err
	}

	if linkCode.Platform() == platform {
		return errs.ErrCannotLinkSamePlatform
	}

	existingAccount, err := s.db.GetAccountByPlatform(ctx, platform, platformUserID)
	if err != nil && err != errs.ErrAccountNotFound {
		log.Printf("Error checking existing account: %v\n", err)

		return err
	}

	if existingAccount != nil && existingAccount.InternalUserID() == linkCode.InternalUserID() {
		return errs.ErrAccountAlreadyLinked
	}

	if existingAccount != nil {
		oldInternalUserID := existingAccount.InternalUserID()

		if err := s.db.UpdateAccountInternalUserID(ctx, existingAccount.ID(), linkCode.InternalUserID()); err != nil {
			log.Printf("Error updating account internal_user_id: %v\n", err)

			return err
		}

		if err := s.db.UpdateTrainingSessionsUserID(ctx, oldInternalUserID, linkCode.InternalUserID()); err != nil {
			log.Printf("Error updating training sessions user_id: %v\n", err)

			_ = s.db.UpdateAccountInternalUserID(ctx, existingAccount.ID(), oldInternalUserID)

			return err
		}

		log.Printf("Linked account and migrated training data: %s -> %s\n", oldInternalUserID, linkCode.InternalUserID())
	} else {
		account := entity.NewUserAccount(
			entity.WithUserAccountRestoreSpec(entity.UserAccountRestoreSpecification{
				ID:             uuid.New().String(),
				InternalUserID: linkCode.InternalUserID(),
				Platform:       platform,
				PlatformUserID: platformUserID,
				Username:       username,
				LinkedAt:       linkCode.CreatedAt(),
				CreatedAt:      linkCode.CreatedAt(),
			}),
		)
		if err := s.db.CreateAccount(ctx, account); err != nil {
			log.Printf("Error creating linked account: %v\n", err)
			return err
		}
	}

	if err := s.db.MarkLinkCodeAsUsed(ctx, linkCode.ID()); err != nil {
		log.Printf("Error marking code as used: %v\n", err)

		return err
	}

	return nil
}

func (s *AccountService) GetLinkedAccounts(ctx context.Context, platform entity.Platform, platformUserID string) ([]entity.UserAccount, error) {
	account, err := s.db.GetAccountByPlatform(ctx, platform, platformUserID)
	if err != nil {
		log.Printf("Error getting account: %v\n", err)

		return nil, err
	}

	return s.db.GetAccountsByInternalUserID(ctx, account.InternalUserID())
}
