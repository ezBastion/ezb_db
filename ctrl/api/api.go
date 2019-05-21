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

package api

import (
	"ezb_db/models"
	"ezb_db/tools"
	"fmt"
	"net/http"
	"regexp"
	s "strings"

	"github.com/gin-gonic/gin"
)

func Find(c *gin.Context) {
	var Raw []models.EzbApi

	db, err := tools.Getdbconn(c)
	if err != "" {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := db.Order("version asc, ctrl").Find(&Raw).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	type api struct {
		ID  int    `json:"id"`
		URL string `json:"api"`
		RGX string `json:"regex"`
	}
	var Ret []api
	for _, r := range Raw {
		var a api
		a.ID = r.ID
		a.URL = fmt.Sprintf("%s %s/v%d/%s/%s", r.Access, r.Bastion, r.Version, r.Ctrl, r.Action)
		if len(r.Path) > 0 {
			a.URL = fmt.Sprintf("%s/%s", a.URL, r.Path)
		}
		if len(r.Query) > 0 {
			a.URL = fmt.Sprintf("%s?%s", a.URL, r.Query)
		}
		a.RGX = fmt.Sprintf("^/v%d/%s/%s", r.Version, r.Ctrl, r.Action)
		if len(r.Path) > 0 {
			var re = regexp.MustCompile(`\{([a-z0-9A-Z-]+)\|([is]{1})\}`)

			subpath := s.Split(r.Path, "/")
			for _, str := range subpath {
				if re.MatchString(str) {
					t := re.FindStringSubmatch(str)[2]
					if t == "i" {
						a.RGX = fmt.Sprintf("%s/([0-9]+)", a.RGX)
					} else if t == "s" {
						a.RGX = fmt.Sprintf("%s/([a-z0-9A-Z-]+)", a.RGX)
					} else {
						// error
					}

				} else {
					a.RGX = fmt.Sprintf("%s/%s", a.RGX, str)
				}
			}
			a.URL = fmt.Sprintf("%s/%s", a.URL, r.Path)
		}
		a.RGX = fmt.Sprintf("%s$", a.RGX)
		Ret = append(Ret, a)
	}
	c.JSON(http.StatusOK, &Ret)

}
