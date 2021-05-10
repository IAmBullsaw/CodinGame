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

// fmt.Fprintln(os.Stderr, "Debug messages...")

/*
	--------------- cell
*/

type cell struct {
	index, richness int
	neighbours      []int
}

/*
	--------------- tree
*/

type tree struct {
	index, size               int
	isMine, isDormant, isSeed bool
}

/*
	--------------- shadow
*/

type shadow struct {
	index, originIndex, size, direction int
}

/*
	--------------- action
*/

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

/*
	--------------- game
*/
type game struct {
	day, nutrients, mySun, myScore, oppSun, oppScore int
	oppIsWaiting                                     bool
	cells                                            map[int]cell
	shadows                                          map[int][]shadow
	trees                                            map[int]tree
	possActions                                      map[actionType][]action
	costs                                            map[int]int
}

/*
	----- utils
*/

func (g *game) clear() {
	g.trees = make(map[int]tree)
	g.possActions = make(map[actionType][]action)
	g.shadows = make(map[int][]shadow)
}

func (g *game) printPossActions() {
	for _, pa := range g.possActions {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", pa))
	}
}

func (g *game) printShadows() {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", g.shadows))
}

/*
	----- getters
*/

func (g *game) getSunDirection() int {
	return g.day % 6
}

func (g *game) getSunDirectionTomorrow() int {
	return (g.day + 1) % 6
}

func (g *game) getCellAt(index int) cell {
	return g.cells[index]
}

func (g *game) getTreeAt(index int) tree {
	return g.trees[index]
}

func (g *game) getShadowsAt(index int) []shadow {
	return g.shadows[index]
}

func (g *game) getMyTrees() (trees map[int]tree) {
	for _, t := range g.trees {
		if t.isMine {
			trees[t.index] = t
		}
	}
	return
}

func (g *game) getGrowCost(currentSize int) int {
	return g.costs[currentSize+1]
}

func (g *game) getSeedCost() int {
	return g.costs[0]
}

func (g *game) getDefaultAction() action {
	return action{action: wait, debugMessage: "Default Action"}
}

func (g *game) getNumberOfSeeds() (seeds int) {
	for _, t := range g.trees {
		if t.isMine && t.isSeed {
			seeds++
		}
	}
	return
}

func (g *game) getNumberOfTrees() (number int) {
	for _, t := range g.trees {
		if t.isMine {
			number++
		}
	}
	return
}

func (g *game) getNumberOfTreesSize(size int) (number int) {
	for _, t := range g.trees {
		if t.isMine && t.size == size {
			number++
		}
	}
	return
}

func (g *game) getShadows(t tree) (shadows []shadow) {
	if t.isSeed {
		return shadows
	}
	for d, index := range g.getCellAt(t.index).neighbours {
		if index == -1 {
			continue
		}
		s := shadow{index: index, originIndex: t.index, direction: d, size: t.size}
		shadows = append(shadows, s)

		n := g.getCellAt(index)
		for i := 1; i < t.size; i++ {
			si := n.neighbours[d]
			if si == -1 {
				continue
			}
			s := shadow{index: si, originIndex: t.index, direction: d, size: t.size}
			shadows = append(shadows, s)
		}
	}
	return
}

/*
	----- issers, hassers and canners
*/

func (g *game) isShadowed(index int) bool {
	return len(g.shadows[index]) > 0
}

func (g *game) isNeighbourToTree(index int) bool {
	for _, t := range g.trees {
		if t.isMine {
			for _, i := range g.getCellAt(t.index).neighbours {
				if index == i {
					return true
				}
			}
		}
	}
	return false
}

func (g *game) canSeed() bool {
	return len(g.possActions[seed]) > 0
}

func (g *game) canGrow() bool {
	return len(g.possActions[grow]) > 0
}

func (g *game) canComplete() bool {
	return len(g.possActions[complete]) > 0
}

/*
	----- updaters
*/

func (g *game) updateShadows() {
	for _, t := range g.trees {
		if t.isSeed {
			continue
		}
		for _, sh := range g.getShadows(t) {
			g.shadows[sh.index] = append(g.shadows[sh.index], sh)
		}
	}
}

/*
	updates the cost map with the costs for growing a tree from size X to Y
	Growing a seed into a size 1 tree costs 1 sun point + the number of size 1 trees you already own.
	Growing a size 1 tree into a size 2 tree costs 3 sun points + the number of size 2 trees you already own.
	Growing a size 2 tree into a size 3 tree costs 7 sun points + the number of size 3 trees you already own.
*/
func (g *game) updateGrowCosts() {
	costs := map[int]int{0: 0, 1: 1, 2: 3, 3: 7}
	for _, t := range g.trees {
		costs[t.size]++
	}
	g.costs = costs
}

func (g *game) nextAction() action {
	g.printPossActions()
	if g.getNumberOfSeeds() < 1 && g.canSeed() {
		for _, s := range g.possActions[seed] {
			if !g.isNeighbourToTree(s.targetCellIndex) {
				return s
			}
		}
		return g.possActions[seed][0]
	}
	if g.canGrow() && g.getNumberOfTreesSize(3) < 3 {
		return g.possActions[grow][0]
	}
	if g.canComplete() && g.getNumberOfTrees() > 1 {
		return g.possActions[complete][0]
	}
	return g.getDefaultAction()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)

	// numberOfCells: 37
	var numberOfCells int
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &numberOfCells)

	var g game
	g.cells = map[int]cell{}

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
		g.cells[newCell.index] = newCell
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
				isSeed:    size == 0,
			}
			g.trees[newTree.index] = newTree
		}
		// numberOfPossibleActions: all legal actions
		var numberOfPossibleActions int
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &numberOfPossibleActions)

		for i := 0; i < numberOfPossibleActions; i++ {
			scanner.Scan()
			possibleAction := scanner.Text()
			a := parseActionString(possibleAction)
			g.possActions[a.action] = append(g.possActions[a.action], a)
		}

		g.updateGrowCosts()
		g.updateShadows()

		g.printShadows()

		if na := g.nextAction(); na.debugMessage == "" {
			fmt.Println(na.String() + " " + na.debugMessage)
		} else {
			fmt.Println(na.String())
		}

	}
}
