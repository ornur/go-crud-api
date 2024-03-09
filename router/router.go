package router

import (
	"net/http"

	db "example.com/go-crud-api/db"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
	
		// Retrieve user from the database
		user, err := db.GetUser(username)
		if err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error":   true,
				"message": "User does not exist",
			})
			return
		}

		// Check if password is correct
		if user.Password != password {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"error":   true,
				"message": "Invalid password",
			})
			return
		}
	
		// If username and password are correct, render the success template
		c.HTML(http.StatusOK, "login_success.html", gin.H{
			"username": username,
		})
	})
	
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	
	r.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
	
		// Check if username or password is empty
		if username == "" || password == "" {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error":   true,
				"message": "Username and password cannot be empty",
			})
			return
		}
	
		// Check if username already exists
		_, err := db.GetUser(username)
		if err == nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"error":   true,
				"message": "Username already exists",
			})
			return
		}
	
		// Create user
		newUser := db.User{
			Username: username,
			Password: password,
		}
	
		_, err = db.CreateUser(&newUser)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error":   true,
				"message": "Error creating user",
			})
			return
		}
	
		c.HTML(http.StatusOK, "register_success.html", gin.H{
			"username": username,
		})
	})
	
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.POST("/users", postUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
	return r
}
func postUser(ctx *gin.Context) {
	var user db.User
	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := db.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"user": res,
	})
}
func getUsers(ctx *gin.Context) {
	res, err := db.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users": res,
	})
}

func getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := db.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"user": res,
	})
}
func updateUser(ctx *gin.Context) {
	var updatedUser db.User
	err := ctx.Bind(&updatedUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	id := ctx.Param("id")
	dbUser, err := db.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	dbUser.Username = updatedUser.Username
	dbUser.Password = updatedUser.Password

	res, err := db.UpdateUser(dbUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"task": res,
	})
}
func deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := db.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}
