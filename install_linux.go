package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Install(src, dest string) {
	destForPath := strings.Join(strings.Split(dest, string(os.PathSeparator))[:len(strings.Split(dest, string(os.PathSeparator)))-1], string(os.PathSeparator))

	homeDir, _ := os.UserHomeDir()

	bashAliases, _ := os.ReadFile(homeDir + string(os.PathSeparator) + ".bash_aliases")
	strBashAliases := string(bashAliases)
	if !strings.Contains(strBashAliases, destForPath) {
		strBashAliases += "\n\nPATH=\"$PATH:" + destForPath + "\"\n"
	}

	_ = os.WriteFile(homeDir+string(os.PathSeparator)+".bash_aliases", []byte(strBashAliases), 0777)

	_ = exec.Command("source", homeDir+string(os.PathSeparator)+".bash_aliases").Run()

	if err := exec.Command("bash", "-c", "mv \""+src+"\" \""+dest+"\"").Start(); err != nil {
		err = fmt.Errorf("erreur lors de l'exécution: %w", err)
	}
}
