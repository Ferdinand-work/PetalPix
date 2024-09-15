package models

type User struct {
	UserId         string   `json:"userId,omitempty" bson:"user_id"`
	Name           string   `json:"name,omitempty" bson:"user_name"`
	ContactNo      string   `json:"contactNo,omitempty" bson:"user_contact_no"`
	Email          string   `json:"email,omitempty" bson:"user_email"`
	Password       string   `json:"password,omitempty" bson:"user_password"`
	FollowingCount int64    `json:"followingCount,omitempty" bson:"user_following_count"`
	Following      []string `json:"following,omitempty" bson:"user_following"`
	FollowersCount int64    `json:"followersCount,omitempty" bson:"user_followers_count"`
	Followers      []string `json:"followers,omitempty" bson:"user_followers"`
}
