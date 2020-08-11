package middleware

import (
	"encoding/json"
	"fmt"
	"gindriver/models"
	"gindriver/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		tmp := strings.Split(authHeader, " ")
		if len(tmp) != 2 {
			c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad authentication token")))
			return
		}
		jwtToken := tmp[1]
		payload := strings.Split(jwtToken, ".")
		if len(payload) != 3 {
			c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad authentication token")))
			return
		}
		jsonPayload := &models.JWTClaims{}
		err := json.Unmarshal(utils.AtobURLSafe(payload[1]), &jsonPayload)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad authentication token")))
			return
		}

		token, err := jwt.ParseWithClaims(jwtToken, &models.JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return utils.ReadFile(fmt.Sprintf("./public/pubkeys/%s.pub", jsonPayload.Name)), nil
			})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("bad credential")))
			return
		}

		if _, err := token.Claims.(*models.JWTClaims); err && token.Valid {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("bad credential")))
			return
		}
	}
}
