// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package timeperiodcustomvar

import (
	"github.com/N-o-X/icingadb/configobject"
	"github.com/N-o-X/icingadb/connection"
	"github.com/N-o-X/icingadb/utils"
)

var (
	ObjectInformation configobject.ObjectInformation
	Fields            = []string{
		"id",
		"timeperiod_id",
		"customvar_id",
		"environment_id",
	}
)

type TimeperiodCustomvar struct {
	Id           string `json:"id"`
	TimeperiodId string `json:"object_id"`
	CustomvarId  string `json:"customvar_id"`
	EnvId        string `json:"environment_id"`
}

func NewTimeperiodCustomvar() connection.Row {
	c := TimeperiodCustomvar{}
	return &c
}

func (c *TimeperiodCustomvar) InsertValues() []interface{} {
	v := c.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(c.Id)}, v...)
}

func (c *TimeperiodCustomvar) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(c.TimeperiodId),
		utils.EncodeChecksum(c.CustomvarId),
		utils.EncodeChecksum(c.EnvId),
	)

	return v
}

func (c *TimeperiodCustomvar) GetId() string {
	return c.Id
}

func (c *TimeperiodCustomvar) SetId(id string) {
	c.Id = id
}

func (c *TimeperiodCustomvar) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{c}, nil
}

func init() {
	name := "timeperiod_customvar"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 "timeperiod:customvar",
		PrimaryMySqlField:        "id",
		Factory:                  NewTimeperiodCustomvar,
		HasChecksum:              false,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "timeperiod",
	}
}
