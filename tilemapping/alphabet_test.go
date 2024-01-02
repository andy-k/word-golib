package tilemapping

import (
	"testing"

	"github.com/matryer/is"
)

func TestToMachineLetters(t *testing.T) {
	is := is.New(t)
	cat, err := GetDistribution(DefaultConfig, "catalan")
	is.NoErr(err)
	// AL·LOQUIMIQUES is only 10 tiles despite being 14 codepoints long.
	// A L·L O QU I M I QU E S
	mls, err := ToMachineLetters("AL·LOQUIMIQUES", cat.TileMapping())
	is.NoErr(err)
	is.Equal(mls, []MachineLetter{
		1, 13, 17, 19, 10, 14, 10, 19, 6, 21,
	})
	mls, err = ToMachineLetters("Al·lOQUIMIquES", cat.TileMapping())
	is.NoErr(err)
	is.Equal(mls, []MachineLetter{
		1, 13 | 0x80, 17, 19, 10, 14, 10, 19 | 0x80, 6, 21,
	})
	mls, err = ToMachineLetters("ARQUEGESSIU", cat.TileMapping())
	is.NoErr(err)
	is.Equal(mls, []MachineLetter{
		1, 20, 19, 6, 8, 6, 21, 21, 10, 23,
	})
	mls, err = ToMachineLetters("L·L", cat.TileMapping())
	is.NoErr(err)
	is.Equal(mls, []MachineLetter{13})
	mls, err = ToMachineLetters("L", cat.TileMapping())
	is.NoErr(err)
	is.Equal(mls, []MachineLetter{12})
}

func TestUV(t *testing.T) {
	is := is.New(t)
	cat, err := GetDistribution(DefaultConfig, "catalan")
	is.NoErr(err)

	uv := MachineWord([]MachineLetter{
		1, 13, 17, 19, 10, 14, 10, 19, 6, 21,
	}).UserVisible(cat.TileMapping())
	is.Equal(uv, "AL·LOQUIMIQUES")

	uv = MachineWord([]MachineLetter{
		1, 13 | 0x80, 17, 19, 10, 14, 10, 19 | 0x80, 6, 21,
	}).UserVisible(cat.TileMapping())
	is.Equal(uv, "Al·lOQUIMIquES")
}

func TestCts(t *testing.T) {
	is := is.New(t)
	eng, err := GetDistribution(DefaultConfig, "english")
	is.NoErr(err)
	is.Equal(eng.TileMapping().NumLetters(), uint8(27))
}

func TestIsVowel(t *testing.T) {
	is := is.New(t)
	eng, err := GetDistribution(DefaultConfig, "english")
	is.NoErr(err)
	is.True(MachineLetter(5).IsVowel(eng))
	is.True(MachineLetter(9).IsVowel(eng))
	is.True(!MachineLetter(0).IsVowel(eng))
	is.True(MachineLetter(1).IsVowel(eng))
	is.True(!MachineLetter(2).IsVowel(eng))
	is.True(!MachineLetter(25).IsVowel(eng))
	is.True(!MachineLetter(26).IsVowel(eng))
	is.True(MachineLetter(21).IsVowel(eng))
}
