package domain

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusResolved, StatusClosed:
		return true
	default:
		return false
	}
}

var allowedTransitions = map[Status][]Status{
	StatusOpen:       {StatusInProgress, StatusClosed},
	StatusInProgress: {StatusResolved, StatusClosed},
	StatusResolved:   {StatusClosed, StatusInProgress},
	StatusClosed:     {},
}

func (s Status) CanTransitionTo(next Status) bool {
	for _, allowed := range allowedTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}
