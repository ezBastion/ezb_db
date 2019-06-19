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

package configuration

import "github.com/ezbastion/ezb_lib/confmanager"

type Configuration struct {
	ListenJWT       string             `json:"listenjwt"`
	ListenPKI       string             `json:"listenpki"`
	PrivateKey      string             `json:"privatekey"`
	PublicCert      string             `json:"publiccert"`
	CaCert          string             `json:"cacert"`
	DB              string             `json:"db"`
	SQLITE          SQLite             `json:"sqlite"`
	ServiceName     string             `json:"servicename"`
	ServiceFullName string             `json:"servicefullname"`
	EzbPki          string             `json:"ezb_pki"`
	SAN             []string           `json:"san"`
	STA             string             `json:"default_sta"`
	Logger          confmanager.Logger `json:"logger"`
}

type SQL struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Instance string `json:"instance"`
}

type SQLite struct {
	DBPath string `json:"dbpath"`
}
type MYSql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
