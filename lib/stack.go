package lib

type Stack struct {
	Items []interface{}
}

func (s Stack) Peek() interface{} {
	if len(s.Items) == 0 {
		return nil
	}

	lastItem := s.Items[len(s.Items)-1]
	return lastItem
}

func (s *Stack) Pop() interface{} {
	if len(s.Items) == 0 {
		return nil
	}

	lastItem := s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	return lastItem
}

func (s *Stack) Push(item interface{}) {
	s.Items = append(s.Items, item)
}
