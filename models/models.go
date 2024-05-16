package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Posts      []Post    `json:"posts" gorm:"foreignKey:UserID"`
	Comments   []Comment `json:"comments"`
	Followers  []*User   `json:"followers" gorm:"many2many:user_followers;joinForeignKey:FollowingID;JoinReferences:FollowerID"`
	Followings []*User   `json:"followings" gorm:"many2many:user_followers;joinForeignKey:FollowerID;JoinReferences:FollowingID"`
}

type Post struct {
    gorm.Model
    Content    string    `json:"content"`
    UserID     uint      `json:"user_id"`
    Username   string    `json:"username" gorm:"-"`
    User       User      `json:"-" gorm:"foreignKey:UserID"`
    Comments   []Comment `json:"comments" gorm:"foreignKey:PostID"`
    LikesCount int       `json:"likes_count" gorm:"-"`
}

type Comment struct {
    gorm.Model
    Content    string `json:"content"`
    UserID     uint   `json:"user_id"`
    Username   string `json:"username" gorm:"-"`
	User       User      `json:"-" gorm:"foreignKey:UserID"`
    PostID     uint   `json:"post_id"`
    LikesCount int    `json:"likes_count" gorm:"-"`
}

type PostLike struct {
	UserID uint `gorm:"index:idx_user_post,uniqueIndex" json:"user_id"`
	PostID uint `gorm:"index:idx_user_post,uniqueIndex" json:"post_id"`
}

type CommentLike struct {
	UserID    uint `gorm:"index:idx_user_comment,uniqueIndex" json:"user_id"`
	CommentID uint `gorm:"index:idx_user_comment,uniqueIndex" json:"comment_id"`
}
