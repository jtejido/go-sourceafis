package sourceafis

import (
	"fmt"
	"sort"
	"sourceafis/features"
	"sourceafis/matcher"
	"sourceafis/transparency"

	"github.com/fxamacker/cbor/v2"
)

type Transparency interface {
	Accepts(key string) bool
	Accept(key, mime string, data []byte) error
}

type DefaultTransparencyLogger struct {
	spi Transparency
}

func NewTransparencyLogger(taker Transparency) *DefaultTransparencyLogger {
	return &DefaultTransparencyLogger{
		spi: taker,
	}
}
func (t *DefaultTransparencyLogger) log(key, mime string, supplier func() ([]byte, error)) error {
	// t.logVersion(); dont version this yet
	if t.spi.Accepts(key) {
		data, err := supplier()
		if err != nil {
			return err
		}
		t.spi.Accept(key, mime, data)
	}

	return nil
}

func (t *DefaultTransparencyLogger) Log(key string, data interface{}) error {
	return t.log(key, "application/cbor", func() ([]byte, error) {
		return cbor.Marshal(data)
	})
}

func (t *DefaultTransparencyLogger) LogSkeleton(keyword string, skeleton *features.Skeleton) error {
	return t.Log(skeleton.T.String()+keyword, transparency.NewConsistentSkeleton(skeleton))
}

func (t *DefaultTransparencyLogger) LogRootPairs(count int, roots []*matcher.MinutiaPair) error {
	return t.Log("roots", transparency.Roots(count, roots))
}
func (t *DefaultTransparencyLogger) LogPairing(pairing *matcher.PairingGraph) error {
	return t.Log("pairing", transparency.NewConsistentPairingGraph(pairing))
}

func (t *DefaultTransparencyLogger) LogBestPairing(pairing *matcher.PairingGraph) error {
	return t.Log("best-pairing", transparency.NewConsistentPairingGraph(pairing))
}

func (t *DefaultTransparencyLogger) LogScore(score *matcher.ScoringData) error {
	return t.Log("score", score)
}

func (t *DefaultTransparencyLogger) LogBestScore(score *matcher.ScoringData) error {
	return t.Log("best-score", score)
}

func (t *DefaultTransparencyLogger) LogEdgeHash(hash map[int][]*features.IndexedEdge) error {
	keys := make([]int, 0, len(hash))
	for key := range hash {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	entries := make([]*transparency.ConsistentHashEntry, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, &transparency.ConsistentHashEntry{
			Key:   key,
			Edges: hash[key],
		})
	}

	return t.Log("edge-hash", entries)
}

func (t *DefaultTransparencyLogger) LogBestMatch(nth int) error {
	t.spi.Accept("best-match", "text/plain", []byte(fmt.Sprintf("%d", nth)))
	return nil
}
