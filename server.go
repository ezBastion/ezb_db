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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/ezBastion/ezb_db/configuration"
	"github.com/ezBastion/ezb_db/routes"
	"github.com/ezBastion/ezb_db/setup"
	"github.com/ezbastion/ezb_lib/logmanager"
	"golang.org/x/sync/errgroup"

	"github.com/ezBastion/ezb_db/Middleware"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

var g errgroup.Group

func routerJWT(db *gorm.DB, lic configuration.License, conf configuration.Configuration) http.Handler {
	loggerJWT := log.WithFields(log.Fields{"module": "jwt", "type": "http"})
	r := gin.Default()
	r.Use(ginrus.Ginrus(loggerJWT, time.RFC3339, true))
	r.Use(Middleware.AddHeaders)
	r.OPTIONS("*a", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})
	r.Use(Middleware.DBMiddleware(db))
	r.Use(Middleware.AuthJWT(db, conf))
	r.Use(Middleware.LicenseMiddleware(lic))
	routes.Routes(r)
	return r
}

func routerPKI(db *gorm.DB, lic configuration.License) http.Handler {
	loggerPKI := log.WithFields(log.Fields{"module": "pki", "type": "http"})
	r := gin.Default()
	r.Use(ginrus.Ginrus(loggerPKI, time.RFC3339, true))
	r.Use(Middleware.AddHeaders)
	r.OPTIONS("*a", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})
	r.Use(Middleware.DBMiddleware(db))
	r.Use(Middleware.LicenseMiddleware(lic))
	routes.Routes(r)
	return r
}
func mainGin(serverchan *chan bool) {
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)
	conf, err := setup.CheckConfig()
	if err != nil {
		panic(err)
	}
	// logmanager.SetLogLevel(conf.Logger.LogLevel, exPath)
	logmanager.SetLogLevel(conf.Logger.LogLevel, exPath, path.Join(exPath, "log/ezb_db.log"), conf.Logger.MaxSize, conf.Logger.MaxBackups, conf.Logger.MaxAge)
	// log.SetOutput(&lumberjack.Logger{
	// 	Filename:   path.Join(exPath, "log/ezb_db.log"),
	// 	MaxSize:    conf.Logger.MaxSize,
	// 	MaxBackups: conf.Logger.MaxBackups,
	// 	MaxAge:     conf.Logger.MaxAge,
	// })

	log.WithFields(log.Fields{"module": "main", "type": "log"})
	log.Debug("loglevel: ", conf.Logger.LogLevel)

	lic := configuration.License{}

	gin.SetMode(gin.ReleaseMode)
	db, err := configuration.InitDB(conf, exPath)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = configuration.InitLic(&lic, db)
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := ioutil.ReadFile(path.Join(exPath, conf.CaCert))
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	/* listner jwt */
	if conf.ListenJWT == "" {
		conf.ListenJWT = "localhost:6000"
	}
	tlsConfigJWT := &tls.Config{}
	serverJWT := &http.Server{
		Addr:      conf.ListenJWT,
		TLSConfig: tlsConfigJWT,
		Handler:   routerJWT(db, lic, conf),
	}
	/* listner jwt */
	/* listner pki */
	if conf.ListenPKI == "" {
		conf.ListenPKI = "localhost:6001"
	}

	tlsConfigPKI := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}
	tlsConfigPKI.BuildNameToCertificate()
	serverPKI := &http.Server{
		Addr:      conf.ListenPKI,
		TLSConfig: tlsConfigPKI,
		Handler:   routerPKI(db, lic),
	}
	/* listner pki */

	g.Go(func() error {
		return serverJWT.ListenAndServeTLS(path.Join(exPath, conf.PublicCert), path.Join(exPath, conf.PrivateKey))
	})

	g.Go(func() error {
		return serverPKI.ListenAndServeTLS(path.Join(exPath, conf.PublicCert), path.Join(exPath, conf.PrivateKey))
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = serverJWT.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
