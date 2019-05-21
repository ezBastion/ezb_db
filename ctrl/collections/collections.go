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

package collections

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

	var Collections []models.EzbCollections
	if err := db.
		Preload("Accounts").
		Preload("Actions").
		Preload("Groups").
		Order("name asc").Find(&Collections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Collections)
}
func Findone(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	name := c.Param("name")

	var Collections models.EzbCollections
	t := "name"
	if tools.StrIsInt(c.Param("name")) {
		t = "id"
	}
	if err := db.
		Preload("Accounts").
		Preload("Actions").
		Preload("Groups").
		Where(fmt.Sprintf("%s = ?", t), name).Find(&Collections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Collections)

}

func Add(c *gin.Context) {
	var Raw models.EzbCollections
	tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbCollections
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbCollections
	tools.Removeraw(c, &Raw)
}

func Unlink(c *gin.Context) {
	var Collections models.EzbCollections
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	id := c.Param("id")
	if tools.StrIsInt(id) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}
	obj := c.Param("obj")
	if tools.StrIsInt(id) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}

	if err := db.Exec("Delete from ezb_actions_has_ezb_collections WHERE ezb_collections_id = ? AND ezb_actions_id = ?;", id, obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.
		Preload("Accounts").
		Preload("Actions").
		Preload("Groups").
		Where("id = ?", id).Find(&Collections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Collections)
}
