package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Install(src, dest string) {
	destForPath := strings.Join(strings.Split(dest, string(os.PathSeparator))[:len(strings.Split(dest, string(os.PathSeparator)))-1], string(os.PathSeparator))

	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.READ|registry.WRITE)
	if err != nil {
		panic(err)
	}
	defer func(k registry.Key) {
		err := k.Close()
		if err != nil {
			panic(err)
		}
	}(k)

	// Lire la valeur Path existante
	pathValue, _, err := k.GetStringValue("Path")
	if err != nil {
		// Si Path n'existe pas encore, on l'initialise vide
		pathValue = ""
	}

	// Vérifier que le chemin n'existe pas déjà
	paths := strings.Split(pathValue, ";")
	for _, p := range paths {
		if strings.EqualFold(p, destForPath) {
			fmt.Println("Le chemin existe déjà dans Path utilisateur.")
			return
		}
	}

	// Ajouter le nouveau chemin
	updatedPath := pathValue
	if len(updatedPath) > 0 && !strings.HasSuffix(updatedPath, ";") {
		updatedPath += ";"
	}
	updatedPath += destForPath

	// Ecrire la nouvelle valeur dans le registre
	err = k.SetStringValue("Path", updatedPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Le chemin a été ajouté avec succès au Path utilisateur.")

	if err := exec.Command("powershell", "-c", "Move-Item", "-Path", "\""+src+"\"", "-Destination", "\""+dest+"\"", "-Force").Start(); err != nil {
		err = fmt.Errorf("erreur lors de l'exécution: %w", err)
	}
}
