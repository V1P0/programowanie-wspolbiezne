package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
    "strings"
)

type Traveler struct {
    ID    int
    Pos   [2]int
}

type Grid struct {
    Width, Height int
    Travelers     []Traveler
    GridMap       [][]int
    mu            sync.Mutex
    visitedEdges  map[[2]int][4]bool
}



func (g *Grid) MoveTraveler(t *Traveler) {
    moves := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

    rand.Seed(time.Now().UnixNano())
    moveDir := rand.Intn(len(moves))
    move := moves[moveDir]
    newPos := [2]int{t.Pos[0] + move[0], t.Pos[1] + move[1]}
    g.mu.Lock()
    if newPos[0] >= 0 && newPos[0] < g.Width && newPos[1] >= 0 && newPos[1] < g.Height && g.GridMap[newPos[0]][newPos[1]] == 0 {
        g.GridMap[t.Pos[0]][t.Pos[1]] = 0
        g.GridMap[newPos[0]][newPos[1]] = t.ID
        t.Pos = newPos
        temp := g.visitedEdges[t.Pos]
        temp[moveDir] = true
        g.visitedEdges[t.Pos] = temp
    }
    g.mu.Unlock()
}
func (g *Grid) AddTraveler() {
	for {
		pos := [2]int{rand.Intn(g.Width), rand.Intn(g.Height)}
		g.mu.Lock()
		if g.GridMap[pos[0]][pos[1]] == 0 {
			i := len(g.Travelers) + 1
			t := Traveler{ID: i, Pos: pos}
			g.GridMap[t.Pos[0]][t.Pos[1]] = t.ID
			g.Travelers = append(g.Travelers, t)
			go func(tr *Traveler) {
				for {
					time.Sleep(1 * time.Second)
					g.MoveTraveler(tr)
				}
			}(&g.Travelers[i-1])
            g.mu.Unlock()
			break
		}else{
            g.mu.Unlock()
        }
	}
}

func (g *Grid) TakePhoto() {
    // Print grid
    g.mu.Lock()
    fmt.Println("Grid:")
    fmt.Println(" 00 01 02 03 04 05 06 07 08 09")
    for i := 0; i < 2*g.Width-1; i++ {
        if(i%2 == 0){
            fmt.Print(i/2)
        }else{
            fmt.Print(" ")
        }
        for j := 0; j < 2*g.Height-1; j++ {
            if(i%2 != 0 && j%2 != 0){
                fmt.Print("+")
            }else if(i%2 != 0){
                if(g.visitedEdges[[2]int{(i+1)/2, (j)/2}][1] || g.visitedEdges[[2]int{(i-1)/2, (j)/2}][0]){
                    fmt.Print("==")
                }else{
                    fmt.Print("--")
                }
            }else if(j%2 != 0){
                if(g.visitedEdges[[2]int{(i)/2, (j+1)/2}][3] || g.visitedEdges[[2]int{(i)/2, (j-1)/2}][2]){
                    fmt.Print(";")
                }else{
                    fmt.Print("|")
                }
            }else{
                if(g.GridMap[i/2][j/2] == 0){
                    fmt.Print("  ")
                }else{
                    fmt.Print(strings.Repeat(" ", 2-len(fmt.Sprint(g.GridMap[i/2][j/2]))), g.GridMap[i/2][j/2])
                }
            }
        }
        fmt.Println()
    }
    for i := 0; i < g.Height; i++ {
        for j := 0; j < g.Width; j++ {
            g.visitedEdges[[2]int{i, j}] = [4]bool{false, false, false, false}
        }
    }
    g.mu.Unlock()
}

func main() {
    // Initialize grid
    g := &Grid{
        Width:  10,
        Height: 10,
        GridMap: make([][]int, 10),
        visitedEdges: make(map[[2]int][4]bool),
    }
    for i := range g.GridMap {
        g.GridMap[i] = make([]int, 10)
    }

    // Generate random travelers
    rand.Seed(time.Now().UnixNano())


	go func() {
		for {
			g.AddTraveler()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			if(g.Width*g.Height == len(g.Travelers)){
				break
			}
		}
	}()


    // Take initial snapshot
    g.TakePhoto()

    // Loop: take snapshot at intervals
    for {
        time.Sleep(5 * time.Second)
        g.TakePhoto()
    }
}