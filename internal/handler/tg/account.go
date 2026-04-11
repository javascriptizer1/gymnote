package tg

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gymnote/internal/entity"
	"gymnote/internal/errs"
)

func (a *API) LinkAccountHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)
	username := message.From.UserName

	text := strings.TrimSpace(message.Text)
	if strings.HasPrefix(text, "/") {
		args := strings.Fields(text)
		if len(args) > 1 {
			code := args[1]
			a.linkWithCode(chatID, userID, username, code)

			return
		}
	}

	linkCode, err := a.accountService.CreateLinkCode(a.ctx, entity.PlatformTelegram, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, errLinkingFailed)
		msg.ParseMode = parseMode

		a.bot.Send(msg)

		return
	}

	text = fmt.Sprintf(linkAccountText, linkCode.Code(), linkCode.Code())
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parseMode

	a.bot.Send(msg)
}

func (a *API) linkWithCode(chatID int64, userID, username, code string) {
	err := a.accountService.LinkAccountWithCode(a.ctx, code, entity.PlatformTelegram, userID, username)
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

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = parseMode

		a.bot.Send(msg)

		return
	}

	msg := tgbotapi.NewMessage(chatID, accountLinkedText)
	msg.ParseMode = parseMode

	a.bot.Send(msg)
}

func (a *API) MyAccountsHandler(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userID := strconv.FormatInt(message.From.ID, 10)

	accounts, err := a.accountService.GetLinkedAccounts(a.ctx, entity.PlatformTelegram, userID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, errInternal)
		msg.ParseMode = parseMode

		a.bot.Send(msg)

		return
	}

	if len(accounts) <= 1 {
		msg := tgbotapi.NewMessage(chatID, noLinkedAccountsText)
		msg.ParseMode = parseMode

		a.bot.Send(msg)

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
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parseMode
	a.bot.Send(msg)
}

func (a *API) getInternalUserID(platformUserID string, username string) (string, error) {
	account, err := a.accountService.GetOrCreateAccount(a.ctx, entity.PlatformTelegram, platformUserID, username)
	if err != nil {
		return "", err
	}

	return account.InternalUserID(), nil
}

func (a *API) getUserID(message *tgbotapi.Message) string {
	platformUserID := strconv.FormatInt(message.From.ID, 10)
	username := message.From.UserName

	internalUserID, err := a.getInternalUserID(platformUserID, username)
	if err != nil {
		log.Printf("Error getting internal user ID: %v", err)

		return platformUserID
	}

	return internalUserID
}

func (a *API) getUserIDFromCallback(callback *tgbotapi.CallbackQuery) string {
	platformUserID := strconv.FormatInt(callback.From.ID, 10)
	username := callback.From.UserName

	internalUserID, err := a.getInternalUserID(platformUserID, username)
	if err != nil {
		log.Printf("Error getting internal user ID: %v", err)

		return platformUserID
	}

	return internalUserID
}
