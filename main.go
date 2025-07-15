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
		fmt.Println(colorText("\nOptions:", green))
		fmt.Println(colorText("[1] Add Category", green))
		fmt.Println(colorText("[2] Add Task", green))
		fmt.Println(colorText("[3] Sort Category", green))
		fmt.Println(colorText("[4] Save & Exit", green))
		fmt.Println(colorText("[5] View/Edit Task", green))
		fmt.Print(colorText("Choose: ", green))

		switch readLine() {
		case "1":
			fmt.Print("Enter new category name: ")
			name := readLine()
			categories = append(categories, Category{Name: name})
		case "2":
			addTask()
		case "3":
			sortCategory()
		case "4":
			saveData()
			fmt.Println("Data saved.")
			return
		case "5":
			viewOrEditTask()
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

		// Column headers
		fmt.Printf("%s %s %s %s %s %s %s %s %s\n",
			colorText(fmt.Sprintf("%-12s", "Name"), magenta),
			colorText("|", cyan),
			colorText(fmt.Sprintf("%-30s", "URL"), magenta),
			colorText("|", cyan),
			colorText(fmt.Sprintf("%-2s", "Pr"), magenta),
			colorText("|", cyan),
			colorText(fmt.Sprintf("%-30s", "Note"), magenta),
			colorText("|", cyan),
			colorText(" ✓", magenta),
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

		// Print each task with colors

		for _, t := range cat.Tasks {
			check := "[ ]"
			if t.Completed {
				check = colorText("[✔]", green)
			}

			// Truncate name and color green if completed
			nameDisplay := t.Name
			if len(nameDisplay) > 11 {
				nameDisplay = nameDisplay[:9] + "..."
			}
			nameColor := white
			if t.Completed {
				nameColor = green
			}
			name := colorText(fmt.Sprintf("%-12s", nameDisplay), nameColor)

			// Truncate and hyperlink URL, colored blue
			displayURL := t.URL
			if len(displayURL) > 27 {
				displayURL = displayURL[:27] + "..."
			}
			displayURL = fmt.Sprintf("%-30s", displayURL)
			url := colorText(hyperlink(displayURL, t.URL), blue)

			noteLines := wrapText(t.Note, 30)
			for i, line := range noteLines {
				noteColor := white
				if t.Completed {
					noteColor = green
				}
				noteFormatted := colorText(fmt.Sprintf("%-30s", line), noteColor)

				if i == 0 {
					fmt.Printf("%s %s %s %s %s %s %s %s %s\n",
						name,
						colorText("|", cyan),
						url,
						colorText("|", cyan),
						colorPriority(t.Priority),
						colorText("|", cyan),
						noteFormatted,
						colorText("|", cyan),
						check,
					)
				} else {
					fmt.Printf("%-12s %s %-30s %s %-3s %s %s %s %s\n",
						strings.Repeat(" ", 12), colorText("|", cyan),
						strings.Repeat(" ", 30), colorText("|", cyan),
						strings.Repeat(" ", 3), colorText("|", cyan),
						noteFormatted, colorText("|", cyan), "",
					)
				}
			}
		}

	}
	fmt.Println("\n" + colorText(strings.Repeat("=", 90), cyan))
}

func colorPriority(priority int) string {
	switch {
	case priority <= 2:
		return colorText(fmt.Sprintf("%d ", priority), red)
	case priority < 6:
		return colorText(fmt.Sprintf("%d ", priority), yellow)
	default:
		return colorText(fmt.Sprintf("%d ", priority), green)
	}
}

// --- View/Edit ---

func viewOrEditTask() {
	cat := selectCategory()
	if cat == nil {
		return
	}
	if len(cat.Tasks) == 0 {
		fmt.Println("No tasks in this category.")
		return
	}
	task := selectTask(cat)
	if task == nil {
		return
	}

	fmt.Println("View or Modify? [v/m]: ")
	choice := strings.ToLower(readLine())

	switch choice {
	case "v":
		displayFullTask(*task)
		promptToContinue()
	case "m":
		modifyTask(task)
	default:
		fmt.Println("Invalid option.")
	}
}

func displayFullTask(t Task) {
	fmt.Println(colorText(strings.Repeat("-", 60), cyan))
	fmt.Println(colorText("Full Task Details", magenta))
	fmt.Println(colorText(strings.Repeat("-", 60), cyan))
	fmt.Printf("%s: %s\n", colorText("Name", yellow), t.Name)
	fmt.Printf("%s: %s\n", colorText("URL", yellow), t.URL)
	fmt.Printf("%s: %d\n", colorText("Priority", yellow), t.Priority)
	fmt.Printf("%s: %s\n", colorText("Note", yellow), t.Note)
	fmt.Printf("%s: %v\n", colorText("Completed", yellow), t.Completed)
	fmt.Println(colorText(strings.Repeat("-", 60), cyan))
}

func promptToContinue() {
	fmt.Println("\nPress Enter to continue...")
	readLine()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func wrapNote(note string, width int) []string {
	words := strings.Fields(note)
	var lines []string
	var line string

	for _, word := range words {
		if len(line)+len(word)+1 > width {
			lines = append(lines, padRight(line, width))
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}
	if line != "" {
		lines = append(lines, padRight(line, width))
	}

	if len(lines) == 0 {
		lines = []string{padRight("", width)}
	}

	return lines
}

func truncateWithDots(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max] // can't safely add "..."
	}
	return s[:max-3] + "..."
}

func displayTaskRow(t Task) {
	check := "[ ]"
	if t.Completed {
		check = colorText("[✔]", green)
	}

	name := truncateWithDots(t.Name, 12)
	nameColor := white
	if t.Completed {
		nameColor = green
	}
	nameColored := colorText(name, nameColor)

	url := truncateWithDots(t.URL, 30)
	urlColored := colorText(hyperlink(fmt.Sprintf("%-30s", url), t.URL), white)

	pr := fmt.Sprintf("%-3s", strconv.Itoa(t.Priority))
	if t.Priority <= 2 {
		pr = colorText(pr, red)
	}

	noteLines := wrapNote(t.Note, 30)
	for i, line := range noteLines {
		var nameCell, urlCell, prCell, noteCell, checkCell string
		if i == 0 {
			nameCell = nameColored
			urlCell = urlColored
			prCell = pr
			noteCell = colorText(padRight(line, 30), white)
			if t.Completed {
				noteCell = colorText(padRight(line, 30), green)
			}
			checkCell = check
		} else {
			nameCell = padRight("", 12)
			urlCell = padRight("", 30)
			prCell = padRight("", 3)
			noteCell = colorText(padRight(line, 30), white)
			checkCell = ""
		}

		fmt.Printf("%s %s %s %s %s %s %s %s %s\n",
			nameCell,
			colorText("|", cyan),
			urlCell,
			colorText("|", cyan),
			prCell,
			colorText("|", cyan),
			noteCell,
			colorText("|", cyan),
			checkCell,
		)
	}
}

func modifyTask(t *Task) {
	fmt.Println("\nWhich field to modify?")
	fmt.Println("[1] Name")
	fmt.Println("[2] URL")
	fmt.Println("[3] Priority")
	fmt.Println("[4] Note")
	fmt.Println("[5] Completed (toggle)")
	fmt.Print("Choose: ")

	switch readLine() {
	case "1":
		fmt.Print("New Name: ")
		t.Name = readLine()
	case "2":
		fmt.Print("New URL: ")
		t.URL = readLine()
	case "3":
		fmt.Print("New Priority: ")
		p, _ := strconv.Atoi(readLine())
		t.Priority = p
	case "4":
		fmt.Print("New Note: ")
		t.Note = readLine()
	case "5":
		t.Completed = !t.Completed
	default:
		fmt.Println("Invalid choice.")
	}
}

// --- Helpers ---

func addTask() {
	if len(categories) == 0 {
		fmt.Println("No categories. Add one first.")
		return
	}
	cat := selectCategory()
	if cat == nil {
		return
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
	cat.Tasks = append(cat.Tasks, newTask)
}

func sortCategory() {
	cat := selectCategory()
	if cat == nil {
		return
	}
	fmt.Print("Sort by [priority/completed/none]: ")
	sortBy := strings.ToLower(readLine())
	if sortBy != "priority" && sortBy != "completed" && sortBy != "none" {
		fmt.Println("Invalid sort option.")
	} else {
		cat.SortBy = sortBy
	}
}

func selectCategory() *Category {
	fmt.Println("Select category index:")
	for i, c := range categories {
		fmt.Printf("[%d] %s\n", i, c.Name)
	}
	fmt.Print("Index: ")
	ci, err := strconv.Atoi(readLine())
	if err != nil || ci < 0 || ci >= len(categories) {
		fmt.Println("Invalid index.")
		return nil
	}
	return &categories[ci]
}

func selectTask(cat *Category) *Task {
	fmt.Println("Select task index:")
	for i, t := range cat.Tasks {
		fmt.Printf("[%d] %s\n", i, t.Name)
	}
	fmt.Print("Index: ")
	ti, err := strconv.Atoi(readLine())
	if err != nil || ti < 0 || ti >= len(cat.Tasks) {
		fmt.Println("Invalid index.")
		return nil
	}
	return &cat.Tasks[ti]
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	var lines []string
	var line string

	for _, word := range words {
		if len(line)+len(word)+1 > width {
			lines = append(lines, line)
			line = word
		} else {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
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
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&categories)
	if err != nil {
		fmt.Println("Load error:", err)
	}
}

// --- ANSI Colors & Hyperlinks ---

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
