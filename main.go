package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Task struct {
	Name      string
	URL       string
	Priority  int
	Note      string
	Completed bool
}

type Category struct {
	Name   string
	Tasks  []Task
	SortBy string
}

var categories []Category

const dataFile = "/Users/jarjohns/git/2learn/data.txt"

func main() {
	loadData()

	for {
		renderUI()
		fmt.Println("\nOptions:")
		fmt.Println("[1] Add Category")
		fmt.Println("[2] Add Task")
		fmt.Println("[3] Sort Category")
		fmt.Println("[4] Save & Exit")
		fmt.Print("Choose: ")

		switch readLine() {
		case "1":
			fmt.Print("Enter new category name: ")
			name := readLine()
			categories = append(categories, Category{Name: name})
		case "2":
			if len(categories) == 0 {
				fmt.Println("No categories. Add one first.")
				continue
			}
			fmt.Println("Select category index:")
			for i, c := range categories {
				fmt.Printf("[%d] %s\n", i, c.Name)
			}
			fmt.Print("Index: ")
			ci, _ := strconv.Atoi(readLine())
			if ci < 0 || ci >= len(categories) {
				fmt.Println("Invalid index.")
				continue
			}
			newTask := Task{}
			fmt.Print("Task Name: ")
			newTask.Name = readLine()
			fmt.Print("Task URL: ")
			newTask.URL = readLine()
			fmt.Print("Priority (1 = high): ")
			p, _ := strconv.Atoi(readLine())
			newTask.Priority = p
			fmt.Print("Note: ")
			newTask.Note = readLine()
			fmt.Print("Completed? (y/n): ")
			newTask.Completed = strings.ToLower(readLine()) == "y"

			categories[ci].Tasks = append(categories[ci].Tasks, newTask)
		case "3":
			fmt.Println("Select category index:")
			for i, c := range categories {
				fmt.Printf("[%d] %s\n", i, c.Name)
			}
			fmt.Print("Index: ")
			ci, _ := strconv.Atoi(readLine())
			if ci < 0 || ci >= len(categories) {
				fmt.Println("Invalid index.")
				continue
			}
			fmt.Print("Sort by [priority/completed/none]: ")
			sortBy := strings.ToLower(readLine())
			if sortBy != "priority" && sortBy != "completed" && sortBy != "none" {
				fmt.Println("Invalid sort option.")
			} else {
				categories[ci].SortBy = sortBy
			}
		case "4":
			saveData()
			fmt.Println("Data saved.")
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func renderUI() {
	fmt.Print("\033[H\033[2J") // Clear screen

	fmt.Println("2learn - What to Learn")
	fmt.Println(colorText(strings.Repeat("=", 90), cyan))

	for _, cat := range categories {
		fmt.Printf("\n%s  [Sort: %s]\n", colorText(cat.Name, yellow), cat.SortBy)
		fmt.Println(colorText(strings.Repeat("-", 90), cyan))

		// Column headers in magenta, separators in cyan
		fmt.Printf("%s %s %s %s %s %s %s %s %s\n",
			colorText(fmt.Sprintf("%-12s", "Name"), magenta),
			colorText("|", cyan),
			colorText(fmt.Sprintf("%-30s", "URL"), magenta),
			colorText("|", cyan),
			colorText(fmt.Sprintf("%-3s", "Pr"), magenta),
			colorText("|", cyan),
			colorText(fmt.Sprintf("%-30s", "Note"), magenta),
			colorText("|", cyan),
			colorText("✓", magenta),
		)

		fmt.Println(colorText(strings.Repeat("-", 90), cyan))

		// Sort tasks if needed
		switch cat.SortBy {
		case "priority":
			sort.Slice(cat.Tasks, func(i, j int) bool {
				return cat.Tasks[i].Priority < cat.Tasks[j].Priority
			})
		case "completed":
			sort.Slice(cat.Tasks, func(i, j int) bool {
				return !cat.Tasks[i].Completed && cat.Tasks[j].Completed
			})
		}

		for _, t := range cat.Tasks {
			check := "[ ]"
			name := t.Name
			note := t.Note
			nameDisplay := fmt.Sprintf("%-12s", t.Name)
			noteDisplay := fmt.Sprintf("%-30s", t.Note)

			if t.Completed {
				check = colorText("[✔]", green)
				name = colorText(nameDisplay, green)
				note = colorText(noteDisplay, green)
			} else {
				name = colorText(nameDisplay, white)
				note = colorText(noteDisplay, white)
			}

			// Fixed-width, colored priority
			prStrRaw := fmt.Sprintf("%-3s", strconv.Itoa(t.Priority))
			prStr := prStrRaw
			if t.Priority <= 2 {
				prStr = colorText(prStrRaw, red)
			}

			// Truncate long URLs for display, keep full link
			displayText := t.URL
			if len(displayText) > 27 {
				displayText = displayText[:27] + "..."
			}
			displayText = fmt.Sprintf("%-30s", displayText)
			url := colorText(hyperlink(displayText, t.URL), white)

			// Print the row
			fmt.Printf("%s %s %s %s %s %s %s %s %s\n",
				name,
				colorText("|", cyan),
				url,
				colorText("|", cyan),
				prStr,
				colorText("|", cyan),
				note,
				colorText("|", cyan),
				check,
			)
		}
	}
	fmt.Println("\n" + colorText(strings.Repeat("=", 90), cyan))
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func saveData() {
	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(categories)
	if err != nil {
		fmt.Println("Encode error:", err)
	}
}

func loadData() {
	file, err := os.Open(dataFile)
	if err != nil {
		// File doesn't exist yet
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&categories)
	if err != nil {
		fmt.Println("Load error:", err)
	}
}

// ----- Coloring & Hyperlink helpers -----

const (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

func colorText(text, color string) string {
	return color + text + reset
}

func hyperlink(text, url string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}
