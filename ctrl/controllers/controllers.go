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

package controllers

import (
	"fmt"
	"net/http"

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

	var Controllers []models.EzbControllers
	if err := db.
		Preload("Accounts").
		// Preload("Actions").
		Preload("Groups").
		Order("name asc").Find(&Controllers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Controllers)
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
	var Controllers models.EzbControllers
	if err := db.
		Preload("Accounts").
		// Preload("Actions").
		Preload("Groups").
		Where(fmt.Sprintf("%s = ?", t), name).Find(&Controllers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Controllers)

}

func Add(c *gin.Context) {
	var Raw models.EzbControllers
	tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbControllers
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbControllers
	tools.Removeraw(c, &Raw)
}

func Enable(c *gin.Context) {
	var Raw models.EzbControllers
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	tools.Enableraw(c, &Raw, Raw.Enable)
}
