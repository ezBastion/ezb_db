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

package jobs

import (
	"ezb_db/models"
	"ezb_db/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	var Raw []models.EzbJobs
	tools.Findraw(c, &Raw, "name asc")
}

func Findone(c *gin.Context) {
	t := "name"
	if tools.StrIsInt(c.Param("name")) {
		t = "id"
	}
	var Raw models.EzbJobs
	tools.Findoneraw(c, &Raw, t)
}
func Add(c *gin.Context) {
	var Raw models.EzbJobs
	tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbJobs
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbJobs
	tools.Removeraw(c, &Raw)
}

func Enable(c *gin.Context) {
	var Raw models.EzbJobs
	if err := c.BindJSON(&Raw); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	tools.Enableraw(c, &Raw, Raw.Enable)
}
