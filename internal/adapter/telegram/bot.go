package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
	"github.com/kitanoyoru/kitaDriveBot/internal/usecase"
	"github.com/rs/zerolog"
)

type Bot struct {
	ownerID       int64
	defaultPath   string
	uploadService *usecase.UploadService
	log           zerolog.Logger
	telegramBot   *bot.Bot
}

func NewBot(
	token string,
	ownerID int64,
	defaultPath string,
	uploadService *usecase.UploadService,
	log zerolog.Logger,
) (*Bot, error) {
	b := &Bot{
		ownerID:       ownerID,
		defaultPath:   defaultPath,
		uploadService: uploadService,
		log:           log,
	}

	tgBot, err := bot.New(token, bot.WithDefaultHandler(b.defaultHandler))
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	b.telegramBot = tgBot
	return b, nil
}

func (b *Bot) Start(ctx context.Context) {
	b.log.Info().Msg("starting telegram bot")
	b.telegramBot.Start(ctx)
}

func (b *Bot) defaultHandler(ctx context.Context, tgBot *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if !IsOwner(update.Message.From, b.ownerID) {
		if update.Message.From != nil {
			b.log.Debug().Int64("user_id", update.Message.From.ID).Msg("ignored message from non-owner")
		}
		return
	}

	if update.Message.Document == nil {
		if strings.TrimSpace(update.Message.Text) == "/start" {
			_, _ = tgBot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Send me a PDF with an optional caption as the Drive folder path (for example: Work/Invoices/2026).",
			})
		}
		return
	}

	b.handleDocument(ctx, tgBot, update.Message)
}

func (b *Bot) handleDocument(ctx context.Context, tgBot *bot.Bot, message *models.Message) {
	doc := message.Document
	if doc == nil {
		return
	}

	file := MapDocument(doc)
	path := domain.ParseFolderPath(message.Caption, b.defaultPath)

	if err := file.Validate(); err != nil {
		b.replyValidationError(ctx, tgBot, message.Chat.ID, err)
		return
	}

	b.reply(ctx, tgBot, message.Chat.ID, fmt.Sprintf("Uploading %s to %s...", file.FileName(), path.String()))

	result, err := b.uploadService.StorePDF(ctx, file, path)
	if err != nil {
		b.log.Error().Err(err).Str("path", path.String()).Msg("store pdf")
		b.reply(ctx, tgBot, message.Chat.ID, "Failed to upload the PDF to Google Drive.")
		return
	}

	b.reply(ctx, tgBot, message.Chat.ID, fmt.Sprintf("Stored in %s\n%s", result.FolderPath, result.WebViewLink))
}

func (b *Bot) replyValidationError(ctx context.Context, tgBot *bot.Bot, chatID int64, err error) {
	switch {
	case errors.Is(err, domain.ErrNotPDF):
		b.reply(ctx, tgBot, chatID, "Only PDF files are supported.")
	case errors.Is(err, domain.ErrFileTooLarge):
		b.reply(ctx, tgBot, chatID, "File is too large. Telegram Bot API supports downloads up to 20 MB.")
	default:
		b.reply(ctx, tgBot, chatID, "Invalid file.")
	}
}

func (b *Bot) reply(ctx context.Context, tgBot *bot.Bot, chatID int64, text string) {
	_, err := tgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		b.log.Error().Err(err).Str("text", text).Msg("send telegram message")
	}
}

func MapDocument(doc *models.Document) domain.IncomingFile {
	if doc == nil {
		return domain.IncomingFile{}
	}

	return domain.IncomingFile{
		ID:       doc.FileID,
		Name:     doc.FileName,
		Size:     doc.FileSize,
		MIMEType: doc.MimeType,
	}
}

func IsOwner(user *models.User, ownerID int64) bool {
	if user == nil {
		return false
	}
	return user.ID == ownerID
}
