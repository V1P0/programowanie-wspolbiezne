package main

import (
	"math/rand"
	"time"
)

//moves (0, 1) - down (0, -1) - up ...
const possibleMoves = 

type Traveler struct {
	ID  int
	Pos [2]int
}

type Node struct {
	Pos               [2]int
	visitedUp         bool
	visitedRight      bool
	currentTraveler   chan Traveler
}

type Grid struct {
	Width, Height int
	Travelers     []Traveler
	GridMap       [][]Node
    Stopped       bool
}

func (g *Grid) runNode(n *Node){
    for{
        if(g.Stopped){
            //wait a second
            
            continue
        }
        select{
            case t:=<-n.currentTraveler:
                move := possibleMoves[rand.Intn(len(possibleMoves))]
                moveX := move[0]
                moveY := move[1]
                if()
        }
    }
}

func prepareGrid(width int, height int) Grid {
	g := Grid{
		Width:     width,
		Height:    height,
		Travelers: make([]Traveler, 0),
		GridMap:   make([][]Node, width),
	}
	for i := 0; i < width; i++ {
		g.GridMap[i] = make([]Node, height)
		for j := 0; j < height; j++ {
			g.GridMap[i][j] = Node{
				Pos:             [2]int{i, j},
				visitedUp:       false,
				visitedRight:    false,
				currentTraveler: nil,
				currentTraveler: make(chan Traveler),
			}
		}
	}
	return g
}

// func (g *Grid) TakePhoto() {
//     fmt.Println("Grid:")
//     fmt.Println(" 00 01 02 03 04 05 06 07 08 09")
//     for i := 0; i < 2*g.Width-1; i++ {
//         if(i%2 == 0){
//             fmt.Print(i/2)
//         }else{
//             fmt.Print(" ")
//         }
//         for j := 0; j < 2*g.Height-1; j++ {
//             if(i%2 != 0 && j%2 != 0){
//                 fmt.Print("+")
//             }else if(i%2 != 0){
//                 if(g.visitedEdges[[2]int{(i+1)/2, (j)/2}][1] || g.visitedEdges[[2]int{(i-1)/2, (j)/2}][0]){
//                     fmt.Print("==")
//                 }else{
//                     fmt.Print("--")
//                 }
//             }else if(j%2 != 0){
//                 if(g.visitedEdges[[2]int{(i)/2, (j+1)/2}][3] || g.visitedEdges[[2]int{(i)/2, (j-1)/2}][2]){
//                     fmt.Print(";")
//                 }else{
//                     fmt.Print("|")
//                 }
//             }else{
//                 if(g.GridMap[i/2][j/2] == 0){
//                     fmt.Print("  ")
//                 }else{
//                     fmt.Print(strings.Repeat(" ", 2-len(fmt.Sprint(g.GridMap[i/2][j/2]))), g.GridMap[i/2][j/2])
//                 }
//             }
//         }
//         fmt.Println()
//     }
//     for i := 0; i < g.Height; i++ {
//         for j := 0; j < g.Width; j++ {
//             g.visitedEdges[[2]int{i, j}] = [4]bool{false, false, false, false}
//         }
//     }
// }

func main() {
	g := prepareGrid(10, 10)

	rand.Seed(time.Now().UnixNano())
}
