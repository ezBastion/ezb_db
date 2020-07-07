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

import (
	"crypto/sha256"
	"fmt"
	"path"

	m "github.com/ezBastion/ezb_db/models"
	uuid "github.com/gofrs/uuid"

	"github.com/ezBastion/ezb_db/tools"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.WithFields(log.Fields{"module": "gorm", "type": "sql"}).Print(v[3])
	}
	if v[0] == "log" {
		log.WithFields(log.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

/// InitDB create database schema.
func InitDB(conf Configuration, exPath string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	log.Debug("db: ", conf.DB)
	switch DB := conf.DB; DB {
	case "sqlite":
		db, err = gorm.Open("sqlite3", path.Join(exPath, conf.SQLITE.DBPath))

		if err != nil {
			fmt.Printf("sql.Open err: %s\n", err)
			return nil, err
		}
		db.Exec("PRAGMA foreign_keys = OFF")

	default:
		log.Fatal("unknow db type.")
		panic("unknow db type.")
	}
	db.SetLogger(&GormLogger{})

	db.SingularTable(true)
	if !db.HasTable(&m.EzbAccess{}) {
		db.CreateTable(&m.EzbAccess{})
		db.Model(&m.EzbAccess{}).AddUniqueIndex("idx_access_id", "id")
		var Get = m.EzbAccess{}
		db.Where(m.EzbAccess{Name: "GET"}).Attrs(m.EzbAccess{Comment: "Read data", Enable: true}).FirstOrCreate(&Get)
		var Put = m.EzbAccess{}
		db.Where(m.EzbAccess{Name: "PUT"}).Attrs(m.EzbAccess{Comment: "Edit data", Enable: false}).FirstOrCreate(&Put)
		var Post = m.EzbAccess{}
		db.Where(m.EzbAccess{Name: "POST"}).Attrs(m.EzbAccess{Comment: "Add data", Enable: false}).FirstOrCreate(&Post)
		var Delete = m.EzbAccess{}
		db.Where(m.EzbAccess{Name: "DELETE"}).Attrs(m.EzbAccess{Comment: "Remove data", Enable: false}).FirstOrCreate(&Delete)
		var Patch = m.EzbAccess{}
		db.Where(m.EzbAccess{Name: "PATCH"}).Attrs(m.EzbAccess{Comment: "Partial edit data", Enable: false}).FirstOrCreate(&Patch)
	}
	if !db.HasTable(&m.EzbAccounts{}) {
		db.CreateTable(&m.EzbAccounts{})
		db.Model(&m.EzbAccounts{}).AddUniqueIndex("idx_accounts_id", "id")
		var Adm m.EzbAccounts
		salt := tools.RandString(5, "")
		defpwd := fmt.Sprintf("%x", sha256.Sum256([]byte("ezbastion"+salt)))
		db.Where(m.EzbAccounts{Name: "admin"}).Attrs(m.EzbAccounts{Enable: true, Comment: "ezBastion admin", Salt: salt, Password: defpwd, Type: "i", Isadmin: true}).FirstOrCreate(&Adm)
	}
	if !db.HasTable(&m.EzbActions{}) {
		db.CreateTable(&m.EzbActions{})
		db.Model(&m.EzbActions{}).AddUniqueIndex("idx_actions_id", "id")
	}
	if !db.HasTable(&m.EzbCollections{}) {
		db.CreateTable(&m.EzbCollections{})
		db.Model(&m.EzbCollections{}).AddUniqueIndex("idx_collections_id", "id")
	}
	if !db.HasTable(&m.EzbControllers{}) {
		db.CreateTable(&m.EzbControllers{})
		db.Model(&m.EzbControllers{}).AddUniqueIndex("idx_controllers_id", "id")
	}
	if !db.HasTable(&m.EzbGroups{}) {
		db.CreateTable(&m.EzbGroups{})
		db.Model(&m.EzbGroups{}).AddUniqueIndex("idx_groups_id", "id")
	}
	if !db.HasTable(&m.EzbJobs{}) {
		db.CreateTable(&m.EzbJobs{})
		db.Model(&m.EzbJobs{}).AddUniqueIndex("idx_jobs_id", "id")
	}
	if !db.HasTable(&m.EzbTags{}) {
		db.CreateTable(&m.EzbTags{})
		db.Model(&m.EzbTags{}).AddUniqueIndex("idx_tags_id", "id")
	}
	if !db.HasTable(&m.EzbWorkers{}) {
		db.CreateTable(&m.EzbWorkers{})
		db.Model(&m.EzbWorkers{}).AddUniqueIndex("idx_workers_id", "id")
	}
	if !db.HasTable(&m.EzbLogs{}) {
		db.CreateTable(&m.EzbLogs{})
	}
	if !db.HasTable(&m.EzbStas{}) {
		db.CreateTable(&m.EzbStas{})
		db.Model(&m.EzbStas{}).AddUniqueIndex("idx_stas_id", "id")
		firstIAM := m.EzbStas{Name: "Default", Enable: true, Type: 0, Comment: "First IAM", EndPoint: conf.STA, Issuer: "changeME", Default: true}
		db.FirstOrCreate(&firstIAM)
	}
	if !db.HasTable(&m.EzbBastions{}) {
		db.CreateTable(&m.EzbBastions{})
		db.Model(&m.EzbBastions{}).AddUniqueIndex("idx_bastions_id", "id")
		Bastion := m.EzbBastions{Name: "changeme"}
		db.FirstOrCreate(&Bastion)
	}
	if !db.HasTable(&m.EzbLicense{}) {
		db.CreateTable(&m.EzbLicense{})
		serial, err := uuid.NewV4()
		if err != nil {
			serial, _ = uuid.FromString(tools.RandString(32, ""))
		}
		// msg := []byte(fmt.Sprintf("ENT 1970-01-01 2 20 %s", serial))
		// key := []byte(t.EncryptDecrypt(fmt.Sprintf("%s", serial)))
		// mac := hmac.New(sha256.New, key)
		// mac.Write(msg)
		// expectedMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		// Serial := m.EzbLicense{UUID: fmt.Sprintf("%s", serial), Level: "ENT", WKS: 2, API: 20, Sign: expectedMAC, SA: "1970-01-01"}
		Serial := m.EzbLicense{UUID: fmt.Sprintf("%s", serial), Level: "LTE", WKS: 0, API: 0}
		db.FirstOrCreate(&Serial)
	}
	if !db.HasTable(&m.EzbApi{}) {
		var viewApi = `
		CREATE VIEW ezb_api
		AS
		SELECT ezb_actions.name AS action, ezb_actions.comment AS actioncomment, ezb_actions.path, ezb_actions.query, ezb_controllers.name AS ctrl, ezb_controllers.comment AS ctrlcomment, ezb_controllers.version, 
		ezb_access.name AS access, ezb_actions.id, ezb_bastions.fqdn AS bastion, ezb_actions.deprecated, ezb_jobs.name AS job, ezb_jobs.comment AS jobcomment
		FROM ezb_actions INNER JOIN
		ezb_controllers ON ezb_actions.ezb_controllers_id = ezb_controllers.id INNER JOIN
		ezb_access ON ezb_actions.ezb_access_id = ezb_access.id INNER JOIN
		ezb_jobs ON ezb_actions.ezb_jobs_id = ezb_jobs.id CROSS JOIN
		ezb_bastions
		WHERE (ezb_actions.ezb_jobs_id > 0) AND (ezb_actions.ezb_access_id > 0) AND (ezb_actions.ezb_controllers_id > 0) AND (ezb_controllers.enable = 1) AND (ezb_access.enable = 1) AND 
		(ezb_actions.enable = 1) AND (ezb_jobs.enable = 1)		
		`
		db.Exec(viewApi)
	}
	if !db.HasTable(&m.EzbAccountsActions{}) {

		var eaccact = `
		CREATE VIEW ezb_accounts_actions
		AS
		SELECT        acc.name AS account, acc.id AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'account,action' AS path
		FROM            ezb_accounts AS acc INNER JOIN
								ezb_accounts_has_ezb_actions ON acc.id = ezb_accounts_has_ezb_actions.ezb_accounts_id INNER JOIN
								ezb_actions AS act ON ezb_accounts_has_ezb_actions.ezb_actions_id = act.id INNER JOIN
								ezb_controllers AS ctrl ON act.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE        (acc.enable = 1) AND (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (access.enable = 1)
		UNION
		SELECT        acc.name AS account, acc.id AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'account,controller' AS path
		FROM            ezb_accounts AS acc INNER JOIN
								ezb_accounts_has_ezb_controllers ON acc.id = ezb_accounts_has_ezb_controllers.ezb_accounts_id INNER JOIN
								ezb_controllers AS ctrl ON ezb_accounts_has_ezb_controllers.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_actions AS act ON ctrl.id = act.ezb_controllers_id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE        (acc.enable = 1) AND (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (access.enable = 1)
		UNION
		SELECT        acc.name AS account, acc.id AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'group,action' AS path
		FROM            ezb_accounts AS acc INNER JOIN
								ezb_accounts_has_ezb_groups AS acc_grp ON acc.id = acc_grp.ezb_accounts_id INNER JOIN
								ezb_groups AS grp ON acc_grp.ezb_groups_id = grp.id INNER JOIN
								ezb_groups_has_ezb_actions AS grp_act ON grp.id = grp_act.ezb_groups_id INNER JOIN
								ezb_actions AS act ON grp_act.ezb_actions_id = act.id INNER JOIN
								ezb_controllers AS ctrl ON act.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE        (acc.enable = 1) AND (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (grp.enable = 1) AND (access.enable = 1)
		UNION
		SELECT        acc.name AS account, acc.id AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'group,controller' AS path
		FROM            ezb_accounts AS acc INNER JOIN
								ezb_accounts_has_ezb_groups AS acc_grp ON acc.id = acc_grp.ezb_accounts_id INNER JOIN
								ezb_groups AS grp ON acc_grp.ezb_groups_id = grp.id INNER JOIN
								ezb_groups_has_ezb_controllers AS grp_ctrl ON grp.id = grp_ctrl.ezb_groups_id INNER JOIN
								ezb_controllers AS ctrl ON grp_ctrl.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_actions AS act ON act.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE        (acc.enable = 1) AND (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (grp.enable = 1) AND (access.enable = 1)
		UNION
		SELECT        acc.name AS account, acc.id AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'account,collection' AS path
		FROM            ezb_accounts AS acc INNER JOIN
								ezb_accounts_has_ezb_collections AS acc_coll ON acc.id = acc_coll.ezb_accounts_id INNER JOIN
								ezb_collections AS coll ON acc_coll.ezb_collections_id = coll.id INNER JOIN
								ezb_actions_has_ezb_collections AS act_coll ON coll.id = act_coll.ezb_collections_id INNER JOIN
								ezb_actions AS act ON act.id = act_coll.ezb_actions_id INNER JOIN
								ezb_controllers AS ctrl ON act.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE        (acc.enable = 1) AND (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (coll.enable = 1) AND (access.enable = 1)
		UNION
		SELECT        acc.name AS account, acc.id AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'group,collection' AS path
		FROM            ezb_accounts AS acc INNER JOIN
								ezb_accounts_has_ezb_groups AS acc_grp ON acc_grp.ezb_accounts_id = acc.id INNER JOIN
								ezb_groups AS grp ON acc_grp.ezb_groups_id = grp.id INNER JOIN
								ezb_groups_has_ezb_collections AS grp_coll ON grp_coll.ezb_groups_id = grp.id INNER JOIN
								ezb_collections AS coll ON grp_coll.ezb_collections_id = coll.id INNER JOIN
								ezb_actions_has_ezb_collections AS act_coll ON act_coll.ezb_collections_id = coll.id INNER JOIN
								ezb_actions AS act ON act.id = act_coll.ezb_actions_id INNER JOIN
								ezb_controllers AS ctrl ON act.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE        (acc.enable = 1) AND (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (coll.enable = 1) AND (grp.enable = 1) AND (access.enable = 1)
		UNION
		SELECT        'anonymous' AS account, 0 AS accountid, ctrl.name AS ctrl, ctrl.id AS ctrlid, ctrl.version as ctrlver, act.name AS action, act.id AS actionid, job.name AS job, job.id AS jobid, access.name AS access, access.id AS accessid, 'anonymous' AS path
		FROM            ezb_controllers AS ctrl INNER JOIN
								ezb_actions AS act ON act.ezb_controllers_id = ctrl.id INNER JOIN
								ezb_jobs AS job ON act.ezb_jobs_id = job.id INNER JOIN
								ezb_access AS access ON act.ezb_access_id = access.id
		WHERE         (act.enable = 1) AND (job.enable = 1) AND (ctrl.enable = 1) AND (access.enable = 1) AND (act.anonymous = 1);
		`
		db.Exec(eaccact)

	}
	tables := []interface{}{&m.EzbAccess{}, &m.EzbAccounts{}, &m.EzbActions{}, &m.EzbCollections{},
		&m.EzbControllers{}, &m.EzbGroups{}, &m.EzbJobs{}, &m.EzbTags{}, &m.EzbWorkers{}, &m.EzbLogs{},
		&m.EzbStas{}, &m.EzbBastions{}, &m.EzbLicense{}}
	db.AutoMigrate(tables...)

	// sandbox(db)
	// if conf.DB != "sqlite" {
	// 	addFK(db)
	// }

	return db, nil
}

func sandbox(db *gorm.DB) {
	// var Tags = m.EzbTags{Name: "tag1", Comment: "nada"}
	// db.FirstOrCreate(&Tags)
	var Jobs = m.EzbJobs{Name: "job1", Enable: true, Comment: "nada", Checksum: "XXXX", Cache: 10, Path: "/xa/applications/getapp.ps1"}
	db.FirstOrCreate(&Jobs)
	var Groups = m.EzbGroups{Name: "group1", Comment: "nada", Enable: true}
	db.FirstOrCreate(&Groups)
	var Account = m.EzbAccounts{Name: "user1", Type: "internal", Enable: true}
	db.FirstOrCreate(&Account)
	var Collection = m.EzbCollections{Name: "collection1", Enable: true, Comment: "nada"}
	db.FirstOrCreate(&Collection)
	var Controller = m.EzbControllers{Name: "ctrl1", Comment: "nada", Enable: true}
	db.FirstOrCreate(&Controller)
	var Action = m.EzbActions{Name: "action1", Enable: true, Comment: "nada"}
	Action.EzbAccessID = 1
	Action.Groups = []m.EzbGroups{Groups}
	Action.EzbJobsID = Jobs.ID
	Action.EzbControllersID = Controller.ID
	// Action.Tags = []m.EzbTags{Tags}
	Action.Accounts = []m.EzbAccounts{Account}
	db.FirstOrCreate(&Action)
	db.Save(&Groups)
	Account.Groups = []m.EzbGroups{Groups}
	db.Save(&Account)
	Controller.Accounts = []m.EzbAccounts{Account}
	Controller.Groups = []m.EzbGroups{Groups}
	// Controller.Actions = []m.EzbActions{Action}
	db.Save(&Controller)
}
