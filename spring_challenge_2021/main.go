package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cell struct {
	index, richness int
	neighbours      []int
}

type tree struct {
	index, size       int
	isMine, isDormant bool
}

type actionType int

const (
	wait actionType = iota
	seed
	grow
	complete
)

func (a actionType) String() string {
	return [...]string{"WAIT", "SEED", "GROW", "COMPLETE"}[a]
}

type action struct {
	action                           actionType
	targetCellIndex, originCellIndex int
	debugMessage                     string
}

func parseActionString(actionString string) (a action) {
	s := strings.Split(actionString, " ")
	switch s[0] {
	case wait.String():
		a.action = wait
	case seed.String():
		a.action = seed
		a.originCellIndex, _ = strconv.Atoi(s[1])
		a.targetCellIndex, _ = strconv.Atoi(s[2])
	case grow.String():
		a.action = grow
		a.targetCellIndex, _ = strconv.Atoi(s[1])
	case complete.String():
		a.action = complete
		a.targetCellIndex, _ = strconv.Atoi(s[1])
	default:
		panic("Not a valid actionType")
	}

	return
}

func (a action) String() string {
	switch a.action {
	case wait:
		return "WAIT"
	case seed:
		return fmt.Sprintf("SEED %d %d", a.originCellIndex, a.targetCellIndex)
	default:
		return fmt.Sprintf("%s %d", a.action.String(), a.targetCellIndex)
	}
}

type game struct {
	day, nutrients, mySun, myScore, oppSun, oppScore int
	oppIsWaiting                                     bool
	board                                            []cell
	trees                                            []tree
	possActions                                      []action
}

func (g *game) removeActionAt(i int) {
	g.possActions[i] = g.possActions[len(g.possActions)-1]
	g.possActions[len(g.possActions)-1] = action{}
	g.possActions = g.possActions[:len(g.possActions)-1]
}

func (g *game) clear() {
	g.trees = nil
	g.possActions = nil
}

func (g *game) nextAction() action {
	g.printPossActions()
	for _, pa := range g.possActions {
		if pa.action == complete && len(g.trees) > 1 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Comp: %v", pa))
			return pa
		} else if pa.action == grow && len(g.trees) > 3 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Grow: %v", pa))
			return pa
		} else if pa.action == seed {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Seed: %v", pa))
			return pa
		}
	}
	g.possActions[0].debugMessage = "Took first Action"
	return g.possActions[0]
}

func (g *game) printPossActions() {
	for _, pa := range g.possActions {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%v", pa))
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// numberOfCells: 37
	var numberOfCells int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &numberOfCells)

	var g game

	for i := 0; i < numberOfCells; i++ {
		// index: 0 is the center cell, the next cells spiral outwards
		// richness: 0 if the cell is unusable, 1-3 for usable cells
		// neigh0: the index of the neighbouring cell for each direction
		var index, richness, neigh0, neigh1, neigh2, neigh3, neigh4, neigh5 int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &index, &richness, &neigh0, &neigh1, &neigh2, &neigh3, &neigh4, &neigh5)
		newCell := cell{
			index:      index,
			richness:   richness,
			neighbours: []int{neigh0, neigh1, neigh2, neigh3, neigh4, neigh5},
		}
		g.board = append(g.board, newCell)
	}
	for {
		g.clear()
		// day: the game lasts 24 days: 0-23
		var day int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &day)
		g.day = day

		// nutrients: the base score you gain from the next COMPLETE action
		var nutrients int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &nutrients)
		g.nutrients = nutrients

		// sun: your sun points
		// score: your current score
		var sun, score int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &sun, &score)
		g.mySun = sun
		g.myScore = score

		// oppSun: opponent's sun points
		// oppScore: opponent's score
		// oppIsWaiting: whether your opponent is asleep until the next day
		var oppSun, oppScore int
		var oppIsWaiting bool
		var _oppIsWaiting int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &oppSun, &oppScore, &_oppIsWaiting)
		oppIsWaiting = _oppIsWaiting != 0
		g.oppSun = oppSun
		g.oppScore = oppScore
		g.oppIsWaiting = oppIsWaiting

		// numberOfTrees: the current amount of trees
		var numberOfTrees int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &numberOfTrees)

		for i := 0; i < numberOfTrees; i++ {
			// cellIndex: location of this tree
			// size: size of this tree: 0-3
			// isMine: 1 if this is your tree
			// isDormant: 1 if this tree is dormant
			var cellIndex, size int
			var isMine, isDormant bool
			var _isMine, _isDormant int
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &cellIndex, &size, &_isMine, &_isDormant)
			isMine = _isMine != 0
			isDormant = _isDormant != 0

			newTree := tree{
				index:     cellIndex,
				size:      size,
				isMine:    isMine,
				isDormant: isDormant,
			}
			g.trees = append(g.trees, newTree)
		}
		// numberOfPossibleActions: all legal actions
		var numberOfPossibleActions int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &numberOfPossibleActions)

		for i := 0; i < numberOfPossibleActions; i++ {
			scanner.Scan()
			possibleAction := scanner.Text()
			g.possActions = append(g.possActions, parseActionString(possibleAction))
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		if na := g.nextAction(); na.debugMessage == "" {
			fmt.Println(na.String() + " " + na.debugMessage)
		} else {
			fmt.Println(na.String())
		}

	}
}
