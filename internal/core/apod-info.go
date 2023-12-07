package core

type APODInfo struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Date        Date   `json:"date"`
	Image       string `json:"image"`
}
