package main

import (
	"github.com/talhaHavadar/kezban"
	"github.com/talhaHavadar/kezban/examples/models"
	"fmt"
	"time"
)

func main() {
	kezban.Initialize("mongodb://localhost:27017", "KezbanMetric")

	layout := "2006-01-02"

	hard := new(models.Kezban)
	hard.Name = "Kezban"
	t, _ := time.Parse(layout, "1992-03-23")
	hard.Birthdate = t
	hard.GrumpyLevel = 10
	hard.SetItself(hard)
	model, err := hard.Save()

	fmt.Println("Kezban || Model:", model, "Error:", err)

	lovely := new(models.Kezban)

	lovely.Name = "Ceyda"
	lovely.GrumpyLevel = 2
	t2, err := time.Parse(layout, "1990-10-05")
	lovely.Birthdate = t2
	lovely.SetItself(lovely)

	model, err = lovely.Save()
	fmt.Println("Kezban || Model:", model, "Error:", err)

	fmt.Println("########### Listing\n--------------")
	var kezbans []models.Kezban
	m := new(models.Kezban)
	m.SetItself(m)
	m.FindAll(kezban.KezQu{}, &kezbans)
	fmt.Println("Kezbans || ", kezbans)

	fetchOne := new(models.Kezban)
	fetchOne.Name = "Kezban"
	fetchOne.SetItself(fetchOne)

	fetchOne.FindOne(fetchOne, fetchOne)
	fmt.Println("Fetched Kezban || ", fetchOne)

	newK := new(models.Kezban)
	newK.Name = "Kezban"
	newK.SetItself(newK)
	if _, err := newK.Save(); err != nil {
		fmt.Println(err.Error())
		if err.Error() == "unique field duplicate" {
			// Give some error message to user bla bla bla
		}
	}
}