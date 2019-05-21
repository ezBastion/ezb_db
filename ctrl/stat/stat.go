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

package stat

import (
	"bytes"
	"ezb_db/tools"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type History struct {
	Year  int
	Month int
	Count int
}

type ElmAccess struct {
	Elm    string
	Access int
}

type AllAccess struct {
	Status     string `json:"status"`
	Token      string `json:"token"`
	Controller string `json:"controller"`
	Action     string `json:"action"`
	Bastion    string `json:"bastion"`
	Worker     string `json:"worker"`
	Issuer     string `json:"issuer"`
	Methode    string `json:"methode"`
	// Duration   int64     `json:"duration"`
	// Size       int       `json:"size"`
}

func GetAccess(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var Access []History
	var buffer bytes.Buffer
	buffer.WriteString("SELECT strftime('%Y',date) year , strftime('%m',date) month , count(*) count ")
	buffer.WriteString("FROM ezb_logs ")
	buffer.WriteString("WHERE controller != 'internal' ")
	buffer.WriteString("GROUP BY strftime('%Y',date),  strftime('%m',date)  ")
	buffer.WriteString("ORDER BY strftime('%Y',date) DESC,  strftime('%m',date) DESC ")
	buffer.WriteString("limit 12; ")
	// buffer.WriteString("SELECT TOP 12 YEAR([date]) year, MONTH([date]) month, COUNT(*) count ")
	// buffer.WriteString("FROM [ezb_logs] ")
	// buffer.WriteString("GROUP BY MONTH([date]), YEAR([date]) ")
	// buffer.WriteString("ORDER BY YEAR([date]) DESC, MONTH([date]) DESC;")
	if err := db.Table("ezb_logs").Raw(buffer.String()).Scan(&Access).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Access)
}

func GetError(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var Access []History
	var buffer bytes.Buffer
	buffer.WriteString("SELECT strftime('%Y',date) year , strftime('%m',date) month , count(*) count ")
	buffer.WriteString("FROM ezb_logs ")
	buffer.WriteString("WHERE controller != 'internal' AND  status > 399  ")
	buffer.WriteString("GROUP BY strftime('%Y',date),  strftime('%m',date)  ")
	buffer.WriteString("ORDER BY strftime('%Y',date) DESC,  strftime('%m',date) DESC ")
	buffer.WriteString("limit 12; ")
	// buffer.WriteString("SELECT TOP 12 YEAR([date]) year, MONTH([date]) month, COUNT(*) count ")
	// buffer.WriteString("FROM [view_ERROR] ")
	// buffer.WriteString("GROUP BY MONTH([date]), YEAR([date]) ")
	// buffer.WriteString("ORDER BY YEAR([date]) DESC, MONTH([date]) DESC;")
	if err := db.Table("ezb_logs").Raw(buffer.String()).Scan(&Access).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Access)
}

func GetElm(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	elm := c.Param("elm")
	year := c.Param("year")
	month := c.Param("month")
	var Access []ElmAccess
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("SELECT %s elm, COUNT(*) access ", elm))
	buffer.WriteString("FROM ezb_logs ")
	buffer.WriteString(fmt.Sprintf("where (strftime('%%Y',date) = '%s') and (strftime('%%m',date) = '%s') and %s != '' ", year, month, elm))
	buffer.WriteString(fmt.Sprintf("GROUP BY %s ", elm))
	buffer.WriteString("ORDER BY access DESC limit 10;")
	if err := db.Table("ezb_logs").Raw(buffer.String()).Scan(&Access).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Access)
}
func GetMonth(c *gin.Context) {
	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	year := c.Param("year")
	month := c.Param("month")
	var Access []AllAccess
	var buffer bytes.Buffer
	buffer.WriteString("SELECT status, token, controller, action, bastion, worker, issuer, methode ")
	buffer.WriteString("FROM ezb_logs ")
	buffer.WriteString(fmt.Sprintf("where (strftime('%%Y',date) = '%s') and (strftime('%%m',date) = '%s') and controller != 'internal' ", year, month))
	buffer.WriteString("ORDER BY id DESC ;")
	if err := db.Table("ezb_logs").Raw(buffer.String()).Scan(&Access).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Access)
}
