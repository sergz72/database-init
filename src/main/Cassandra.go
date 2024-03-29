package main

import (
	"github.com/gocql/gocql"
	"time"
)

type cassandraDriver struct {
	host    string
	session *gocql.Session
}

func newCassandraDriver(host string) *cassandraDriver {
	return &cassandraDriver{host: host}
}

func (d *cassandraDriver) Connect() error {
	cluster := gocql.NewCluster(d.host)
	cluster.Timeout = time.Minute
	cluster.Consistency = gocql.Quorum
	var err error
	d.session, err = cluster.CreateSession()
	return err
}

func (d *cassandraDriver) Disconnect() {
	d.session.Close()
}

func (d *cassandraDriver) Exec(sql string) error {
	return d.session.Query(sql).Exec()
}

func (d *cassandraDriver) CreateDatabase(dbName string, dbUser string, dbPass string) error {
	if err := d.Exec(preprocess(dbName, dbUser, dbPass, "DROP KEYSPACE IF EXISTS ${DB_NAME}")); err != nil {
		return err
	}
	return d.Exec(preprocess(dbName, dbUser, dbPass, "CREATE KEYSPACE ${DB_NAME} WITH REPLICATION = { 'class': 'SimpleStrategy','replication_factor': 1} AND DURABLE_WRITES =  true"))
}
