package main

import "fmt"
import "github.com/Ho7/dialoger"

type State struct {
	actionOneDid  bool
	actionTwoDid  bool
	dialogManager dialoger.DialogManager
}

func main() {
	d := dialoger.NewDialoger("/Users/nyuokunev/go/src/Dialoger/example/dialogs", "/Users/nyuokunev/go/src/Dialoger/example/k.png", "Uno")
	s := State{false, false, d}
	s.Loop()
}

func (s State) Loop() {
	if !s.actionOneDid {
		s.actionOneDid = true
		s.dialogManager.StartDialog("pika_dialog1", func() { s.Loop() })
	} else if !s.actionTwoDid {
		s.actionTwoDid = true
		s.a1()
	}

	fmt.Println("Game over")
}

func (s State) a1() {
	fmt.Println("Action one")
}
