package chain

type Collect struct {
	list  []ICollect
	ch    chan string
	index int
	data  *InputData
}

func NewCollect(ch chan string, lastVersion string, collect ...ICollect) *Collect {
	return &Collect{
		list: collect,
		ch:   ch,
		data: &InputData{
			Version: lastVersion,
		},
	}
}

func (l *Collect) Collect() {
	l.list[l.index].Title(l)

	for {
		select {
		case txt := <-l.ch:
			if l.Handle(txt) {
				return
			}
		}
	}
}

func (l *Collect) Handle(txt string) bool {
	var ok bool

	defer func() {
		if ok {
			l.index++
		} else {
			l.list[l.index].Title(l)
		}
	}()

	v := l.list[l.index]

	ok = v.Collect(txt, l.data)
	if !ok {
		return false
	}

	return l.index >= len(l.list)-1
}

func (l *Collect) GetData() [][]string {
	return [][]string{
		{
			"git", "tag", "-a", "v" + l.data.Version, "-m", l.data.Version,
		},
		{
			"git", "push", "origin", "v" + l.data.Version,
		},
	}
}
