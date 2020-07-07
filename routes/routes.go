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

package routes

import (
	"github.com/ezBastion/ezb_db/ctrl/access"
	"github.com/ezBastion/ezb_db/ctrl/accountactions"
	"github.com/ezBastion/ezb_db/ctrl/accounts"
	"github.com/ezBastion/ezb_db/ctrl/actions"
	"github.com/ezBastion/ezb_db/ctrl/api"
	"github.com/ezBastion/ezb_db/ctrl/bastions"
	"github.com/ezBastion/ezb_db/ctrl/collections"
	"github.com/ezBastion/ezb_db/ctrl/controllers"
	"github.com/ezBastion/ezb_db/ctrl/groups"
	"github.com/ezBastion/ezb_db/ctrl/jobs"
	"github.com/ezBastion/ezb_db/ctrl/license"
	"github.com/ezBastion/ezb_db/ctrl/logs"
	"github.com/ezBastion/ezb_db/ctrl/stas"
	"github.com/ezBastion/ezb_db/ctrl/stat"
	"github.com/ezBastion/ezb_db/ctrl/tags"
	"github.com/ezBastion/ezb_db/ctrl/workers"

	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {

	License := route.Group("/license")
	{
		License.GET("", license.Find)
		License.PUT("", license.Update)
	}
	Access := route.Group("/access")
	{
		Access.GET("", access.Find)
		Access.GET("/:name", access.Findone)
		Access.POST("", access.Add)
		Access.DELETE("/:id", access.Remove)
		Access.PUT("", access.Update)
	}
	Accounts := route.Group("/accounts")
	{
		Accounts.GET("", accounts.Find)
		Accounts.GET("/:name", accounts.Findone)
		Accounts.POST("", accounts.Add)
		Accounts.DELETE("/:id", accounts.Remove)
		Accounts.DELETE("/:id/actions/:obj", accounts.UnlinkActions)
		Accounts.DELETE("/:id/groups/:obj", accounts.UnlinkGroups)
		Accounts.DELETE("/:id/controllers/:obj", accounts.UnlinkControllers)
		Accounts.DELETE("/:id/collections/:obj", accounts.UnlinkCollections)
		Accounts.PUT("", accounts.Update)
		Accounts.PUT("/enable", accounts.Enable)
	}
	Actions := route.Group("/actions")
	{
		Actions.GET("", actions.Find)
		Actions.GET("/:name", actions.Findone)
		Actions.POST("", actions.Add)
		Actions.DELETE("/:id", actions.Remove)
		Actions.DELETE("/:id/job", actions.RemoveJob)
		Actions.DELETE("/:id/tag/:obj", actions.RemoveTag)
		Actions.PUT("", actions.Update)
		Actions.PUT("/rename", actions.Rename)
		Actions.PUT("/enable", actions.Enable)
		Actions.PUT("/tag/:obj", actions.AddTag)
	}
	Bastions := route.Group("/bastions")
	{
		Bastions.GET("", bastions.Find)
		Bastions.GET("/:name", bastions.Findone)
		Bastions.POST("", bastions.Add)
		Bastions.DELETE("/:id", bastions.Remove)
		Bastions.PUT("", bastions.Update)
	}
	Collections := route.Group("/collections")
	{
		Collections.GET("", collections.Find)
		Collections.GET("/:name", collections.Findone)
		Collections.POST("", collections.Add)
		Collections.DELETE("/:id", collections.Remove)
		Collections.DELETE("/:id/:obj", collections.Unlink)
		Collections.PUT("", collections.Update)
	}
	Controllers := route.Group("/controllers")
	{
		Controllers.GET("", controllers.Find)
		Controllers.GET("/:name", controllers.Findone)
		Controllers.POST("", controllers.Add)
		Controllers.DELETE("/:id", controllers.Remove)
		Controllers.PUT("", controllers.Update)
		Controllers.PUT("/enable", controllers.Enable)

	}
	Groups := route.Group("/groups")
	{
		Groups.GET("", groups.Find)
		Groups.GET("/:name", groups.Findone)
		Groups.POST("", groups.Add)
		Groups.DELETE("/:id", groups.Remove)
		Groups.PUT("", groups.Update)
	}
	Jobs := route.Group("/jobs")
	{
		Jobs.GET("", jobs.Find)
		Jobs.GET("/xtrack/:name", jobs.Findone)
		Jobs.POST("", jobs.Add)
		Jobs.DELETE("/:id", jobs.Remove)
		Jobs.PUT("", jobs.Update)
		Jobs.PUT("/enable", jobs.Enable)
	}
	Logs := route.Group("/logs")
	{
		// Logs.GET("", logs.Find)
		Logs.GET("/todayerror", logs.TodayError)
		Logs.GET("/xtrack/:name", logs.Findone)
		Logs.GET("/lasterror/:nb", logs.LastError)
		// Logs.GET("/:name", logs.Findone)
		Logs.POST("", logs.Add)
		Logs.PUT("", logs.Update)
	}
	Stas := route.Group("/stas")
	{
		Stas.GET("", stas.Find)
		Stas.GET("/:name", stas.Findone)
		Stas.POST("", stas.Add)
		Stas.DELETE("/:id", stas.Remove)
		Stas.PUT("", stas.Update)
	}
	Tags := route.Group("/tags")
	{
		Tags.GET("", tags.Find)
		Tags.GET("/:name", tags.Findone)
		Tags.POST("", tags.Add)
		Tags.DELETE("/:id", tags.Remove)
		Tags.PUT("", tags.Update)
	}
	Workers := route.Group("/workers")
	{
		Workers.GET("", workers.Find)
		Workers.GET("/:name", workers.Findone)
		Workers.POST("", workers.Add)
		Workers.DELETE("/:id", workers.Remove)
		Workers.DELETE("/:id/tag/:obj", workers.Removetag)
		Workers.PUT("", workers.Update)
		Workers.PUT("/inc/:id", workers.IncRequest)
		Workers.PUT("/tag/:obj", workers.Addtag)
	}
	AccountActions := route.Group("/accountactions")
	{
		AccountActions.GET("", accountactions.Find)
		AccountActions.GET("/:name", accountactions.Findone)
	}
	Api := route.Group("/api")
	{
		Api.GET("", api.Find)
	}

	Stat := route.Group("/stat")
	{
		Stat.GET("/access", stat.GetAccess)
		Stat.GET("/error", stat.GetError)
		Stat.GET("/all/:year/:month", stat.GetMonth)
		Stat.GET("/elm/:elm/:year/:month", stat.GetElm)
	}

}
