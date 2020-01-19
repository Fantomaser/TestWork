package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

)

var port string = ":8080"
var takeMSGAdr string = "/GetText"
var UserDB []User

type User struct {
	Key     string
	Message string
	Adres   string
}

func main() {

	r := gin.Default()
	r.Use(LiberalCORS)
	r.POST("/makeText", makeText)
	r.POST("/getText", getText)

	go GmailStream()

	r.Run(port)

	fmt.Println("ALL is Ok")
}

func makeText(c *gin.Context) {

	tmp := struct {
		Msg string `json: "msg"`
	}{}

	c.BindJSON(&tmp)

	str := []byte(tmp.Msg)

	if len(str) <= 0 {
		return
	}

	msg, key := Encrypt(str)

	user := User{string(key), string(msg), takeMSGAdr}
	UserDB = append(UserDB, user)

	fmt.Println("Send: ", user)

	c.JSON(200, gin.H{"user": user})
}

func getText(c *gin.Context) {
	tmp := struct {
		Key string `json: "key"`
	}{}

	c.BindJSON(&tmp)

	msg, err := Decrypt([]byte(tmp.Key))

	if err != nil {
		c.JSON(200, gin.H{"status": "deny"})
	}

	c.JSON(200, gin.H{"status": "pass", "message": string(msg)})

}

func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	if c.Request.Method == "OPTIONS" {

		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}

}
