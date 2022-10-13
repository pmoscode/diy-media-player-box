package utils

import (
	"fmt"
	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"
	"io"
	"log"
	"os"
)

func GetAudioInformation(audioFilename string) (string, string) {
	file, _ := os.Open(audioFilename)

	// Read mp3TagInfo
	m, err := tag.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}

	// Get Track length
	t := 0.0

	d := mp3.NewDecoder(file)
	var f mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			t = 0.0
			break
		}

		t = t + f.Duration().Seconds()
	}

	return m.Title(), fmt.Sprintf("%.2f", t)
}
