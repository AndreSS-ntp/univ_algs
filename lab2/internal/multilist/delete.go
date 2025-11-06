package multilist

import (
	"strings"
)

type nextGetter func(*Node) *Node
type nextSetter func(*Node, *Node)

func removeFromList(headPtr **Node, target *Node, getNext nextGetter, setNext nextSetter) (removed bool, wasTail bool) {
	if headPtr == nil || *headPtr == nil || target == nil {
		return false, false
	}

	if *headPtr == target {
		*headPtr = getNext(target)
		return true, getNext(target) == nil
	}

	prev := *headPtr
	for prev != nil && getNext(prev) != target {
		prev = getNext(prev)
	}
	if prev == nil {
		return false, false
	}
	// prev.Next = target.Next
	nextOfTarget := getNext(target)
	setNext(prev, nextOfTarget)
	return true, nextOfTarget == nil
}

func updateTail(ld *ListDescriptor, getNext nextGetter) {
	if ld.First == nil {
		ld.Last = nil
		return
	}

	cur := ld.First
	var prev *Node
	for cur != nil {
		prev = cur
		cur = getNext(cur)
	}
	ld.Last = prev
}

func (ml *MultiList) findNodeByLastName(last string) *Node {
	cur := ml.All.First
	lastTrim := strings.TrimSpace(strings.ToLower(last))
	for cur != nil {
		if strings.TrimSpace(strings.ToLower(cur.Data.LastName)) == lastTrim {
			return cur
		}
		cur = cur.NextAll
	}
	return nil
}

func (ml *MultiList) DeleteByLastName(last string) bool {
	target := ml.findNodeByLastName(last)
	if target == nil {
		return false
	}

	_, wasTail := removeFromList(&ml.All.First, target, func(n *Node) *Node { return n.NextAll }, func(n *Node, nx *Node) { n.NextAll = nx })
	if wasTail {
		updateTail(&ml.All, func(n *Node) *Node { return n.NextAll })
	} else {
		if ml.All.First == nil {
			ml.All.Last = nil
		}
	}

	removed, wasTail := removeFromList(&ml.Excellent.First, target, func(n *Node) *Node { return n.NextExcellent }, func(n *Node, nx *Node) { n.NextExcellent = nx })
	if removed {
		if wasTail {
			updateTail(&ml.Excellent, func(n *Node) *Node { return n.NextExcellent })
		} else if ml.Excellent.First == nil {
			ml.Excellent.Last = nil
		}
	}

	removed, wasTail = removeFromList(&ml.Distinct.First, target, func(n *Node) *Node { return n.NextDistinct }, func(n *Node, nx *Node) { n.NextDistinct = nx })
	if removed {
		if wasTail {
			updateTail(&ml.Distinct, func(n *Node) *Node { return n.NextDistinct })
		} else if ml.Distinct.First == nil {
			ml.Distinct.Last = nil
		}
	}

	removed, wasTail = removeFromList(&ml.OutOfTown.First, target, func(n *Node) *Node { return n.NextOutOfTown }, func(n *Node, nx *Node) { n.NextOutOfTown = nx })
	if removed {
		if wasTail {
			updateTail(&ml.OutOfTown, func(n *Node) *Node { return n.NextOutOfTown })
		} else if ml.OutOfTown.First == nil {
			ml.OutOfTown.Last = nil
		}
	}

	removed, wasTail = removeFromList(&ml.NeedsDorm.First, target, func(n *Node) *Node { return n.NextNeedsDorm }, func(n *Node, nx *Node) { n.NextNeedsDorm = nx })
	if removed {
		if wasTail {
			updateTail(&ml.NeedsDorm, func(n *Node) *Node { return n.NextNeedsDorm })
		} else if ml.NeedsDorm.First == nil {
			ml.NeedsDorm.Last = nil
		}
	}

	return true
}

func (ml *MultiList) DeleteAllByLastName(last string) int {
	count := 0

	for {
		if ml.findNodeByLastName(last) == nil {
			break
		}
		ok := ml.DeleteByLastName(last)
		if !ok {
			break
		}
		count++
	}
	return count
}

func (ml *MultiList) DeleteAll() {
	ml.All.First, ml.All.Last = nil, nil
	ml.Excellent.First, ml.Excellent.Last = nil, nil
	ml.Distinct.First, ml.Distinct.Last = nil, nil
	ml.OutOfTown.First, ml.OutOfTown.Last = nil, nil
	ml.NeedsDorm.First, ml.NeedsDorm.Last = nil, nil
}
