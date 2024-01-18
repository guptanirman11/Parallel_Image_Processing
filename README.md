
# Parallel Project \: An Image Processing System
[Worked on the project from September 2023 to December 2023]

Go encourages explicit communication
between threads (for example using channels), which is often easier to
understand and makes data races easier to avoid than programming models
where multiple threads share large amounts of data that they work on at
the same time. The model encouraged by Go is sometimes referred to as
"CSP style programming", for "Communicating Sequential Processes". If
you want to learn more, here's a brief explanation of the CSP acronym:
<https://levelup.gitconnected.com/communicating-sequential-processes-csp-for-go-developer-in-a-nutshell-866795eb879d>
In this project, I used CSP as well as more advance work
distribution techniques.

## Preliminaries

Many algorithms in image
processing benefit from parallelization (especially those that run on
GPUs). In this project, we will create an image processing system
that runs on a CPU, and reads in a series of images and applies certain
effects to them using image convolution. 

## Assignment: Image Processing System
For this
project, we will create an image editor that will apply image effects
on series of images using 2D image convolution. The program will read
in from a series of JSON strings, where each string represents an image
along with the effects that should be applied to that image. Each string
will have the following format,

``` json
{ 
  "inPath": string, 
  "outPath": string, 
  "effects": [string] 
}
```

where each key-value is described in the table below,

| Key-Value                     | Description |
|-------------------------------|-------------|
| ``"inPath":"sky.png"``        | The ``"inPath"`` pairing represents the file path of the image to read in. Images in  this assignment will always be PNG files. All images are relative to the ``data`` directory inside the ``proj2`` folder. |
| ``"outPath:":"sky_out.png"``  | The ``"outPath"`` pairing represents the file path to save the image after applying the effects. All images are relative to the ``data`` directory inside the ``proj2`` folder. |
| ``"effects":["S"\,"B"\,"E"]`` | The ``"effects"`` pairing  represents the image effects to apply to the image. You must apply these in the order they are listed. If no effects are specified (e.g.\, ``[]``) then the out image is the same as the input image. |

The program will read in the images, apply the effects associated with
an image, and save the images to their specified output file paths. How
the program processes this file is described in the **Program
Specifications** section.

### Image Effects

The sharpen, edge-detection, and blur image effects are required to use
image convolution to apply their effects to the input image. Again, we
can read about how to perform image convolution here:

-   [Two Dimensional
    Convolution](http://www.songho.ca/dsp/convolution/convolution2d_example.html)

As stated in the above article, the size of the input and output image
are fixed (i.e., they are the same). Thus, results around the border
pixels will not be fully accurate because we will need to pad zeros
where inputs are not defined. We are required to use the a zero-padding
when working with pixels that are not defined. 

Each effect is identified by a single character that is described below,

| Image Effect | Description |
| -------------|-------------|
| ``"S"`` | Performs a sharpen effect with the following kernel (provided as a flat go array): ``[9]float6 {0,-1,0,-1,5,-1,0,-1,0}``. |
| ``"E"`` | Performs an edge detection effect with the following kernel (provided as a flat go array): ``[9]float64{-1,-1,-1,-1,8,-1,-1,-1,-1}``. |
| ``"B"`` | Performs a blur effect with the following kernel (provided as a flat go array): ``[9]float64{1 / 9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}``. |
| ``"G"`` | Performs a grayscale effect on the image. This is done by averaging the values of all three color numbers for a pixel, the red, green and blue, and then replacing them all by that average. So if the three colors were 25, 75 and 250, the average would be 116, and all three numbers would become 116. |
### The `data` Directory


The Data directory was placed inside the `proj` directory that contains the
subdirectories: `editor` and `png`. **I HAVE NOT COMMITTED THIS DIRECTORY TO
 REPOSITORY**. These are very large files!

### Working with Images in Go

I decided to use `image` package which is provided by Go that
makes it easy to load,read,and save PNG images. I recommend looking at
the examples from these links:

-   [Go PNG docs](https://golang.org/pkg/image/png/)
-   A [helpful
    tutorial](https://www.devdungeon.com/content/working-images-go) for
    working on png images. I used the
    `At()` and `Set()` methods as specified by the Go PNG documentation.

> **Note**:
> As the image package only allows you to read an image data and not modify
> it in-place. I created a separate out buffer to represent
> the modified pixels. I have done this in the `Image`
> struct as follows:

``` go
type Image struct {
  in  *image.RGBA64  
  out *image.RGBA64 
  Bounds  image.Rectangle
  ... 
} 
```

As these are
**pointers** so I only need to swap the pointers to make the old out
buffer the new in buffer when applying one effect after another effect.
This process is less expensive than copying pixel data after apply each
effect.

### Program Specifications

For this project, I implemented two versions of this image
processing system. The versions will include a sequential version and
one parallel version(pipeline).

The program has the following
usage statement:

    Usage: editor data_dir [mode] [number_of_threads]
    data_dir = The data directories to use to load the images.
    mode     = (bsp) run the BSP mode, (pipeline) run the pipeline mode
    number_of_threads = Runs the parallel version of the program with the specified number of threads (i.e., goroutines).

The `data_dir` argument will always be either `big`, `small`, or
`mixture` or a combination between them. The program will always read
from the `data/effects.txt` file; however, the `data_dir` argument
specifies which directory to use. The user can also add a `+` to perform
the effects on multiple directories. For example, `big` will apply the
`effects.txt` file on the images coming from the `big` directory. The
argument `big+small` will apply the `effects.txt` file on both the `big`
and `small` directory. The program must always prepend the `data_dir`
identifier to the beginning of the `outPath`. For example, running the
program as follows:

    $: go run editor.go big pipleine 4 

will produce inside the `out` directory the following files:

    big_IMG_2020_Out.png 
    big_IMG_2724_Out.png 
    big_IMG_3695_Out.png 
    big_IMG_3696_Out.png 
    big_IMG_3996_Out.png 
    big_IMG_4061_Out.png 
    big_IMG_4065_Out.png
    big_IMG_4066_Out.png 
    big_IMG_4067_Out.png
    big_IMG_4069_Out.png

Here's an example of a combination run:

    $: go run editor.go big+small pipeline 2

will produce inside the `out` directory the following files:

    big_IMG_2020_Out.png 
    big_IMG_2724_Out.png 
    big_IMG_3695_Out.png 
    big_IMG_3696_Out.png 
    big_IMG_3996_Out.png 
    big_IMG_4061_Out.png 
    big_IMG_4065_Out.png
    big_IMG_4066_Out.png 
    big_IMG_4067_Out.png
    big_IMG_4069_Out.png
    small_IMG_2020_Out.png 
    small_IMG_2724_Out.png 
    small_IMG_3695_Out.png 
    small_IMG_3696_Out.png 
    small_IMG_3996_Out.png 
    small_IMG_4061_Out.png 
    small_IMG_4065_Out.png
    small_IMG_4066_Out.png 
    small_IMG_4067_Out.png
    small_IMG_4069_Out.png

The `mode` and `number_of_threads` arguments will be used to run one of
the parallel versions. 

The scheduling (i.e., running) of the various implementations is handled
by the `scheduler` package defined in `proj/scheduler` directory. The
`editor.go` program will create a configuration object (similar to
project 1) using the following struct:

``` go
type Config struct {
  DataDirs string //Represents the data directories to use to load the images.
  Mode     string // Represents which scheduler scheme to use
  // If Mode == "s" run the sequential version
  // If Mode == "pipeline" run the pipeline version
  // If Mode == "bsp" run the pipeline version
  // These are the only values for Version
  ThreadCount int // Runs the parallel version of the program with the
  // specified number of threads (i.e., goroutines)
}
```

The `Schedule` function inside the `proj/scheduler/scheduler.go` file
will then call the correct version to run based on the `Mode` field of
the configuration value. Each of the functions to begin running the
various implementation will be explained in the following sections.

## Part 1: Sequential Implementation


The sequential version is ran by default when executing the `editor`
program when the `mode` and `number_of_threads` are both not provided.
The sequential program is relatively straightforward. This version
run through the images specified by the strings coming in from
`effects.txt`, apply their effects and save the modified images to their
output files inside the `data/out` directory. 


## Part 2: Pipeline + BSP Implementation

This parallel implementation will use channels and is
implemented as follows:

1.  For this version, all synchronization between the goroutines is
    done using channels.

2.  We implemented the  **fan-in/fan-out** scheme.

    -   **Image Task Generator**: As stated earlier, the program will
        read in the images to process via `effects.txt`. Reading is done
        by a single generator goroutine. The image task generator will
        read in the JSON strings and do any preparation needed before
        applying their effects. 

    -   **ImageTask**: A value that holds everything needed to do
        filtering for a specific JSON string. Again, its up to you how
        you define the `ImageTask` struct.

    -   **Workers**: The workers are the goroutines that are performing
        the filtering effects on the images. The number of workers is
        static and is equal to the `number_of_threads` command line
        argument. A worker use a pipeline pattern. 
        Each `Worker` of the pipeline, have an internal *data
        decomposition* component, do the following:

        -   Spawn `N` number of goroutines, where
            `N = number_of_threads`. We will call these "mini-workers".
        -   Each mini-worker goroutine is given a section of the image
            to work on.
        -   Each mini-worker goroutine will apply the effect for that
            stage to its assigned section.
        -   You should give approximately equal portions to all
            mini-workers.

        Visually the splitting could look something like this if
        `number_of_threads=6`:

        This means that if there will be a total of 6 `Workers` with
        each having 6 mini-workers running, which totals to 36 "worker"
        goroutines running in parallel. 
        The output of a worker is an `ImageResult`.

    -   **ImageResult**: the final image after applying its effects.
        
    -   **Results Aggregator**: The results aggregator gorountine reads
        from the channel that holds the `ImageResults` and saves the
        filtered image to its `"outPath"` file.

3.  If all the images have been processed then the main goroutine can
    exit the program. The main goroutine does not handle the
    Image Task Generator and/or the Results Aggregator.

4.  The `mode` command line argument value for executing this version is
    `"pipeline"`.

Another way of stating the above breakdown is the following:

-   *X* images to process
-   *P* number of threads (which is supplied by the `number_of_threads`
    flag)
-   *N* workers (*N = P*)

1.  `ImageTaskGenerator` produces *X* `ImageTasks` and dumps them all
    into a channel.
2.  A worker *w1* tries to grab a single `ImageTask` *x1* from the
    channel. This worker is SOLELY responsible for performing all *E*
    effects for this one image.
3.  The worker *w1* splits up this single `ImageTask` *x1* into *P*
    roughly equal portions and spawns *P* goroutines to apply effect
    *e1* to the image. Once all goroutines have applied effect *e1*,
    they can begin applying *e2*, *e3*, ... , until all *E* effects have
    been applied in order.
4.  Meanwhile, other workers *w2*, *w3*, ... all the way up to *N* are
    concurrently grabbing images of their own, and performing step (3)
    on their own images.
5.  If *X \> N*, then a worker *w* will go back up to step (2), grab
    another `ImageTask` *x*, and perform step (3) on it. We repeat this
    process until all *X* images have been processed.



   

## Part 3: Performance Measurements and Speedup Graphs

I ran timing measurements on both the sequential and parallel
versions of the `editor.go` program. The `data` directory will be used
to measure the performance of the parallel versions versus the
sequential version. We will keep things simple and only look at
measuring single data directories: `small`, `mixture`, and `big`. The
measurements gathered will be used to create speedup graphs. Each speedup graph is based around a single parallel version
(e.g., `pipeline`) where each line represents running a specific data
directory. The set of threads will be `{2,4,6,8,12}` 


1.  Each line in the graph represents a `data` directory size (i.e.,
    `small`, `mixture`, and `big`) that you will run for each thread
    number in the set of threads (i.e., `{2,4,6,8,12}`).



3.  For each speedup graph, the y-axis lists the speedup measurement
    and the x-axis lists the number of threads.

4.  All  work for this section is placed in the `benchmark`
    directory along with the generated speedup graphs. 

## Part 4: Performance Analysis

### Project Summary

In my project report, I provided a comprehensive summary of the experiment results and the corresponding conclusions drawn from them. The report incorporated essential graphs, along with a detailed analysis to facilitate a clear understanding of the developed code, conducted experiments, and the supporting data. Here are the key components included in the report:

### Project Description
I began with a concise paragraph outlining the project, summarizing its core objectives and scope.

### Testing Script Instructions
Clear instructions were provided on how to execute the testing script. Users could seamlessly run the script with a simple command, such as `sbatch benchmark-proj.sh`.

### Graph Analysis
I delved into a thorough analysis of the graphs, addressing key questions:
- Identified hotspots and bottlenecks in the sequential program.
- Discussed the superior performance of a specific parallel implementation and provided insights into the reasons behind its efficiency.
- Explored the impact of problem size (data size) on performance.
- Speculated on potential differences in performance measurements if the Go runtime scheduler utilized a `1:1` or `N:1` scheduler.

### Performance Improvement Hypotheses
I highlighted hypothetical areas in the implementation that could see performance improvements. A rationale was provided for anticipating these enhancements.

In essence, the project report aimed to encapsulate the entire project experience, from its foundational description to the practical aspects of running experiments and drawing insightful conclusions.
