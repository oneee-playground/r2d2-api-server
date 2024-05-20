package section_module

import "github.com/oneee-playground/r2d2-api-server/internal/domain"

// changeIndex changes each section's field index.
// sections should be orderd by its field index.
// It assumes that from and to is valid index of sections
func changeIndex(sections []domain.Section, from, to int) {
	if from == to {
		return
	}

	// First change the current index into desird index.
	sections[from].Index = uint8(to)

	// Shift all the elements between from and to
	isIncreased := from < to
	if isIncreased {
		for i := from + 1; i <= to; i++ {
			sections[i].Index--
		}
	} else {
		for i := from - 1; i >= to; i-- {
			sections[i].Index++
		}
	}
}
