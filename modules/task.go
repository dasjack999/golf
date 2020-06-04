package modules

import (
	"../repository"
)

type Task struct {
	//
	repo.BaseEntity `json:",inline" bson:",inline"`
	//
	Name string `json:"name" bson:"name"`
	//
	CreateTime int64 `json:"ctime" bson:"ctime"`
	//
	TableId string `json:"tid" bson:"tid"`
	//
	OwnerId string `json:"oid" bson:"oid"`
	//
	//Props map[string]interface{} `json:"props" bson:"props"`
	//
	owner *Account
}

//
func NewTask(source string) *Task {
	a := &Task{}

	ent := repo.GetDataSource(source)
	if err := a.Init(ent, map[string]interface{}{
		"CollName": "task",
	}); err != nil {
		return nil
	}

	return a
}

//
func (t *Task) LoadById(id string) bool {
	return t.BaseEntity.LoadById(id, t)
}

//
func (t *Task) Save() bool {
	return t.BaseEntity.Save(t)
}

//
func (t *Task) SetOwner(a *Account) {
	t.owner = a
	t.OwnerId = a.Id
}

//
func (t *Task) GetOwner() *Account {
	if t.owner == nil {
		ac := NewAccount("redis")
		if !ac.LoadById(t.OwnerId, ac) {
			return nil
		}
		t.owner = ac
	}
	return t.owner
}
