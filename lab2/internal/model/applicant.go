package model

type Applicant struct {
	LastName  string
	Exam1     int
	Exam2     int
	Exam3     int
	Distinct  bool
	City      string
	NeedsDorm bool
}
