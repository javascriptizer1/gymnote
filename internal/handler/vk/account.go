package vk

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/v2/events"

	"gymnote/internal/entity"
	"gymnote/internal/errs"
)

func (a *API) LinkAccountHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := strconv.Itoa(obj.Message.FromID)
	username := ""
	text := strings.TrimSpace(obj.Message.Text)

	if strings.HasPrefix(text, "/") {
		args := strings.Fields(text)
		if len(args) > 1 {
			code := args[1]
			a.linkWithCode(peerID, userID, username, code)

			return
		}
	}

	linkCode, err := a.accountService.CreateLinkCode(a.ctx, entity.PlatformVK, userID)
	if err != nil {
		a.sendMessage(peerID, errLinkingFailed, nil)

		return
	}

	text = fmt.Sprintf(linkAccountText, linkCode.Code(), linkCode.Code())

	a.sendMessage(peerID, text, nil)
}

func (a *API) linkWithCode(peerID int, userID, username, code string) {
	err := a.accountService.LinkAccountWithCode(a.ctx, code, entity.PlatformVK, userID, username)
	if err != nil {
		var text string
		switch {
		case errors.Is(err, errs.ErrLinkCodeNotFound):
			text = errInvalidLinkCode
		case errors.Is(err, errs.ErrCannotLinkSamePlatform):
			text = errCannotLinkSamePlatform
		case errors.Is(err, errs.ErrAccountAlreadyLinked):
			text = accountAlreadyLinkedText
		default:
			text = errLinkingFailed
		}

		a.sendMessage(peerID, text, nil)

		return
	}

	a.sendMessage(peerID, accountLinkedText, nil)
}

func (a *API) MyAccountsHandler(obj *events.MessageNewObject) {
	peerID := obj.Message.PeerID
	userID := strconv.Itoa(obj.Message.FromID)

	accounts, err := a.accountService.GetLinkedAccounts(a.ctx, entity.PlatformVK, userID)
	if err != nil {
		a.sendMessage(peerID, errInternal, nil)

		return
	}

	if len(accounts) <= 1 {
		a.sendMessage(peerID, noLinkedAccountsText, nil)

		return
	}

	var accountsList strings.Builder
	for _, acc := range accounts {
		platformName := "Неизвестно"
		accountInfo := acc.PlatformUserID()

		switch acc.Platform() {
		case entity.PlatformTelegram:
			platformName = "Telegram"
			if acc.Username() != "" {
				accountInfo = "@" + acc.Username()
			}
		case entity.PlatformVK:
			platformName = "VK"
			if acc.Username() != "" {
				accountInfo = acc.Username()
			} else {
				accountInfo = fmt.Sprintf("vk.com/id%s", acc.PlatformUserID())
			}
		}

		fmt.Fprintf(&accountsList, "• %s: %s\n", platformName, accountInfo)
	}

	text := fmt.Sprintf(myAccountsText, accountsList.String())

	a.sendMessage(peerID, text, nil)
}

func (a *API) getInternalUserID(platformUserID string, username string) (string, error) {
	account, err := a.accountService.GetOrCreateAccount(a.ctx, entity.PlatformVK, platformUserID, username)
	if err != nil {
		return "", err
	}

	return account.InternalUserID(), nil
}

func (a *API) getUserID(obj *events.MessageNewObject) string {
	platformUserID := strconv.Itoa(obj.Message.FromID)
	username := ""

	internalUserID, err := a.getInternalUserID(platformUserID, username)
	if err != nil {
		log.Printf("Error getting internal user ID: %v", err)

		return platformUserID
	}

	return internalUserID
}

func (a *API) getUserIDFromCallback(obj *events.MessageEventObject) string {
	platformUserID := strconv.Itoa(obj.UserID)
	username := ""

	internalUserID, err := a.getInternalUserID(platformUserID, username)
	if err != nil {
		log.Printf("Error getting internal user ID: %v", err)

		return platformUserID
	}

	return internalUserID
}
