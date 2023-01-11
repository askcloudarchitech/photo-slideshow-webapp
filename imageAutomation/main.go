package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"

	fswatch "github.com/andreaskoch/go-fswatch"
)

var (
	sourceDir = getEnvDefault("SOURCE_DIR", "../IMAGES")
	dest      = getEnvDefault("DEST_URL", "http://localhost:8080/api/upload?p=defaultPassword")
)

func main() {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic("could not initialize cookiejar")
	}

	client := &http.Client{
		Jar: jar,
	}

	recurse := false // include all sub directories

	skipDotFilesAndFolders := func(path string) bool {
		return strings.HasPrefix(filepath.Base(path), ".")
	}

	checkIntervalInSeconds := 1

	folderWatcher := fswatch.NewFolderWatcher(
		sourceDir,
		recurse,
		skipDotFilesAndFolders,
		checkIntervalInSeconds,
	)

	folderWatcher.Start()

	for folderWatcher.IsRunning() {

		changes := <-folderWatcher.ChangeDetails()

		fmt.Printf("%s\n", changes.String())
		fmt.Printf("New: %#v\n", changes.New())
		fmt.Printf("Modified: %#v\n", changes.Modified())
		fmt.Printf("Moved: %#v\n", changes.Moved())

		for _, v := range changes.New() {
			go func(file string) {
				openFile, err := mustOpen(file)
				if err != nil {
					fmt.Println(err)
					return
				}
				values := map[string]io.Reader{
					"file": openFile,
				}
				err = Upload(client, dest, values)
				if err != nil {
					fmt.Println(err)
				}
			}(v)
		}

	}
}

func getEnvDefault(env string, defaultVal string) string {
	val := os.Getenv(env)
	if val == "" {
		return defaultVal
	}
	return val
}

// stolen from https://stackoverflow.com/questions/20205796/post-data-using-the-content-type-multipart-form-data
func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("bad status: %s", res.Status)
	} else {
		fmt.Printf("success response: %s", res.Status)
	}
	return
}

func mustOpen(f string) (*os.File, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return r, nil
}
