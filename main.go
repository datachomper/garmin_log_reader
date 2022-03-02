package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	header := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Paths</name>
    <description>Examples of paths. Note that the tessellate tag is by default
      set to 0. If you want to create tessellated lines, they must be authored
      (or edited) directly in KML.</description>
    <Style id="yellowLineGreenPoly">
      <LineStyle>
        <color>ffff00ff</color>
        <width>4</width>
      </LineStyle>
      <PolyStyle>
        <color>ffff00ff</color>
      </PolyStyle>
    </Style>
    <Placemark>
      <name>Absolute Extruded</name>
      <description>Transparent green wall with yellow outlines</description>
      <styleUrl>#yellowLineGreenPoly</styleUrl>
      <LineString>
        <coordinates>`

	footer := `</coordinates>
	      </LineString>
	          </Placemark>
		    </Document>
		    </kml>`

	outfilename := flag.String("out", "output.kml", "don't forget the output filename")
	flag.Parse()

	outfile, err := os.Create(*outfilename)
	if err != nil {
		fmt.Println("Unable to open output file: ", *outfilename)
		return
	}
	outfile.WriteString(header)

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Useage: ./main -out output_filename.kml <log_file_1.csv> <log_file_2.csv> ...")
	}

	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}

		csvreader := csv.NewReader(f)
		csvreader.Comment = '#'
		data, err := csvreader.ReadAll()
		if err != nil {
			fmt.Println(err)
			continue
		}

		for i, line := range data {
			// Data is 1hz, only take 60hz samples to reduce output file size
			if i > 1 && (i%60) == 0 {
				// Empty values bork up the map
				if len(line[5]) == 0 || len(line[4]) == 0 {
					continue
				}
				outfile.WriteString(fmt.Sprintf("%v,%v\n", line[5], line[4]))
			}
		}

		f.Close()
	}

	outfile.WriteString(footer)
	outfile.Close()
}
