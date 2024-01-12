package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"proj3/png"
	"strings"
)

func RunSequential(config Config) {
	effectsPathFile := fmt.Sprintf("../data/effects.txt")
	effectsFile, _ := os.Open(effectsPathFile)

	dataDirs := strings.Split(config.DataDirs, "+")

	// Traversing through each directory and creating an Entry struct.
	for _, dataDir := range dataDirs {

		outDir := fmt.Sprintf("../data/out/%s", dataDir)
		reader := json.NewDecoder(effectsFile)

		for {
			var Entry ImageTask

			if err := reader.Decode(&Entry); err != nil {
				break
			}

			inPath := fmt.Sprintf("../data/in/%s/%s", dataDir, Entry.InPath)
			outPath := fmt.Sprintf("%s"+"_"+"%s", outDir, Entry.OutPath)

			// Load the image from inPath.
			img, err := png.Load(inPath)
			if err != nil {
				fmt.Println("Error in image loading part")
				fmt.Println(err)
			}
			// Applying each effect to the specific image.
			for i, effect := range Entry.Effects {

				if effect == "E" {
					img.EdgeDetect(img.Bounds.Min.Y, img.Bounds.Max.Y)
				}
				if effect == "S" {
					img.Sharpen(img.Bounds.Min.Y, img.Bounds.Max.Y)
				}
				if effect == "B" {
					img.Blur(img.Bounds.Min.Y, img.Bounds.Max.Y)
				}
				if effect == "G" {
					img.Grayscale(img.Bounds.Min.Y, img.Bounds.Max.Y)
				}
				// Swapping the images before applying the next effect.
				if i != len(Entry.Effects)-1 {
					img.In = img.Out
					img.Out = image.NewRGBA64(img.Bounds)

				}

			}
			// saving the image
			err2 := img.Save(outPath)

			//Checks to see if there were any errors when saving.
			if err2 != nil {
				panic(err2)
			}

		}

	}

}
