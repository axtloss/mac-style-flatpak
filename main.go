package main

import (
	"errors"
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

func uninstall(filepath string, path string) {
	cmd := exec.Command("flatpak", "remove", "-y", strings.Split(strings.Replace(filepath, ".desktop", "", 1), "/")[1])
	stdout, err := cmd.Output()
	if err != nil {
		beeeper := beeep.Notify("Error", "Unable to uninstall application "+strings.Split(strings.Replace(filepath, ".desktop", "", 1), "/")[1], "assets/notif.png")
		if beeeper != nil {
			panic(beeeper)
		}
	} else {
		beeper := beeep.Notify("Mac Style Uninstall", strings.Split(strings.Replace(filepath, ".desktop", "", 1), "/")[1]+" uninstalled succesfully", "assets/notif.png")
		if beeper != nil {
			panic(beeper)
		}
	}
	fmt.Println(string(stdout))
}

func install(filepath string, path string) {
	cmd := exec.Command("flatpak", "install", "--user", "-y", filepath)
	stdout, err := cmd.Output()
	if err != nil {
		beeperr := beeep.Notify("Mac Style Install", strings.Split(strings.Replace(filepath, ".flatpakref", "", 1), "/")[1]+" failed to install!", "assets/notif.png")
		if beeperr != nil {
			panic(beeperr)
		}
	} else {
		symlink(filepath, path)

		symlink(filepath, path)

		if _, err := os.Stat(filepath); err == nil {
			if err := os.Remove(filepath); err != nil {
				log.Fatal(err)
			}
			beeperr := beeep.Notify("Mac Style Install", strings.Split(strings.Replace(filepath, ".flatpakref", "", 1), "/")[1]+" installed succesfully", "assets/notif.png")
			if beeperr != nil {
				panic(beeperr)
			}
		} else if errors.Is(err, os.ErrNotExist) {
			beeperr := beeep.Notify("Mac Style Install", strings.Split(strings.Replace(filepath, ".flatpakref", "", 1), "/")[1]+" installed succesfully", "assets/notif.png")
			if beeperr != nil {
				panic(beeperr)
			}
		} else {
			if err := os.Remove(filepath); err != nil {
				log.Fatal(err)
			}
			beeperr := beeep.Notify("Mac Style Install", strings.Split(strings.Replace(filepath, ".flatpakref", "", 1), "/")[1]+" installed succesfully", "assets/notif.png")
			if beeperr != nil {
				panic(beeperr)
			}
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
				} else if strings.Contains(event.String(), "RENAME") && strings.Contains(event.String(), ".desktop") {
					fmt.Println("event:", event)
					fmt.Println("filepath:", event.Name)
					uninstall(event.Name, path)
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
