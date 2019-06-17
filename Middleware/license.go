package Middleware

import (
	"github.com/ezBastion/ezb_db/configuration"
	"github.com/gin-gonic/gin"
)

func LicenseMiddleware(lic configuration.License) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("wksLimit", lic.WksLimit)
		c.Set("apiLimit", lic.ApiLimit)
		// log.Printf("wksLimit: %d\n", lic.WksLimit)
		// log.Printf("apiLimit: %d\n", lic.ApiLimit)
		c.Next()
	}
}
