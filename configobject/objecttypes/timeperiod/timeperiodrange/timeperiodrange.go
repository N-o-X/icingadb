// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package timeperiodrange

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
		"range_key",
		"range_value",
		"environment_id",
	}
)

type TimeperiodRange struct {
	Id           string `json:"id"`
	TimeperiodId string `json:"timeperiod_id"`
	RangeKey     string `json:"range_key"`
	RangeValue   string `json:"range_value"`
	EnvId        string `json:"environment_id"`
}

func NewTimeperiodRange() connection.Row {
	t := TimeperiodRange{}
	return &t
}

func (t *TimeperiodRange) InsertValues() []interface{} {
	v := t.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(t.Id)}, v...)
}

func (t *TimeperiodRange) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(t.TimeperiodId),
		t.RangeKey,
		t.RangeValue,
		utils.EncodeChecksum(t.EnvId),
	)

	return v
}

func (t *TimeperiodRange) GetId() string {
	return t.Id
}

func (t *TimeperiodRange) SetId(id string) {
	t.Id = id
}

func (t *TimeperiodRange) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{t}, nil
}

func init() {
	name := "timeperiod_range"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 "timeperiod:range",
		PrimaryMySqlField:        "id",
		Factory:                  NewTimeperiodRange,
		HasChecksum:              false,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "timeperiod",
	}
}
