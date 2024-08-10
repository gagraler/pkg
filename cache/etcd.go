package cache

import (
	"context"
	"github.com/gagraler/pkg/logger"
	"go.etcd.io/etcd/clientv3"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/8/10 18:37
 * @file: etcd.go
 * @description: etcd operation
 */

type Etcd struct {
	EndPoints []string
	cli       *clientv3.Client
}

var logs = logger.SugaredLogger()

func NewEtcd(endpoints []string) *Etcd {

	client, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		logs.Errorf("etcd client error: %v", err)
	}

	defer func(client *clientv3.Client) {
		err := client.Close()
		if err != nil {
			logs.Errorf("etcd client close error: %v", err)
		}
	}(client)
	return &Etcd{
		EndPoints: endpoints,
		cli:       client,
	}
}

// Set key value
func (e *Etcd) Set(key, value string) error {

	_, err := e.cli.Put(context.Background(), key, value)
	if err != nil {
		return err
	}
	return nil
}

// HSet key field value
func (e *Etcd) HSet(key string, fields map[string]string) error {

	for field, value := range fields {
		_, err := e.cli.Put(context.Background(), key+"/"+field, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Etcd) grant(ttl int64) clientv3.LeaseID {
	grant, err := e.cli.Grant(context.Background(), ttl)
	if err != nil {
		logs.Errorf("etcd grant error: %v", err)
	}
	return grant.ID
}

// LessSet key value
func (e *Etcd) LessSet(key, value string, ttl int64) error {

	grant, err := e.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	_, err = e.cli.Put(context.Background(), key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return err
	}
	return nil
}

// LessHSet key field value
func (e *Etcd) LessHSet(key string, fields map[string]string, ttl int64) error {

	grant, err := e.cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	for field, value := range fields {
		_, err := e.cli.Put(context.Background(), key+"/"+field, value, clientv3.WithLease(grant.ID))
		if err != nil {
			return err
		}
	}
	return nil
}

// Get key value
func (e *Etcd) Get(key string) (string, error) {

	resp, err := e.cli.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// HGet key field value
func (e *Etcd) HGet(key, field string) (string, error) {

	resp, err := e.cli.Get(context.Background(), key+"/"+field)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// Del key
func (e *Etcd) Del(key string) error {

	_, err := e.cli.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}

// Watch key
func (e *Etcd) Watch(key string) error {

	rch := e.cli.Watch(context.Background(), key, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			logs.Infof("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
	return nil
}
