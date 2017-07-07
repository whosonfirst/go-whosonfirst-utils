package main

import (
	"flag"
	"github.com/mmcloughlin/globe"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"image/color"
	"io"		
	"log"
	"strconv"
	"time"		
)

func main() {

	outfile := flag.String("out", "", "Where to write globe")
	size := flag.Int("size", 1600, "The size of the globe (in pixels)")
	// mode := flag.String("mode", "meta", "... (default is 'meta' for one or more meta files)")

	center_lat := flag.Float64("latitude", 37.755244, "")
	center_lon := flag.Float64("longitude", -122.447777, "")	
	
	flag.Parse()

	green := color.NRGBA{0x00, 0x64, 0x3c, 192}
	g := globe.New()
	g.DrawGraticule(10.0)

	t1 := time.Now()
	
	for _, path := range flag.Args() {

		reader, err := csv.NewDictReaderFromPath(path)

		if err != nil {
			log.Fatal(err)
		}

		for {
			row, err := reader.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			str_lat, ok := row["geom_latitude"]

			if !ok {
				continue
			}

			str_lon, ok := row["geom_longitude"]

			if !ok {
				continue
			}

			lat, err := strconv.ParseFloat(str_lat, 64)

			if err != nil {
				log.Println(err, str_lat)
				continue
			}

			lon, err := strconv.ParseFloat(str_lon, 64)
			
			if err != nil {
				log.Println(err, str_lon)
				continue
			}

			g.DrawDot(lat, lon, 0.01, globe.Color(green))
		}
	}

	t2 := time.Since(t1)

	log.Printf("time to read all the things %v\n", t2)

	t3 := time.Now()
	
	g.CenterOn(*center_lat, *center_lon)
	err := g.SavePNG(*outfile, *size)

	t4 := time.Since(t3)

	log.Printf("time to draw all the things %v\n", t4)
	
	if err != nil {
		log.Fatal(err)
	}

}
