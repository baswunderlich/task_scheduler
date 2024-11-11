package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	name             string
	deadline         int
	computation_time int
	period           int
	comp_done        int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing cmdLine args")
		return
	}
	fmt.Print("Reading file...")
	tasks := read_conf_file()
	fmt.Println("OK")

	// fmt.Println("\nRate Monotonic Scheduling")
	// tasks01 := make([]Task, len(tasks))
	// copy(tasks01, tasks)
	// rate_monotonic_scheduling(tasks01)

	// fmt.Println("\nDeadline Monotonic Scheduling")
	// tasks02 := make([]Task, len(tasks))
	// copy(tasks02, tasks)
	// deadline_monotonic_scheduling(tasks02)

	// fmt.Println("\nResponse Time Analysis")
	// tasks03 := make([]Task, len(tasks))
	// copy(tasks03, tasks)
	// response_time_analysis(tasks03)

	// fmt.Println("\nOptimal Solution")
	// tasks04 := make([]Task, len(tasks))
	// copy(tasks04, tasks)
	// optimal_solution(tasks04)

	fmt.Println("\nEarliest Deadline First (EDF)")
	tasks05 := make([]Task, len(tasks))
	copy(tasks05, tasks)
	earliest_deadline_first_scheduling(tasks05)
}

func earliest_deadline_first_scheduling(tasks []Task) {
	time := 0
	for {
		time += 1
		chosen_task := task_not_done_with_smallest_remaining_time(tasks, time)
		chosen_task.comp_done += 1
		fmt.Print(chosen_task.name)
		if critical_instant_reached(tasks, time) {
			break
		}
		reset_tasks_when_period_reached(tasks, time)
	}
}

func optimal_solution(tasks []Task) {
	tasks = rate_monotonic_scheduling_sort(tasks)
	for k := len(tasks) - 1; k > -1; k-- {
		ok := false
		for next := k; next > -1; next-- {
			//Swap
			tasks[k], tasks[next] = tasks[next], tasks[k]
			//Is RTA ok?
			response_time := calc_response_time(k, tasks)
			ok = response_time <= tasks[k].deadline
			if ok {
				break
			}
		}
		if !ok {
			fmt.Println("Es konnte keine optimale Lösung gefunden werden")
			fmt.Println(tasks)
			return
		}
	}
	fmt.Println("Die optimale Prioritätenvergabe lautet:")
	fmt.Println(tasks)
}

func response_time_analysis(tasks []Task) {
	sorted_tasks := rate_monotonic_scheduling_sort(tasks)
	response_times := calc_response_times(sorted_tasks)
	fmt.Println(response_times)
}

func calc_response_time(tasksIndex int, tasks []Task) int {
	response_time := 0
	//Calc response time
	for {
		new_response_time := tasks[tasksIndex].computation_time
		//break when response time value doesnt change
		//New iteration => calc formula
		for j := tasksIndex - 1; j > -1; j-- {
			new_response_time += int(math.Ceil(float64(response_time)/float64(tasks[j].period))) * tasks[j].computation_time
		}
		if new_response_time == response_time {
			fmt.Println("Final response time:", new_response_time)
			break
		}
		response_time = new_response_time
	}
	return response_time
}

func calc_response_times(tasks []Task) []int {
	response_times := make([]int, len(tasks))
	for i := 0; i < len(tasks); i++ {
		response_times[i] = calc_response_time(i, tasks)
	}
	return response_times
}

// The task with the highest priority is has the index '0'
func rate_monotonic_scheduling_sort(tasks []Task) []Task {
	sorted_tasks := make([]Task, 0, len(tasks))
	for range tasks {
		var task *Task
		for i := range tasks {
			if task == nil && !do_sorted_tasks_contain(sorted_tasks, tasks[i]) {
				task = &tasks[i]
			}
		}
		for i := range tasks {
			if tasks[i].period < task.period && !do_sorted_tasks_contain(sorted_tasks, tasks[i]) {
				task = &tasks[i]
			}
		}
		sorted_tasks = append(sorted_tasks, *task)
	}
	return sorted_tasks
}

func do_sorted_tasks_contain(sorted_tasks []Task, task Task) bool {
	for i := range sorted_tasks {
		if strings.Compare(sorted_tasks[i].name, task.name) == 0 {
			return true
		}
	}
	return false
}

func read_conf_file() []Task {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
	}
	scanner := bufio.NewScanner(file)
	tasks := make([]Task, 0)
	for scanner.Scan() {
		split_line := strings.Split(scanner.Text(), ",")
		task_name := split_line[0]

		task_comp_time, err := strconv.Atoi(split_line[1])
		if err != nil {
			fmt.Println(err.Error())
		}

		task_period, err := strconv.Atoi(split_line[2])
		if err != nil {
			fmt.Println(err.Error())
		}
		task_deadline := 0
		if len(split_line) < 4 {
			task_deadline = task_period
		} else {
			task_deadline, err = strconv.Atoi(split_line[3])
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		tasks = append(tasks, Task{name: task_name, computation_time: task_comp_time, period: task_period, deadline: task_deadline})
	}
	return tasks
}

func rate_monotonic_scheduling(tasks []Task) {
	time := 0
	for {
		time += 1
		chosen_task := task_not_done_with_smallest_period(tasks)
		chosen_task.comp_done += 1
		fmt.Print(chosen_task.name)
		if critical_instant_reached(tasks, time) {
			break
		}
		reset_tasks_when_period_reached(tasks, time)
	}
}

func deadline_monotonic_scheduling(tasks []Task) {
	time := 0
	for {
		time += 1
		chosen_task := task_not_done_with_smallest_distance_to_deadline(tasks, time)
		chosen_task.comp_done += 1
		fmt.Print(chosen_task.name)
		if critical_instant_reached(tasks, time) {
			break
		}
		reset_tasks_when_period_reached(tasks, time)
	}
}

func task_not_done_with_smallest_distance_to_deadline(tasks []Task, time int) *Task {
	smallest_task := &Task{name: "_", computation_time: math.MaxInt32, deadline: math.MaxInt32}
	for i := 0; i < len(tasks); i++ {
		if smallest_task.deadline > tasks[i].deadline && tasks[i].comp_done < tasks[i].computation_time {
			smallest_task = &tasks[i]
		}
	}
	//fmt.Printf("Result: %s \n", smallest_task.name)
	return smallest_task
}

func task_not_done_with_smallest_period(tasks []Task) *Task {
	smallest_task := &Task{name: "_", computation_time: 1000000000, deadline: 10000000000}
	for i := 0; i < len(tasks); i++ {
		if smallest_task.deadline > tasks[i].deadline && tasks[i].comp_done < tasks[i].computation_time {
			smallest_task = &tasks[i]
		}
	}
	return smallest_task
}

// Earliest deadline first
func task_not_done_with_smallest_remaining_time(tasks []Task, time int) *Task {
	smallest_task := &Task{name: "_", computation_time: 1000000000, deadline: 10000000000}
	for i := 0; i < len(tasks); i++ {
		if smallest_task.deadline > tasks[i].deadline && tasks[i].comp_done < tasks[i].computation_time {
			smallest_task = &tasks[i]
		}
	}
	return smallest_task
}

func critical_instant_reached(tasks []Task, time int) bool {
	for i := range tasks {
		if time%tasks[i].period != 0 {
			return false
		}
	}
	return true
}

func reset_tasks_when_period_reached(tasks []Task, time int) {
	for i := range tasks {
		if time%tasks[i].period == 0 {
			if tasks[i].comp_done < tasks[i].computation_time {
				fmt.Println("\nERROR! Scheduling failed")
			}
			//Offset the deadline to match the next period
			tasks[i].deadline += tasks[i].period
			tasks[i].comp_done = 0
		}
	}
}
