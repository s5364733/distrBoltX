package serv

import (
	"fmt"
	rpc "github.com/s5364733/distrBoltX/rpc/proto"
	"github.com/s5364733/distrBoltX/web"
	"io"
	"time"
)

type AckSyncDialerService struct {
	rpc.UnimplementedAckSyncDialerServer
	s *web.Server
}

func NewAckSyncDialerService(serv *web.Server) *AckSyncDialerService {
	return &AckSyncDialerService{
		s: serv,
	}
}

func (c *AckSyncDialerService) Dial(stream rpc.AckSyncDialer_DialServer) error {
	for {
		//1. 拿到最新KEY
		key, v, err2 := c.s.Db.GetNextKeyForReplication()
		if err2 != nil || key == nil || v == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		//2.发送到副本数据同步到副本节点bucket
		err2 = stream.Send(&rpc.NextKeyValue{
			Key:   string(key),
			Value: string(v),
		})
		if err2 != nil {
			fmt.Errorf("err  %v", err2)
		}
		fmt.Printf("Data sent to the replica is synchronized to the replica node key = %q,value=%q", key, v)
		//3.副本节点同步成功后发送ACK 标识
		ack, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("The ACK identifier of the replica node synchronization is completed ack = %q", ack)
		//4 删除主节点副本
		if ack != nil && ack.Ack {
			fmt.Printf("The key asynchronized from the master ,which has been deleted key %q val %q ", string(key), string(v))
			err2 := c.s.Db.DeleteReplicationKey(key, v)
			if err2 != nil {
				fmt.Errorf("%v", err2)
			}
			fmt.Printf("Data replica sync have been done,deleting local key =  %q", key)
		}
	}
}
