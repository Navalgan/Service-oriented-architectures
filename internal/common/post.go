package common

const (
	Like = 0
	View = 1
)

type PostText struct {
	Text string `json:"text"`
}

type PostStatistic struct {
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Operation int8   `json:"like_or_view"`
}
