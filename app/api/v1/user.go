package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"chatapp/app/biz/models"

	"github.com/gin-gonic/gin"
)

// verifyUser
func (api *Apiv1) verifyUser(c *gin.Context) {
	token, err := c.Cookie("ca-token")
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"msg": "API token required"})
		return
	}
	user := &models.User{}
	err = json.Unmarshal([]byte(token), user)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"msg": "API token required"})
		return
	}
	if user.ID == 0 || user.Email == "" || user.Name == "" {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"msg": "API token required"})
		return
	}
	c.Next()
}

// checkComplete ...
func (api *Apiv1) checkComplete(c *gin.Context) {
	token, err := c.Cookie("ca-token")
	if err != nil {
		c.Next()
		return
	}
	user := &models.User{}
	err = json.Unmarshal([]byte(token), user)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}
	c.HTML(http.StatusOK, "home.html", gin.H{"user": user})
}

// Register ...
// @Summary Register
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param username formData string true "username"
// @Param email formData string true "email of the user"
// @Param password formData string true "some password"
// @Success 200 {string} text/html "HTML page to login"
// @Success 400 {string} text/html "HTML page to register with the error"
// @Success 500 {string} text/html "HTML page to register with the error"
// @Router /register [post]
func (api *Apiv1) Register(c *gin.Context) {
	name := c.Request.PostFormValue("username")
	email := c.Request.PostFormValue("email")
	pass := c.Request.PostFormValue("password")

	// validate entry
	if name == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"msg": "check your body request"})
		return
	}
	if email == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"msg": "check your body request"})
		return
	}
	if pass == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"msg": "check your body request"})
		return
	}

	newUser := &models.User{
		ID:       0,
		Name:     name,
		Email:    email,
		Password: pass,
	}

	id, err := api.bizLayer.CreateUser(newUser)
	if err != nil {
		api.logger.Printf("Error creating the user: %+v\n", err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"msg": err.Error()})
		return
	}

	api.logger.Printf("New user created id: %dn", id)
	c.HTML(http.StatusOK, "login.html", gin.H{"msg": "user created"})
}

// LogIn ...
// @Summary LogIn
// @Description LogIn to an account
// @Tags Users
// @Accept json
// @Produce json
// @Param email formData string true "email of the user"
// @Param password formData string true "some password"
// @Success 200 {string} text/html "HTML page to home"
// @Success 400 {string} text/html "HTML page to login with the error"
// @Success 500 {string} text/html "HTML page to login with the error"
// @Router /login [post]
func (api *Apiv1) LogIn(c *gin.Context) {
	email := c.Request.PostFormValue("email")
	pass := c.Request.PostFormValue("password")

	// validate entry
	if email == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"msg": "check your body request"})
		return
	}
	if pass == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"msg": "check your body request"})
		return
	}

	user, err := api.bizLayer.GetUser(email, pass)
	if err != nil {
		api.logger.Printf("Error creating the user: %+v\n", err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"msg": err.Error()})
		return
	}

	if user == nil {
		api.logger.Println("User not found.")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"msg": "User not found."})
		return
	}

	user.Password = "******"
	api.logger.Printf("get user: %+v\n", user)

	xuser, err := json.Marshal(user)
	if err != nil {
		str := fmt.Sprintf("Error adding user token: %+v\n", err)
		api.logger.Printf(str)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"msg": str})
		return
	}

	fmt.Println("xuser: ", string(xuser))
	// add token to the header of the response
	c.SetCookie("ca-token", string(xuser), 60*60*24, "/", "", true, true)

	c.HTML(http.StatusOK, "home.html", gin.H{"user": user})
}

// Logout ...
// @Summary Logout
// @Description Logout from an account
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {string} text/html "HTML page to login"
// @Router /logout [post]
func (api *Apiv1) Logout(c *gin.Context) {
	c.SetCookie("ca-token", "", -1, "", "", false, true)
	c.HTML(http.StatusOK, "login.html", nil)
}
