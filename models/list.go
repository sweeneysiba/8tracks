package models

//Item ...
type List struct {
	ID        int64       `gorm:"primary_key" json:"id"`
	Name      string      `json:"name"`
	Likes     int64       `json:"likes"`
	Plays     int64       `json:"plays"`
	Songs     []SongLists `json:"songs"`
	CreatedAt int64       `json:"created_at"`
	Tag       string      `json:"tag"`
}

//UpdateList ...
type UpdateList struct {
	ID        int64             `gorm:"primary_key" json:"id"`
	Name      string            `json:"name"`
	Likes     int64             `json:"likes"`
	Songs     []UpdateSongLists `json:"songs"`
	Plays     int64             `json:"plays"`
	CreatedAt int64             `json:"created_at"`
	Tag       string            `json:"tag"`
}
