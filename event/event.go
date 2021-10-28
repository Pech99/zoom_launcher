package event

import (
	"encoding/xml"
	"errors"
	"fmt"

	"../myTime"
)

//Event è un singolo evento in calendario
type Event = *event
type event struct {
	id     string
	psw    string
	name   string
	day    int
	hour   int
	min    int
	closed bool
}
type export struct {
	ID   string
	Psw  string
	Name string
	Day  int
	Hour int
	Min  int
}

/*################ CREAZIONE ################*/

//New crea un nuovo eveto
func New() Event {
	var E event
	E.day, E.hour, E.min = -1, -1, -1
	return &E
}

//SetID imposta l'ID
func (E Event) SetID(ID string) Event {
	if !E.closed {
		E.id = ID
	}
	return E
}

//SetPsw imposta la password (può essere anche vuto)
func (E Event) SetPsw(Psw string) Event {
	if !E.closed {
		E.psw = Psw
	}
	return E
}

//SetName imposta il nome associato alla riunione
func (E Event) SetName(Name string) Event {
	if !E.closed {
		E.name = Name
	}
	return E
}

//SetDate imposta la data e l'ora della riunione
func (E Event) SetDate(D, H, M int) Event {
	if !E.closed {
		E.day = D
		E.hour = H
		E.min = M
	}
	return E
}

//Build crea la riunione controllando che sia tutto ok
func (E Event) Build() (Event, error) {
	if E.day < 0 || E.hour < 0 || E.min < 0 {
		return nil, errors.New("date mast be set")
	}
	if E.id == "" {
		return nil, errors.New("ID mast be set")
	}
	if E.name == "" {
		return nil, errors.New("name mast be set")
	}

	E.closed = true
	return E, nil
}

/*################ INTERROGAZIONE ################*/

//GetID restituisce l'id ella riunuione in questione
func (E Event) GetID() string {
	return E.id
}

//GetPsw restituisce la password della riunione in qestione
func (E Event) GetPsw() string {
	return E.psw
}

//GetName restituisce il nome della riunione in questione
func (E Event) GetName() string {
	return E.name
}

//GetDate restituisce la data della riunione in questione
func (E Event) GetDate() int {
	return myTime.DayToInt(E.day, E.hour, E.min)
}

//GetTS restituisce il time stamp della riunione
func (E Event) GetTS() string {
	return fmt.Sprintf("%d:%2.d:%2.d", E.day, E.hour, E.min)
}

//IsClosed valuta se l'evento sul quale è invocato è stato buildato o meno
func (E Event) IsClosed() bool {
	return E.closed
}

/*################ EXPORT ################*/

//Xml esporta l'oggetto sul quale è stato invoato in xml
func (E Event) Xml() (string, error) {

	if E == nil || !E.closed {
		return "", errors.New("invalid event")
	}

	var ex export
	ex.Name = E.name
	ex.ID = E.id
	ex.Psw = E.psw
	ex.Day = E.day
	ex.Hour = E.hour
	ex.Min = E.min
	S, err := xml.Marshal(ex)
	return string(S) + "\n", err
}

//XmlToEvent passatol'xml ritorna l'oggetto
func XmlToEvent(xmlRow string) (Event, error) {

	var ex export
	if err := xml.Unmarshal([]byte(xmlRow), &ex); err != nil {
		return nil, err
	}
	E, err := New().SetName(ex.Name).SetID(ex.ID).SetPsw(ex.Psw).SetDate(ex.Day, ex.Hour, ex.Min).Build()

	return E, err
}
