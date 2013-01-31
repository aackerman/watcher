package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"flag"
	"time"
	"os/exec"
)

type WatchFile struct {
	Name string
	LastModified time.Time
}

func (wf *WatchFile) watch(c string) {
	wf.updateLastModified()
	for {
		if wf.lastModifiedChanged() {
			wf.updateLastModified()
			wf.exec(c)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (wf *WatchFile) mtime() (time.Time) {
	fileinfo, err := os.Stat(wf.Name);
	if err != nil {
		log.Fatal(err)
	}
	return fileinfo.ModTime()
}

func (wf *WatchFile) updateLastModified() {
	wf.LastModified = wf.mtime()
}

func (wf *WatchFile) lastModifiedChanged() (bool) {
	if wf.mtime().Equal(wf.LastModified) {
		return false
	}
	return true
}

func (wf *WatchFile) exec(c string) {
	cmd := exec.Command(c, wf.Name)
	stdout, err := cmd.StdoutPipe();

	if err != nil {
		log.Fatal(err)
		return
	}

	stderr, err := cmd.StderrPipe();

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return
	}

	go func(){
		io.Copy(os.Stdout, stdout)
		io.Copy(os.Stderr, stderr)
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	command := flag.String("c", "false", "input command to be run on each file")
	flag.Parse()
	files := flag.Args()

	if len(files) == 0 {
		fmt.Println("Please specify a file or list of files")
		return
	}

	if (*command == "false") {
		fmt.Println("Please specify a command")
		return
	}

	for _, file := range files {
		watchfile := new(WatchFile)
		watchfile.Name = file
		go watchfile.watch(*command)
	}
	<-make(chan bool)
}
