package kwg

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/domino14/word-golib/cache"
	"github.com/domino14/word-golib/tilemapping"
)

const (
	CacheKeyPrefix = "kwg:"
)

// CacheLoadFunc is the function that loads an object into the global cache.
func CacheLoadFunc(cfg map[string]any, key string) (interface{}, error) {
	lexiconName := strings.TrimPrefix(key, CacheKeyPrefix)
	if dataPath, ok := cfg["data-path"].(string); !ok {
		return nil, errors.New("could not find data-path in the configuration")
	} else {
		return LoadKWG(cfg, filepath.Join(dataPath, "lexica", "gaddag", lexiconName+".kwg"))
	}
}

func LoadKWG(cfg map[string]any, filename string) (*KWG, error) {
	log.Debug().Msgf("Loading %v ...", filename)
	file, filesize, err := cache.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// KWG is a simple map of nodes. There is no alphabet information in it,
	// so we must derive it from the filename, for now.
	lexfile := filepath.Base(filename)
	lexname, found := strings.CutSuffix(lexfile, ".kwg")
	if !found {
		return nil, errors.New("filename not in correct format")
	}

	kwg, err := ScanKWG(file, filesize)
	if err != nil {
		return nil, err
	}
	kwg.lexiconName = lexname
	lexname = strings.ToLower(lexname)
	var alphabetName string
	switch {
	case strings.HasPrefix(lexname, "nwl") ||
		strings.HasPrefix(lexname, "nswl") ||
		strings.HasPrefix(lexname, "twl") ||
		strings.HasPrefix(lexname, "owl") ||
		strings.HasPrefix(lexname, "csw") ||
		strings.HasPrefix(lexname, "america") ||
		strings.HasPrefix(lexname, "cel") ||
		strings.HasPrefix(lexname, "ecwl"):

		alphabetName = "english"

	// more cases here
	case strings.HasPrefix(lexname, "osps"):
		alphabetName = "polish"
	case strings.HasPrefix(lexname, "nsf"):
		alphabetName = "norwegian"
	case strings.HasPrefix(lexname, "fra"):
		alphabetName = "french"
	case strings.HasPrefix(lexname, "rd"):
		alphabetName = "german"
	case strings.HasPrefix(lexname, "disc"):
		alphabetName = "catalan"
	case strings.HasPrefix(lexname, "fise"):
		alphabetName = "spanish"
	default:
		return nil, errors.New("cannot determine alphabet from lexicon name " + lexname)
	}

	ld, err := tilemapping.NamedLetterDistribution(cfg, alphabetName)
	if err != nil {
		return nil, err
	}
	// we don't care about the distribution right now, just the tilemapping.
	kwg.alphabet = ld.TileMapping()

	return kwg, nil

}

// Get loads a named KWG from the cache or from a file
func Get(cfg map[string]any, name string) (*KWG, error) {

	key := CacheKeyPrefix + name
	obj, err := cache.Load(cfg, key, CacheLoadFunc)
	if err != nil {
		return nil, err
	}
	ret, ok := obj.(*KWG)
	if !ok {
		return nil, errors.New("could not read kwg from file")
	}
	return ret, nil
}
