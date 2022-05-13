package pkg

import "fmt"

type FileContent struct {
	Name    string
	Date    string
	Content string
	Err     error
}

func (f FileContent) String() string {
	if f.Err != nil {
		return fmt.Sprintf("---------------\nName: %s\nError: %s\n---------------", f.Name, f.Err)
	}
	return fmt.Sprintf("---------------\nName: %s\nDate: %s\nContent: %s\n---------------", f.Name, f.Date, f.Content)
}

func SprintFileList(scheduleds []FileContent) string {
	result := "\n####################################\nSelected files\n"
	for _, scheduled := range scheduleds {
		result += fmt.Sprintln(scheduled)
	}
	result += "####################################\n"
	return result
}
