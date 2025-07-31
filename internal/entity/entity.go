package entity

type Task struct {
	Link       []string
	Status     string
	ArchiveUrl string
	Err        map[string]error
}
