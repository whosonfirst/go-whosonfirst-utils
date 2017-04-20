package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func HashId(wofid int64, sources map[string]string) (map[string]string, error) {

	hashes := make(map[string]string)

	rel_path, err := uri.Id2RelPath(int(wofid))	// OH GOD FIX ME...

	if err != nil {
		return hashes, err
	}

	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)

	for src, root := range sources {

		wg.Add(1)

		go func(src string, root string, rel_path string, wg *sync.WaitGroup, mu *sync.Mutex) {

			defer wg.Done()

			hash, err := HashRecord(root, rel_path)

			mu.Lock()

			if err != nil {
				log.Println(src, err)
				hashes[src] = ""
			} else {
				hashes[src] = hash
			}

			mu.Unlock()

		}(src, root, rel_path, wg, mu)
	}

	wg.Wait()

	return hashes, nil
}

func HashRecord(root string, rel_path string) (string, error) {

	uri := root + rel_path

	rsp, err := http.Get(uri)

	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return "", err
	}

	var stub interface{}

	err = json.Unmarshal(body, &stub)

	if err != nil {
		return "", err
	}

	body, err = json.Marshal(stub)

	if err != nil {
		return "", err
	}

	hash := md5.Sum(body)
	str_hash := hex.EncodeToString(hash[:])

	return str_hash, nil
}

func CompareHashes(hashes map[string]string) bool {

	last := "-"

	for _, hash := range hashes {

		if last != "-" && last != hash {
			return false
		}

		last = hash
	}

	return true
}

func main() {

	// wofid := flag.Int64("wofid", 0, "A valid WOF ID")

	flag.Parse()

	wofids := make([]int64, 0)

	args := flag.Args()

	for _, id := range args {

		wofid, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		wofids = append(wofids, wofid)
	}

	if len(wofids) == 0 {
		log.Fatal("Missing WOF ID")
	}

	sources := map[string]string{
		"wof":    "https://whosonfirst.mapzen.com/data/",
		"github": "https://raw.githubusercontent.com/whosonfirst-data/whosonfirst-data/master/data/",
		"s3":     "https://s3.amazonaws.com/whosonfirst.mapzen.com/data/",
	}

	var writer *csv.DictWriter

	for _, wofid := range wofids {

		hashes, err := HashId(wofid, sources)

		if err != nil {
			log.Fatal(err)
		}

		match := "MATCH"

		if !CompareHashes(hashes) {
			match = "MISMATCH"
		}

		out := map[string]string{
			"wofid": strconv.FormatInt(wofid, 10),
			"match": match,
		}

		for src, hash := range hashes {
			out[src] = hash
		}

		if writer == nil {

			fieldnames := make([]string, 0)

			for k, _ := range out {
				fieldnames = append(fieldnames, k)
			}

			writer, err = csv.NewDictWriter(os.Stdout, fieldnames)

			if err != nil {
				log.Fatal(err)
			}

			writer.WriteHeader()
		}

		writer.WriteRow(out)
	}

	os.Exit(0)
}
