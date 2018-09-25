package score

import (
	"github.com/faiface/pixel/text"
)

type ScoreKeeper map[string]int

func (k ScoreKeeper) Add(key string, val int) {
	k[key] += val
}

func (k ScoreKeeper) Get(key string) int {
	s, ok := k[key]
	if !ok {
		return 0
	}
	return s
}

type ScoreRenderer struct {
	renderFuncs []func(*text.Text, ScoreKeeper)
}

func (s *ScoreRenderer) Render(context *text.Text, scoreKeeper ScoreKeeper) {
	for _, f := range s.renderFuncs {
		f(context, scoreKeeper)
	}
}

func (s *ScoreRenderer) RenderFunc(f func(*text.Text, ScoreKeeper)) {
	s.renderFuncs = append(s.renderFuncs, f)
}

func NewScoreRenderer() *ScoreRenderer {
	return &ScoreRenderer{}
}
