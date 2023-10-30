package main

import (
	"math/rand"
	"time"
    "sync"
    "fmt"
    "strings"
)

//moves (0, 1) - down (0, -1) - up ...
var possibleMoves = [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type Traveler struct {
	ID  int
    moving bool
    currentNode *Node
}

type Node struct {
	Pos               [2]int
	visitedUp         bool
	visitedRight      bool
    visitedDown       bool
    visitedLeft       bool
    currentTraveler   *Traveler
	addTraveler   chan *Traveler
}

type Grid struct {
	Width, Height int
	GridMap       [][]Node
    Stopped       bool
    travelerCount int
    travelerCountLock sync.Mutex
}

func (g *Grid) runNode(n *Node){
    for{
        if(g.Stopped){
            //wait a second
            time.Sleep(time.Second)
            continue
        }
        //handle incoming traveler
        select{
            case t:=<-n.addTraveler:
                if(n.currentTraveler == nil){
                    oldX := t.currentNode.Pos[0]
                    oldY := t.currentNode.Pos[1]
                    newX := n.Pos[0]
                    newY := n.Pos[1]
                    fmt.Println("Traveler", t.ID, "moved from", oldX, oldY, "to", newX, newY)
                    if(oldX == newX+1){
                        n.visitedRight = true
                    }
                    if(oldX == newX-1){
                        n.visitedLeft = true
                    }
                    if(oldY == newY+1){
                        n.visitedUp = true
                    }
                    if(oldY == newY-1){
                        n.visitedDown = true
                    }
                    n.currentTraveler = t
                    t.currentNode.currentTraveler = nil
                    t.currentNode = n

                }else{
                    t.moving = false
                }
            default:
        }
        //handel current traveler
        shouldMove := rand.Intn(2) == 0
        if n.currentTraveler != nil && shouldMove{
            move := possibleMoves[rand.Intn(4)]
            moveX := n.Pos[0] + move[0]
            moveY := n.Pos[1] + move[1]
            if(moveX >= 0 && moveX < g.Width && moveY >= 0 && moveY < g.Height && !n.currentTraveler.moving){
                n.currentTraveler.moving = true
                select{
                case g.GridMap[moveX][moveY].addTraveler <- n.currentTraveler:
                default:
                    n.currentTraveler.moving = false
                }
            }
        }
        if(n.currentTraveler == nil){
            addTraveler := rand.Intn(10) == 0
            if(!addTraveler){
                continue
            }
            g.travelerCountLock.Lock()
            if(g.travelerCount == 25){
                g.travelerCountLock.Unlock()
                continue
            }
            t:=Traveler{ID: g.travelerCount, moving: false, currentNode: n}
            n.currentTraveler = &t
            g.travelerCount++
            g.travelerCountLock.Unlock()
        }else{
            time.Sleep(time.Millisecond*100)
        }
    }
}

func (g *Grid) TakePhoto() {
    g.Stopped = true
    time.Sleep(time.Millisecond*10)
    fmt.Println("Grid:")
    fmt.Println("  00 01 02 03 04 05 06 07 08 09")
    for i := 0; i < 2*g.Height-1; i++ {
        if(i%2 == 0){
            fmt.Print(i/2, " ")
        }else{
            fmt.Print("  ")
        }
        for j := 0; j < 2*g.Width-1; j++ {
            if(i%2 != 0 && j%2 != 0){
                fmt.Print("+")
            }else if(i%2 != 0){
                if(g.GridMap[(j+1)/2][i/2].visitedUp || g.GridMap[(j-1)/2][i/2].visitedDown){
                    fmt.Print("==")
                }else{
                    fmt.Print("--")
                }
            }else if(j%2 != 0){
                if(g.GridMap[j/2][(i+1)/2].visitedRight || g.GridMap[j/2][(i-1)/2].visitedLeft){
                    fmt.Print(";")
                }else{
                    fmt.Print("|")
                }
            }else{
                if(g.GridMap[j/2][i/2].currentTraveler == nil){
                    fmt.Print("  ")
                }else{
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
		Width:     width,
		Height:    height,
		GridMap:   make([][]Node, width),
        Stopped:   false,
        travelerCount: 0,
        travelerCountLock: sync.Mutex{},
	}
	var wg sync.WaitGroup
    for i := 0; i < width; i++ {
        g.GridMap[i] = make([]Node, height)
        for j := 0; j < height; j++ {
            g.GridMap[i][j] = Node{
                Pos:         [2]int{i, j},
                visitedUp:   false,
                visitedRight:false,
                visitedDown: false,
                visitedLeft: false,
                currentTraveler: nil,
                addTraveler: make(chan *Traveler, 1),
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
    for i:=0; i<10; i++{
        g.TakePhoto()
        time.Sleep(time.Second)
    }
    wg.Wait() // wait for all goroutines to finish
}
