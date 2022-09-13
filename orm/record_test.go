package orm

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/orm/dialect"
	"github.com/oldbai555/lb/orm/session"
	"testing"
)

var (
	user1 = &User{
		Name: "Tom",
		Age:  18,
	}
	user2 = &User{
		Name: "Sam",
		Age:  25,
	}
	user3 = &User{
		Name: "Jack",
		Age:  25,
	}
)

func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	engine, err := NewEngine(dialect.DMYSQL, "root:xxx@tcp(xxx:3306)/orm")
	if err != nil {
		log.Errorf("err:%v", err)
		return nil
	}
	s := engine.NewSession().Model(&User{})
	err1 := dropTable(s)
	err2 := createTable(s)
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
