package main

import (
	"errors"
	"flag"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/globe"			// for to make DrawPreparedPaths public
	"github.com/whosonfirst/go-whosonfirst-uri"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func DrawFeature(feature []byte, gl *globe.Globe) error {

	geom_type := gjson.GetBytes(feature, "geometry.type")

	if !geom_type.Exists() {
		return errors.New("Geometry is missing a type property")
	}

	coords := gjson.GetBytes(feature, "geometry.coordinates")

	if !coords.Exists() {
		return errors.New("Geometry is missing a coordinates property")
	}

	switch geom_type.String() {

	case "Point":

	/*
		lonlat := coords.Array()
		lat := lonlat[1].Float()
		lon := lonlat[0].Float()

		gl.DrawDot(lat, lon, 0.01, globe.Color(green))
		*/
		
	case "Polygon":

		paths := make([][]*globe.Point, 0)

		for _, ring := range coords.Array() {

			path := make([]*globe.Point, 0)
			
			for _, r := range ring.Array() {

				lonlat := r.Array()
				lat := lonlat[1].Float()
				lon := lonlat[0].Float()

				pt := globe.NewPoint(lat, lon)
				path = append(path, &pt) 
			}

			paths = append(paths, path)
		}

		gl.DrawPaths(paths)
		
	case "MultiPolygon":
		// log.Println("Can't process MultiPolygon")

	default:
		return errors.New("Unsupported geometry type")
	}

	return nil
}

func main() {

	outfile := flag.String("out", "", "Where to write globe")
	size := flag.Int("size", 1600, "The size of the globe (in pixels)")
	mode := flag.String("mode", "meta", "... (default is 'meta' for one or more meta files)")

	center := flag.String("center", "", "")
	center_lat := flag.Float64("latitude", 37.755244, "")
	center_lon := flag.Float64("longitude", -122.447777, "")

	flag.Parse()

	if *center != "" {

		latlon := strings.Split(*center, ",")

		lat, err := strconv.ParseFloat(latlon[0], 64)

		if err != nil {
			log.Fatal(err)
		}

		lon, err := strconv.ParseFloat(latlon[1], 64)

		if err != nil {
			log.Fatal(err)
		}

		*center_lat = lat
		*center_lon = lon
	}

	green := color.NRGBA{0x00, 0x64, 0x3c, 192}
	g := globe.New()
	g.DrawGraticule(10.0)

	t1 := time.Now()

	if *mode == "meta" {

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
					log.Println(err, path)
					break
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

	} else if *mode == "repo" {

		for _, path := range flag.Args() {

			var cb = func(path string, info os.FileInfo) error {

				if info.IsDir() {
					return nil
				}

				is_wof, err := uri.IsWOFFile(path)

				if err != nil {
					log.Printf("unable to determine whether %s is a WOF file, because %s\n", path, err)
					return err
				}

				if !is_wof {
					return nil
				}

				is_alt, err := uri.IsAltFile(path)

				if err != nil {
					log.Printf("unable to determine whether %s is an alt (WOF) file, because %s\n", path, err)
					return err
				}

				if is_alt {
					return nil
				}

				fh, err := os.Open(path)

				if err != nil {
					log.Printf("failed to open %s, because %s\n", path, err)
					return err
				}

				defer fh.Close()

				feature, err := ioutil.ReadAll(fh)

				if err != nil {
					log.Printf("failed to read %s, because %s\n", path, err)
					return err
				}

				return DrawFeature(feature, g)
			}

			cr := crawl.NewCrawler(path)
			cr.Crawl(cb)
		}

	} else {

		log.Fatal("Invalid mode")
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
