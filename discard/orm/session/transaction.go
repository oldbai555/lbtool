package session

import "github.com/oldbai555/lbtool/log"

func (s *Session) Begin() (err error) {
	log.Infof("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Errorf("err is %v", err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Infof("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Errorf("err is %v", err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Infof("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Errorf("err is %v", err)
	}
	return
}
