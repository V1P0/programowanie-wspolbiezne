package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

//moves (0, 1) - down (0, -1) - up ...
var possibleMoves = [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type GridObject struct {
	Pos			[2]int
	ID          int
	moving      bool
	typeID      int
	lifeTime    int
	currentNode *Node
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

func (g* Grid) runTraveler(t *GridObject){
	for{
		if(t.isToBeRemoved){
			break
		}
		if(g.Stopped || t.moving){
			time.Sleep(time.Second)
			continue
		}else{
			if(rand.Intn(10) == 0){
				direction := rand.Int(4)
				move := possibleMoves[direction]
				moveX := move[0] + t.Pos[0]
				moveY := move[1] + t.Pos[1]

				if(moveX < 0 || moveX >= g.Width || moveY < 0 || moveY >= g.Height){
					continue
				}
				t.moving = travelerCount
				select{
				case g.GridMap[moveX][moveY].addTraveler <- t:
				default:
					t.moving = false
				}
				
			}
			time.Sleep(time.Second)
		}
	}
}

func (g *Grid) runNode(n *Node) {
	for {
		if(n.currentTraveler == nil){
			if(g.Stopped){
				//wait a second
				time.Sleep(time.Second)
				continue
			}
			select{
			case t<-n.addTraveler:
				n.currentTraveler = t
				n.currentTraveler.Pos = n.Pos
				n.currentTraveler.moving = false
				n.currentTraveler.currentNode.currentTraveler = nil
				n.currentTraveler.currentNode = n
			default:
				if(rand.Intn(10)==0){
					g.addTraveler(n)
				}
			}
			
		}else{
			select{
			case t<-n.addTraveler:
				if(n.currentTraveler.typeID == TravelerType){
					t.moving = false
				}else if(n.currentTraveler.typeID == WildType){
					
				}else{
					t.isToBeRemoved = true
					t.currentNode.currentTraveler = nil
					n.currentTraveler = nil
				}
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

func (g *Grid) addTraveler(n *Node){
	//random typeID
	typeID := rand.Intn(3) + 1
	if(typeID == 1){
		g.travelerCountLock.Lock()
		g.travelerCount++
		n.currentTraveler = &GridObject{
			Pos:         n.Pos,
			ID:          g.travelerCount,
			moving:      false,
			typeID:      TravelerType,
			lifeTime:    0,
			currentNode: n,
			isToBeRemoved: false
		}
		g.travelerCountLock.Unlock()
		go g.runTraveler(n.currentTraveler)
	}else if(typeID == 2){
		lifetime := rand.Intn(10) + 1
		n.currentTraveler = &GridObject{
			Pos:         n.Pos,
			ID:          -1,
			moving:      false,
			typeID:      WildType,
			lifeTime:    lifetime,
			currentNode: n,
			isToBeRemoved: false
		}
	}else{
		lifetime := rand.Intn(10) + 1
		n.currentTraveler = &GridObject{
			Pos:         n.Pos,
			ID:          0,
			moving:      false,
			typeID:      Danger,
			lifeTime:    lifetime,
			currentNode: n,	
			isToBeRemoved: false
		}
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
				if g.GridMap[(j+1)/2][i/2].visitedUp || g.GridMap[(j-1)/2][i/2].visitedDown {
					fmt.Print("==")
				} else {
					fmt.Print("--")
				}
			} else if j%2 != 0 {
				if g.GridMap[j/2][(i+1)/2].visitedRight || g.GridMap[j/2][(i-1)/2].visitedLeft {
					fmt.Print(";")
				} else {
					fmt.Print("|")
				}
			} else {
				if g.GridMap[j/2][i/2].currentTraveler == nil {
					fmt.Print("  ")
				} else {
					fmt.Print(strings.Repeat(" ", 2-len(fmt.Sprint(g.GridMap[j/2][i/2].currentTraveler.ID))), g.GridMap[j/2][i/2].currentTraveler.ID)
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
		GridMap:           make([][]Node, width),
		Stopped:           false,
		travelerCount:     0,
		travelerCountLock: sync.Mutex{},
	}
	var wg sync.WaitGroup
	for i := 0; i < width; i++ {
		g.GridMap[i] = make([]Node, height)
		for j := 0; j < height; j++ {
			g.GridMap[i][j] = Node{
				Pos:             [2]int{i, j},
				visitedUp:       false,
				visitedRight:    false,
				visitedDown:     false,
				visitedLeft:     false,
				currentTraveler: nil,
				addTraveler:     make(chan *Traveler, 1),
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				g.runNode(&g.GridMap[i][j])
			}(i, j)
		}
	}
	for i := 0; i < 10; i++ {
		g.TakePhoto()
		time.Sleep(time.Second)
	}
	wg.Wait() // wait for all goroutines to finish
}
