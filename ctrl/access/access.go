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

package access

import (
	"ezb_db/models"
	"ezb_db/tools"

	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	var Raw []models.EzbAccess
	tools.Findraw(c, &Raw, "name asc")
}

func Findone(c *gin.Context) {

	t := "name"
	if tools.StrIsInt(c.Param("name")) {
		t = "id"
	}
	var Raw models.EzbAccess
	tools.Findoneraw(c, &Raw, t)
}
func Add(c *gin.Context) {
	var Raw models.EzbAccess
	tools.Addraw(c, &Raw)
}

func Update(c *gin.Context) {
	var Raw models.EzbAccess
	tools.Updateraw(c, &Raw)
}

func Remove(c *gin.Context) {
	var Raw models.EzbAccess
	tools.Removeraw(c, &Raw)
}
