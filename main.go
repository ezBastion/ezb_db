//go:generate  goversioninfo -64 -platform-specific=false

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
	"fmt"
	"log"
	"os"

	"github.com/ezBastion/ezb_db/configuration"
	"github.com/ezBastion/ezb_db/setup"

	"github.com/urfave/cli"
	"golang.org/x/sys/windows/svc"
)

var (
	exPath string
	conf   configuration.Configuration
)

func main() {

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}
	if !isIntSess {
		conf, err := setup.CheckConfig()
		if err == nil {
			runService(conf.ServiceName, false)
		}
		return
	}
	app := cli.NewApp()
	app.Name = "ezb_db"
	app.Version = "0.3.0"
	app.Usage = "Manage ezBastion database."

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Genarate config file and PKI certificat.",
			Action: func(c *cli.Context) error {
				err := setup.Setup(true)
				return err
			},
		}, {
			Name:  "debug",
			Usage: "Start ezb_db in console.",
			Action: func(c *cli.Context) error {
				conf, _ := setup.CheckConfig()
				runService(conf.ServiceName, true)
				return nil
			},
		}, {
			Name:  "install",
			Usage: "Add ezb_db deamon windows service.",
			Action: func(c *cli.Context) error {
				conf, _ := setup.CheckConfig()
				return installService(conf.ServiceName, conf.ServiceFullName)
			},
		}, {
			Name:  "remove",
			Usage: "Remove ezb_db deamon windows service.",
			Action: func(c *cli.Context) error {
				conf, _ := setup.CheckConfig()
				return removeService(conf.ServiceName)
			},
		}, {
			Name:  "start",
			Usage: "Start ezb_db deamon windows service.",
			Action: func(c *cli.Context) error {
				conf, _ := setup.CheckConfig()
				return startService(conf.ServiceName)
			},
		}, {
			Name:  "stop",
			Usage: "Stop ezb_db deamon windows service.",
			Action: func(c *cli.Context) error {
				conf, _ := setup.CheckConfig()
				return controlService(conf.ServiceName, svc.Stop, svc.Stopped)
			},
		},
		{
			Name:  "newadmin",
			Usage: "Add an admin account.",
			Action: func(c *cli.Context) error {
				err := setup.ResetPWD()
				return err
			},
		},
		{
			Name:  "backup",
			Usage: "Dump db in file.",
			Action: func(c *cli.Context) error {
				err := setup.DumpDB()
				return err
			},
		},
		{
			Name:  "restore",
			Usage: "Restore db from file.",
			Action: func(c *cli.Context) error {
				err := setup.RestoreDB()
				return err
			},
		},
	}

	cli.AppHelpTemplate = fmt.Sprintf(`

	███████╗███████╗██████╗  █████╗ ███████╗████████╗██╗ ██████╗ ███╗   ██╗
	██╔════╝╚══███╔╝██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██║██╔═══██╗████╗  ██║
	█████╗    ███╔╝ ██████╔╝███████║███████╗   ██║   ██║██║   ██║██╔██╗ ██║
	██╔══╝   ███╔╝  ██╔══██╗██╔══██║╚════██║   ██║   ██║██║   ██║██║╚██╗██║
	███████╗███████╗██████╔╝██║  ██║███████║   ██║   ██║╚██████╔╝██║ ╚████║
	╚══════╝╚══════╝╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝
																		   
							██████╗ ██████╗                                
							██╔══██╗██╔══██╗                               
							██║  ██║██████╔╝                               
							██║  ██║██╔══██╗                               
							██████╔╝██████╔╝                               
							╚═════╝ ╚═════╝               

%s
INFO:
		http://www.ezbastion.com		
		support@ezbastion.com
		`, cli.AppHelpTemplate)
	app.Run(os.Args)
}
