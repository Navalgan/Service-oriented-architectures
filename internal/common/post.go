package common

type PostText struct {
	Text string `json:"text"`
}

type PostStatistic struct {
	PostID string `json:"post_id"`
	Author string `json:"author"`
	Count  uint64 `json:"count"`
}

type TopPosts struct {
	Posts []PostStatistic `json:"posts,omitempty"`
}

type UserStatistic struct {
	Login string `json:"login"`
	Count uint64 `json:"count"`
}

type TopUsers struct {
	Users []UserStatistic `json:"users,omitempty"`
}
