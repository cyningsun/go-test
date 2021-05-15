package mutex

import "sync"

type Game struct {
	mtx       sync.Mutex
	bestScore int
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) HandlePlayer(p Player) error {
	for {
		score, err := p.NextScore()
		if err != nil {
			return err
		}
		g.mtx.Lock()
		if g.bestScore < score {
			g.bestScore = score
		}
		g.mtx.Unlock()
	}
}
