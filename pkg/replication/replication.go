package replication

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/s5364733/distrBoltX/api"
	"github.com/s5364733/distrBoltX/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// NextKeyValue contains the response for GetNextKeyForReplication.
type NextKeyValue struct {
	Key   string
	Value string
	Err   error
}

type client struct {
	db         *db.Database
	leaderAddr string
}

// ClientLoop continuously downloads new keys from the master and applies them.
func ClientLoop(db *db.Database, leaderAddr string) {
	c := &client{db: db, leaderAddr: leaderAddr}
	for {
		present, err := c.loop()
		if err != nil {
			log.Printf("Loop error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if !present {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// ClientGrpcLoop continuously stream rpc for grpc sync  data's
func ClientGrpcLoop(db *db.Database, leaderAddr string) {
	c := &client{db: db, leaderAddr: leaderAddr}
	//for {
	err := c.grpcLoop()
	if err != nil {
		log.Printf("grpcLoop error: %v", err)
		time.Sleep(time.Second)
		//continue
	}
	//}
}

// grpcLoop grpc
// the default keepalive tpc syn ack is opened for this link
func (c *client) grpcLoop() (err error) {
	conn, err := grpc.Dial("127.0.0.2:50030", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return err
	}
	defer conn.Close()
	dialerClient := api.NewAckSyncDialerClient(conn)
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	stream, err := dialerClient.Dial(context.TODO())
	if err != nil {
		return err
	}
	waitc := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			//set to current replication
			c.syncReplicationBolt(res)
			err = stream.Send(&api.SyncD{
				Ack: true,
			})
			if err != nil {
				fmt.Errorf("err %v", err)
			}
		}
	}()

	select {
	case <-waitc:
		stream.CloseSend()
	}
	return nil
}

func (c *client) syncReplicationBolt(res *api.NextKeyValue) {
	//设置到当前节点
	if err := c.db.SetKeyOnReplica(res.Key, []byte(res.Value)); err != nil {
		fmt.Errorf("err for operation of sync , key = %q,value = %q, err = %+v", res.Key, res.Value, err)
	}
	log.Printf("The key asynchronized from the master , which has been loaded (key:%q,value:%q)", res.Key, res.Value)
}

// loop Return false used to do wait 100 millis
// the default keepalive tpc syn ack is opened for this link
func (c *client) loop() (present bool, err error) {
	//Sync
	//拿到主分片的副本分片数据
	resp, err := http.Get("http://" + c.leaderAddr + "/next-replication-key")
	if err != nil {
		return false, err
	}

	var res NextKeyValue
	//解析成功
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	defer resp.Body.Close()

	//报错
	if res.Err != nil {
		fmt.Sprintf("The value of sync to master which occurs err ")
		return false, err
	}
	//没有这个KEY
	if res.Key == "" {
		fmt.Sprintf("The value of sync to master which is nil ")
		return false, nil
	}
	//errors.New()
	//设置到当前节点
	if err := c.db.SetKeyOnReplica(res.Key, []byte(res.Value)); err != nil {
		err := errors.New("error")
		panic(err) //throws error
		return false, err
	}

	log.Printf("The key asynchronized from the master , which has been loaded (key:%q,value:%q)", res.Key, res.Value)
	//Deletes the key of replica's bucket of master
	if err := c.deleteFromReplicationQueue(res.Key, res.Value); err != nil {
		log.Printf("DeleteKeyFromReplication failed: %v", err)
	}

	return true, nil
}

func (c *client) deleteFromReplicationQueue(key, value string) error {
	u := url.Values{}
	u.Set("key", key)
	u.Set("value", value)

	log.Printf("Deleting key=%q, value=%q from replication queue on %q", key, value, c.leaderAddr)

	resp, err := http.Get("http://" + c.leaderAddr + "/delete-replication-key?" + u.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !bytes.Equal(result, []byte("ok")) {
		return errors.New(string(result))
	}

	return nil
}
