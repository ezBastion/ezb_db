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

package accounts

import (
	"github.com/ezBastion/ezb_db/models"
	"github.com/ezBastion/ezb_db/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var Account []models.EzbAccounts
	if err := db.
		Preload("Actions").
		Preload("Groups").
		Preload("Controllers").
		Preload("Collections").
		Preload("STA").
		Order("name asc").Find(&Account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Account)
}

func Findone(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var Account models.EzbAccounts
	name := c.Param("name")
	// t := "name"
	if tools.StrIsInt(name) {
		// t = "id"
		if err := db.
			Where("id = ?", name).
			Preload("Actions").
			Preload("Groups").
			Preload("Controllers").
			Preload("Collections").
			Preload("STA").
			Find(&Account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := db.
			Where("name LIKE ? OR real LIKE ?", name, name).
			Preload("Actions").
			Preload("Groups").
			Preload("Controllers").
			Preload("Collections").
			Preload("STA").
			Find(&Account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, Account)
}

func Add(c *gin.Context) {
	var Raw models.EzbAccounts
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	Raw.Salt = tools.RandString(5)
	db.NewRecord(&Raw)
	if err := db.Create(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbAccounts
	tools.Updateraw(c, &Raw)
}

func Enable(c *gin.Context) {
	var Raw models.EzbAccounts
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	tools.Enableraw(c, &Raw, Raw.Enable)
}

func Remove(c *gin.Context) {
	var Raw models.EzbAccounts
	tools.Removeraw(c, &Raw)
}

func UnlinkActions(c *gin.Context) {
	var Account models.EzbAccounts
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

	if err := db.Exec("Delete from ezb_accounts_has_ezb_actions WHERE ezb_accounts_id = ? AND ezb_actions_id = ?;", id, obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.
		Preload("Actions").
		Preload("Groups").
		Preload("Controllers").
		Preload("Collections").
		Where("id = ?", id).Find(&Account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Account)
}

func UnlinkGroups(c *gin.Context) {
	var Account models.EzbAccounts
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

	if err := db.Exec("Delete from ezb_accounts_has_ezb_groups WHERE ezb_accounts_id = ? AND ezb_groups_id = ?;", id, obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.
		Preload("Actions").
		Preload("Groups").
		Preload("Controllers").
		Preload("Collections").
		Where("id = ?", id).Find(&Account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Account)

}
func UnlinkControllers(c *gin.Context) {
	var Account models.EzbAccounts
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

	if err := db.Exec("Delete from ezb_accounts_has_ezb_controllers WHERE ezb_accounts_id = ? AND ezb_controllers_id = ?;", id, obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.
		Preload("Actions").
		Preload("Groups").
		Preload("Controllers").
		Preload("Collections").
		Where("id = ?", id).Find(&Account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Account)

}
func UnlinkCollections(c *gin.Context) {
	var Account models.EzbAccounts
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

	if err := db.Exec("Delete from ezb_accounts_has_ezb_collections WHERE ezb_accounts_id = ? AND ezb_collections_id = ?;", id, obj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.
		Preload("Actions").
		Preload("Groups").
		Preload("Controllers").
		Preload("Collections").
		Where("id = ?", id).Find(&Account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Account)

}
