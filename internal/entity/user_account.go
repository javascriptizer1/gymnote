package entity

import (
	"time"

	"github.com/google/uuid"
)

type Platform string

const (
	PlatformTelegram Platform = "telegram"
	PlatformVK       Platform = "vk"
)

type UserAccountOption func(o *UserAccount)

type UserAccount struct {
	id             string
	internalUserID string
	platform       Platform
	platformUserID string
	username       string
	linkedAt       time.Time
	createdAt      time.Time
}

func (ua *UserAccount) ID() string {
	return ua.id
}

func (ua *UserAccount) InternalUserID() string {
	return ua.internalUserID
}

func (ua *UserAccount) Platform() Platform {
	return ua.platform
}

func (ua *UserAccount) PlatformUserID() string {
	return ua.platformUserID
}

func (ua *UserAccount) Username() string {
	return ua.username
}

func (ua *UserAccount) LinkedAt() time.Time {
	return ua.linkedAt
}

func (ua *UserAccount) CreatedAt() time.Time {
	return ua.createdAt
}

func NewUserAccount(opts ...UserAccountOption) *UserAccount {
	account := &UserAccount{}

	for _, opt := range opts {
		opt(account)
	}

	return account
}

type UserAccountInitSpecification struct {
	Platform       Platform
	PlatformUserID string
	Username       string
}

func WithUserAccountInitSpec(spec UserAccountInitSpecification) UserAccountOption {
	return func(ua *UserAccount) {
		ua.id = uuid.New().String()
		ua.internalUserID = uuid.New().String()
		ua.platform = spec.Platform
		ua.platformUserID = spec.PlatformUserID
		ua.username = spec.Username
		ua.linkedAt = time.Now()
		ua.createdAt = time.Now()
	}
}

type UserAccountRestoreSpecification struct {
	ID             string
	InternalUserID string
	Platform       Platform
	PlatformUserID string
	Username       string
	LinkedAt       time.Time
	CreatedAt      time.Time
}

func WithUserAccountRestoreSpec(spec UserAccountRestoreSpecification) UserAccountOption {
	return func(ua *UserAccount) {
		ua.id = spec.ID
		ua.internalUserID = spec.InternalUserID
		ua.platform = spec.Platform
		ua.platformUserID = spec.PlatformUserID
		ua.username = spec.Username
		ua.linkedAt = spec.LinkedAt
		ua.createdAt = spec.CreatedAt
	}
}

type LinkCodeOption func(o *LinkCode)

type LinkCode struct {
	id             string
	code           string
	internalUserID string
	platform       Platform
	expiresAt      time.Time
	createdAt      time.Time
	used           bool
}

func (lc *LinkCode) ID() string {
	return lc.id
}

func (lc *LinkCode) Code() string {
	return lc.code
}

func (lc *LinkCode) InternalUserID() string {
	return lc.internalUserID
}

func (lc *LinkCode) Platform() Platform {
	return lc.platform
}

func (lc *LinkCode) ExpiresAt() time.Time {
	return lc.expiresAt
}

func (lc *LinkCode) CreatedAt() time.Time {
	return lc.createdAt
}

func (lc *LinkCode) Used() bool {
	return lc.used
}

func NewLinkCode(opts ...LinkCodeOption) *LinkCode {
	linkCode := &LinkCode{}

	for _, opt := range opts {
		opt(linkCode)
	}

	return linkCode
}

type LinkCodeInitSpecification struct {
	Code           string
	InternalUserID string
	Platform       Platform
	ExpiresAt      time.Time
}

func WithLinkCodeInitSpec(spec LinkCodeInitSpecification) LinkCodeOption {
	return func(lc *LinkCode) {
		lc.id = uuid.New().String()
		lc.code = spec.Code
		lc.internalUserID = spec.InternalUserID
		lc.platform = spec.Platform
		lc.expiresAt = spec.ExpiresAt
		lc.createdAt = time.Now()
		lc.used = false
	}
}

type LinkCodeRestoreSpecification struct {
	ID             string
	Code           string
	InternalUserID string
	Platform       Platform
	ExpiresAt      time.Time
	CreatedAt      time.Time
	Used           bool
}

func WithLinkCodeRestoreSpec(spec LinkCodeRestoreSpecification) LinkCodeOption {
	return func(lc *LinkCode) {
		lc.id = spec.ID
		lc.code = spec.Code
		lc.internalUserID = spec.InternalUserID
		lc.platform = spec.Platform
		lc.expiresAt = spec.ExpiresAt
		lc.createdAt = spec.CreatedAt
		lc.used = spec.Used
	}
}
