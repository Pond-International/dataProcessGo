package models

type User struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	CreatedAt       string `json:"created_at"`
	Username        string `json:"username"`
	Name            string `json:"name"`
	PublicMetrics   struct {
		FollowersCount int `json:"followers_count"`
		FollowingCount int `json:"following_count"`
		TweetCount     int `json:"tweet_count"`
		ListedCount    int `json:"listed_count"`
		LikeCount      int `json:"like_count"`
	} `json:"public_metrics"`
}

type Response struct {
	Data []User `json:"data"`
}
