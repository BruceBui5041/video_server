package usertransport

import (
	"net/http"
	"video_server/component"
	"video_server/model/user/userbiz"
	user "video_server/model/user/usermodel"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input user.CreateUser

		// Decode the request body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userstore := userstore.NewSQLStore(appCtx.GetMainDBConnection())
		userBiz := userbiz.NewRegisterBusiness(userstore)

		// Create the user
		newUser, err := userBiz.RegisterUser(c.Request.Context(), &input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Prepare and send the response
		c.JSON(http.StatusCreated, gin.H{
			"username": newUser.Username,
			"email":    newUser.Email,
		})
	}
}
