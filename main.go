package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

)

var port string = ":9090"
var takeMSGAdr string = "/GetText"
var userDB []User

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
	r.Run(":8080")

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

	btMask := criptKey()

	for it := 0; len(userDB) != 0; it++ {
		if userDB[it].Key == string(btMask) {
			btMask = criptKey()
			it = 0
		}

		if it == (len(userDB) - 1) {
			break
		}

	}

	for i, _ := range []byte(str) {
		str[i] = str[i] ^ btMask[i%len(btMask)]
	}

	user := User{string(btMask), string(str), takeMSGAdr}
	userDB = append(userDB, user)

	fmt.Println("Send: ", user)

	c.JSON(200, gin.H{"user": user})
}

func getText(c *gin.Context) {
	tmp := struct {
		Key string `json: "key"`
	}{}

	c.BindJSON(&tmp)

	var Userid int = -1

	for i, _ := range userDB {
		if userDB[i].Key == tmp.Key {
			Userid = i
			break
		}
	}

	if Userid == -1 {
		c.JSON(200, gin.H{"status": "deny"})
	}

	str := []byte(userDB[Userid].Message)

	for i, _ := range str {
		str[i] = str[i] ^ tmp.Key[i%len(tmp.Key)]
	}

	c.JSON(200, gin.H{"status": "pass", "message": string(str)})

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

func criptKey() []byte {

	btMask := make([]byte, 8)

	for i, _ := range btMask {
		btMask[i] = byte(rand.Int63() % 126)
	}

	return btMask
}
