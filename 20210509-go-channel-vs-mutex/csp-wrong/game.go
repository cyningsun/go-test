package csp

type Game struct {
	bestScore int
	scores    chan int
}

func (g *Game) run() {
	for score := range g.scores {
		if g.bestScore < score {
			g.bestScore = score
		}
	}
}

func (g *Game) HandlePlayer(p Player) error {
	for {
		score, err := p.NextScore()
		if err != nil {
			return err
		}
		g.scores <- score
	}
}

func NewGame() (g *Game) {
	g = &Game{
		bestScore: 0,
		scores:    make(chan int),
	}
	go g.run()
	return g
}
