package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var possibleMoves = [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type GridObject struct {
	Pos           [2]int
	ID            int
	moving        bool
	typeID        int
	lifeTime      int
	currentNode   *Node
	isToBeRemoved bool
}

const (
	TravelerType = 1
	WildType     = 2
	Danger       = 3
)

type Node struct {
	Pos             [2]int
	visitedUp       bool
	visitedRight    bool
	visitedDown     bool
	visitedLeft     bool
	currentTraveler *GridObject
	addTraveler     chan *GridObject
}

type Grid struct {
	Width, Height     int
	GridMap           [][]Node
	Stopped           bool
	travelerCount     int
	travelerCountLock sync.Mutex
}

func (g *Grid) runWild(w *GridObject) {
	creationTime := time.Now()
	for {
		if g.Stopped {
			time.Sleep(time.Second)
			continue
		}
		if creationTime.Add(time.Second * time.Duration(w.lifeTime)).Before(time.Now()) {
			w.currentNode.currentTraveler = nil
			break
		}
		time.Sleep(time.Millisecond)
	}
}
func (g *Grid) moveWild(w *GridObject, diffPos [2]int) {
	w.moving = true
	moved := false
	if w.Pos[0] > 0 && w.Pos[0]-1 != diffPos[0] {
		select {
		case g.GridMap[w.Pos[0]-1][w.Pos[1]].addTraveler <- w:
			moved = true
		default:
		}
	}
	if w.Pos[0] < g.Height-1 && w.Pos[0]+1 != diffPos[0] {
		select {
		case g.GridMap[w.Pos[0]+1][w.Pos[1]].addTraveler <- w:
			moved = true
		default:
		}
	}
	if w.Pos[1] > 0 && w.Pos[1]-1 != diffPos[1] {
		select {
		case g.GridMap[w.Pos[0]][w.Pos[1]-1].addTraveler <- w:
			moved = true
		default:
		}
	}
	if w.Pos[1] < g.Width-1 && w.Pos[1]+1 != diffPos[1] {
		select {
		case g.GridMap[w.Pos[0]][w.Pos[1]+1].addTraveler <- w:
			moved = true
		default:
		}
	}
	if !moved {
		w.moving = false
	}
}
func (g *Grid) runDanger(d *GridObject) {
	creationTime := time.Now()
	for {
		if g.Stopped {
			time.Sleep(time.Second)
			continue
		}
		if creationTime.Add(time.Second * time.Duration(d.lifeTime)).Before(time.Now()) {
			d.currentNode.currentTraveler = nil
			break
		}
		time.Sleep(time.Millisecond)
	}
}
func (g *Grid) runTraveler(t *GridObject) {
	for {
		if t.isToBeRemoved {
			break
		}
		if g.Stopped || t.moving {
			time.Sleep(time.Second)
			continue
		} else {
			if rand.Intn(3) == 0 {
				direction := rand.Intn(4)
				move := possibleMoves[direction]
				moveX := move[0] + t.Pos[0]
				moveY := move[1] + t.Pos[1]

				if moveX < 0 || moveX >= g.Height || moveY < 0 || moveY >= g.Width {
					continue
				}
				t.moving = true
				select {
				case g.GridMap[moveX][moveY].addTraveler <- t:
				default:
					t.moving = false
				}

			}
		}
		time.Sleep(time.Millisecond)
	}
}
func (g *Grid) runNode(n *Node) {
	for {
		if n.currentTraveler == nil {
			if g.Stopped {
				//wait a second
				time.Sleep(time.Second)
				continue
			}
			select {
			case t := <-n.addTraveler:
				if !t.moving {
					continue
				}
				if t.Pos[0] < n.Pos[0] {
					n.visitedUp = true
				}
				if t.Pos[0] > n.Pos[0] {
					n.visitedDown = true
				}
				if t.Pos[1] < n.Pos[1] {
					n.visitedLeft = true
				}
				if t.Pos[1] > n.Pos[1] {
					n.visitedRight = true
				}
				t.Pos = n.Pos
				t.moving = false
				t.currentNode.currentTraveler = nil
				t.currentNode = n
				n.currentTraveler = t
			default:
				if rand.Intn(10) == 0 {
					g.addTraveler(n)
				}
			}

		} else {
			select {
			case t := <-n.addTraveler:
				if n.currentTraveler == nil {
					n.addTraveler <- t
					continue
				}
				typeID := n.currentTraveler.typeID
				if typeID == TravelerType {
					t.moving = false
				} else if typeID == WildType {
					g.moveWild(n.currentTraveler, t.Pos)
					select {
					case n.addTraveler <- t:
					default:
						t.moving = false
					}
				} else {
					t.isToBeRemoved = true
					t.currentNode.currentTraveler = nil
					n.currentTraveler = nil
				}
			default:
			}
		}
		time.Sleep(time.Millisecond)
	}
}

func (g *Grid) addTraveler(n *Node) {
	//random typeID
	typeID := rand.Intn(3) + 1
	if typeID == 1 {
		g.travelerCountLock.Lock()
		if g.travelerCount == 25 {
			g.travelerCountLock.Unlock()
			return
		}
		g.travelerCount++
		n.currentTraveler = &GridObject{
			Pos:           n.Pos,
			ID:            g.travelerCount,
			moving:        false,
			typeID:        TravelerType,
			lifeTime:      0,
			currentNode:   n,
			isToBeRemoved: false,
		}
		g.travelerCountLock.Unlock()
		go g.runTraveler(n.currentTraveler)
	} else if typeID == 2 {
		lifetime := rand.Intn(5) + 1
		n.currentTraveler = &GridObject{
			Pos:           n.Pos,
			ID:            -1,
			moving:        false,
			typeID:        WildType,
			lifeTime:      lifetime,
			currentNode:   n,
			isToBeRemoved: false,
		}
		go g.runWild(n.currentTraveler)
	} else {
		lifetime := rand.Intn(5) + 1
		n.currentTraveler = &GridObject{
			Pos:           n.Pos,
			ID:            -2,
			moving:        false,
			typeID:        Danger,
			lifeTime:      lifetime,
			currentNode:   n,
			isToBeRemoved: false,
		}
		go g.runDanger(n.currentTraveler)
	}
}

func (g *Grid) TakePhoto() {
	g.Stopped = true
	time.Sleep(time.Millisecond * 10)
	fmt.Println("Grid:")
	fmt.Println("  00 01 02 03 04 05 06 07 08 09")
	for i := 0; i < 2*g.Height-1; i++ {
		if i%2 == 0 {
			fmt.Print(i/2, " ")
		} else {
			fmt.Print("  ")
		}
		for j := 0; j < 2*g.Width-1; j++ {
			if i%2 != 0 && j%2 != 0 {
				fmt.Print("+")
			} else if i%2 != 0 {
				if g.GridMap[(i+1)/2][j/2].visitedUp && g.GridMap[(i-1)/2][j/2].visitedDown {
					fmt.Print("--")
				} else if g.GridMap[(i+1)/2][j/2].visitedUp {
					fmt.Print("\\/")
				} else if g.GridMap[(i-1)/2][j/2].visitedDown {
					fmt.Print("/\\")
				} else {
					fmt.Print("  ")
				}

			} else if j%2 != 0 {
				if g.GridMap[i/2][(j+1)/2].visitedLeft && g.GridMap[i/2][(j-1)/2].visitedRight {
					fmt.Print("|")
				} else if g.GridMap[i/2][(j+1)/2].visitedLeft {
					fmt.Print(">")
				} else if g.GridMap[i/2][(j-1)/2].visitedRight {
					fmt.Print("<")
				} else {
					fmt.Print(" ")
				}
			} else {
				if g.GridMap[i/2][j/2].currentTraveler == nil {
					fmt.Print("  ")
				} else if g.GridMap[i/2][j/2].currentTraveler.typeID == WildType {
					fmt.Print("*", g.GridMap[i/2][j/2].currentTraveler.lifeTime)
				} else if g.GridMap[i/2][j/2].currentTraveler.typeID == Danger {
					fmt.Print("#", g.GridMap[i/2][j/2].currentTraveler.lifeTime)
				} else {
					fmt.Print(strings.Repeat(" ", 2-len(fmt.Sprint(g.GridMap[i/2][j/2].currentTraveler.ID))), g.GridMap[i/2][j/2].currentTraveler.ID)
				}
			}
		}
		fmt.Println()
	}
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			g.GridMap[i][j].visitedUp = false
			g.GridMap[i][j].visitedRight = false
			g.GridMap[i][j].visitedDown = false
			g.GridMap[i][j].visitedLeft = false
		}
	}
	g.Stopped = false
}

func main() {
	width := 10
	height := 10
	g := Grid{
		Width:             width,
		Height:            height,
		GridMap:           make([][]Node, height),
		Stopped:           false,
		travelerCount:     0,
		travelerCountLock: sync.Mutex{},
	}
	var wg sync.WaitGroup
	for i := 0; i < height; i++ {
		g.GridMap[i] = make([]Node, width)
		for j := 0; j < width; j++ {
			g.GridMap[i][j] = Node{
				Pos:             [2]int{i, j},
				visitedUp:       false,
				visitedRight:    false,
				visitedDown:     false,
				visitedLeft:     false,
				currentTraveler: nil,
				addTraveler:     make(chan *GridObject, 1),
			}
		}
	}
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				g.runNode(&g.GridMap[i][j])
			}(i, j)
		}
	}
	for {
		g.TakePhoto()
		time.Sleep(time.Second)
	}
}
