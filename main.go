package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Name string
	Age  int
	jwt.StandardClaims
}

var jwt_secret = []byte("reallysecret")

func CreateNewToken() (string, error) {
	c := &Claims{
		Name: "CONSTANT NAME",
		Age:  22,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tokenClaims.SignedString(jwt_secret)
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt_secret, nil
	})

	return tokenClaims.Claims.(*Claims), err
}

func main() {
	r := gin.Default()

	r.GET("/api", func(ctx *gin.Context) {
		remote_token := ctx.GetHeader("remote-token")

		if remote_token == "" {
			new_token, err := CreateNewToken()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"token": new_token,
			})
			return
		}

		// parse token
		claims, err := ParseToken(remote_token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"name": claims.Name,
			"age":  claims.Age,
		})
	})

	r.Run(":3000")
}
