package models

type User struct {
	Id int64 `db:"id" json:"id"`
	Name string `sql:"size:60" db:"name" json:"name"`
	Password string `sql:"size:60" db:"password" json:"password"`
	Hobby string `sql:"size:60" db:"hobby" json:"hobby"`
}