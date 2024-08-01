package tilemapping

import (
	"testing"

	"github.com/domino14/word-golib/config"
	"github.com/matryer/is"
	"github.com/stretchr/testify/assert"
)

func EnglishAlphabet() *TileMapping {
	ld, err := GetDistribution(config.DefaultConfig, "english")
	if err != nil {
		panic(err)
	}
	return ld.TileMapping()
}

func TestScoreOn(t *testing.T) {
	is := is.New(t)

	ld, err := EnglishLetterDistribution(config.DefaultConfig)
	is.NoErr(err)
	type racktest struct {
		rack string
		pts  int
	}
	testCases := []racktest{
		{"ABCDEFG", 16},
		{"XYZ", 22},
		{"??", 0},
		{"?QWERTY", 21},
		{"RETINAO", 7},
	}
	for _, tc := range testCases {
		r := RackFromString(tc.rack, ld.TileMapping())
		score := r.ScoreOn(ld)
		if score != tc.pts {
			t.Errorf("For %v, expected %v, got %v", tc.rack, tc.pts, score)
		}
	}
}

func TestRackFromString(t *testing.T) {
	alph := EnglishAlphabet()
	rack := RackFromString("AENPPSW", alph)

	expected := make([]int, alph.NumLetters())
	expected[1] = 1
	expected[5] = 1
	expected[14] = 1
	expected[16] = 2
	expected[19] = 1
	expected[23] = 1

	assert.Equal(t, expected, rack.LetArr)

}

func TestRackTake(t *testing.T) {
	alph := EnglishAlphabet()
	rack := RackFromString("AENPPSW", alph)
	rack.Take(MachineLetter(16))
	expected := make([]int, alph.NumLetters())
	expected[1] = 1
	expected[5] = 1
	expected[14] = 1
	expected[16] = 1
	expected[19] = 1
	expected[23] = 1

	assert.Equal(t, expected, rack.LetArr)

	rack.Take(MachineLetter(16))
	expected[16] = 0
	assert.Equal(t, expected, rack.LetArr)
}

func TestRackTakeAll(t *testing.T) {
	alph := EnglishAlphabet()
	rack := RackFromString("AENPPSW", alph)

	rack.Take(MachineLetter(16))
	rack.Take(MachineLetter(16))
	rack.Take(MachineLetter(1))
	rack.Take(MachineLetter(5))
	rack.Take(MachineLetter(14))
	rack.Take(MachineLetter(19))
	rack.Take(MachineLetter(23))
	expected := make([]int, alph.NumLetters())

	assert.Equal(t, expected, rack.LetArr)
}

func TestRackTakeAndAdd(t *testing.T) {
	alph := EnglishAlphabet()
	rack := RackFromString("AENPPSW", alph)

	rack.Take(MachineLetter(16))
	rack.Take(MachineLetter(16))
	rack.Take(MachineLetter(1))
	rack.Add(MachineLetter(1))

	expected := make([]int, alph.NumLetters())
	expected[1] = 1
	expected[5] = 1
	expected[14] = 1
	expected[19] = 1
	expected[23] = 1

	assert.Equal(t, expected, rack.LetArr)

}
