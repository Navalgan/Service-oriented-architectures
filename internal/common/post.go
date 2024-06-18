package common

type PostText struct {
	Text string `json:"text"`
}

type PostStatistic struct {
	PostID string `json:"post_id"`
	UserID string `json:"user_id"`
}
