package tea

type screen string

func (m *TeaModel) currentScreen() screen {
	return m.screens[len(m.screens)-1]
}

func (m *TeaModel) screenPush(screen screen) {
	m.screens = append(m.screens, screen)
}

func (m *TeaModel) screenPop() {
	if len(m.screens) > 1 {
		m.screens = m.screens[:len(m.screens)-1]
	}
}
