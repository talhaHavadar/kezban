package models

import (
	"github.com/talhaHavadar/kezban"
	"time"
)

type Kezban struct {
	kezban.Model			`,inline`
	Name		string		`bson:"name,omitempty" json:"name" kezban:"unique"`
	Birthdate	time.Time	`bson:"birthdate,omitempty" json:"birthdate"`
	GrumpyLevel	int		`bson:"grumpy_level,omitempty" json:"grumyLevel"`
}

func (self *Kezban) GetCollectionName() string {
	return "kezbans"
}