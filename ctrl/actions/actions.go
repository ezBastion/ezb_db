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

package actions

import (
	"fmt"
	"net/http"

	"github.com/ezBastion/ezb_db/models"
	"github.com/ezBastion/ezb_db/tools"

	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	al, _ := c.MustGet("apiLimit").(int)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var Action []models.EzbActions
	if err := db.
		Preload("Tags").
		Preload("Access").
		Preload("Jobs").
		Preload("Controllers").
		Preload("Accounts").
		Preload("Groups").
		Order("name asc").Find(&Action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if al == 0 || al > len(Action) {
		c.JSON(http.StatusOK, Action)
	} else {
		c.JSON(http.StatusOK, Action[0:al])
	}
}

func Findone(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	name := c.Param("name")
	t := "name"
	if tools.StrIsInt(name) {
		t = "id"
	}
	var Action models.EzbActions
	if err := db.
		Set("gorm:auto_preload", true).
		Where(fmt.Sprintf("%s = ?", t), name).
		Find(&Action).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var Workers []models.EzbWorkers
	var proc = `
	SELECT w.*
	FROM ezb_workers AS w INNER JOIN (
	select wt.ezb_workers_id from ezb_workers_has_ezb_tags AS wt
	where wt.ezb_tags_id IN (
		select t.id from ezb_tags AS t inner join 
		ezb_actions_has_ezb_tags AS at ON at.ezb_tags_id = t.id INNER JOIN 
		ezb_actions a on a.id = at.ezb_actions_id 
		WHERE a.id = ? )
	group by wt.ezb_workers_id
	having COUNT(*) = (
		SELECT COUNT(*) FROM ezb_tags AS t inner join 
		ezb_actions_has_ezb_tags AS at ON at.ezb_tags_id = t.id INNER JOIN 
		ezb_actions AS a ON a.id = at.ezb_actions_id 
		WHERE a.id = ?)
	) AS f ON f.ezb_workers_id = w.id
	`
	if err := db.Raw(proc, Action.ID, Action.ID).Scan(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Action.Workers = Workers
	c.JSON(http.StatusOK, Action)
}

func Add(c *gin.Context) {
	var Raw models.EzbActions
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var nbAction int
	db.Model(&models.EzbActions{}).Count(&nbAction)
	db.NewRecord(&Raw)
	if err := db.Create(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, &Raw)

	// tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbActions
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbActions
	tools.Removeraw(c, &Raw)
}

func Rename(c *gin.Context) {
	var Raw models.EzbActions
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	tools.Renameraw(c, &Raw, Raw.Name)
}

func Enable(c *gin.Context) {
	var Raw models.EzbActions
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	tools.Enableraw(c, &Raw, Raw.Enable)
}

func RemoveJob(c *gin.Context) {
	var Raw models.EzbActions
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

	if err := db.First(&Raw, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Raw.EzbJobsID = 0
	if err := db.Save(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Raw)
}

func RemoveTag(c *gin.Context) {
	var Raw models.EzbActions
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	actionID := c.Param("id")
	if tools.StrIsInt(actionID) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}
	tagID := c.Param("obj")
	if tools.StrIsInt(tagID) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}
	if err := db.Exec("Delete from ezb_actions_has_ezb_tags WHERE ezb_actions_id = ? AND ezb_tags_id = ?;", actionID, tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.
		Set("gorm:auto_preload", true).
		Where("id = ?", actionID).
		Find(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var Workers []models.EzbWorkers
	var proc = `
	SELECT w.*
	FROM ezb_workers AS w INNER JOIN (
	select wt.ezb_workers_id from ezb_workers_has_ezb_tags AS wt
	where wt.ezb_tags_id IN (
		select t.id from ezb_tags AS t inner join 
		ezb_actions_has_ezb_tags AS at ON at.ezb_tags_id = t.id INNER JOIN 
		ezb_actions a on a.id = at.ezb_actions_id 
		WHERE a.id = ? )
	group by wt.ezb_workers_id
	having COUNT(*) = (
		SELECT COUNT(*) FROM ezb_tags AS t inner join 
		ezb_actions_has_ezb_tags AS at ON at.ezb_tags_id = t.id INNER JOIN 
		ezb_actions AS a ON a.id = at.ezb_actions_id 
		WHERE a.id = ?)
	) AS f ON f.ezb_workers_id = w.id
	`
	if err := db.Raw(proc, actionID, actionID).Scan(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Raw.Workers = Workers

	c.JSON(http.StatusOK, Raw)
}

func AddTag(c *gin.Context) {
	var Raw models.EzbActions
	// var Tags []models.EzbTags
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	actionID := Raw.ID
	tagID := c.Param("obj")
	if tools.StrIsInt(tagID) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}

	if err := db.Exec("INSERT INTO ezb_actions_has_ezb_tags (ezb_actions_id, ezb_tags_id ) VALUES (?, ?);", actionID, tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.
		Set("gorm:auto_preload", true).
		Where("id = ?", actionID).
		Find(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var Workers []models.EzbWorkers
	var proc = `
	SELECT w.*
	FROM ezb_workers AS w INNER JOIN (
	select wt.ezb_workers_id from ezb_workers_has_ezb_tags AS wt
	where wt.ezb_tags_id IN (
		select t.id from ezb_tags AS t inner join 
		ezb_actions_has_ezb_tags AS at ON at.ezb_tags_id = t.id INNER JOIN 
		ezb_actions a on a.id = at.ezb_actions_id 
		WHERE a.id = ? )
	group by wt.ezb_workers_id
	having COUNT(*) = (
		SELECT COUNT(*) FROM ezb_tags AS t inner join 
		ezb_actions_has_ezb_tags AS at ON at.ezb_tags_id = t.id INNER JOIN 
		ezb_actions AS a ON a.id = at.ezb_actions_id 
		WHERE a.id = ?)
	) AS f ON f.ezb_workers_id = w.id
	`
	if err := db.Raw(proc, actionID, actionID).Scan(&Workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Raw.Workers = Workers

	c.JSON(http.StatusOK, Raw)
}
