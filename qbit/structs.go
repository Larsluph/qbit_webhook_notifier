package qbit

type TorrentProps struct {
	AdditionDate   int64  `json:"addition_date"`
	Comment        string `json:"comment"`
	CompletionDate int64  `json:"completion_date"`
	CreationDate   int64  `json:"creation_date"`
	DlSpeedAvg     int64  `json:"dl_speed_avg"`
	Hash           string `json:"hash"`
	Name           string `json:"name"`
	TotalSize      int64  `json:"total_size"`
}
