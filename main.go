package Dialoger

import (
	"github.com/naoina/toml"
	"os"
)

type DialogManager struct {
	Dialogs TOML
}

type TOML struct {
	Title string
	Dialogs []Dialog
}

type Dialog struct {
	Id int
	PersonName string
	PersonImgPath string
	Phase string
	AnswerOne Answer
	AnswerTwo Answer
	AnswerThree Answer
	AnswerFour Answer
}

type Answer struct {
	Phase string
	DialogId int
}

func NewDialoger(file string) DialogManager{
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var dialogs TOML
	if err := toml.NewDecoder(f).Decode(&dialogs); err != nil {
		panic(err)
	}

	return DialogManager{Dialogs:dialogs}
}
