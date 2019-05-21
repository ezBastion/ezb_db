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

package workers

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
	var Workers []models.EzbWorkers
	if err := db.
		Preload("Tags").
		Order("name asc").Find(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Workers)
}
func Findone(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	name := c.Param("name")

	by := "name"
	if tools.StrIsInt(name) {
		by = "id"
	}
	var Workers models.EzbWorkers
	if err := db.
		Preload("Tags").
		Where(fmt.Sprintf("%s = ?", by), name).Find(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Workers)

}
func Add(c *gin.Context) {
	var Raw models.EzbWorkers

	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var nbWorker int
	db.Model(&models.EzbWorkers{}).Count(&nbWorker)
	db.NewRecord(&Raw)
	if err := db.Create(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbWorkers
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbWorkers
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

	if err := db.First(&Raw, id).Association("Tags").Clear().Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Delete(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, Raw)
}

func Removetag(c *gin.Context) {
	var Workers models.EzbWorkers
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

	if err := db.Exec("Delete from ezb_workers_has_ezb_tags WHERE ezb_workers_id = ? AND ezb_tags_id = ?;", id, obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.
		Preload("Tags").
		Where("id = ?", id).Find(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, Workers)
}

func Addtag(c *gin.Context) {
	var Workers models.EzbWorkers
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// workerID := c.Param("id")
	// if tools.StrIsInt(workerID) == false {
	// 	c.JSON(http.StatusConflict, "WRONG PARAMETER")
	// 	return
	// }
	if err := c.BindJSON(&Workers); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	workerID := Workers.ID
	tagID := c.Param("obj")
	if tools.StrIsInt(tagID) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}

	if err := db.Exec("INSERT INTO ezb_workers_has_ezb_tags (ezb_workers_id, ezb_tags_id ) VALUES (?, ?);", workerID, tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Preload("Tags").Where("id = ?", workerID).Find(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Workers)
}

func IncRequest(c *gin.Context) {
	var Workers models.EzbWorkers
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	workerID := c.Param("id")
	if err := db.Where("id = ?", workerID).Find(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Workers.Request++
	db.Save(&Workers)
	c.JSON(http.StatusOK, Workers)
}
