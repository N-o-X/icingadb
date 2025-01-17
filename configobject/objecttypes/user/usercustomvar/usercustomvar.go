// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package usercustomvar

import (
	"github.com/N-o-X/icingadb/configobject"
	"github.com/N-o-X/icingadb/connection"
	"github.com/N-o-X/icingadb/utils"
)

var (
	ObjectInformation configobject.ObjectInformation
	Fields            = []string{
		"id",
		"user_id",
		"customvar_id",
		"environment_id",
	}
)

type UserCustomvar struct {
	Id          string `json:"id"`
	UserId      string `json:"object_id"`
	CustomvarId string `json:"customvar_id"`
	EnvId       string `json:"environment_id"`
}

func NewUserCustomvar() connection.Row {
	c := UserCustomvar{}
	return &c
}

func (c *UserCustomvar) InsertValues() []interface{} {
	v := c.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(c.Id)}, v...)
}

func (c *UserCustomvar) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(c.UserId),
		utils.EncodeChecksum(c.CustomvarId),
		utils.EncodeChecksum(c.EnvId),
	)

	return v
}

func (c *UserCustomvar) GetId() string {
	return c.Id
}

func (c *UserCustomvar) SetId(id string) {
	c.Id = id
}

func (c *UserCustomvar) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{c}, nil
}

func init() {
	name := "user_customvar"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 "user:customvar",
		PrimaryMySqlField:        "id",
		Factory:                  NewUserCustomvar,
		HasChecksum:              false,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "user",
	}
}
