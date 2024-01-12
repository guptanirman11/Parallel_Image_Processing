## Some extra information about files present apart from the pdf for part5

This directory presents the graphs with both the variations when only workers which will process images are varies and second when both workers which will process images and second which will save the images.

Additionally it contains the graphs from project1 both parslices and parfiles to compare those results with my combination of pipelines and BSP along with stealing algorithm.

The speed up graphs can plotted directly by calling sbatch benchmark-proj1.sh which internally calls the python file plot_script.py. This will add the graphs starting with speedup-grade..


### pdf for performance report and project overview is present [here](https://github.com/mpcs-jh/project-3-guptanirman11/blob/main/proj3/benchmark/project3%20performance%20report.pdf) with name project3 performance report


Before running please download tha data directory from here. https://www.dropbox.com/s/cwse3i736ejcxpe/data.zip?dl=0


The values for peanut cluster are present for both the variations in slurm/out directory, you can refer to the directory to verify that the project and the plot script both works on the cluster.



### To generate the speed up plots one can directly run sbatch benchmark-proj1.sh with the specified time so that it doesn’t fail because of time run out error. 
## However there is a catch: the first will be produced when there are fixed (2) result aggregators so You need to do 4 modifications in total .

### 3 in scheduler/pipeline.go 
### Line 51,52 and 63 from 2 in the both for loops and defining the capacity of the result channel to config.Threadcount.

### 1 in benchmark/plot_script.py

### Change the name of plot to speed_up_vary2. Which is line 72 
