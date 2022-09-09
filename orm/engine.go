package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/orm/dialect"
	"github.com/oldbai555/lb/orm/session"
)

type Engine struct {
	db   *sqlx.DB
	dial dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sqlx.Open(driver, source)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// make sure the specific dialect exists
	dial, err := dialect.GetDialect(driver)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	e = &Engine{
		db:   db,
		dial: dial,
	}
	log.Infof("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Errorf("Failed to close database,err is %v", err)
		return
	}
	log.Infof("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.NewSession(engine.db, engine.dial)
}
