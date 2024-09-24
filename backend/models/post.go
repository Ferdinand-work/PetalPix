package models

import "time"

type Post struct {
	UserId      string    `json:"userId" bson:"user_id"`
	Title       string    `json:"title" bson:"post_title"`
	Content     string    `json:"content" bson:"post_content"`
	DateCreated time.Time `json:"dateCreated" bson:"post_date_created"`
	DateEdited  time.Time `json:"dateEdited" bson:"post_date_edited"`
	Tags        []string  `json:"tags,omitempty" bson:"post_tags"`
	ImagePath   string    `json:"imagePath,omitempty" bson:"post_image_path"`
}
