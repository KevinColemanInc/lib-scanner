/*
 */

package handle

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"github.com/KevinColemanInc/lib-crawl/src/report"
	"sync"

	"github.com/tucnak/climax"
)

const TARGET_EXT = ".rb"
const NUMBER_OF_DIGESTERS = 10000

func digester(done <-chan struct{}, paths <-chan string, c chan<- report.Warning) {
	for path := range paths { // HLpaths
		RubyScan(done, path, c)
	}
}

func Scan(ctx climax.Context) int {
	warnings, err := scanAll(ctx.Args[0])
	if err != nil {
		fmt.Println("Unexpected error", err)
	}
	for _, w := range warnings {
		w.ToCLI("verbose")
	}
	return 0
}

func scanAll(root string) ([]report.Warning, error) {
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan report.Warning) // HLc
	var wg sync.WaitGroup
	wg.Add(NUMBER_OF_DIGESTERS)
	for i := 0; i < NUMBER_OF_DIGESTERS; i++ {
		go func() {
			digester(done, paths, c) // HLc
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c) // HLc
	}()
	// End of pipeline. OMIT

	warnings := make([]report.Warning, 0)
	for w := range c {
		if w.Err == nil {
			warnings = append(warnings, w)
		}
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return nil, err
	}
	return warnings, nil

}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)
		counter := 0
		errc <- filepath.Walk(root,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() || len(info.Name()) < 4 || filepath.Ext(info.Name()) != TARGET_EXT {
					return nil
				}

				select {
				case paths <- path: // HL
					counter++
				case <-done: // HL
					return errors.New("walk canceled")
				}
				return nil
			})
		fmt.Printf("Found paths: %v\n", counter)
	}()
	return paths, errc
}
