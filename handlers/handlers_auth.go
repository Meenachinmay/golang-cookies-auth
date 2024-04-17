package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-cookies/handlers/models"
	"golang-cookies/utils"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token  string    `json:"string"`
	Expire time.Time `json:"expires"`
}

type SessionData struct {
	Token  string    `json:"string"`
	UserId uuid.UUID `json:"userId"`
}

func (lac *LocalApiConfig) SignInHandler(c *gin.Context) {
	var userToAuth models.UserToAuth

	if err := c.ShouldBindJSON(&userToAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// do validation here
	validationErrors := utils.ValidateUserAuth(userToAuth)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": validationErrors,
		})
		return
	}

	//validationResult := utils.ValidateEmail(userToAuth.Email)
	//if validationResult.Error != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": validationResult.Error.Error(),
	//	})
	//	return
	//}
	//
	//if !validationResult.IsValid {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": validationResult.Error.Error(),
	//	})
	//	return
	//}

	// fetch the user from database to match
	foundUser, err := lac.DB.FindUserByEmail(c, userToAuth.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No user found with this email or password.",
		})
		return
	}

	if foundUser.Password != userToAuth.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Email: userToAuth.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// till here you are having token
	sessionID := uuid.New().String()

	// prepare the session data in redis
	sessionData := map[string]interface{}{
		"token":  tokenString,
		"userId": foundUser.ID,
	}

	sessionDataJSON, err := json.Marshal(sessionData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode the json data",
		})
		return
	}

	// saving the session data in redis
	lac.RedisClient.Set(c, sessionID, sessionDataJSON, time.Until(expirationTime)).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save token in session",
		})
		return
	}

	// set the session in cookies for the client side
	c.SetCookie("session_id", sessionID, int(time.Until(expirationTime).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"expires": expirationTime,
	})

}

func (lac *LocalApiConfig) LogoutHandler(c *gin.Context) {
	// Retrieve the session id from the session
	sessionID, err := c.Cookie("session_id")

	if err != nil {
		// Handle missing or invalid cookie
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized - No session to logout",
		})
		return
	}

	// Delete the session data from redis
	err = lac.RedisClient.Del(c, sessionID).Err()
	if err != nil {
		// Handle potential redis err
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to end session",
		})
		return
	}

	// Delete the cookies by setting the maximum age of -1
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	// Respond to the client that logout successful
	c.JSON(http.StatusOK, gin.H{
		"error": "Logged Out successfully",
	})
}
