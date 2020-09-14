package redis

import (
	"errors"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type DbClient interface {
    Get(key string) (value string, err error)
    Set(key string, value string) (err error)
    Delete(key string) (err error)
}

type db struct {
	pool *redis.Pool
}

func (d *db) Get(key string) (value string, err error){
	conn := d.pool.Get()
	defer conn.Close()
	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return
	} else if value == "" {
		err = errors.New("empty value")
		return
	}
	return
}

func (d *db) Set(key string, value string) (err error){
	conn := d.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return
	}
	//untilMorrow := time.Until()
	_, err = conn.Do("EXPIRE", key, 3600)
	if err != nil {
		return
	}
	return
}

func (d *db) Delete(key string) (err error){
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