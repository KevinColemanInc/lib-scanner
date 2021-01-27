package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"fmt"
)

const tmpDir = "/tmp/ruby-gems-crawl/tmp"

/*
 Step 2 - download and extract the gems in latest_gems.csv
*/

type gemStruct struct {
	name, version string
}

func (m gemStruct) folder_name() string {
	return m.name + "-" + m.version
}

func (m gemStruct) file_name() string {
	return m.folder_name() + ".gem"
}

func unpackGem(job *gemStruct) {
	if _, err := os.Stat(tmpDir + "/gems/" + job.folder_name()); !os.IsNotExist(err) {
		// fmt.Println("skipping", tmpDir+"/gems/"+job.folder_name())
		return
	}
	os.Mkdir(tmpDir+"/uncompressed/"+job.folder_name(), 0755)
	cmd := exec.Command("tar", "xf", job.file_name(), "-C", "./uncompressed/"+job.folder_name())
	cmd.Dir = tmpDir
	cmd.Run()
	os.Mkdir(tmpDir+"/gems/"+job.folder_name(), 0755)
	uncomCmd := exec.Command("tar", "xzf", tmpDir+"/uncompressed/"+job.folder_name()+"/data.tar.gz", "-C", tmpDir+"/gems/"+job.folder_name())
	uncomCmd.Run()
}

func gemFileGems() (<-chan *gemStruct, <-chan error) {
	ch := make(chan *gemStruct)
	errc := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errc)
		reader, err := os.Open("latest_gems.csv")
		if err != nil {
			errc <- err
		}
		defer reader.Close()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			txt := scanner.Text()
			gem_name := strings.Split(txt, ",")

			if len(gem_name) == 2 {
					ch <- &gemStruct{name: gem_name[0], version: gem_name[1]}
			}
		}
	}()

	return ch, errc
}

func fetchGem(job *gemStruct) {
	if _, err := os.Stat(tmpDir + "/" + job.file_name()); !os.IsNotExist(err) {
		return
	}

	DownloadFile(tmpDir+"/"+job.file_name(), "https://rubygems.org/gems/"+job.file_name())
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func process(job *gemStruct) {
	fetchGem(job)
	unpackGem(job)
}

func digester(ch <-chan *gemStruct) {
	for j := range ch {
		process(j)
	}
}

func main() {
	jobs, errc := gemFileGems()

	var wg sync.WaitGroup
	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(jobs) // HLc
			wg.Done()
		}()
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		fmt.Println("failure", err)
	}
	wg.Wait()
}
