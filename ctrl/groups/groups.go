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

package groups

import (
	"ezb_db/models"
	"ezb_db/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var Groups []models.EzbGroups
	if err := db.
		Preload("Accounts").
		Preload("Actions").
		Preload("Collections").
		Preload("Controllers").
		Order("name asc").Find(&Groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Groups)
}

func Findone(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	name := c.Param("name")
	t := "name"
	if tools.StrIsInt(c.Param("name")) {
		t = "id"
	}
	var Groups models.EzbGroups
	if err := db.
		Preload("Accounts").
		Preload("Actions").
		Preload("Collections").
		Preload("Controllers").
		Where(fmt.Sprintf("%s = ?", t), name).Find(&Groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Groups)

}

func Add(c *gin.Context) {
	var Raw models.EzbGroups
	tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbGroups
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbGroups
	tools.Removeraw(c, &Raw)
}
