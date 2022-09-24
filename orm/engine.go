package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/orm/dialect"
	"github.com/oldbai555/lbtool/orm/session"
	"github.com/oldbai555/lbtool/pkg/exception"
	"strings"
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

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != any(nil) {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}

// difference returns a - b
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (re interface{}, err error) {
		_, err = doDescTable(s)
		if err != nil {
			log.Errorf("err:%v", err)
			if exception.GetErrCode(err) == exception.ErrOrmTableNotExist {
				return nil, createTable(s)
			}
			return nil, err
		}
		table := s.Table
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.DbName, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
