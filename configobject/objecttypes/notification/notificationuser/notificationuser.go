// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package notificationuser

import (
	"github.com/N-o-X/icingadb/configobject"
	"github.com/N-o-X/icingadb/connection"
	"github.com/N-o-X/icingadb/utils"
)

var (
	ObjectInformation configobject.ObjectInformation
	Fields            = []string{
		"id",
		"notification_id",
		"user_id",
		"environment_id",
	}
)

type NotificationUser struct {
	Id             string `json:"id"`
	NotificationId string `json:"notification_id"`
	UserId         string `json:"user_id"`
	EnvId          string `json:"environment_id"`
}

func NewNotificationUser() connection.Row {
	n := NotificationUser{}
	return &n
}

func (n *NotificationUser) InsertValues() []interface{} {
	v := n.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(n.Id)}, v...)
}

func (n *NotificationUser) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(n.NotificationId),
		utils.EncodeChecksum(n.UserId),
		utils.EncodeChecksum(n.EnvId),
	)

	return v
}

func (n *NotificationUser) GetId() string {
	return n.Id
}

func (n *NotificationUser) SetId(id string) {
	n.Id = id
}

func (n *NotificationUser) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{n}, nil
}

func init() {
	name := "notification_user"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 "notification:user",
		PrimaryMySqlField:        "id",
		Factory:                  NewNotificationUser,
		HasChecksum:              false,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "notification",
	}
}
