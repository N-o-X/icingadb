package customvarflat

import (
	"encoding/json"
	"fmt"
	"git.icinga.com/icingadb/icingadb-main/configobject"
	"git.icinga.com/icingadb/icingadb-main/connection"
	"git.icinga.com/icingadb/icingadb-main/utils"
)

var (
	ObjectInformation configobject.ObjectInformation
	Fields         = []string{
		"id",
		"env_id",
		"customvar_id",
		"flatname_checksum",
		"flatname",
		"flatvalue",
	}
)

type CustomvarFlat struct {
	Id                   string `json:"id"`
	EnvId                string `json:"env_id"`
	NameChecksum         string `json:"name_checksum"`
	Name                 string `json:"name"`
	Value                string `json:"value"`
}

func NewCustomvarFlat() connection.Row {
	c := CustomvarFlat{}

	return &c
}

func (c *CustomvarFlat) InsertValues() []interface{} {
	return nil
}

func (c *CustomvarFlat) UpdateValues() []interface{} {
	return nil
}

func (c *CustomvarFlat) GetId() string {
	return c.Id
}

func (c *CustomvarFlat) SetId(id string) {
	c.Id = id
}

func (c *CustomvarFlat) GetFinalRows() ([]connection.Row, error) {
	var values interface{} = nil
	if err := json.Unmarshal([]byte(c.Value), &values); err != nil {
		return nil, err
	}

	return CollectScalarVars(c, values, c.Name, make([]string, 0)), nil
}

func CollectScalarVars(c *CustomvarFlat, value interface{}, name string, path []string) []connection.Row {
	path = append(path, name)
	switch v := value.(type) {
	case map[string]interface{}:
		var rows = []connection.Row{}
		for flatName, flatValue := range v {
			rows = append(rows, CollectScalarVars(c, flatValue, flatName, path)...)
		}

		return rows
	case []interface{}:
		var rows = []connection.Row{}
		for i, flatValue := range v {
			rows = append(rows, CollectScalarVars(c, flatValue, fmt.Sprintf("%d", i), path)...)
		}

		return rows
	default:
		flatName := fmt.Sprintf("%v", path)
		flatValue := fmt.Sprintf("%v", v)
		return []connection.Row{
			&CustomvarFlatFinal{
				Id:               utils.StringToSha1String(c.EnvId + flatName +  flatValue),
				EnvId:            c.EnvId,
				CustomvarId:      c.Id,
				FlatNameChecksum: utils.StringToSha1String(flatName),
				FlatName:         flatName,
				FlatValue:        flatValue,
			},
		}
	}
}

type CustomvarFlatFinal struct {
	Id          		string
	EnvId       		string
	CustomvarId			string
	FlatNameChecksum    string
	FlatName        	string
	FlatValue       	string
}

func (c *CustomvarFlatFinal) InsertValues() []interface{} {
	v := c.UpdateValues()

	return append([]interface{}{utils.Checksum(c.Id)}, v...)
}

func (c *CustomvarFlatFinal) UpdateValues() []interface{} {
	v := make([]interface{}, 0)

	v = append(
		v,
		utils.Checksum(c.EnvId),
		utils.Checksum(c.CustomvarId),
		utils.Checksum(c.FlatNameChecksum),
		c.FlatName,
		c.FlatValue,
	)

	return v
}

func (c *CustomvarFlatFinal) GetId() string {
	return c.Id
}

func (c *CustomvarFlatFinal) SetId(id string) {
	c.Id = id
}

func (c *CustomvarFlatFinal) GetFinalRows() ([]connection.Row, error) {
	return []connection.Row{c}, nil
}

func init() {
	name := "customvar_flat"
	ObjectInformation = configobject.ObjectInformation{
		ObjectType: name,
		RedisKey: "customvar",
		PrimaryMySqlField: "customvar_id",
		Factory: NewCustomvarFlat,
		HasChecksum: false,
		BulkInsertStmt: connection.NewBulkInsertStmt(name, Fields, "id"),
		BulkDeleteStmt: connection.NewBulkDeleteStmt(name,  "customvar_id"),
		BulkUpdateStmt: connection.NewBulkUpdateStmt(name, Fields),
	}
}