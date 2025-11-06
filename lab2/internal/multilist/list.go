package multilist

import "multilist/internal/model"

type Node struct {
	Data *model.Applicant

	NextAll       *Node
	NextExcellent *Node
	NextDistinct  *Node
	NextOutOfTown *Node
	NextNeedsDorm *Node
}

type ListDescriptor struct {
	First *Node
	Last  *Node
}

type MultiList struct {
	All       ListDescriptor
	Excellent ListDescriptor
	Distinct  ListDescriptor
	OutOfTown ListDescriptor
	NeedsDorm ListDescriptor

	UniversityCity string
}
