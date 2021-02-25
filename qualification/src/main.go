package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Simulation struct {
	context         SimulationContext
	intersectionMap map[int]*Intersection
	streetMap       map[string]*Street
	carMap          map[int]*Car
}

type SimulationContext struct {
	duration          int
	intersectionCount int
	streetCount       int
	carCount          int
	bonusPoints       int
}

type Intersection struct {
	id              int
	incomingStreets []*Street
	outgoingStreets []*Street
	trafficLights   []TrafficLight
}

func (intersection *Intersection) getCarCount() int {
	carCount := 0

	for _, street := range intersection.incomingStreets {
		carCount += street.carCount
	}

	return carCount
}

func (intersection *Intersection) getIncomingMaxLength() int {
	maxLength := 0

	for _, street := range intersection.incomingStreets {
		if street.length > maxLength {
			maxLength = street.length
		}
	}

	return maxLength
}

type Street struct {
	name              string
	length            int
	startIntersection *Intersection
	endIntersection   *Intersection
	carCount          int
	carsAtEnd         []*Car
}

type Car struct {
	id          int
	streetCount int
	streets     []*Street
}

type TrafficLight struct {
	street   *Street
	duration int
}

func main() {
	inputFileName := os.Args[1]
	inputFilePath := fmt.Sprintf("./input/%v", inputFileName)
	simulation := parseInputFile(inputFilePath)

	for _, intersection := range simulation.intersectionMap {
		sort.Slice(intersection.incomingStreets, func(i, j int) bool {
			return len(intersection.incomingStreets[i].carsAtEnd) > len(intersection.incomingStreets[j].carsAtEnd)
		})

		for _, street := range intersection.incomingStreets {
			if street.carCount == 0 {
				continue
			}

			duration := int(math.Min(float64(simulation.context.duration), math.Max(1., float64(len(street.carsAtEnd)))))

			trafficLight := TrafficLight{
				street,
				duration,
			}
			intersection.trafficLights = append(intersection.trafficLights, trafficLight)
		}
	}

	writeOutputFile(simulation, inputFileName)
}

func writeOutputFile(simulation *Simulation, fileName string) {
	outputFilePath := fmt.Sprintf("./output/%v", fileName)
	outputFile, _ := os.Create(outputFilePath)

	defer outputFile.Close()

	intersectionCount := 0

	for _, intersection := range simulation.intersectionMap {
		if len(intersection.trafficLights) > 0 {
			intersectionCount++
		}
	}

	outputFile.WriteString(fmt.Sprintf("%v\n", intersectionCount))

	for _, intersection := range simulation.intersectionMap {
		if len(intersection.trafficLights) == 0 {
			continue
		}

		outputFile.WriteString(fmt.Sprintf("%v\n%v\n", intersection.id, len(intersection.trafficLights)))

		for _, trafficLight := range intersection.trafficLights {
			outputFile.WriteString(fmt.Sprintf("%v %v\n", trafficLight.street.name, trafficLight.duration))
		}
	}

	outputFile.Sync()
	fmt.Println(fmt.Sprintf("%v: Done!", fileName))
}

func parseInputFile(filePath string) *Simulation {
	lines := readFile(filePath)
	context := parseInputHeader(lines[0])
	simulation := &Simulation{
		context,
		make(map[int]*Intersection),
		make(map[string]*Street),
		make(map[int]*Car),
	}

	for i := 0; i < context.streetCount; i++ {
		line := lines[1+i]
		parseStreet(simulation, line)
	}

	for i := 0; i < context.carCount; i++ {
		line := lines[1+context.streetCount+i]
		parseCar(simulation, line)
	}

	return simulation
}

func parseInputHeader(header string) SimulationContext {
	parts := strings.Split(header, " ")

	duration, _ := strconv.Atoi(parts[0])
	intersectionCount, _ := strconv.Atoi(parts[1])
	streetCount, _ := strconv.Atoi(parts[2])
	carCount, _ := strconv.Atoi(parts[3])
	bonusPoints, _ := strconv.Atoi(parts[4])

	return SimulationContext{
		duration, intersectionCount, streetCount, carCount, bonusPoints,
	}
}

func parseStreet(simulation *Simulation, line string) {
	parts := strings.Split(line, " ")

	startIntersectionId, _ := strconv.Atoi(parts[0])
	endIntersectionId, _ := strconv.Atoi(parts[1])
	name := parts[2]
	length, _ := strconv.Atoi(parts[3])

	startIntersection := getOrCreateIntersection(simulation, startIntersectionId)
	endIntersection := getOrCreateIntersection(simulation, endIntersectionId)

	street := &Street{
		name, length, startIntersection, endIntersection, 0, make([]*Car, 0),
	}
	simulation.streetMap[name] = street

	startIntersection.outgoingStreets = append(startIntersection.outgoingStreets, street)
	endIntersection.incomingStreets = append(endIntersection.incomingStreets, street)
}

func parseCar(simulation *Simulation, line string) {
	parts := strings.Split(line, " ")
	streetCount, _ := strconv.Atoi(parts[0])
	car := createCar(simulation, streetCount)

	for i, streetName := range parts[1:] {
		street := simulation.streetMap[streetName]
		car.streets[i] = street
		street.carCount++

		if i == 0 {
			street.carsAtEnd = append(street.carsAtEnd, car)
		}
	}
}

func getOrCreateIntersection(simulation *Simulation, id int) *Intersection {
	intersection, hasIntersection := simulation.intersectionMap[id]

	if !hasIntersection {
		intersection = &Intersection{
			id,
			make([]*Street, 0),
			make([]*Street, 0),
			make([]TrafficLight, 0),
		}
		simulation.intersectionMap[id] = intersection
	}

	return intersection
}

func createCar(simulation *Simulation, streetCount int) *Car {
	id := len(simulation.carMap)

	return &Car{
		id, streetCount, make([]*Street, streetCount),
	}
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
