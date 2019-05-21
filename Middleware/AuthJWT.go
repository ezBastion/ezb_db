// This file is part of ezBastion.

//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

package Middleware

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ezBastion/ezb_db/configuration"
	"github.com/ezBastion/ezb_db/models"

	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Payload struct {
	JTI string `json:"jti"`
	ISS string `json:"iss"`
	SUB string `json:"sub"`
	AUD string `json:"aud"`
	EXP int    `json:"exp"`
	IAT int    `json:"iat"`
}

func AuthJWT(db *gorm.DB, conf configuration.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {

		logg := log.WithFields(log.Fields{"Middleware": "jwt"})
		var err error
		authHead := c.GetHeader("Authorization")
		bearer := strings.Split(authHead, " ")
		if len(bearer) != 2 {
			logg.Error("bad Authorization #J0001: " + authHead)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0001"))
			return
		}
		if strings.Compare(strings.ToLower(bearer[0]), "bearer") != 0 {
			logg.Error("bad Authorization #J0002: " + authHead)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0002"))
			return
		}
		tokenString := bearer[1]
		ex, _ := os.Executable()
		exPath := filepath.Dir(ex)
		parts := strings.Split(tokenString, ".")
		p, err := base64.RawStdEncoding.DecodeString(parts[1])
		if err != nil {
			logg.Error("Unable to decode payload: ", err)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0009"))
			return
		}
		var payload Payload
		err = json.Unmarshal(p, &payload)
		if err != nil {
			logg.Error("Unable to parse payload: ", err)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0011"))
			return
		}
		jwtkeyfile := fmt.Sprintf("%s.crt", payload.ISS)
		jwtpubkey := path.Join(exPath, "cert", jwtkeyfile)

		if _, err := os.Stat(jwtpubkey); os.IsNotExist(err) {
			logg.Error("Unable to load sta public certificat: ", err)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0010"))
			return
		}

		key, _ := ioutil.ReadFile(jwtpubkey)
		var ecdsaKey *ecdsa.PublicKey
		if ecdsaKey, err = jwt.ParseECPublicKeyFromPEM(key); err != nil {
			logg.Error("Unable to parse ECDSA public key: ", err)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0003"))
		}
		methode := jwt.GetSigningMethod("ES256")
		// parts := strings.Split(tokenString, ".")
		err = methode.Verify(strings.Join(parts[0:2], "."), parts[2], ecdsaKey)
		if err != nil {
			logg.Error("Error while verifying key: ", err)
			c.AbortWithError(http.StatusForbidden, errors.New("#J0004"))
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwt.ParseECPublicKeyFromPEM(key)
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// log.Println(claims["iss"], claims["sub"])
			var Account models.EzbAccounts
			if err := db.Where("name LIKE ? OR real LIKE ?", claims["sub"], claims["sub"]).Find(&Account).Error; err != nil {
				c.AbortWithError(http.StatusForbidden, errors.New("#J0006"))
				return
			}
			if !Account.Isadmin {
				c.AbortWithError(http.StatusForbidden, errors.New("#J0007"))
				return
			}
		} else {
			c.AbortWithError(http.StatusForbidden, errors.New("#J0005"))
			logg.Error(err)
			return
		}
		c.Next()
	}
}
