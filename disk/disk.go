package disk

import (
	"errors"
	"os"
)

/*################ FUNZIONI GENERICHE ################*/

func itob(d int64) []byte {
	b := make([]byte, 8)
	b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7] = 0, 0, 0, 0, 0, 0, 0, 0

	for i := 7; d > 0; i-- {
		b[i] = byte(d % 16)
		d = d / 16
	}

	return b
}

func btoi(b []byte) int64 {
	var d int64
	for _, esa := range b {
		d = d*16 + int64(esa)
	}
	return d
}

/*################ DISK ################*/

type Disk = *disk
type disk struct {
	recNun  int64    //numeo di record presenti nel file
	fileDir string   //percorso del file
	filePnt *os.File //puntatore al file
}

//getOfset restituisce il peso del file
func (D Disk) getOfset() int64 {
	fi, _ := D.filePnt.Stat()
	return fi.Size()
}

//Open crea e dinizializza una nuova struttura disk
func Open() Disk {
	var D disk
	D.fileDir = os.Args[0]
	D.filePnt, _ = os.OpenFile(D.fileDir, os.O_APPEND|os.O_RDWR, 0770)

	header := make([]byte, 16)
	n, _ := D.filePnt.ReadAt(header, D.getOfset()-16)
	recnum := btoi(header[0:8])
	recnum++
	if recnum < 0 || recnum > 256 || n < 16 {
		D.DellAll()
		recnum = 0
	}
	D.recNun = recnum

	return &D
}

//Write aggiunge un novo record al file
func (D Disk) Write(S string) error {
	S += string(itob(D.recNun)) + string(itob(D.getOfset()-16))
	_, err := D.filePnt.Write([]byte(S))
	D.recNun++
	return err
}

//DellAll "cancella" tutti i record nel file
func (D Disk) DellAll() error {
	header := make([]byte, 8)
	header[0], header[1], header[2], header[3], header[4], header[5], header[6], header[7] = 15, 15, 15, 15, 15, 15, 15, 15
	_, err := D.filePnt.Write([]byte(string(header) + string(itob(D.getOfset()))))
	return err
}

//WriteAll aggiunge una serie di record al file "cancellando" quelli esistenti
func (D Disk) WriteAll(S []string) error {

	if err := D.DellAll(); err != nil {
		return err
	}

	D.recNun = 0

	for _, c := range S {
		if err := D.Write(c); err != nil {
			return err
		}
	}

	return nil
}

//Close chiude il file
func (D Disk) Close() {
	D.filePnt.Close()
}

/*################ DISK ITERATOR ################*/

type DiskIterator = *diskIterator
type diskIterator struct {
	actualOfset int64
	recordInx   int64
	filePointer *os.File
}

//readHeader legge l'header del record cosiderato
func (D DiskIterator) readHeader() (int64, int64) {

	header := make([]byte, 16)

	n, err := D.filePointer.ReadAt(header, D.actualOfset)
	if n < 16 || err != nil {
		return -1, -1
	}

	return btoi(header[8:]), btoi(header[0:8])
}

//readN legge n byte, restituisce errore se non riesce a leggerli tutti
func (D DiskIterator) readN(nByte int) ([]byte, error) {

	cnt := make([]byte, nByte)
	n, err := D.filePointer.ReadAt(cnt, D.actualOfset+16)
	if err != nil {
		return cnt, err
	}
	if n < nByte {
		return cnt, errors.New("non tutti i byte sono stati letti")
	}
	return cnt, nil
}

//NewIter crea un nuovo iteratore di disk
func (D Disk) NewIter() DiskIterator {
	var I diskIterator

	fi, _ := D.filePnt.Stat()

	I.actualOfset = fi.Size() - 16
	I.filePointer = D.filePnt
	I.recordInx = D.recNun - 1

	return &I
}

//HasNext restituisce un buleano rappresentativo della presenza di un prossimo record
func (D DiskIterator) HasNext() bool {
	return D.recordInx > 0
}

//Next restituisce un nuovo record
func (D DiskIterator) Next() (string, error) {

	of, ind := D.readHeader()
	if of < 0 || ind > D.recordInx || ind < 0 {
		return "", errors.New("errore nella lettura del heder")
	}
	len := D.actualOfset - of
	D.actualOfset = of
	D.recordInx = ind
	cnt, err := D.readN(int(len))
	if err != nil {
		return "", errors.New("errore nella lettura del record")
	}

	return string(cnt), nil
}
