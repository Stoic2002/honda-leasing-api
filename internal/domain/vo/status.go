package vo

type TaskStatus string

const (
	StatusDraft      TaskStatus = "draft"
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "inprogress"
	StatusCompleted  TaskStatus = "completed"
	StatusCanceled   TaskStatus = "canceled"
	StatusActive     TaskStatus = "active"
	StatusApproved   TaskStatus = "approved"
)

func (s TaskStatus) String() string {
	return string(s)
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusPending, StatusInProgress, StatusCompleted, StatusCanceled, StatusActive:
		return true
	}
	return false
}
