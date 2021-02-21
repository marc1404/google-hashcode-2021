package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var sampleSize = 100
var xorBinaryCache = make(map[string]int)

func main() {
	inputFileName := os.Args[1]
	inputFilePath := fmt.Sprintf("./input/%s", inputFileName)
	lines := readFile(inputFilePath)
	ingredientToIndex := make(map[string]int)
	ingredientIndex := 0
	pizzas := make([]*Pizza, len(lines)-1)
	headerLine := lines[0]
	headerParts := strings.Split(headerLine, " ")
	teams := make([]*Team, 0)

	for headerIndex, part := range headerParts[1:] {
		teamCount, _ := strconv.Atoi(part)
		teamSize := headerIndex + 2

		for i := 0; i < teamCount; i++ {
			teams = append(teams, &Team{teamSize, make([]*Pizza, 0), []int{}, false})
		}
	}

	for pizzaId, line := range lines[1:] {
		lineParts := strings.Split(line, " ")
		pizzas[pizzaId] = &Pizza{pizzaId, lineParts[1:], []int{}}

		for _, ingredient := range lineParts[1:] {
			_, hasIngredient := ingredientToIndex[ingredient]

			if hasIngredient {
				continue
			}

			ingredientToIndex[ingredient] = ingredientIndex
			ingredientIndex++
		}
	}

	for _, pizza := range pizzas {
		binaryCode := make([]int, len(ingredientToIndex))

		for _, ingredient := range pizza.ingredients {
			ingredientIndex := ingredientToIndex[ingredient]
			binaryCode[ingredientIndex] = 1
		}

		pizza.binaryCode = binaryCode
	}

	maxPizzaCount := len(pizzas)
	progress := float64(0)
	teamCount := len(teams)

	for i := 0; i < teamCount; i++ {
		team := teams[teamCount-1-i]
		pizzas = deliverPizzas(team, pizzas)

		if len(pizzas) == 0 {
			break
		}

		newProgress := math.Floor(float64(maxPizzaCount-len(pizzas)) / float64(maxPizzaCount) * 100)

		if newProgress-progress >= 10 {
			progress = newProgress

			fmt.Println(fmt.Sprintf("%v: %v%%", inputFileName, progress))
		}
	}

	outputFilePath := fmt.Sprintf("./output/%s", inputFileName)
	outputFile, _ := os.Create(outputFilePath)

	defer outputFile.Close()

	completeTeamsCount := 0

	for _, team := range teams {
		if team.isComplete {
			completeTeamsCount++
		}
	}

	outputFile.WriteString(fmt.Sprintf("%v\n", completeTeamsCount))

	for _, team := range teams {
		if !team.isComplete {
			continue
		}

		outputFile.WriteString(fmt.Sprintf("%v", team.size))

		for _, pizza := range team.pizzas {
			outputFile.WriteString(fmt.Sprintf(" %v", pizza.id))
		}

		outputFile.WriteString("\n")
	}

	outputFile.Sync()
	fmt.Println(fmt.Sprintf("%v: Done!", inputFileName))
}

func deliverPizzas(team *Team, pizzas []*Pizza) []*Pizza {
	if len(team.pizzas) == 0 {
		pizzas = deliverFirstPizza(team, pizzas)

		if len(pizzas) == 0 {
			return pizzas
		}
	}

	for {
		pizzas = deliverNextPizza(team, pizzas)

		if team.isComplete || len(pizzas) == 0 {
			return pizzas
		}
	}
}

func deliverFirstPizza(team *Team, pizzas []*Pizza) []*Pizza {
	pizza := pizzas[0]
	team.pizzas = append(team.pizzas, pizza)
	team.binaryCode = pizza.binaryCode

	return removePizza(pizzas, 0)
}

func deliverNextPizza(team *Team, pizzas []*Pizza) []*Pizza {
	bestPizza := pizzas[0]
	bestPizzaIndex := 0
	maxDifference := 0

	for i := 0; i < sampleSize; i++ {
		pizza, index := getRandomPizza(pizzas)
		difference := xorBinary(team.binaryCode, pizza.binaryCode)

		if difference > maxDifference {
			bestPizza = pizza
			bestPizzaIndex = index
			maxDifference = difference
		}
	}

	team.pizzas = append(team.pizzas, bestPizza)
	team.binaryCode = andBinary(team.binaryCode, bestPizza.binaryCode)

	if len(team.pizzas) == team.size {
		team.isComplete = true
	}

	return removePizza(pizzas, bestPizzaIndex)
}

func getRandomPizza(pizzas []*Pizza) (*Pizza, int) {
	randomIndex := rand.Intn(len(pizzas))

	return pizzas[randomIndex], randomIndex
}

func removePizza(pizzas []*Pizza, indexToRemove int) []*Pizza {
	pizzaCount := len(pizzas)
	pizzas[indexToRemove] = pizzas[pizzaCount-1]

	return pizzas[:pizzaCount-1]
}

func xorBinary(leftBinaryCode, rightBinaryCode []int) int {
	cacheKey := fmt.Sprintf("%v%v", leftBinaryCode, rightBinaryCode)
	cacheResult, hasCache := xorBinaryCache[cacheKey]

	if hasCache {
		return cacheResult
	}

	difference := 0

	for i, leftBit := range leftBinaryCode {
		if leftBit != rightBinaryCode[i] {
			difference++
		}
	}

	xorBinaryCache[cacheKey] = difference

	return difference
}

func andBinary(leftBinaryCode, rightBinaryCode []int) []int {
	for i, leftBit := range rightBinaryCode {
		if leftBit == rightBinaryCode[i] {
			leftBinaryCode[i] = 1
		}
	}

	return leftBinaryCode
}

type Pizza struct {
	id          int
	ingredients []string
	binaryCode  []int
}

type Team struct {
	size       int
	pizzas     []*Pizza
	binaryCode []int
	isComplete bool
}

func readFile(filePath string) []string {
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}
