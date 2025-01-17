// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package usergroup

import (
	"github.com/N-o-X/icingadb/configobject"
	"github.com/N-o-X/icingadb/connection"
	"github.com/N-o-X/icingadb/utils"
)

var (
	ObjectInformation configobject.ObjectInformation
	Fields            = []string{
		"id",
		"environment_id",
		"name_checksum",
		"properties_checksum",
		"customvars_checksum",
		"name",
		"name_ci",
		"display_name",
		"zone_id",
	}
)

type Usergroup struct {
	Id                 string  `json:"id"`
	EnvId              string  `json:"environment_id"`
	NameChecksum       string  `json:"name_checksum"`
	PropertiesChecksum string  `json:"checksum"`
	CustomvarsChecksum string  `json:"customvars_checksum"`
	Name               string  `json:"name"`
	NameCi             *string `json:"name_ci"`
	DisplayName        string  `json:"display_name"`
	ZoneId             string  `json:"zone_id"`
}

func NewUsergroup() connection.Row {
	u := Usergroup{}
	u.NameCi = &u.Name

	return &u
}

func (u *Usergroup) InsertValues() []interface{} {
	v := u.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(u.Id)}, v...)
}

func (u *Usergroup) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(u.EnvId),
		utils.EncodeChecksum(u.NameChecksum),
		utils.EncodeChecksum(u.PropertiesChecksum),
		utils.EncodeChecksum(u.CustomvarsChecksum),
		u.Name,
		u.NameCi,
		u.DisplayName,
		utils.EncodeChecksum(u.ZoneId),
	)

	return v
}

func (u *Usergroup) GetId() string {
	return u.Id
}

func (u *Usergroup) SetId(id string) {
	u.Id = id
}

func (u *Usergroup) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{u}, nil
}

func init() {
	name := "usergroup"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 name,
		PrimaryMySqlField:        "id",
		Factory:                  NewUsergroup,
		HasChecksum:              true,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "usergroup",
	}
}
