package dialoger

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/naoina/toml"
	"github.com/qeesung/image2ascii/convert"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
)

type DialogManager struct {
	Dialogs          map[string]Dialog
	CharacterImgPath string
	CharacterName    string
}

type Dialog struct {
	Name   string
	Scenes []Scene
}

type Scene struct {
	Id            int
	PersonName    string
	PersonImgPath string
	Phase         string
	AnswerOne     Answer
	AnswerTwo     Answer
	AnswerThree   Answer
	AnswerFour    Answer
}

type Answer struct {
	Phase   string
	SceneId int
}

func NewDialoger(dialogDirPath string, characterImgPath string, characterName string) DialogManager {
	dialogs := make(map[string]Dialog)

	files, err := ioutil.ReadDir(dialogDirPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		var dialog Dialog
		f, err := os.Open(fmt.Sprintf("%s/%s", dialogDirPath, file.Name()))
		if err != nil {
			panic(err)
		}

		if err := toml.NewDecoder(f).Decode(&dialog); err != nil {
			panic(err)
		}

		dialogs[dialog.Name] = dialog

		err = f.Close()
		if err != nil {
			panic(err)
		}
	}

	return DialogManager{Dialogs: dialogs, CharacterImgPath: characterImgPath, CharacterName: characterName}
}

func (d DialogManager) StartDialog(dialogName string, callback func()) {

	dialog, ok := d.Dialogs[dialogName]
	if ok {
		d.Render(dialog.Scenes, 0)
	} else {
		panic("Диалог не найден")
	}

	callback()
}

func (d DialogManager) Render(scenes []Scene, sceneId int) {
	scene := scenes[sceneId]
	characterImgAscii := imageToAscii(d.CharacterImgPath)
	personImgAscii := imageToAscii(scene.PersonImgPath)

	app := tview.NewApplication()

	newScene := func(sceneId int) {
		if sceneId > 0 {
			d.Render(scenes, sceneId)
		} else {
			app.Stop()
		}
	}

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	list := tview.NewList()
	list.AddItem(scene.AnswerOne.Phase, "", '1', func() { newScene(scene.AnswerOne.SceneId) })
	list.AddItem(scene.AnswerTwo.Phase, "", '2', func() { newScene(scene.AnswerTwo.SceneId) })
	list.AddItem(scene.AnswerThree.Phase, "", '3', func() { newScene(scene.AnswerThree.SceneId) })
	list.AddItem(scene.AnswerFour.Phase, "", '4', func() { newScene(scene.AnswerFour.SceneId) })

	list.SetShortcutColor(tcell.Color18)
	list.SetBackgroundColor(tcell.ColorDimGray)

	grid := tview.NewGrid().
		SetRows(20, 1, 7).
		SetColumns(10, 0, 10).
		SetBorders(true)

	grid.AddItem(newPrimitive(fmt.Sprintf(characterImgAscii)), 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(newPrimitive(personImgAscii), 0, 2, 1, 2, 0, 0, false)

	grid.AddItem(newPrimitive(d.CharacterName), 1, 0, 1, 2, 0, 0, false)
	grid.AddItem(newPrimitive(scene.PersonName), 1, 2, 1, 2, 0, 0, false)

	grid.AddItem(list, 2, 0, 1, 2, 0, 0, true)
	grid.AddItem(newPrimitive(scene.Phase), 2, 2, 1, 2, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}

func imageToAscii(imgPath string) string {
	convertOptions := convert.DefaultOptions
	convertOptions.Ratio = 5
	convertOptions.FitScreen = true
	convertOptions.FixedWidth = 65
	convertOptions.FixedHeight = 23
	convertOptions.Colored = false

	converter := convert.NewImageConverter()
	return converter.ImageFile2ASCIIString(imgPath, &convertOptions)
}
