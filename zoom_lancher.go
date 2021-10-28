package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"./event"
	"./schedule"
	"./zoomLib"
	"github.com/AllenDang/w32"
)

const vers string = "\tZoom Launcher\n\t\tVersion:1.3.3\n\t\tThe Pecce"

func parsTS(ts string) (int, int, int) {
	var d, h, m int
	var err error

	cmp := strings.Split(ts, ":")
	if len(cmp) != 3 {
		return -1, -1, -1
	}

	if d, err = strconv.Atoi(cmp[0]); err != nil {
		d = -1
	}
	if h, err = strconv.Atoi(cmp[1]); err != nil {
		d = -1
	}
	if m, err = strconv.Atoi(cmp[2]); err != nil {
		d = -1
	}

	return d, h, m
}

func loop(S schedule.Schedule) {
	if Ev, Dis := S.GetEvent(); Ev != nil {
		if dialog(Ev.GetName(), Dis) {
			zoomLib.Comando(Ev.GetID(), Ev.GetPsw()).Start()
		}
	}
	time.Sleep(10 * time.Minute)
}

//mostra una finestra di dialogo per informare l'utente dell'apertura di zoom
func dialog(evName string, min int) bool {
	var ric string
	if min >= 0 {
		ric = "L'evento " + evName + " inizierà tra " + fmt.Sprint(min) + " minuti\nVoi avviare zoom?"
	} else {
		ric = "L'evento " + evName + " è iniziato da " + fmt.Sprint(-min) + " minuti\nVoi avviare zoom?"
	}
	// si = 6; no = 7
	return w32.MessageBox(0, ric, "Zoom Launcher", 4) == 6
}

//HideConsole Nasconde la console
func HideConsole(Acc bool) {
	console := w32.GetConsoleWindow()
	if console == 0 {
		return
	}

	if Acc {
		w32.ShowWindow(console, 0)
	} else {
		w32.ShowWindow(console, 1)
	}
}

func scan() []string {

	var arg string
	var args []string
	var quot bool = false

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	txt := scanner.Text()

	for _, r := range txt {
		if !quot && (r == ' ' || r == '	') {
			args = append(args, arg)
			arg = ""
		} else if r == '"' {
			quot = !quot
		} else {
			arg += string(r)
		}
	}

	return append(args, arg)
}

func cmd(S schedule.Schedule) {

	for {
		cmd := scan()

		if len(cmd) == 1 && cmd[0] == "loop" { //Loop
			loop(S)

		} else if len(cmd) == 4 && cmd[0] == "add" { //add [nome evento] [G:HH:MM] [Link Evento]
			ID, Psw := zoomLib.Parse(cmd[3])
			if E, err := event.New().SetName(cmd[1]).SetID(ID).SetPsw(Psw).SetDate(parsTS(cmd[2])).Build(); err != nil {
				fmt.Println(err)
			} else {
				S.AddEvent(E)
			}

		} else if len(cmd) == 1 && cmd[0] == "list" { //list
			fmt.Println(S.ListEvent())

		} else if len(cmd) == 2 && cmd[0] == "dell" { //dell [# evento]
			indx, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println(err)
			} else {
				S.DellEvent(indx)
			}

		} else if cmd[0] == "clear" { //clear
			S.Clear()
		} else if len(cmd) == 2 && cmd[0] == "start" { //start [link]
			zoomLib.Comando(zoomLib.Parse(cmd[1])).Start()
		} else if len(cmd) == 2 && cmd[0] == "save" { //save
			S.Save()
		} else if len(cmd) == 1 && cmd[0] == "version" { //version
			fmt.Println(vers)
		}

	}

}

func main() {

	if len(os.Args) < 2 || os.Args[1] != "w" {
		HideConsole(true)
	}

	S := schedule.New()
	defer S.Close()
	err := S.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	go cmd(S)

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "loop" { //Loop
			loop(S)

		} else if len(os.Args) > i+3 && os.Args[i] == "add" { //add [nome evento] [G:HH:MM] [Link Evento]
			ID, Psw := zoomLib.Parse(os.Args[i+3])
			E, _ := event.New().SetName(os.Args[i+1]).SetID(ID).SetPsw(Psw).SetDate(parsTS(os.Args[i+2])).Build()
			S.AddEvent(E)
			i += 3

		} else if os.Args[i] == "list" { //list
			fmt.Println(S.ListEvent())

		} else if len(os.Args) > i+1 && os.Args[i] == "dell" { //dell [# evento]
			indx, err := strconv.Atoi(os.Args[i+1])
			if err != nil {
				fmt.Println(err)
			} else {
				S.DellEvent(indx)
				S.Save()
				i += 1
			}

		} else if os.Args[i] == "start" { //start
			zoomLib.Comando(zoomLib.Parse(os.Args[i])).Start()
			return
		} else if os.Args[i] == "clear" { //clear
			S.Clear()
			return
		} else if os.Args[i] == "version" { //version
			fmt.Println(vers)
		}

	}

	if len(os.Args) < 2 {
		if Ev, Dis := S.GetEvent(); Ev != nil {
			if dialog(Ev.GetName(), Dis) {
				zoomLib.Comando(Ev.GetID(), Ev.GetPsw()).Start()
			}
		} else {
			fmt.Println("Nulla da visualizzare")
			w32.MessageBox(0, "Nessuna riunione programmata nei prossimi 15 minuti", "Zoom Launcher", 0)
		}
	}

}
