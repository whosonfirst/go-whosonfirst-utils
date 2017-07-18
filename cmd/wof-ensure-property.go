package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func main() {

	repo := flag.String("repo", "/usr/local/data/whosonfirst-data", "The WOF repo whose files you want to test.")
	procs := flag.Int("processes", runtime.NumCPU()*2, "The number of concurrent processes to use")
	prop := flag.String("property", "", "The dotted notation for the property whose existence you want to test.")
	proptype := flag.String("type", "", "IF a property exists it MUST be the specified type. Valid: 'string', 'number', 'boolean', 'null'. If set, missing properties are ignored.")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	_, err := os.Stat(*repo)

	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	data := filepath.Join(*repo, "data")

	_, err = os.Stat(data)

	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	if *prop == "" {
		log.Fatal("You forgot to specify anything to check for.")
	}

	fieldnames := []string{"id", "path", "details"}
	writer, err := csv.NewDictWriter(os.Stdout, fieldnames)

	writer.WriteHeader()

	mu := new(sync.Mutex)

	callback := func(path string, info os.FileInfo) error {

		var details = ""

		if info.IsDir() {
			return nil
		}

		is_wof, err := uri.IsWOFFile(path)

		if err != nil {
			return err
		}

		if !is_wof {
			return nil
		}

		is_alt, err := uri.IsAltFile(path)

		if err != nil {
			return err
		}

		if is_alt {
			return nil
		}

		fh, err := os.Open(path)
		defer fh.Close()

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		var jpath string

		if strings.HasPrefix(*prop, "properties") {
			jpath = *prop
		} else {
			jpath = fmt.Sprintf("properties.%s", *prop)
		}

		result := gjson.GetBytes(body, jpath)

		if *proptype != "" {

			if result.Exists() != true {
				return nil
			} else if *proptype == "string" && result.Type.String() == "String" {
				return nil
			} else if *proptype == "number" && result.Type.String() == "Number" {
				return nil
			} else if *proptype == "boolean" && (result.Type.String() == "False" || result.Type.String() == "True") {
				return nil
			} else if *proptype == "null" && result.Type.String() == "Null" {
				return nil
			} else if *proptype == "object" && result.String()[0:1] == "{" {
				return nil
			} else if *proptype == "array" && result.String()[0:1] == "[" {
				return nil
			}

			resulttype := result.Type.String()
			if resulttype == "JSON" {
				if result.String()[0:1] == "{" {
					resulttype = "object"
				} else if result.String()[0:1] == "[" {
					resulttype = "array"
				}
			}

			details = "unexpected type '%s': " + resulttype
		} else if result.Exists() {
			return nil
		} else {
			details = "missing '%s'"
		}

		id, err := uri.IdFromPath(path)

		if err != nil {
			return err
		}

		str_id := strconv.FormatInt(id, 10)

		details = fmt.Sprintf(details, jpath)

		mu.Lock()
		defer mu.Unlock()

		row := make(map[string]string)
		row["id"] = str_id
		row["path"] = path
		row["details"] = details

		writer.WriteRow(row)
		return nil
	}

	cr := crawl.NewCrawler(data)
	cr.Crawl(callback)
}
