package models

//UpdateSongLists ...
type UpdateSongLists struct {
	ID     int64  `json:"id"`
	SongID int64  `json:"song_id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	ListID int64  `json:"list_id"`
	Type   string `json:"type"`
}

//Songs ...
type Songs struct {
	ID        int64  `gorm:"primary_key;autoIncrement" json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	CreatedAt int64  `json:"created"`
}

//InsertSongs ...
type InsertSongs struct {
	Name      string `json:"name"`
	Author    string `json:"author"`
	CreatedAt int64  `json:"created"`
}

//SongList ...
type SongList struct {
	ID     int64  `gorm:"primary_key;autoIncrement" json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	ListID int64  `json:"list_id"`
	SongID int64  `json:"song_id"`
}

//SongLists ...
type SongLists struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	ListID int64  `json:"list_id"`
	SongID int64  `json:"song_id"`
}
