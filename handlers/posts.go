package handlers

import (
	"fmt"

	"strconv"

	"errors"

	"github.com/almirpernen/database"
	"github.com/almirpernen/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(cp *fiber.Ctx) error {
	post := new(models.Post)
	if err := cp.BodyParser(post); err != nil {
		return cp.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error parsing request body"})
	}

	userIDInterface := cp.Locals("userID")
	fmt.Printf("Extracted userIDInterface: %#v\n", userIDInterface)

	if userIDInterface == nil {
		fmt.Println("userIDInterface is nil. User ID was not set in context.")
		return cp.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {
		fmt.Printf("Failed to assert userIDInterface to uint or userID is 0. userIDInterface actual value: %#v\n", userIDInterface)
		return cp.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID format or userID is 0"})
	}

	var user models.User
	if err := database.DB.Db.First(&user, userID).Error; err != nil {
		return cp.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User does not exist", "userID": userID})
	}

	post.UserID = userID
	if err := database.DB.Db.Create(&post).Error; err != nil {
		return cp.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving the post to the database", "error": err.Error()})
	}

	database.DB.Db.First(&post.User, post.UserID)

	return cp.Status(200).JSON(post)
}

func ListPosts(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

    sortField := c.Query("sortField", "created_at")
    sortOrder := c.Query("sortOrder", "desc")

    validSortFields := map[string]bool{
        "created_at":  true,
        "content":     true,
        "likes_count": true,
    }

    if _, ok := validSortFields[sortField]; !ok {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid sort field"})
    }

    var posts []models.Post

    offset := (page - 1) * pageSize
    query := database.DB.Db.Model(&models.Post{}).Offset(offset).Limit(pageSize).Order(fmt.Sprintf("%s %s", sortField, sortOrder)).Preload("User")

    if err := query.Find(&posts).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching comments"})
    }

    for i, post := range posts {
        var count int64
        database.DB.Db.Model(&models.PostLike{}).Where("post_id = ?", post.ID).Count(&count)
        posts[i].LikesCount = int(count)


        posts[i].Username = post.User.Username
    }

    return c.Status(200).JSON(posts)
}

func GetPost(c *fiber.Ctx) error {
    postID := c.Params("id")
    var post models.Post

    if err := database.DB.Db.Preload("User").First(&post, postID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Post not found"})
    }

    var count int64
    database.DB.Db.Model(&models.PostLike{}).Where("post_id = ?", post.ID).Count(&count)
    post.LikesCount = int(count)


    post.Username = post.User.Username

    return c.Status(fiber.StatusOK).JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	var post models.Post
	if err := database.DB.Db.Preload("Comments").First(&post, postID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Post not found"})
	}

	for _, comment := range post.Comments {
		if err := database.DB.Db.Delete(&comment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting associated comments", "error": err.Error()})
		}
	}

	if err := database.DB.Db.Delete(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting the post", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post and associated comments deleted successfully"})
}

func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	var existingPost models.Post

	if err := database.DB.Db.First(&existingPost, postID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Post not found",
		})
	}

	newPost := new(models.Post)
	if err := c.BodyParser(newPost); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Model(&existingPost).Updates(newPost)
	database.DB.Db.First(&existingPost.User, existingPost.UserID)

	return c.Status(fiber.StatusOK).JSON(existingPost)
}

func LikePost(c *fiber.Ctx) error {
	postIDParam := c.Params("id")
	userIDInterface := c.Locals("userID")

	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid post ID"})
	}

	var existingLike models.PostLike
	result := database.DB.Db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike)

	if result.Error == nil {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{"message": "User has already liked this post"})
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error checking for existing like", "error": result.Error.Error()})
	}

	newLike := models.PostLike{
		UserID: userID,
		PostID: uint(postID),
	}

	if err := database.DB.Db.Create(&newLike).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not like the post", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post liked successfully"})
}

func UnlikePost(c *fiber.Ctx) error {
	postIDParam := c.Params("id")
	userIDInterface := c.Locals("userID")

	if userIDInterface == nil {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid post ID"})
	}

	var postLike models.PostLike

	if err := database.DB.Db.Where("user_id = ? AND post_id = ?", userID, postID).First(&postLike).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Like not found"})
	}

	if err := database.DB.Db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.PostLike{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not unlike the post", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Post unliked successfully"})
}
