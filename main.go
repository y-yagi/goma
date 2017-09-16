package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const taskDuration = 25 * time.Minute
const restDuration = 5 * time.Minute
const longRestDuration = 15 * time.Minute

func formatMinutes(timeLeft time.Duration) string {
	minutes := int(timeLeft.Minutes())
	seconds := int(timeLeft.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func countDown(target time.Time) {
	for range time.Tick(100 * time.Millisecond) {
		timeLeft := -time.Since(target)
		if timeLeft < 0 {
			fmt.Print("Countdown: ", formatMinutes(0), "   \r")
			return
		}
		fmt.Fprint(os.Stdout, "Countdown: ", formatMinutes(timeLeft), "   \r")
		os.Stdout.Sync()
	}
}

func task() error {
	start := time.Now()
	finish := start.Add(taskDuration)
	fmt.Printf("Start task.\n")

	countDown(finish)

	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nTag: ")

	if !scanner.Scan() {
		return nil
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	tag := scanner.Text()
	createTomato(tag)

	return nil
}

func rest(duration time.Duration) {
	start := time.Now()
	finish := start.Add(duration)
	fmt.Printf("\nStart rest.\n")

	countDown(finish)

	_ = exec.Command("mpg123", "data/ringing.mp3").Start()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\nPlease press the Enter key for start next tomato.\n")
	scanner.Scan()
}

func run(args []string, outStream, errStream io.Writer) int {
	var show string

	flags := flag.NewFlagSet("goma", flag.ExitOnError)
	flags.SetOutput(errStream)
	flags.StringVar(&show, "s", "", "Show your tomatoes. You can specify range, 'today', 'week', 'month' or 'all'.")
	flags.Parse(args[1:])

	err := initDB()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	if len(show) != 0 {
		if !contains([]string{"today", "week", "month", "all"}, show) {
			fmt.Printf("'%s' is invalid argument. Please specify 'today', 'week', 'month' or 'all'.\n", show)
			return 1
		}

		err = showTomatoes(show)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return 1
		}
		return 0
	}

	for i := 1; ; i++ {
		err = task()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return 1
		}

		if i%4 == 0 {
			rest(longRestDuration)
		} else {
			rest(restDuration)
		}
	}
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
