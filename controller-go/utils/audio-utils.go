package utils

import (
	"fmt"
	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetAudioInformation(audioFilename string) (string, string) {
	file, _ := os.Open(audioFilename)

	title := "no title"

	// Read mp3TagInfo
	m, err := tag.ReadFrom(file)
	if err != nil {
		log.Println(err)
	} else {
		title = m.Title()
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

		t = t + f.Duration().Minutes()
	}

	time := fmt.Sprintf("%.2f", t)

	return title, convertMinutes(time)
}

func convertMinutes(min string) string {
	split := strings.Split(min, ".")

	secs, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		log.Println(err)
		return min
	}

	return split[0] + ":" + fmt.Sprintf("%.0f", secs*0.6)
}
