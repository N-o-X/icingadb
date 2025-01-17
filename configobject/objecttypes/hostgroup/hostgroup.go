// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package hostgroup

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

type Hostgroup struct {
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

func NewHostgroup() connection.Row {
	h := Hostgroup{}
	h.NameCi = &h.Name

	return &h
}

func (h *Hostgroup) InsertValues() []interface{} {
	v := h.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(h.Id)}, v...)
}

func (h *Hostgroup) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(h.EnvId),
		utils.EncodeChecksum(h.NameChecksum),
		utils.EncodeChecksum(h.PropertiesChecksum),
		utils.EncodeChecksum(h.CustomvarsChecksum),
		h.Name,
		h.NameCi,
		h.DisplayName,
		utils.EncodeChecksum(h.ZoneId),
	)

	return v
}

func (h *Hostgroup) GetId() string {
	return h.Id
}

func (h *Hostgroup) SetId(id string) {
	h.Id = id
}

func (h *Hostgroup) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{h}, nil
}

func init() {
	name := "hostgroup"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 name,
		PrimaryMySqlField:        "id",
		Factory:                  NewHostgroup,
		HasChecksum:              true,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "hostgroup",
	}
}
