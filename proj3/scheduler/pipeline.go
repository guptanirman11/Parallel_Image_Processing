package scheduler

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math/rand"
	"os"
	"proj3/png"
	"strings"
)

// ImageResult represents the result of processing an ImageTask
type ImageResult *ImageTask

// RunPipeline implements the pipeline
func RunPipeline(config Config) {
	effectsPathFile := fmt.Sprintf("../data/effects.txt")
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		fmt.Println("Error opening effects file:", err)
		return
	}

	dataDirs := strings.Split(config.DataDirs, "+")

	// Creating channels for communication between stages
	imageTaskChan := make(chan ImageTask, config.ThreadCount)
	imageResultChan := make(chan ImageResult, config.ThreadCount)

	workerdoneChan := make(chan struct{}, config.ThreadCount)

	// Creating worker queues

	workerQueues := make([]DEQueue, config.ThreadCount)
	for i := range workerQueues {
		workerQueues[i] = NewBDEQueue()

	}

	// Starting the ImageTaskGenerator goroutine
	go ImageTaskGenerator(effectsFile, dataDirs, imageTaskChan)

	// Starting worker goroutines
	for i := 0; i < config.ThreadCount; i++ {
		go Worker(i, imageTaskChan, imageResultChan, workerQueues[i], workerQueues, config.ThreadCount, workerdoneChan)
	}

	//Starting resul aggreagtor go routines
	doneChan := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go ResultsAggregator(imageResultChan, doneChan)
	}

	for i := 0; i < config.ThreadCount; i++ {
		<-workerdoneChan
	}

	close(imageResultChan)

	// Wait for ResultsAggregator to finish
	for i := 0; i < 2; i++ {
		<-doneChan
	}

}

// ImageTaskGenerator reads from effectsFile and generates ImageTasks
func ImageTaskGenerator(effectsFile *os.File, dataDirs []string, outChan chan<- ImageTask) {
	defer close(outChan)

	for _, dataDir := range dataDirs {
		effectsFile.Seek(0, 0)

		// Reinitialize the reader for each directory
		reader := json.NewDecoder(effectsFile)
		for {
			var entry ImageTask
			if err := reader.Decode(&entry); err != nil {
				if err == io.EOF {
					// End of file is expected
					break
				}
				fmt.Println("Error reading the image task:", err)
				break
			}

			entry.InPath = fmt.Sprintf("../data/in/%s/%s", dataDir, entry.InPath)
			entry.OutPath = fmt.Sprintf("../data/out/%s"+"_"+"%s", dataDir, entry.OutPath)
			entry.DataDir = dataDir

			outChan <- entry

		}
	}

}

// Worker processes ImageTasks and sends the results to the result channel
func Worker(workerID int, inChan <-chan ImageTask, outChan chan<- ImageResult, ownQueue DEQueue, allQueues []DEQueue, threadCount int, done chan<- struct{}) {
	defer func() {
		done <- struct{}{} // Signal completion of the worker goroutine
	}()
	for {
		select {
		case task, ok := <-inChan:
			if !ok {
				// In case the channel is closed
				// Process remaining tasks in the queue before returning
				for ownQueue.Size() > 1 {
					task := ownQueue.popBottom()

					if task != nil {
						ProcessImageTask(task, workerID, threadCount)

						result := ImageResult(task)
						outChan <- result
					}
				}
				task := ownQueue.popBottom()

				if task != nil {
					ProcessImageTask(task, workerID, threadCount)

					result := ImageResult(task)
					outChan <- result
				}

				return
			}

			// Add the task to the worker's queue

			ownQueue.pushBottom(&task)

		default:

			if ownQueue.Size() > 1 {
				task := ownQueue.popBottom()

				// Process the task and apply effects
				if task != nil {
					ProcessImageTask(task, workerID, threadCount)

					// fmt.Println("before adding image result channel", task)
					result := ImageResult(task)
					outChan <- result
				}

			} else {

				// Try to steal work if the worker's queue is empty
				randomInt := rand.Intn(threadCount)
				stealQueue := allQueues[randomInt]
				StolenTask := stealQueue.popTop()
				if StolenTask != nil {

					// Adding the stolen task to the worker's queue
					ownQueue.pushBottom(StolenTask)
				}
			}
		}
	}
}

// ResultsAggregator reads from the result channel and saves the filtered images
func ResultsAggregator(inChan <-chan ImageResult, done chan<- struct{}) {
	for result := range inChan {
		// Save the filtered image
		err := result.ProcessedImage.Save(result.OutPath)

		if err != nil {
			fmt.Printf("Error saving image to %s: %v\n", result.OutPath, err)
		}
	}
	done <- struct{}{}
}

// ProcessImageTask applies effects to the image task
func ProcessImageTask(task *ImageTask, wokrerId int, threadCount int) {

	img, err := png.Load(task.InPath)

	if err != nil {
		fmt.Printf("Error loading image from %s: %v\n", task.InPath, err)
	}

	rectangleHeight := img.Out.Bounds().Max.Y / threadCount

	for i, effect := range task.Effects {

		// Helper function to break images in slices and apply effects parallely
		// and return after applying effects to the image

		processChunk(threadCount, rectangleHeight, effect, img)

		if i != len(task.Effects)-1 {
			img.In = img.Out
			img.Out = image.NewRGBA64(img.Bounds)

		}

	}

	task.ProcessedImage = img

}

// Applies each effect to image by breaking into several chunks.
func processChunk(threadCount int, rectangleHeight int, effect string, img *png.Image) {

	// var completed int32
	var done = make(chan struct{}, threadCount)

	interval := img.Out.Bounds().Max.Y

	for i := 0; i < threadCount; i++ {
		recStart := rectangleHeight * i
		if recStart > interval {
			recStart = interval
		}
		recEnd := recStart + rectangleHeight + 100
		if recEnd > interval {
			recEnd = interval + 100
		}

		go func(start, end int) {
			defer func() { done <- struct{}{} }()

			switch effect {
			case "E":
				img.EdgeDetect(start, end)
			case "S":
				img.Sharpen(start, end)
			case "B":
				img.Blur(start, end)
			case "G":
				img.Grayscale(start, end)
			}

		}(recStart, recEnd)
	}

	for i := 0; i < threadCount; i++ {
		<-done
	}

	// Close the done channel
	close(done)

}
