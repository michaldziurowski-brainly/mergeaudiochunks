package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

/*
port from
https://github.com/aws-samples/amazon-chime-media-capture-pipeline-demo/blob/main/src/processLambda/app/app.py
*/
func main() {
	ffmpegCmd := "ffmpeg"
	// ReadDir sorts files by name by default
	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		log.Fatalf("error while reading filest from dir: %s", err.Error())
	}

	f, err := os.OpenFile("chunk_list.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	tmpFileNames := []string{}
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(fileName string) {
			// filename without extension
			basename := fileName[:len(fileName)-len(filepath.Ext(fileName))]
			cmdArgs := fmt.Sprintf("-i ./data/%s -bsf:v h264_mp4toannexb -f mpegts -framerate 15 -c copy ./tmp_data/%s.ts -y", fileName, basename)
			runCmd(ffmpegCmd, strings.Fields(cmdArgs)...)

			fmt.Printf("processed %s\n", basename)

			tmpFileNames = append(tmpFileNames, fmt.Sprintf("%s.ts", basename))

			wg.Done()
		}(file.Name())
	}

	wg.Wait()

	if len(tmpFileNames) == len(files) {
		// sort in correct order (files are named by date so ascending order is ok)
		sort.Strings(tmpFileNames)

		for _, tmpFileName := range tmpFileNames {
			if _, err := fmt.Fprintf(f, "file ./tmp_data/%s\n", tmpFileName); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("written to file %s\n", tmpFileName)
		}

		cmdArgs := "-f concat -safe 0 -i ./chunk_list.txt  -c copy ./merged.mp4 -y"

		runCmd(ffmpegCmd, strings.Fields(cmdArgs)...)
	} else {
		log.Fatal("did not create enough tmp ts files")
	}

	fmt.Println("done")
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("error occured while executing worker: %v\n", err)
	}
}
