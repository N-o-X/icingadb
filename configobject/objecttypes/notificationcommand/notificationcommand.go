// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package notificationcommand

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
		"name",
		"name_ci",
		"zone_id",
		"command",
		"timeout",
	}
)

type NotificationCommand struct {
	Id                 string  `json:"id"`
	EnvId              string  `json:"environment_id"`
	NameChecksum       string  `json:"name_checksum"`
	PropertiesChecksum string  `json:"checksum"`
	Name               string  `json:"name"`
	NameCi             *string `json:"name_ci"`
	ZoneId             string  `json:"zone_id"`
	Command            string  `json:"command"`
	Timeout            float32 `json:"timeout"`
}

func NewNotificationCommand() connection.Row {
	c := NotificationCommand{}
	c.NameCi = &c.Name

	return &c
}

func (c *NotificationCommand) InsertValues() []interface{} {
	v := c.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(c.Id)}, v...)
}

func (c *NotificationCommand) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(c.EnvId),
		utils.EncodeChecksum(c.NameChecksum),
		utils.EncodeChecksum(c.PropertiesChecksum),
		c.Name,
		c.NameCi,
		utils.EncodeChecksum(c.ZoneId),
		c.Command,
		c.Timeout,
	)

	return v
}

func (c *NotificationCommand) GetId() string {
	return c.Id
}

func (c *NotificationCommand) SetId(id string) {
	c.Id = id
}

func (c *NotificationCommand) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{c}, nil
}

func init() {
	name := "notificationcommand"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 name,
		PrimaryMySqlField:        "id",
		Factory:                  NewNotificationCommand,
		HasChecksum:              true,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "notificationcommand",
	}
}
