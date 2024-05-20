package section_module

import (
	"slices"
	"testing"

	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestChangeIndex(t *testing.T) {
	testData := []domain.Section{
		{Index: 0},
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}

	testcases := []struct {
		desc     string
		from, to int
		expected []domain.Section
	}{
		{
			desc: "no change",
			from: 1, to: 1,
			expected: []domain.Section{
				{Index: 0},
				{Index: 1},
				{Index: 2},
				{Index: 3},
			},
		},
		{
			desc: "increased",
			from: 1, to: 3,
			expected: []domain.Section{
				{Index: 0},
				{Index: 3},
				{Index: 1},
				{Index: 2},
			},
		},
		{
			desc: "decreased",
			from: 2, to: 0,
			expected: []domain.Section{
				{Index: 1},
				{Index: 2},
				{Index: 0},
				{Index: 3},
			},
		},
	}

	for _, tc := range testcases {
		tmp := slices.Clone(testData)

		changeIndex(tmp, tc.from, tc.to)

		assert.Equal(t, tc.expected, tmp)
	}
}
