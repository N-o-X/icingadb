// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package supervisor

import (
	"github.com/N-o-X/icingadb/connection"
	"github.com/N-o-X/icingadb/jsondecoder"
	"sync"
)

type Supervisor struct {
	ChErr    chan error
	ChDecode chan *jsondecoder.JsonDecodePackages
	Rdbw     *connection.RDBWrapper
	Dbw      *connection.DBWrapper
	EnvId    []byte
	EnvLock  *sync.Mutex
}
