package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"ot-recorder/app/model"
	"ot-recorder/app/response"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const mapLink = "https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=18/%f/%f"
const two = 2
const locationCMD = "location"
const helpCMD = "help"
const helpDescription = `*/loc<space><username>* - get user last location
*/location<space><username>* - get user last location
*/help* - for a list of commands`

var commands = map[string]string{
	"/location": locationCMD,
	"/loc":      locationCMD,
	"/help":     helpCMD,
}

type locationUsecase struct {
	repo           model.LocationRepository
	contextTimeout time.Duration
}

func NewLocationUsecase(repo model.LocationRepository, timeout time.Duration) model.LocationUsecase {
	return &locationUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *locationUsecase) Ping(c context.Context, l *model.Location) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// store location
	err = u.repo.CreateLocation(ctx, l)
	if err != nil {
		logrus.Errorln(err)

		return response.WrapError(errors.New("internal server error, please report to admin"), http.StatusInternalServerError)
	}

	return nil
}

func (u *locationUsecase) LastLocation(
	c context.Context,
	username string,
) (location *model.LocationDetails, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	l, err := u.repo.GetUserLastLocation(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}

		logrus.Errorln(err)

		return nil, response.WrapError(
			errors.New("internal server error, please report to admin"),
			http.StatusInternalServerError,
		)
	}

	return toLastLocationDetails(&l), nil
}

func (u *locationUsecase) TelegramHook(c context.Context, req *model.TelegramRequest) (res *model.TelegramResponse) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	text := ""
	command, username := extractTelegramCommand(req.Message.Text)

	switch command {
	case locationCMD:
		text = getUserLocationForTelegram(u, ctx, username)
	case helpCMD:
		text = helpDescription
	default:
		text = helpDescription
	}

	return &model.TelegramResponse{
		Method:                "sendMessage",
		ChatID:                req.Message.Chat.ID,
		ReplyToMessageID:      req.Message.MessageID,
		Text:                  text,
		ParseMode:             "Markdown",
		DisableWebPagePreview: true,
	}
}

func extractTelegramCommand(text string) (command, username string) {
	textRaw := strings.Split(text, " ")
	if len(textRaw) >= two {
		command = strings.Trim(textRaw[0], " ")
		username = strings.Trim(textRaw[1], " ")

		if strings.Contains(command, "@") {
			cmdPart := strings.Split(command, "@")
			command = cmdPart[0]
		}

		if v, ok := commands[command]; ok {
			command = v
		} else {
			command = helpCMD
		}

		if username == "" {
			command = helpCMD
		} else {
			username = strings.TrimPrefix(username, "@")
		}
	}

	return command, username
}

func getUserLocationForTelegram(u *locationUsecase, ctx context.Context, name string) string {
	loc, err := u.repo.GetUserLastLocation(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "*username not found!*"
		}

		logrus.Errorln(err)

		return "*internal server error, please report to admin.*"
	}

	return toTelegramMessage(&loc)
}
