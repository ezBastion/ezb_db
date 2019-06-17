package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	t "ezb_priv/tools"

	"github.com/ezBastion/ezb_db/models"
	"github.com/ezBastion/ezb_db/tools"
	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var License models.EzbLicense
	if err := db.First(&License).Error; err == nil {
		c.JSON(http.StatusOK, License)
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}

type serialLic struct {
	Lic string `json:"lic"`
	Sig string `json:"sig"`
}

func Update(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var License models.EzbLicense
	if err := db.First(&License).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var slic serialLic
	if err := c.BindJSON(&slic); err != nil {
		c.JSON(http.StatusBadRequest, "BAD LICENSE FORMAT")
		return
	}
	lic, be := base64.StdEncoding.DecodeString(slic.Lic)
	if be != nil {
		c.JSON(http.StatusBadRequest, "BAD LICENSE FORMAT")
		return
	}
	a := strings.Split(string(lic), " ")
	WKS, _ := strconv.Atoi(a[2])
	API, _ := strconv.Atoi(a[3])
	if License.UUID != a[4] {
		c.JSON(http.StatusBadRequest, "BAD LICENSE FORMAT")
		return
	}

	msg := []byte(lic)
	mac := hmac.New(sha256.New, []byte(t.EncryptDecrypt(a[4])))
	mac.Write(msg)
	expectedMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if expectedMAC != slic.Sig {
		c.JSON(http.StatusBadRequest, "BAD LICENSE FORMAT")
		return
	}
	License.Level = a[0]
	License.SA = a[1]
	License.WKS = WKS
	License.API = API
	License.Sign = slic.Sig
	if err := db.Save(&License).Error; err == nil {
		c.JSON(http.StatusOK, License)
		return
	} else {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
}
