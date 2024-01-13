package discord

import (
	"qbit_webhook/helpers"
	"qbit_webhook/qbit"
	"time"
)

func GenerateAddedEmbed(torrentProps qbit.TorrentProps) WebhookPayload {
	return WebhookPayload{
		Embeds: []Embed{
			{
				Title:       "New torrent added",
				Description: torrentProps.Name,
				Color:       0xFF0000,
				URL:         torrentProps.Comment,
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
				Description: torrentProps.Name,
				Color:       0x00FF00,
				URL:         torrentProps.Comment,
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
			Description: error.Description,
			Color:       0xCA0BE4,
			Footer:      Footer{Text: error.CodeLocation},
			Datetime:    time.Now().Format(time.RFC3339),
		}},
	}
}
