package redis

import (
	"errors"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type DbClient interface {
    Get(key string) (value string, err error)
    Set(key string, value string) (err error)
}

type db struct {
	pool *redis.Pool
}

func (d *db) Get(key string) (value string, err error){
	conn := d.pool.Get()
	defer conn.Close()
	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		log.Warn(err)
		return
	} else if value == "" {
		err = errors.New("empty key")
		log.Trace(err)
		return
	}
	log.Trace("Get ", key, ":", value)
	return
}

func (d *db) Set(key string, value string) (err error){
	conn := d.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	if err != nil {
		log.Warn(err)
	}
	log.Trace("Set ", key, ":", value)
	return
}

// NewDb returns a new Db instance.
func NewDb() DbClient {
	pool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("DATABASE_URL"))
		},
	}

	return &db{pool: pool}
}