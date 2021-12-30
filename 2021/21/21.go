package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/zigdon/adventofcode/common"
)

type die struct {
	roll   func() int
	rolled int
}

func (d *die) cast(n int) (int, []int) {
	r := 0
	rs := []int{}
	for ; n > 0; n-- {
		d.rolled++
		v := d.roll()
		r += v
		rs = append(rs, v)
	}
	return r, rs
}

func deterministict() *die {
	n := 0
	return &die{
		roll: func() int {
			n++
			if n > 100 {
				n = 1
			}
			return n
		}}
}

func dirac(id int64) *die {
	return &die{
		roll: func() int {
			r := id % 3
			id /= 3
			return int(r) + 1
		},
	}
}

type game struct {
	pos    []int
	score  []int
	target int
}

func newGame(players int, pos []int, target int) *game {
	g := &game{
		pos:    pos,
		score:  make([]int, len(pos)),
		target: target,
	}

	return g
}

func (g *game) loserScore() int {
	for _, s := range g.score {
		if s < g.target {
			return s
		}
	}

	return -1
}

func (g *game) winner() int {
	for i, s := range g.score {
		if s >= g.target {
			return i
		}
	}

	return -1
}

func (g *game) play(player, rolls int, d *die) bool {
	v, _ := d.cast(rolls)
	g.pos[player] += v
	for g.pos[player] > 10 {
		g.pos[player] -= 10
	}
	g.score[player] += g.pos[player]
	// log.Printf("%d rolled %d, landing at %d, adding up to a score of %d",
	//	player, v, g.pos[player], g.score[player])
	return g.score[player] >= g.target
}

func (g *game) runGame(d *die) int {
	for {
		if g.play(0, 3, d) {
			return 0
		}
		if g.play(1, 3, d) {
			return 1
		}
	}
}

func (g *game) clone() *game {
	ng := &game{
		pos:    []int{g.pos[0], g.pos[1]},
		score:  []int{g.score[0], g.score[1]},
		target: g.target,
	}

	return ng
}

func readFile(path string) (int, int) {
	data := common.ReadTransformedFile(
		path,
		common.IgnoreBlankLines,
		common.SplitWords,
	)

	return common.MustInt(data[0].([]string)[4]),
		common.MustInt(data[1].([]string)[4])
}

type pos struct {
	a, b int
}

type key struct {
	id             int64
	pA, pB, sA, sB int
}

func (k key) String() string {
	return fmt.Sprintf("%d", k.id)
}

type wins struct {
	a, b int64
}

var cache map[key]wins
var lastPrint time.Time

func allGames(g *game, gameID int64) wins {
	k := key{gameID, g.pos[0], g.pos[1], g.score[0], g.score[1]}
	if w, ok := cache[k]; ok {
		return w
	}

	d := dirac(gameID)

	w := wins{}
	ng := g.clone()
	if ng.play(0, 3, d) || ng.play(1, 3, d) {
		if ng.winner() == 0 {
			w.a++
		} else {
			w.b++
		}
		cache[k] = w

		return w
	}

	// run all possible die rolls
	opts := int64(math.Pow(3, 6))
	for id := int64(0); id < opts; id++ {
		got := allGames(ng, id)
		w.a += got.a
		w.b += got.b
	}
	cache[k] = w
	if time.Now().Sub(lastPrint) > time.Second*5 {
		log.Printf("%s: %d", k, len(cache))
		lastPrint = time.Now()
	}

	return w
}

func main() {
	a, b := readFile("input.txt")
	g := newGame(2, []int{a, b}, 1000)
	d := deterministict()
	g.runGame(d)
	log.Printf("loser score: %d, rolls: %d, res: %d", g.loserScore(), d.rolled, g.loserScore()*d.rolled)

	cache = make(map[key]wins)
	g = newGame(2, []int{a, b}, 21)
	log.Printf("all games: %v", allGames(g, 0))
}
