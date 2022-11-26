package subscriber

type Criteria interface {
	Fits(flat Flat) bool
}
