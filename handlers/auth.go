package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/almirpernen/database"
	"github.com/almirpernen/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("your-secret-key")

func Signup(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}
	user.Password = string(hashedPassword)

	database.DB.Db.Create(&user)

	user.Password = ""
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Signin(c *fiber.Ctx) error {
	loginData := new(models.User)
	if err := c.BodyParser(loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	var user models.User
	if err := database.DB.Db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error querying the database",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["username"] = user.Username
	accessClaims["userID"] = user.ID
	accessClaims["exp"] = time.Now().Add(15 * time.Minute).Unix() // Short expiry for access token

	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating JWT token",
		})
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["userID"] = user.ID
	refreshClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating refresh token",
		})
	}

	// Send both tokens to the client
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	})
}

func ListUsers(c *fiber.Ctx) error {
	var users []models.User

	if err := database.DB.Db.Preload("Posts").Preload("Posts.Comments").Preload("Comments").Preload("Followers").Preload("Followings").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving users"})
	}

	for i := range users {
		for j := range users[i].Posts {
			var likesCount int64
			if err := database.DB.Db.Model(&models.PostLike{}).Where("post_id = ?", users[i].Posts[j].ID).Count(&likesCount).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving likes count for posts"})
			}
			users[i].Posts[j].LikesCount = int(likesCount)

			for k := range users[i].Posts[j].Comments {

				var commentLikesCount int64
				if err := database.DB.Db.Model(&models.CommentLike{}).Where("comment_id = ?", users[i].Posts[j].Comments[k].ID).Count(&commentLikesCount).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving likes count for comments"})
				}
				users[i].Posts[j].Comments[k].LikesCount = int(commentLikesCount)
			}
		}
	}

	for i := range users {
		for j := range users[i].Comments {
			var likesCount int64
			if err := database.DB.Db.Model(&models.CommentLike{}).Where("comment_id = ?", users[i].Comments[j].ID).Count(&likesCount).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving likes count for comments"})
			}
			users[i].Comments[j].LikesCount = int(likesCount)
		}
	}

	return c.Status(200).JSON(users)
}

func GetUsers(c *fiber.Ctx) error {
	userID := c.Params("id")
	var user models.User

	if err := database.DB.Db.Preload("Followers").Preload("Followings").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var user models.User

	if err := database.DB.Db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error querying the database",
		})
	}

	if err := database.DB.Db.Delete(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func FollowUser(c *fiber.Ctx) error {
	followerIDInterface := c.Locals("userID")
	if followerIDInterface == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}
	followerID := followerIDInterface.(uint) // Assuming userID is stored as uint

	followingIDParam := c.Params("id")
	followingID, err := strconv.ParseUint(followingIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	if followerID == uint(followingID) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot follow yourself"})
	}

	var follower, following models.User
	if err := database.DB.Db.First(&follower, followerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Follower not found"})
	}

	if err := database.DB.Db.First(&following, followingID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User to follow not found"})
	}

	database.DB.Db.Model(&follower).Association("Followings").Append(&following)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Followed successfully"})
}

func UnfollowUser(c *fiber.Ctx) error {
	followerID := c.Locals("userID")
	if followerID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User not authenticated"})
	}

	unfollowingID := c.Params("id")
	if followerID == unfollowingID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Cannot unfollow yourself"})
	}

	var follower, unfollowing models.User
	if err := database.DB.Db.First(&follower, followerID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Follower not found"})
	}

	if err := database.DB.Db.First(&unfollowing, unfollowingID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User to unfollow not found"})
	}

	database.DB.Db.Model(&follower).Association("Followings").Delete(&unfollowing)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Unfollowed successfully"})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshTokenString := c.FormValue("refreshToken")
	if refreshTokenString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "No refresh token provided"})
	}

	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid refresh token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userID"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token claims"})
	}

	newAccessToken := jwt.New(jwt.SigningMethodHS256)
	newAccessClaims := newAccessToken.Claims.(jwt.MapClaims)
	newAccessClaims["userID"] = claims["userID"]
	newAccessClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()

	newAccessTokenString, err := newAccessToken.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating new access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken": newAccessTokenString,
	})
}
