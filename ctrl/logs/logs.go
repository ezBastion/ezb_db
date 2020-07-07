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

package logs

import (
	"fmt"
	"net/http"
	"strconv"

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

	var Logs []models.EzbLogs
	if err := db.Find(&Logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Logs)
}

func TodayError(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var Logs []models.EzbLogs
	sql := "select * from ezb_logs where status > 399 and strftime('%Y-%m-%d',date) = date('now') order by id desc; "
	if err := db.Table("ezb_logs").Raw(sql).Scan(&Logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Logs)
}

func LastError(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	nb := c.Param("nb")
	if tools.StrIsInt(nb) == false {
		c.JSON(http.StatusConflict, "WRONG PARAMETER")
		return
	}
	nbr, _ := strconv.Atoi(nb)
	var Logs []models.EzbLogs
	sql := fmt.Sprintf("select * from ezb_logs where status > 399 order by id desc limit %d; ", nbr)
	if err := db.Table("ezb_logs").Raw(sql).Scan(&Logs).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, sql)
		return
	}
	c.JSON(http.StatusOK, Logs)
}

func Findone(c *gin.Context) {
	var Log models.EzbLogs
	t := "xtrack"
	if tools.StrIsInt(c.Param("name")) {
		t = "id"
	}
	tools.Findoneraw(c, &Log, t)
}

func Add(c *gin.Context) {
	var Log models.EzbLogs
	tools.Addraw(c, &Log)
}
func Update(c *gin.Context) {
	var Log models.EzbLogs
	tools.Updateraw(c, &Log)
}
