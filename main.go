package main

import (
	"fmt"
	"math/rand"
	//"log"
	"net/http"

	"github.com/gin-gonic/gin"

)

var port string = ":9090"
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

	btMask := make([]byte, 8)

	for i, _ := range btMask {
		btMask[i] = byte(rand.Int63() % 126)
	}

	for i, _ := range []byte(str) {
		str[i] = str[i] ^ btMask[i%len(btMask)]
	}

	user := User{string(btMask), string(str), "/GetText"}
	userDB = append(userDB, user)

	fmt.Println("Send: ", user)

	c.JSON(200, gin.H{"user": user})
}

func getText(c *gin.Context) {

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
