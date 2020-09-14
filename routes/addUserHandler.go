package routes

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
	gomail "gopkg.in/gomail.v2"
	"zsbrybnik.pl/server-go/utils"
)

type addUserJSON struct {
	Login string `json:"login"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AddUserEmailTemplate struct {
	Password string
	QrCode   string
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// AddUserHandler - Handling add-user route
func AddUserHandler(context *gin.Context) {
	token := context.GetHeader("Authorization")
	err := utils.VerifyToken(token)
	utils.ErrorHandler(err, false)
	if err == nil {
		err := utils.VerifyToken(token)
		utils.ErrorHandler(err, false)
		if err == nil {
			var addUserData addUserJSON
			err := context.Bind(&addUserData)
			utils.ErrorHandler(err, false)
			if err == nil {
				database, ok := context.MustGet("database").(*sql.DB)
				if ok {
					password := randomString(10)
					bcrytedPassword := utils.CreateBCryptHash(password)
					secretLength := 16
					googleAuthCode := gotp.RandomSecret(secretLength)
					query := "INSERT INTO users (login, password, email, role, google_auth_code) VALUES (?, ?, ?, ?, ?)"
					result, err := database.Query(query, addUserData.Login, bcrytedPassword, addUserData.Email, addUserData.Role, googleAuthCode)
					utils.ErrorHandler(err, false)
					defer result.Close()
					if err == nil {
						var zsbEmail string = os.Getenv("EMAIL")
						var zsbEmailPassword string = os.Getenv("EMAIL_PASSWORD")

						authLink := gotp.NewDefaultTOTP(googleAuthCode).ProvisioningUri(addUserData.Login, "ZSBRybnik")
						qrImagePath := "temp/addUser" + addUserData.Login + ".png"
						err = qrcode.WriteFile(authLink, qrcode.Medium, 256, qrImagePath)
						utils.ErrorHandler(err, false)
						templateData := AddUserEmailTemplate{
							Password: password,
							QrCode:   qrImagePath,
						}
						emailTemplate, err := template.ParseFiles("templates/addUser.html")
						utils.ErrorHandler(err, false)
						buffer := new(bytes.Buffer)
						err = emailTemplate.Execute(buffer, templateData)
						utils.ErrorHandler(err, false)
						emailBody := buffer.String()
						email := gomail.NewMessage()
						email.SetHeader("From", zsbEmail)
						email.SetHeader("To", addUserData.Email)
						email.SetHeader("Subject", "Twoje konto ZSB")
						email.Embed(qrImagePath)
						email.SetBody("text/html", emailBody)
						email.Embed("templates/addUserLinks.html")
						planedEmail := gomail.NewPlainDialer("smtp.gmail.com", 587, zsbEmail, zsbEmailPassword)
						err = planedEmail.DialAndSend(email)
						utils.ErrorHandler(err, false)
						err = os.Remove(qrImagePath)
						if err == nil {
							context.Status(200)
						} else {
							context.AbortWithError(500, errors.New("Internal Server Error"))
						}
					} else {
						context.AbortWithError(500, errors.New("Internal Server Error"))
					}
				} else {
					log.Fatalln("Can't find database in gin-gonic context")
					context.AbortWithError(500, errors.New("Internal Server Error"))
				}
			} else {
				context.AbortWithError(400, errors.New("Bad Request"))
			}
		} else {
			context.AbortWithError(401, errors.New("Unauthorized"))
		}
	} else {
		context.AbortWithError(403, errors.New("Forbidden"))
	}
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}
