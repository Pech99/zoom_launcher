package zoomLib

import (
	"net/url"
	"os"
	"os/exec"
	"path"
)

//Comando Passato ID e Pass crea il comando per eseguire zoom
func Comando(ID, Pass string) *exec.Cmd {
	zoomD, _ := os.UserHomeDir()
	zoomD += `\AppData\Roaming\Zoom\bin\Zoom.exe`

	if Pass == "" {
		return exec.Command(zoomD, `--url=zoommtg://zoom.us/join?confno=`+ID, `Disabilitato`)
	}
	return exec.Command(zoomD, `--url=zoommtg://zoom.us/join?confno=`+ID+`&pwd=`+Pass, `Disabilitato`)
}

//Parse fa il parsing di un url estraendo ID e Pass
func Parse(link string) (string, string) {
	u, err := url.Parse(link)
	if err != nil {
		return "", ""
	}

	pw := u.Query()["pwd"]
	if len(pw) > 0 {
		return path.Base(u.Path), pw[0]
	}

	return path.Base(u.Path), ""
}
