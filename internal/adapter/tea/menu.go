package tea

type menu struct {
	groupNames []string
	groupSizes []int
}

func (m *menu) addGroup(group string, size int) {
	m.groupNames = append(m.groupNames, group)
	m.groupSizes = append(m.groupSizes, size)
}

func (m menu) getGroup(cursor int) (group string, item int) {
	var count int
	for i := range m.groupSizes {
		a := count
		b := count + m.groupSizes[i]
		if a <= cursor && cursor < b {
			return m.groupNames[i], cursor - a
		}
		count = b
	}
	return "", 0
}

func (m menu) maxCursor() int {
	var count int
	for i := range m.groupSizes {
		count += m.groupSizes[i]
	}
	return count - 1
}
