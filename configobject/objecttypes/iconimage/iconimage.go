// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package iconimage

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
		"icon_image",
	}
)

type IconImage struct {
	Id        string `json:"id"`
	EnvId     string `json:"environment_id"`
	IconImage string `json:"icon_image"`
}

func NewIconImage() connection.Row {
	a := IconImage{}

	return &a
}

func (a *IconImage) InsertValues() []interface{} {
	v := a.UpdateValues()

	return append([]interface{}{utils.EncodeChecksum(a.Id)}, v...)
}

func (a *IconImage) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.EncodeChecksum(a.EnvId),
		a.IconImage,
	)

	return v
}

func (a *IconImage) GetId() string {
	return a.Id
}

func (a *IconImage) SetId(id string) {
	a.Id = id
}

func (a *IconImage) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{a}, nil
}

func init() {
	name := "icon_image"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType:        name,
		RedisKey:          name,
		PrimaryMySqlField: "id",
		Factory:           NewIconImage,
		HasChecksum:       false,
		BulkInsertStmt:    connection.NewBulkInsertStmt(name, Fields),
		BulkDeleteStmt:    connection.NewBulkDeleteStmt(name, "id"),
		BulkUpdateStmt:    connection.NewBulkUpdateStmt(name, Fields),
	}
}
