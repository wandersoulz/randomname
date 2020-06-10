package randomname

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"github.com/wandersoulz/godes"
)

var funcDist *godes.FunctionalDistr
var names []string
var cd *conditionalDistribution
var contextSize int

type conditionalDistribution []*ngramItem

var probabilitiesLookup map[string][]ngramProb

// Init Function to get the initial conditions for starting the name generation properties
func Init(filename string, initContextSize int) {
	cd = getConditionalDistribution(filename, initContextSize)
	contextSize = initContextSize
}

func findInSlice(arr []ngramProb, val string) int {
	for j := range arr {
		if arr[j].nextChar == val {
			return j
		}
	}
	return -1
}

func (cd *conditionalDistribution) lookUpProbabilities(lookUp string) []ngramProb {
	var results []ngramProb
	numValuesFound := 0.0
	if val, ok := probabilitiesLookup[lookUp]; ok {
		return val
	}
	for i := range *cd {
		item := (*cd)[i]
		if lookUp == item.prev {
			numValuesFound += 1.0
			probIndex := findInSlice(results, item.next)
			if probIndex == -1 {
				value := ngramProb{
					item.next,
					1.0,
				}
				results = append(results, value)
			} else {
				results[probIndex].probability += 1.0
			}
		}
	}

	for j := range results {
		prob := results[j].probability

		results[j].probability = prob / numValuesFound
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].probability < results[j].probability
	})
	probabilitiesLookup[lookUp] = results
	return probabilitiesLookup[lookUp]
}

// GetName Main function to call to get a name
func GetName() string {
	if funcDist == nil {
		funcDist = godes.NewFunctionalDistr(true)
	}
	if probabilitiesLookup == nil {
		probabilitiesLookup = make(map[string][]ngramProb)
	}

	name := strings.Repeat(" ", contextSize)
	nextValue := cd.sampleDistribution(name)
	for nextValue != "#" {
		name = name + nextValue
		length := len(name)
		context := name[length-contextSize : length]
		nextValue = cd.sampleDistribution(context)
	}
	return strings.Title(strings.Trim(name, " "))
}

func sumProbs(arr []ngramProb) float64 {
	sum := 0.0
	for i := 0; i < len(arr); i++ {
		sum += arr[i].probability
	}
	return sum
}

func (cd *conditionalDistribution) sampleDistribution(lookUp string) string {
	probabilities := cd.lookUpProbabilities(lookUp)
	sampleIndex := funcDist.Get(getValues(probabilities), 0.0, 1.0)
	if sampleIndex == -1 {
		return "#"
	}
	return probabilities[sampleIndex].nextChar
}

func getValues(probs []ngramProb) []float64 {
	ret := make([]float64, len(probs))
	for i := range probs {
		ret[i] = probs[i].probability
	}
	return ret
}

type ngramProb struct {
	nextChar    string
	probability float64
}

type ngramItem struct {
	prev string
	next string
}

func getConditionalDistribution(filename string, contextSize int) *conditionalDistribution {
	getNames(filename)
	paramN := contextSize + 1
	pad := strings.Repeat(" ", contextSize)
	nm := make([]string, len(names))
	for i := 0; i < len(names); i++ {
		nm[i] = pad + strings.ToLower(names[i]) + "#"
	}

	allTokens := []string{}
	for i := 0; i < len(nm); i++ {
		tokens := splitWord(nm[i], paramN)
		allTokens = append(allTokens, tokens...)
	}

	conditionalDist := conditionalDistribution{}
	allCount := 0
	for i := 0; i < len(allTokens); i++ {
		prev := allTokens[i][0:contextSize]
		next := string(allTokens[i][paramN-1])
		item := &ngramItem{
			prev,
			next,
		}
		conditionalDist = append(conditionalDist, item)
		allCount++
	}

	return &conditionalDist
}

func splitWord(input string, size int) []string {
	tokens := []string{input[0:size]}
	for i := 1; i < len(input)-size; i++ {
		tokens = append(tokens, input[i:i+size])
	}
	return tokens
}

func getNames(filename string) {
	nameFile, _ := os.Open(filename)
	defer nameFile.Close()
	scanner := bufio.NewScanner(nameFile)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
}
