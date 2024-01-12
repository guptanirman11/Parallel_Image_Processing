import subprocess
import matplotlib
matplotlib.use('Agg')

import matplotlib.pyplot as plt
import os
import numpy as np

threads = np.array([2, 4, 6, 8, 12])

def run_and_capture_output(scheduler, repetitions):
    # Relative path to the editor directory
    editor_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'editor')

    timing_map = {
        'small': [],
        'mixture': [],
        'big': []
    }

    for _ in range(repetitions):
        
        if scheduler == "s":
            for data_dir in timing_map.keys():
                print(f'Running in {data_dir} directory with 1')
                command = ['go', 'run', 'editor.go', data_dir, scheduler, str(1)]
                result = subprocess.run(command, text=True, capture_output=True, check=True, cwd=editor_dir)

                total_time_line = result.stdout.strip().split('\n')[0]
                total_time = float(total_time_line.split()[1])
                timing_map[data_dir].append(total_time)
        else:
            for data_dir in timing_map.keys():
                for n in threads:
                    print(f'Running in {data_dir} directory with {n} threads.')
                    command = ['go', 'run', 'editor.go', data_dir, scheduler, str(n)]
                    result = subprocess.run(command, text=True, capture_output=True, check=True, cwd=editor_dir)

                    total_time_line = result.stdout.strip().split('\n')[0]
                    
                    total_time = float(total_time_line.split()[1])
                    timing_map[data_dir].append(total_time)

    return timing_map

def calculate_average(timing_map, repetitions):
    average_map = {}
    for data_dir, timings in timing_map.items():

        average_timings = [(timings[i] + timings[i +len(threads)] + timings[i + 2*len(threads)]) / repetitions for i in range(0, len(threads))]
        average_map[data_dir] = average_timings
        print(average_map)
    return average_map

def plot_timing_map(average_map, scheduler, sequential_data):
    plt.figure()
    
    x = threads

    for data_dir, y in average_map.items():
        base_time = sequential_data[data_dir]
        speedup = [round(base_time / time, 3) if time != 0 else 0 for time in y]
        plt.plot(x, speedup, label=data_dir)
    
    plt.xlabel('Number of Threads (n)')
    plt.ylabel('Speedup')
    plt.title(f'Editor Speedup ({scheduler})')
    plt.legend()
    plt.grid(True)
    plt.xticks(x)

    plt.savefig(f'speedup-test-{scheduler}.png')
    plt.close()

def main():
    print('Running sequential baseline')
    sequential_data = run_and_capture_output('s', repetitions=3)

    for key, value in sequential_data.items():
        sequential_data[key] = sum(value) / 3

    print(sequential_data)
   

    for x in ['pipeline']:
        print('Scheduler:', x)
        timing_data = run_and_capture_output(x, repetitions=3)
        average_map = calculate_average(timing_data, repetitions=3)
        plot_timing_map(average_map, x, sequential_data)

if __name__ == '__main__':
    main()