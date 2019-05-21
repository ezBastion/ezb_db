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

package tools

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Getdbconn(c *gin.Context) (db *gorm.DB, ret string) {

	db, _ = c.MustGet("db").(*gorm.DB)
	if db == nil {
		dberrmsg, ok := c.MustGet("dberr").(error)
		if ok {
			ret = dberrmsg.Error()
		} else {
			ret = string("unknow database connection error")
		}
		return nil, ret
	}
	return db, ""
}
func StrIsInt(data string) bool {
	match, _ := regexp.MatchString("^([0-9]+)$", data)
	if match {
		return true
	} else {
		return false
	}
}

func Findraw(c *gin.Context, d interface{}, order string) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := db.Order(order).Find(d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, d)
}

func Findoneraw(c *gin.Context, d interface{}, by string) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	name := c.Param("name")
	if err := db.Where(fmt.Sprintf("%s = ?", by), name).Find(d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, d)
}

func Addraw(c *gin.Context, d interface{}) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.BindJSON(d); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	db.NewRecord(d)
	if err := db.Create(d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, d)
}

func Updateraw(c *gin.Context, d interface{}) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.BindJSON(d); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Save(d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, d)
}

func Renameraw(c *gin.Context, d interface{}, newname string) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := db.Model(d).UpdateColumn("name", newname).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, d)
}

func Enableraw(c *gin.Context, d interface{}, enable bool) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := db.Model(d).UpdateColumn("enable", enable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, d)
}

func Removeraw(c *gin.Context, d interface{}) {
	db, err := Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	id := c.Param("id")
	if StrIsInt(id) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}

	if err := db.First(d, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Delete(d).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, d)
}
