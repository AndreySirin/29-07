package entity

type Task struct {
	Link       []string
	Status     string
	ArchiveUrl string
	ErrorLoad  map[string]string
}
