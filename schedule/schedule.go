package schedule

import (
	"encoding/json"
	"errors"
	"fmt"

	"../disk"
	"../event"
	"../myTime"
)

//Schedule contiene tutti gli eventi su zoom della settimana
type Schedule = *schedule
type schedule struct {
	sche *[]event.Event
	mdf  bool
	mem  disk.Disk
}

//New crea una programmazione
func New() Schedule {
	var S schedule
	var sch []event.Event

	S.sche = &sch
	S.mem = disk.Open()
	S.mdf = false

	return &S
}

//Load carica lo schedule
func (S Schedule) Load() error {

	iter := S.mem.NewIter()

	for iter.HasNext() {
		xml, err := iter.Next()
		if err != nil {
			return err
		}

		E, err := event.XmlToEvent(xml)
		if err != nil {
			return err
		}

		err = S.add(E)
		if err != nil {
			return err
		}
	}

	return nil
}

//AddEvent aggiunge un elemento allo schedule
func (P Schedule) add(E event.Event) error {
	S := *P.sche

	if !E.IsClosed() {
		return errors.New("event must be closed")
	}

	tmp := append(S, E)
	P.sche = &tmp

	return nil
}

//AddEvent aggiunge un elemento allo schedule
func (P Schedule) AddEvent(E event.Event) error {

	/*if len(*P.sche) < 1 {
		P.mem.DellAll()
	}*/

	P.add(E)
	xml, _ := E.Xml()

	P.mem.Write(xml)

	return nil
}

//DellEvent cancella un evento passando l'indice
func (P Schedule) DellEvent(ind int) error {
	S := *P.sche

	if ind > len(S)-1 {
		return errors.New("invalid index")
	}
	S[ind] = nil
	P.mdf = true

	return nil
}

//GetEvent retituisce ID e Pass del'evento più vicino a ora
func (P Schedule) GetEvent() (event.Event, int) {
	S := *P.sche

	today := myTime.Now()
	for _, E := range S {
		if dis := myTime.Distance(today, E.GetDate()); dis > -45 && dis < 15 {
			return E, dis
		}
	}
	return nil, 0
}

//Clear "cancella" tutti i dati in uno schedule
func (P Schedule) Clear() {
	P.mem.DellAll()
}

//
/*################ EXPORT ################*/

//ListEvent elenca gli eventi presenti nello schedule
func (P Schedule) ListEvent() string {
	S := *P.sche
	l := ""

	for i, r := range S {
		ts := r.GetTS()
		l += fmt.Sprintln(i, ") ", ts, " - ", r.GetName())
	}

	if l == "" {
		l = "Nessun evento da mostrare"
	}
	return l
}

//Save salva il calendrio
func (P Schedule) Save() error {

	if P.mdf {
		return P.mem.WriteAll(P.Xml())
	}
	return nil
}

func (P Schedule) Close() {
	P.mem.Close()
}

//Xml esporta l'oggetto sul quale è stato invoato in xml
func (P Schedule) Xml() []string {
	S := *P.sche
	var xml []string

	for _, E := range S {
		if E != nil {
			evXml, _ := E.Xml()
			xml = append(xml, evXml)
		}
	}
	return xml
}

//XmlToSchedule legge l'xml e crea l'oggetto
func XmlToSchedule(xmlRow string) (Schedule, error) {
	S := New()

	var ex []string
	if err := json.Unmarshal([]byte(xmlRow), &ex); err != nil {
		return nil, err
	}

	for _, evXml := range ex {
		E, err := event.XmlToEvent(evXml)
		if err != nil {
			return nil, err
		}
		S.AddEvent(E)
	}

	return S, nil
}
