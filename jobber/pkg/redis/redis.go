package redis

import (
	"errors"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type DbClient interface {
    Get(key interface{}) (value string, err error)
    Set(key interface{}, value string) (err error)
    Delete(key interface{}) (err error)
}

type db struct {
	pool *redis.Pool
}

func (d *db) Get(key interface{}) (value string, err error){
	conn := d.pool.Get()
	defer conn.Close()
	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return
	} else if value == "" {
		err = errors.New("empty key")
		return
	}
	return
}

func (d *db) Set(key interface{}, value string) (err error){
	conn := d.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return
	}
	_, err = conn.Do("EXPIRE", key, os.Getenv("CASH_EXPIRE"))
	if err != nil {
		return
	}
	return
}

func (d *db) Delete(key interface{}) (err error){
	conn := d.pool.Get()
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		return
	}
	return
}

// NewDb returns a new Db instance.
func NewDb() DbClient {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS_URL"))
		},
	}

	return &db{pool: pool}
}