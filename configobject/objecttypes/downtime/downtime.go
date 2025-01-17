// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package downtime

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
		"triggered_by_id",
		"object_type",
		"host_id",
		"service_id",
		"name_checksum",
		"properties_checksum",
		"name",
		"author",
		"comment",
		"entry_time",
		"scheduled_start_time",
		"scheduled_end_time",
		"flexible_duration",
		"is_flexible",
		"is_in_effect",
		"start_time",
		"end_time",
		"zone_id",
	}
)

type Downtime struct {
	Id                 string  `json:"id"`
	EnvId              string  `json:"environment_id"`
	TriggeredById      string  `json:"triggered_by_id"`
	ObjectType         string  `json:"object_type"`
	HostId             string  `json:"host_id"`
	ServiceId          string  `json:"service_id"`
	NameChecksum       string  `json:"name_checksum"`
	PropertiesChecksum string  `json:"checksum"`
	Name               string  `json:"name"`
	Author             string  `json:"author"`
	Comment            string  `json:"comment"`
	EntryTime          float64 `json:"entry_time"`
	ScheduledStartTime float64 `json:"scheduled_start_time"`
	ScheduledEndTime   float64 `json:"scheduled_end_time"`
	FlexibleDuration   float64 `json:"flexible_duration"`
	IsFlexible         bool    `json:"is_flexible"`
	IsInEffect         bool    `json:"is_in_effect"`
	StartTime          float64 `json:"start_time"`
	EndTime            float64 `json:"end_time"`
	ZoneId             string  `json:"zone_id"`
}

func NewDowntime() connection.Row {
	d := Downtime{}

	return &d
}

func (d *Downtime) InsertValues() []interface{} {
	v := d.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(d.Id)}, v...)
}

func (d *Downtime) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(d.EnvId),
		utils.EncodeChecksum(d.TriggeredById),
		d.ObjectType,
		utils.EncodeChecksum(d.HostId),
		utils.EncodeChecksum(d.ServiceId),
		utils.EncodeChecksum(d.NameChecksum),
		utils.EncodeChecksum(d.PropertiesChecksum),
		d.Name,
		d.Author,
		d.Comment,
		d.EntryTime,
		d.ScheduledStartTime,
		d.ScheduledEndTime,
		d.FlexibleDuration,
		utils.Bool[d.IsFlexible],
		utils.Bool[d.IsInEffect],
		d.StartTime,
		d.EndTime,
		utils.EncodeChecksum(d.ZoneId),
	)

	return v
}

func (d *Downtime) GetId() string {
	return d.Id
}

func (d *Downtime) SetId(id string) {
	d.Id = id
}

func (d *Downtime) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{d}, nil
}

func init() {
	name := "downtime"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:               name,
		RedisKey:                 "downtime",
		PrimaryMySqlField:        "id",
		Factory:                  NewDowntime,
		HasChecksum:              true,
		BulkInsertStmt:           connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:           connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:           connection.NewBulkUpdateStmt(name, Fields),
		NotificationListenerType: "downtime",
	}
}
