// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package configobject

import (
	"github.com/N-o-X/icingadb/connection"
)

type ObjectInformation struct {
	ObjectType               string
	RedisKey                 string
	PrimaryMySqlField        string
	HasChecksum              bool
	NotificationListenerType string
	Factory                  connection.RowFactory
	BulkInsertStmt           *connection.BulkInsertStmt
	BulkDeleteStmt           *connection.BulkDeleteStmt
	BulkUpdateStmt           *connection.BulkUpdateStmt
}
