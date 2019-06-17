package configuration

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	t "ezb_priv/tools"

	"github.com/ezBastion/ezb_db/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type License struct {
	WksLimit int
	ApiLimit int
}

func InitLic(lic *License, db *gorm.DB) error {
	logg := log.WithFields(log.Fields{"module": "lic", "type": "log"})
	logg.Debug("start init lic")
	var Lic models.EzbLicense
	if err := db.First(&Lic).Error; err != nil {
		return errors.New("License error L0012")
	} else {
		if strings.Compare(Lic.Level, "LTE") == 0 {
			lic.WksLimit = 0
			lic.ApiLimit = 0
		} else {
			msg := []byte(fmt.Sprintf("%s %s %d %d %s", Lic.Level, Lic.SA, Lic.WKS, Lic.API, Lic.UUID))
			key := []byte(t.EncryptDecrypt(Lic.UUID))
			mac := hmac.New(sha256.New, key)
			mac.Write(msg)
			expectedMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			if Lic.Sign == expectedMAC {
				lic.WksLimit = 0
				lic.ApiLimit = 0
			} else {
				return errors.New("License error L0013")
			}
		}
		return nil
	}
}
