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

package models

import (
	"time"
)

type EzbLicense struct {
	ID    int    `json:"id" gorm:"primary_key"`
	UUID  string `json:"uuid"`
	Level string `json:"level"`
	WKS   int    `json:"workers"`
	API   int    `json:"api"`
	SA    string `json:"saexpiry"`
	Sign  string `json:"sign"`
}
type EzbActions struct {
	ID               int              `json:"id" gorm:"primary_key"`
	Name             string           `gorm:"size:250;not null" json:"name"`
	Enable           bool             `gorm:"not null;default:'0'" json:"enable"  `
	Comment          string           `json:"comment"`
	LastRequest      time.Time        `json:"lastrequest"`
	Tags             []*EzbTags       `json:"tags" gorm:"many2many:ezb_actions_has_ezb_tags;"`
	Access           EzbAccess        `json:"access" gorm:"ForeignKey:ID;AssociationForeignKey:EzbAccessID;association_autoupdate:false;association_autocreate:false;association_save_reference:false;"`
	EzbAccessID      int              `gorm:"default:'0'" json:"ezbaccessid" sql:"type:int REFERENCES ezb_access(id)"`
	Jobs             EzbJobs          `json:"jobs" gorm:"ForeignKey:ID;AssociationForeignKey:EzbJobsID;association_autoupdate:false;association_autocreate:false;association_save_reference:false;"`
	EzbJobsID        int              `json:"ezbjobsid" sql:"type:int REFERENCES ezb_jobs(id)"`
	Collections      []EzbCollections `json:"collections" gorm:"many2many:ezb_actions_has_ezb_collections;"`
	Accounts         []EzbAccounts    `json:"accounts" gorm:"many2many:ezb_accounts_has_ezb_actions;"`
	Groups           []EzbGroups      `json:"groups" gorm:"many2many:ezb_groups_has_ezb_actions;"`
	Controllers      EzbControllers   `json:"controllers" gorm:"ForeignKey:ID;AssociationForeignKey:EzbControllersID;association_autoupdate:false;association_autocreate:false;association_save_reference:false;"`
	EzbControllersID int              `gorm:"default:'0'" json:"ezbcontrollersid" sql:"type:int REFERENCES ezb_controllers(id)"`
	Workers          []EzbWorkers     `json:"workers" gorm:"-"`
	Path             string           `json:"path"`
	Query            string           `json:"query" sql:"type:text"`
	Body             string           `json:"body"`
	Constant         string           `json:"constant" sql:"type:text"`
	Deprecated       int              `json:"deprecated" gorm:"not null; default:'0'"` //false or new action id
	Anonymous        bool             `gorm:"not null;default:'0'" json:"anonymous"  `
}

type EzbCollections struct {
	ID       int           `json:"id" gorm:"primary_key"`
	Name     string        `gorm:"size:250;not null;unique" json:"name"`
	Enable   bool          `gorm:"not null;default:'0'" json:"enable"  `
	Comment  string        `json:"comment"`
	Actions  []EzbActions  `json:"actions" gorm:"many2many:ezb_actions_has_ezb_collections;"`
	Accounts []EzbAccounts `json:"accounts" gorm:"many2many:ezb_accounts_has_ezb_collections;"`
	Groups   []EzbGroups   `json:"groups" gorm:"many2many:ezb_groups_has_ezb_collections;"`
}

type EzbAccess struct {
	ID      int    `json:"id" gorm:"primary_key"`
	Name    string `gorm:"size:250;not null;unique" json:"name"`
	Enable  bool   `gorm:"not null;default:'0'" json:"enable"  `
	Comment string `json:"comment"`
}

type EzbAccounts struct {
	ID          int              `json:"id" gorm:"primary_key"`
	Name        string           `gorm:"size:250;not null;unique;index:name" json:"name"`
	Enable      bool             `gorm:"not null;default:'0'" json:"enable"`
	Isadmin     bool             `gorm:"not null;default:'0'" json:"isadmin"`
	Comment     string           `json:"comment"`
	LastRequest time.Time        `json:"lastrequest"`
	Type        string           `json:"type"`
	Real        string           `json:"real"`
	Email       string           `json:"email"`
	Password    string           `json:"password"`
	Salt        string           `json:"salt"`
	Actions     []EzbActions     `json:"actions" gorm:"many2many:ezb_accounts_has_ezb_actions;"`
	Groups      []EzbGroups      `json:"groups" gorm:"many2many:ezb_accounts_has_ezb_groups;"`
	Controllers []EzbControllers `json:"controllers" gorm:"many2many:ezb_accounts_has_ezb_controllers;"`
	Collections []EzbCollections `json:"collections" gorm:"many2many:ezb_accounts_has_ezb_collections;"`
	STA         EzbStas          `json:"sta" gorm:"ForeignKey:ID;AssociationForeignKey:EzbStasID;association_autoupdate:false;association_autocreate:false;association_save_reference:false;"`
	EzbStasID   int              `gorm:"default:'1'" json:"stasid" sql:"type:int REFERENCES ezb_stas(id)"`
}
type EzbGroups struct {
	ID          int              `json:"id" gorm:"primary_key"`
	Name        string           `gorm:"size:250;not null;unique" json:"name"`
	Enable      bool             `gorm:"not null;default:'0'" json:"enable"  `
	Comment     string           `json:"comment"`
	Accounts    []EzbAccounts    `json:"accounts" gorm:"many2many:ezb_accounts_has_ezb_groups;"`
	Actions     []EzbActions     `json:"actions" gorm:"many2many:ezb_groups_has_ezb_actions;"`
	Collections []EzbCollections `json:"collections" gorm:"many2many:ezb_groups_has_ezb_collections;"`
	Controllers []EzbControllers `json:"controllers" gorm:"many2many:ezb_groups_has_ezb_controllers;"`
}

type EzbControllers struct {
	ID       int           `json:"id" gorm:"primary_key"`
	Name     string        `gorm:"size:250;not null;" json:"name"`
	Enable   bool          `gorm:"not null;default:'0'" json:"enable"  `
	Comment  string        `json:"comment"`
	Accounts []EzbAccounts `gorm:"many2many:ezb_accounts_has_ezb_controllers;"`
	Groups   []EzbGroups   `gorm:"many2many:ezb_groups_has_ezb_controllers;"`
	Version  int           `json:"version" gorm:"not null; default:'1'" sql:"type:int" `
	// Actions  []EzbActions  `gorm:"ForeignKey:ID"`
}

type EzbJobs struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Name     string `gorm:"size:250;not null;unique" json:"name"`
	Enable   bool   `gorm:"not null;default:'0'" json:"enable"  `
	Comment  string `json:"comment"`
	Checksum string `json:"checksum"`
	Path     string `json:"path"`
	Cache    int    `json:"cache"`
	Async    bool   `json:"async"`
	Output   string `json:"output" sql:"type:text"`
}
type EzbWorkers struct {
	ID          int        `json:"id" gorm:"primary_key"`
	Name        string     `gorm:"size:250;not null;unique" json:"name"`
	Enable      bool       `gorm:"not null;default:'0'" json:"enable"`
	Comment     string     `json:"comment"`
	Tags        []*EzbTags `json:"tags" gorm:"many2many:ezb_workers_has_ezb_tags;"`
	Fqdn        string     `gorm:"size:250;" json:"fqdn"`
	Register    time.Time  `json:"register"`
	LastRequest time.Time  `json:"lastrequest"`
	Request     int        `gorm:"not null;default:'0'" json:"request"`
}
type EzbTags struct {
	ID      int           `json:"id" gorm:"primary_key"`
	Name    string        `gorm:"size:250;not null;unique" json:"name"`
	Comment string        `json:"comment"`
	Workers []*EzbWorkers `json:"workers" gorm:"many2many:ezb_workers_has_ezb_tags;"`
	Actions []*EzbActions `json:"actions" gorm:"many2many:ezb_actions_has_ezb_tags;"`
}
type EzbBastions struct {
	ID      int    `json:"id" gorm:"primary_key"`
	Name    string `gorm:"size:250;not null;unique" json:"name"`
	Enable  bool   `gorm:"not null;default:'0'" json:"enable"`
	Fqdn    string `gorm:"not null" json:"fqdn"`
	Comment string `json:"comment"`
}
type EzbStas struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Name     string `gorm:"size:250;not null;unique" json:"name"`
	Enable   bool   `gorm:"not null;default:'0'" json:"enable"`
	Type     int    `gorm:"not null;default:'0'" json:"type"` // 0:internal  1:AD  2:oAuth2
	Comment  string `json:"comment"`
	EndPoint string `json:"authorization_endpoint"`
	Issuer   string `json:"issuer"`
	Default  bool   `gorm:"not null;default:'0'" json:"default"`
}
type EzbLogs struct {
	ID         int       `json:"id" gorm:"primary_key"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`
	Ipaddr     string    `json:"ipaddr"`
	Token      string    `json:"token"`
	Account    string    `json:"account"`
	Controller string    `json:"controller"`
	Action     string    `json:"action"`
	URL        string    `json:"url"`
	Bastion    string    `json:"bastion"`
	Worker     string    `json:"worker"`
	Xtrack     string    `json:"xtrack"`
	Deprecated bool      `json:"deprecated"`
	Issuer     string    `json:"issuer"`
	Methode    string    `json:"methode"`
	Duration   int64     `json:"duration"`
	Size       int       `json:"size"`
	Error      string    `json:"error"`
}

type EzbAccountsActions struct {
	Account   string `json:"account"`
	Accountid int    `json:"accountid"`
	Ctrl      string `json:"ctrl"`
	Ctrlid    int    `json:"ctrlid"`
	Ctrlver   int    `json:"ctrlver"`
	Action    string `json:"action"`
	Actionid  int    `json:"actionid"`
	Job       string `json:"job"`
	Jobid     int    `json:"jobid"`
	Access    string `json:"access"`
	Accessid  int    `json:"accessid"`
	Path      string `json:"path"`
}

type EzbApi struct {
	ID            int    `json:"id"`
	Access        string `json:"access"`
	Ctrl          string `json:"ctrl"`
	Ctrlcomment   string `json:"ctrlcomment"`
	Version       int    `json:"version"`
	Action        string `json:"action"`
	Actioncomment string `json:"actioncomment"`
	Path          string `json:"path"`
	Query         string `json:"query"`
	Bastion       string `json:"bastion"`
	Deprecated    int    `json:"deprecated"`
	Job           string `json:"job"`
	Jobcomment    string `json:"jobcomment"`
}
