package cache

import (
	"errors"

	"github.com/etcd-io/bbolt"
)

type Cache struct {
	db    *bbolt.DB

	table []byte
}

func New(db *bbolt.DB, table string) *Cache { return &Cache{db: db, table: []byte(table)} }

func (c *Cache) Get(key string) string {
	var res []byte
	c.db.View(func(tx *bbolt.Tx) error {
		bu := tx.Bucket(c.table)
		// 桶不存在
		if bu == nil {
			return nil
		}
		r := bu.Get([]byte(key))
		// key不存在
		if r != nil {
			res = r
		}
		return nil
	})
	return string(res)
}

func (c *Cache) Set(key string, val []byte) error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		bu, err := tx.CreateBucketIfNotExists(c.table)
		if err != nil {
			return err
		}
		return bu.Put([]byte(key), val)
	})
}

func (c *Cache) IsExist(key string) bool {
	if c.db.View(func(tx *bbolt.Tx) error {
		bu, err := tx.CreateBucketIfNotExists(c.table)
		if err != nil {
			return errors.New("bucket not exists")
		}
		if bu.Get([]byte(key)) == nil {
			return errors.New("not exists")
		}
		return nil
	}) != nil {
		return false
	}
	return true
}

func (c *Cache) Delete(key string) error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		bu := tx.Bucket(c.table)
		return bu.Delete([]byte(key))
	})
}

func (c *Cache) Flush() error {
	err := c.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket(c.table)
	})
	if err != bbolt.ErrBucketNotFound {
		return err
	}
	return nil
}
