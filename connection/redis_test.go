// IcingaDB | (c) 2019 Icinga GmbH | GPLv2+

package connection

import (
	"github.com/go-redis/redis"
	"github.com/N-o-X/icingadb/connection/redisd"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func NewTestRDBW(rdb RedisClient) RDBWrapper {
	dbw := RDBWrapper{Rdb: rdb, ConnectedAtomic: new(uint32), ConnectionLostCounterAtomic: new(uint32)}
	dbw.ConnectionUpCondition = sync.NewCond(&sync.Mutex{})
	return dbw
}

func TestNewRDBWrapper(t *testing.T) {
	var server redisd.Server

	client, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()
	defer client.Close()

	rdbw := NewRDBWrapper(client.Options().Addr, 64)
	assert.True(t, rdbw.CheckConnection(false), "Redis should be connected")

	rdbw = NewRDBWrapper("asdasdasdasdasd:5123", 64)
	assert.False(t, rdbw.CheckConnection(false), "Redis should not be connected")
	//TODO: Add more tests here
}

func TestRDBWrapper_GetConnectionCheckInterval(t *testing.T) {
	rdbw := NewTestRDBW(nil)

	//Should return 15s, if connected - counter doesn't madder
	rdbw.CompareAndSetConnected(true)
	assert.Equal(t, 15*time.Second, rdbw.getConnectionCheckInterval())

	//Should return 5s, if not connected and counter < 4
	rdbw.CompareAndSetConnected(false)
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 0)
	assert.Equal(t, 5*time.Second, rdbw.getConnectionCheckInterval())

	//Should return 10s, if not connected and 4 <= counter < 8
	rdbw.CompareAndSetConnected(false)
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 4)
	assert.Equal(t, 10*time.Second, rdbw.getConnectionCheckInterval())

	//Should return 30s, if not connected and 8 <= counter < 11
	rdbw.CompareAndSetConnected(false)
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 8)
	assert.Equal(t, 30*time.Second, rdbw.getConnectionCheckInterval())

	//Should return 60s, if not connected and 11 <= counter < 14
	rdbw.CompareAndSetConnected(false)
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 11)
	assert.Equal(t, 60*time.Second, rdbw.getConnectionCheckInterval())

	//Should panic, if not connected and counter > 13
	rdbw.CompareAndSetConnected(false)
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 14)
	assert.Panics(t, func() { rdbw.getConnectionCheckInterval() }, "Should panic")
}

func TestRDBWrapper_CheckConnection(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(nil)

	rdbw.Rdb = rdb
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 512312312)
	assert.True(t, rdbw.CheckConnection(false), "DBWrapper should be connected")
	assert.Equal(t, uint32(0), atomic.LoadUint32(rdbw.ConnectionLostCounterAtomic))

	rdbw.Rdb = redis.NewClient(&redis.Options{
		Addr:         "dasdasdasdasdasd:5123",
		DialTimeout:  time.Minute / 2,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	})
	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 0)
	assert.False(t, rdbw.CheckConnection(false), "DBWrapper should not be connected")
	assert.Equal(t, uint32(0), atomic.LoadUint32(rdbw.ConnectionLostCounterAtomic))

	atomic.StoreUint32(rdbw.ConnectionLostCounterAtomic, 10)
	assert.False(t, rdbw.CheckConnection(true), "DBWrapper should not be connected")
	assert.Equal(t, uint32(11), atomic.LoadUint32(rdbw.ConnectionLostCounterAtomic))
}

func TestRDBWrapper_HGetAll(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.Del("herpdaderp")
	rdb.HSet("herpdaderp", "one", 5)
	rdb.HSet("herpdaderp", "two", 11)

	rdbw.CompareAndSetConnected(false)

	var data map[string]string
	var err error
	done := make(chan bool)
	go func() {
		var res = rdbw.HGetAll("herpdaderp")
		data, err = res.Result()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	rdbw.CheckConnection(true)

	<-done

	assert.NoError(t, err)
	assert.Contains(t, data, "one")
	assert.Contains(t, data, "two")
}

func TestRDBWrapper_HKeys(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.Del("firstKey")
	rdb.Del("secondKey")
	rdb.HSet("firstKey", "foo", 5)
	rdb.HSet("firstKey", "abc", 2)
	rdb.HSet("secondKey", "bar", 11)

	assert.Equal(t, []string{"foo", "abc"}, rdbw.HKeys("firstKey").Val())
	assert.Equal(t, []string{"bar"}, rdbw.HKeys("secondKey").Val())
}

func TestRDBWrapper_HMGet(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.Del("firstKey")
	rdb.HSet("firstKey", "foo", "5")
	rdb.HSet("firstKey", "abc", "2")

	assert.Equal(t, []interface{}{"5"}, rdbw.HMGet("firstKey", "foo").Val())
	assert.Equal(t, []interface{}{"2"}, rdbw.HMGet("firstKey", "abc").Val())
}

func TestRDBWrapper_XRead(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.XTrim("teststream", 0)
	rdb.XAdd(&redis.XAddArgs{Stream: "teststream", Values: map[string]interface{}{"one": "5", "two": "11", "herp": "11"}})

	rdbw.CompareAndSetConnected(false)

	var data *redis.XStreamSliceCmd
	done := make(chan bool)
	go func() {
		data = rdbw.XRead(&redis.XReadArgs{Streams: []string{"teststream", "0"}})
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	rdbw.CheckConnection(true)

	<-done

	streams, err := data.Result()
	assert.NoError(t, err)
	value := streams[0].Messages[0].Values

	assert.Contains(t, value, "one")
	assert.Contains(t, value, "two")
}

func TestRDBWrapper_XDel(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.XTrim("teststream", 0)
	adds := rdb.XAdd(&redis.XAddArgs{Stream: "teststream", Values: map[string]interface{}{"one": "5", "two": "11", "herp": "11"}})

	rdbw.CompareAndSetConnected(false)

	done := make(chan bool)
	go func() {
		rdbw.XDel("teststream", adds.Val())
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	rdbw.CheckConnection(true)

	<-done

	data := rdbw.XRead(&redis.XReadArgs{Streams: []string{"teststream", "0"}, Block: -1})
	streams, err := data.Result()
	assert.Error(t, err)
	assert.Len(t, streams, 0)
}

func TestRDBWrapper_Publish(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	var msg *redis.Message
	var err error
	done := make(chan bool)
	go func() {
		msg, err = rdb.Subscribe("testchannel").ReceiveMessage()
		done <- true
	}()

	rdbw.CompareAndSetConnected(false)

	go func() {
		rdbw.Publish("testchannel", "Hello there")
	}()

	time.Sleep(50 * time.Millisecond)
	rdbw.CheckConnection(true)

	<-done

	assert.NoError(t, err)
	assert.Equal(t, "Hello there", msg.Payload)
}

func TestRDBWrapper_TxPipelined(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.Del("firstKey")
	rdb.Del("secondKey")
	rdb.HSet("firstKey", "foo", 5)
	rdb.HSet("secondKey", "bar", 11)

	rdbw.CompareAndSetConnected(false)

	var firstMap *redis.StringStringMapCmd
	var secondMap *redis.StringStringMapCmd
	var err error
	done := make(chan bool)
	go func() {
		_, err = rdbw.TxPipelined(func(pipe redis.Pipeliner) error {
			firstMap = pipe.HGetAll("firstKey")
			secondMap = pipe.HGetAll("secondKey")
			return nil
		})
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	rdbw.CheckConnection(true)

	<-done

	assert.NoError(t, err)
	assert.Contains(t, firstMap.Val(), "foo")
	assert.Contains(t, secondMap.Val(), "bar")
}

func TestRDBWrapper_PipeConfigChunks(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.Del("icinga:config:testkey")
	rdb.Del("icinga:checksum:testkey")

	rdb.HSet("icinga:config:testkey", "123534534fsdf12sdas12312adg23423f", "this-should-be-the-config")
	rdb.HSet("icinga:checksum:testkey", "123534534fsdf12sdas12312adg23423f", "this-should-be-the-checksum")

	chChunk := rdbw.PipeConfigChunks(make(chan struct{}), []string{"123534534fsdf12sdas12312adg23423f"}, "testkey")
	chunk := <-chChunk
	assert.Equal(t, "this-should-be-the-config", chunk.Configs[0])
	assert.Equal(t, "this-should-be-the-checksum", chunk.Checksums[0])
}

func TestRDBWrapper_PipeChecksumChunks(t *testing.T) {
	var server redisd.Server

	rdb, errSrv := server.Start()
	if errSrv != nil {
		t.Fatal(errSrv)
		return
	}

	defer server.Stop()

	rdbw := NewTestRDBW(rdb)

	if !rdbw.CheckConnection(true) {
		t.Fatal("This test needs a working Redis connection")
	}

	rdb.Del("icinga:checksum:testkey")

	rdb.HSet("icinga:checksum:testkey", "123534534fsdf12sdas12312adg23423f", "this-should-be-the-checksum")

	chChunk := rdbw.PipeChecksumChunks(make(chan struct{}), []string{"123534534fsdf12sdas12312adg23423f"}, "testkey")
	chunk := <-chChunk
	assert.Equal(t, "this-should-be-the-checksum", chunk.Checksums[0])
}
