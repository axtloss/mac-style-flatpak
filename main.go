package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gen2brain/beeep"
)

func symlink(filepath string, path string) {
	homedir, err := os.UserHomeDir()
	fmt.Println(homedir + "/.local/share/flatpak/exports/share/applications/" + strings.Split(strings.Replace(filepath, ".flatpakref", ".desktop", 1), "/")[1])
	cmd := exec.Command("ln", "-s", homedir+"/.local/share/flatpak/exports/share/applications/"+strings.Split(strings.Replace(filepath, ".flatpakref", ".desktop", 1), "/")[1], path)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(stdout))
}

func install(filepath string, path string) {
	cmd := exec.Command("flatpak", "install", "--user", "-y", filepath)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		if err := os.Remove("./" + filepath); err != nil {

			beeperr := beeep.Notify("Mac Style Install", strings.Split(strings.Replace(filepath, ".flatpakref", "", 1), "/")[1]+" failed to install!", "assets/notif.png")
			if beeperr != nil {
				panic(beeperr)
			}
			log.Fatal(err)
		}
		symlink(filepath, path)
		beeperr := beeep.Notify("Mac Style Install", strings.Split(strings.Replace(filepath, ".flatpakref", "", 1), "/")[1]+" installed succesfully", "assets/notif.png")
		if beeperr != nil {
			panic(beeperr)
		}
	}
	fmt.Println(string(stdout))

}

func main() {

	path := os.Getenv("APPLICATIONS_PATH")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if strings.Contains(event.String(), "CREATE") && strings.Contains(event.String(), ".flatpakref") {
					fmt.Println("event:", event)
					fmt.Println("filepath:", event.Name)
					install(event.Name, path)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
