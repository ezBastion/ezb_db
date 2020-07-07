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

package tags

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
	var Tags []models.EzbTags
	if err := db.
		// Preload("Workers").
		// Preload("Actions").
		Order("name asc").Find(&Tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Tags)
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
	var Tags models.EzbTags
	if err := db.
		// Preload("Workers").
		// Preload("Actions").
		Where(fmt.Sprintf("%s = ?", t), name).Find(&Tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Tags)

}

func Add(c *gin.Context) {
	var Raw models.EzbTags
	tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbTags
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbTags

	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	tagID := c.Param("id")
	if tools.StrIsInt(tagID) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}

	if err := db.First(&Raw, tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Exec("Delete from ezb_actions_has_ezb_tags WHERE ezb_tags_id = ? ;", tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Exec("Delete from ezb_workers_has_ezb_tags WHERE ezb_tags_id = ? ;", tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Delete(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, Raw)
}
