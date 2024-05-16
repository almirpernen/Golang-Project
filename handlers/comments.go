package handlers

import (
	"fmt"
	"strconv"

	"github.com/almirpernen/database"
	"github.com/almirpernen/models"
	"github.com/gofiber/fiber/v2"
	"errors"
	"gorm.io/gorm"
)

func CreateComment(cp *fiber.Ctx) error {
	comment := new(models.Comment)
	postIDParam := cp.Params("id")
	if err := cp.BodyParser(comment); err != nil {
		return cp.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error parsing request body"})
	}

	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		return cp.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid post ID format"})
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

	var post models.Post
	if err := database.DB.Db.First(&post, postID).Error; err != nil {
		return cp.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Post not found"})
	}

	comment.UserID = userID
	comment.PostID = uint(postID)

	if err := database.DB.Db.Create(&comment).Error; err != nil {
		return cp.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving the comment to the database", "error": err.Error()})
	}

	return cp.Status(200).JSON(comment)
}

func ListComments(c *fiber.Ctx) error {
	// Pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))          // Default to page 1 if not specified
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10")) // Default page size to 10 if not specified

	// Sorting parameters
	sortField := c.Query("sortField", "created_at") // Default sorting by created_at date
	sortOrder := c.Query("sortOrder", "desc")       // Default sort order

	// Validate sort field to prevent SQL injection or errors
	validSortFields := map[string]bool{
		"created_at":  true,
		"content":     true,
		"likes_count": true,
	}
	if _, ok := validSortFields[sortField]; !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid sort field"})
	}

	comments := []models.Comment{}

	offset := (page - 1) * pageSize
	query := database.DB.Db.Model(&models.Comment{}).Offset(offset).Limit(pageSize).Order(fmt.Sprintf("%s %s", sortField, sortOrder)).Preload("User")

	if err := query.Find(&comments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching comments"})
	}

	for i := range comments {
		var likesCount int64
		if err := database.DB.Db.Model(&models.CommentLike{}).Where("comment_id = ?", comments[i].ID).Count(&likesCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching likes count"})
		}
		comments[i].LikesCount = int(likesCount)

		// Assuming Comment struct has a Username field
		comments[i].Username = comments[i].User.Username
	}

	return c.Status(200).JSON(comments)
}

func GetComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var comment models.Comment

	if err := database.DB.Db.Preload("User").First(&comment, commentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Comment not found"})
	}

	var count int64
	database.DB.Db.Model(&models.CommentLike{}).Where("comment_id = ?", comment.ID).Count(&count)
	comment.LikesCount = int(count)

	// Assuming Comment struct has a Username field
	comment.Username = comment.User.Username

	return c.Status(fiber.StatusOK).JSON(comment)
}

func DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var comment models.Comment
	if err := database.DB.Db.First(&comment, commentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Comment not found",
		})
	}

	database.DB.Db.Delete(&comment)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}

func UpdateComment(c *fiber.Ctx) error {
	commentID := c.Params("id")
	var existingComment models.Comment

	if err := database.DB.Db.First(&existingComment, commentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Comment not found",
		})
	}

	newComment := new(models.Comment)
	if err := c.BodyParser(newComment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Model(&existingComment).Updates(newComment)

	return c.Status(fiber.StatusOK).JSON(existingComment)
}

func LikeComment(c *fiber.Ctx) error {
	commentIDParam := c.Params("id")
	userIDInterface := c.Locals("userID")

	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	commentID, err := strconv.ParseUint(commentIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid comment ID"})
	}

	var existingLike models.CommentLike
	result := database.DB.Db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingLike)

	if result.Error == nil {
		return c.Status(fiber.StatusAlreadyReported).JSON(fiber.Map{"message": "User has already liked this comment"})
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error checking for existing like", "error": result.Error.Error()})
	}

	newLike := models.CommentLike{
		UserID:    userID,
		CommentID: uint(commentID),
	}

	if err := database.DB.Db.Create(&newLike).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not like the comment", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment liked successfully"})
}

func UnlikeComment(c *fiber.Ctx) error {
	commentIDParam := c.Params("id")
	userIDInterface := c.Locals("userID")

	if userIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	commentID, err := strconv.ParseUint(commentIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid comment ID"})
	}

	var commentlike models.CommentLike

	if err := database.DB.Db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&commentlike).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Like not found"})
	}

	if err := database.DB.Db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not unlike the comment", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Comment unliked successfully"})
}
