package models

type FollowRequest struct {
	FollowUser  string   `json:"followUser,omitempty"`
	FollowUsers []string `json:"followUsers,omitempty"`
}

type UnfollowRequest struct {
	UnfollowUser  string   `json:"unfollowUser,omitempty"`
	UnfollowUsers []string `json:"unfollowUsers,omitempty"`
}
