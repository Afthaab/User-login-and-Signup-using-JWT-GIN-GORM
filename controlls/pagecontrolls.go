package controlls

import (
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/project_login/config"
	"github.com/project_login/models"
)

func SignupPage(c *gin.Context) {
	//Integrating the html Page
	tmpl := template.Must(template.ParseFiles("views/signup.html"))
	tmpl.Execute(c.Writer, nil)

}

func SignupUser(c *gin.Context) {
	c.Request.ParseForm()
	type Userdata struct {
		fname    string
		lname    string
		uname    string
		mail     string
		password string
	}
	fieldvalidate := models.Errors{Errors: "The fields are empty"}
	usernamevalidate := models.Errors{Errors: "This username is not available"}
	emailvalidate := models.Errors{Errors: "This email is invalid"}
	passwordvalidate := models.Errors{Errors: "Ateast 7 charecters; include atleast 1 number, 1 Uppercase, 1 Speical Charecter"}
	DB := config.DBConn()
	var data Userdata
	var temp_user models.User
	data.fname = c.Request.PostForm["firstname"][0]
	data.lname = c.Request.PostForm["lastname"][0]
	data.uname = c.Request.PostForm["username"][0]
	data.mail = c.Request.PostForm["email"][0]
	data.password = c.Request.PostForm["password"][0]

	if data.fname == "" || data.lname == "" || data.uname == "" || data.mail == "" || data.password == "" {
		tmpl := template.Must(template.ParseFiles("views/signup.html"))
		tmpl.Execute(c.Writer, fieldvalidate)
	} else if !isEmailValid(data.mail) {
		tmpl := template.Must(template.ParseFiles("views/signup.html"))
		tmpl.Execute(c.Writer, emailvalidate)
	} else if !isValid(data.password) {
		tmpl := template.Must(template.ParseFiles("views/signup.html"))
		tmpl.Execute(c.Writer, passwordvalidate)
	} else {
		result1 := DB.First(&temp_user, "username LIKE ?", data.uname)
		if result1.Error != nil {
			user := models.User{First_name: data.fname, Last_name: data.lname, Username: data.uname, Email: data.mail, Password: data.password}
			result2 := DB.Create(&user)
			if result2.Error != nil {
				c.Redirect(http.StatusMovedPermanently, "/signup")
				return
			} else {
				c.Redirect(http.StatusMovedPermanently, "/login")
			}
		} else {

			tmpl := template.Must(template.ParseFiles("views/signup.html"))
			tmpl.Execute(c.Writer, usernamevalidate)
		}

	}

}

func Loginpage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("views/login.html"))
	tmpl.Execute(c.Writer, nil)
}

func Loginuser(c *gin.Context) {
	invaliduser := models.Errors{Errors: "This is user is Invalid"}
	fieldvalidate := models.Errors{Errors: "The fields are empty"}
	type data struct {
		uname    string
		password string
	}
	var userdata data
	c.Request.ParseForm()
	userdata.uname = c.Request.PostForm["username"][0]
	userdata.password = c.Request.PostForm["password"][0]
	if userdata.uname == "" || userdata.password == "" {
		tmpl := template.Must(template.ParseFiles("views/login.html"))
		tmpl.Execute(c.Writer, fieldvalidate)
	} else {
		DB := config.DBConn()
		var temp_user models.User
		result := DB.First(&temp_user, "username LIKE ? AND password LIKE ? AND is_admin LIKE ?", userdata.uname, userdata.password, "no")
		if result.Error != nil {
			tmpl := template.Must(template.ParseFiles("views/login.html"))
			tmpl.Execute(c.Writer, invaliduser)
		} else {
			//creating a JWT Token

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": temp_user.Username,
				"exp":      time.Now().Add(time.Hour * 24 * 3).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "jwt token error",
				})
				return
			}

			//return response
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Authorization", tokenString, 36000*24*30, "", "", false, true)
			c.Redirect(http.StatusMovedPermanently, "/home")
		}

	}

}

func Homepage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("views/home.html"))
	tmpl.Execute(c.Writer, nil)
}

func Adminloginpage(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("views/adminsignin.html"))
	tmpl.Execute(c.Writer, nil)
}
func Adminlogin(c *gin.Context) {
	fieldvalidate := models.Errors{Errors: "The Fields are empty"}
	adminvalidate := models.Errors{Errors: "Wrong Credentials"}
	c.Request.ParseForm()
	type data struct {
		uname    string
		password string
	}
	var userdata data
	var temp_user models.User
	userdata.uname = c.Request.PostForm["username"][0]
	userdata.password = c.Request.PostForm["password"][0]
	if userdata.uname == "" || userdata.password == "" {
		tmpl := template.Must(template.ParseFiles("views/adminsignin.html"))
		tmpl.Execute(c.Writer, fieldvalidate)
	} else {
		DB := config.DBConn()
		result := DB.Find(&temp_user, "username LIKE ? AND password LIKE ? AND is_admin LIKE ?", userdata.uname, userdata.password, "yes")
		if result.Error != nil {
			tmpl := template.Must(template.ParseFiles("views/adminsignin.html"))
			tmpl.Execute(c.Writer, adminvalidate)
		} else {
			c.Redirect(http.StatusMovedPermanently, "/adminpanel")
		}
	}

}

func Adminpanel(c *gin.Context) {

	var temp_user []models.User
	DB := config.DBConn()
	result := DB.Find(&temp_user)
	if result.Error != nil {
		c.Redirect(http.StatusMovedPermanently, "/adminloginpage")
	} else {
		tmpl := template.Must(template.ParseFiles("views/adminpanel.html"))
		tmpl.Execute(c.Writer, temp_user)
	}

}
func Delete(c *gin.Context) {
	c.Request.ParseForm()
	delete_id := c.Request.PostForm["id"]
	DB := config.DBConn()
	var temp_user models.User
	DB.Delete(&temp_user, delete_id)
	c.Redirect(http.StatusMovedPermanently, "/adminpanel")
}
func Search(c *gin.Context) {
	c.Request.ParseForm()
	search_uname := c.Request.PostForm["username"]
	DB := config.DBConn()
	var temp_user []models.User
	result := DB.First(&temp_user, " username LIKE ? ", search_uname)
	if result.Error != nil {
		c.Redirect(http.StatusMovedPermanently, "/adminpanel")

	} else {
		tmpl := template.Must(template.ParseFiles("views/adminpanel.html"))
		tmpl.Execute(c.Writer, temp_user)

	}

}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
func isValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
