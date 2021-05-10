/*
	Plan:

	When to grow a tree?
		- when
			- we have too many of X size trees
			- we have too many trees
			- it increases the amount of sun we have

	When to plant a seed?
		- when
			- we only have 1 tree
			- we can plant in the center
			- we can shadow oppent tree(s)
			- it increases the amount of sun we have

	When to complete a tree?
		- when
			- we are close to finishing
			- it does not give any sun && not blocks opponents trees
			- opponent score is way higher(?)
			- it increases the amount of sun we have

	When to wait?
		- when
			- we can't take an action
			- we can't afford desired action (?)

	Where to plant a seed?
		- where
			- in the center
			- close to the center
			- edges for sun generators?


	Take the action that makes you the most sun gain each round

	--------------- TRASH BELOW

	Order of importance?
		- Grow if you can
		- Complete if you can
		- Plant if you can
*/

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

func (g *game) getCellAt(index int) (bool, cell) {
	for _, c := range g.board {
		if c.index == index {
			return true, c
		}
	}
	return false, cell{}
}

func (g *game) getTreeAt(index int) (bool, tree) {
	for _, t := range g.trees {
		if t.index == index {
			return true, t
		}
	}
	return false, tree{}
}

func (t tree) isSeed() bool {
	return t.size == 0
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

func (g *game) clear() {
	g.trees = nil
	g.possActions = nil
}

func (g *game) sunCost() map[int]int {
	costs := map[int]int{}
	for _, t := range g.trees {
		costs[t.size] += 1
	}
	return costs
}

func (g *game) getComplete() (bool, action) {
	for _, pa := range g.possActions {
		if pa.action == complete {
			return true, pa
		}
	}
	return false, action{}
}

func (g *game) getGrow() (bool, action) {
	grows := []action{}
	for _, pa := range g.possActions {
		if pa.action == grow {
			grows = append(grows, pa)
		}
	}

	// grow the most expensive one
	costs := g.sunCost()
	k := 0
	v := 0
	for kk, vv := range costs {
		if vv > v {
			k = kk
			v = vv
		}
	}
	for _, pa := range grows {
		if ok, t := g.getTreeAt(pa.targetCellIndex); ok && t.size == k {
			return true, pa
		}
	}

	return false, action{}
}

func (g *game) getSeed() (bool, action) {
	seeds := []action{}
	for _, pa := range g.possActions {
		if pa.action == seed {
			seeds = append(seeds, pa)
		}
	}

	// find out which seed has the most potential with nutrients
	pa := action{}
	bestNutrients := -1
	for _, s := range seeds {
		_, c := g.getCellAt(s.targetCellIndex)
		if c.richness > bestNutrients {
			bestNutrients = c.richness
			pa = s
		}
	}
	if bestNutrients != -1 {
		return true, pa
	}

	return false, action{}
}

func (g *game) defaultAction() action {
	return action{action: wait, debugMessage: "Default Action"}
}

func (g *game) numberOfSeeds() (seeds int) {
	for _, t := range g.trees {
		if t.isSeed() {
			seeds++
		}
	}
	return
}

func (g *game) myTrees() (trees []tree) {
	for _, t := range g.trees {
		if t.isMine {
			trees = append(trees, t)
		}
	}
	return
}

func (g *game) nextAction() action {
	g.printPossActions()
	ok, pa := g.getGrow()
	if ok && g.day < 20 {
		return pa
	}
	ok, complete := g.getComplete()
	if ok && len(g.myTrees()) > 3 || ok && g.day > 19 {
		return complete
	}
	ok, seed := g.getSeed()
	if ok && g.numberOfSeeds() < 1 {
		return seed
	}
	return g.defaultAction()
}

func (g *game) printPossActions() {
	for _, pa := range g.possActions {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", pa))
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
