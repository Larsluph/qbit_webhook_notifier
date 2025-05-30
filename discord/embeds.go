package discord

import (
	"qbit_webhook/helpers"
	"qbit_webhook/qbit"
	"regexp"
	"time"
)

func escapeMarkdown(text string) string {
	re := regexp.MustCompile(`([\\*_~\|\(\)\[\]` + "`" + `])`)
	return re.ReplaceAllString(text, `\$1`)
}

func GenerateAddedEmbed(torrentProps qbit.TorrentProps) WebhookPayload {
	return WebhookPayload{
		Embeds: []Embed{
			{
				Title:       "New torrent added",
				Description: escapeMarkdown(torrentProps.Name),
				Color:       0xFF0000,
				Fields: []Field{
					{
						Name:   "Total Size",
						Value:  helpers.FormatByteSize(torrentProps.TotalSize),
						Inline: true,
					},
				},
				Footer: Footer{
					Text: torrentProps.Hash,
				},
				Datetime: time.Unix(torrentProps.AdditionDate, 0).Format(time.RFC3339),
			},
		},
	}
}

func GenerateCompletedEmbed(torrentProps qbit.TorrentProps) WebhookPayload {
	return WebhookPayload{
		Embeds: []Embed{
			{
				Title:       "Torrent completed",
				Description: escapeMarkdown("Dolly_Dyson_-_Dolly_s_full_vacation_movie.mp4"),
				Color:       0x00FF00,
				Fields: []Field{
					{
						Name:   "Total Size",
						Value:  helpers.FormatByteSize(torrentProps.TotalSize),
						Inline: true,
					},
					{
						Name:   "Torrent finished in",
						Value:  helpers.RelativeTimeElapsed(torrentProps.AdditionDate, torrentProps.CompletionDate),
						Inline: true,
					},
					{
						Name:   "Average DL Speed",
						Value:  helpers.FormatByteSpeed(torrentProps.DlSpeedAvg),
						Inline: true,
					},
				},
				Footer: Footer{
					Text: torrentProps.Hash,
				},
				Datetime: time.Unix(torrentProps.CompletionDate, 0).Format(time.RFC3339),
			},
		},
	}
}

func GenerateErrorEmbed(error helpers.ErrorPayload) WebhookPayload {
	return WebhookPayload{
		Content: "",
		Embeds: []Embed{{
			Title:       "Error Occured",
			Description: escapeMarkdown(error.Description),
			Color:       0xCA0BE4,
			Footer:      Footer{Text: error.CodeLocation},
			Datetime:    time.Now().Format(time.RFC3339),
		}},
	}
}
