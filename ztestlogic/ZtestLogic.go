package ztestlogic

import (
	"fmt"
	"math/rand"
)

// global variables; will NOT change
var Adam [2]string = [2]string{"X", "Y"}
var Eve [2]string = [2]string{"X", "X"}
var Lilith [2]string = [2]string{"Z", "Y"}
var Diana [2]string = [2]string{"Z", "X"}

// global variables used to track #s across multiple timelines (aka multiple runs of GenTryFail)
var MaleExtinction int = 0  //counts how many times males died out completely across timelines
var FemExtinction int = 0   //counts how many times females died out completely across timelines
var Zextinction int = 0     //counts how many times both Lilith and Diana (Z chromosom carriers) died out across timelines
var TotalExtinction int = 0 //counts how many times EVERYONE died out across timelines
var MaxPopReached int = 0   //counts how many times the population cap was reached across timelines
var LastGen int = 0         //if Z or men died out, this is the greatest # of generations it took for that to happen across timelines
var PopCapGen int = 0       //if the population cap was reached, this is the greatest # of generations it took for that to happem across timelines

// setup a random population
func RandomPop() [4]int {
	totalRandomPopSize := 0
	fmt.Println("What small population size do you want to start with (ex: 200)? ")
	fmt.Scanln(&totalRandomPopSize)

	A := 0
	E := 0
	L := 0
	D := 0

	for i := 0; i < totalRandomPopSize; i++ {
		var child [2]string = [2]string{Adam[rand.Intn(2)], randomWoman()[rand.Intn(2)]}

		if child[0] == "X" && child[1] == "X" {
			E++
		} else if child[0] == "Z" && child[1] == "Y" {
			L++
		} else if child[0] == "Y" && child[1] == "Z" {
			L++
		} else if child[0] == "Z" && child[1] == "X" {
			D++
		} else if child[0] == "X" && child[1] == "Z" {
			D++
		} else if child[0] == "Y" && child[1] == "Y" {
			//YY not viable, try again
			i--
		} else {
			A++
		}
	}
	var array [4]int = [4]int{A, E, L, D}
	return array
}

// pick a random Eve, Lilith, or Diana
func randomWoman() [2]string {
	rand := rand.Intn(3)
	if rand == 0 {
		return Eve
	} else if rand == 1 {
		return Lilith
	} else {
		return Diana
	}
}

// uses the seed population, birthrate, viability of Y eggs, and the population cap to generate the next population
func nextGen(seedPop [4]int, birthRateELD [3]float64, viableY string, maxPopulation int) [4]int {
	var newFem int
	var nMale int
	var nEve int
	var nLilith int
	var nDiana int

	//Next gen born from Eve
	for i := 0.0; i < (float64(seedPop[1]) * birthRateELD[0]); i++ {
		var kid [2]string = [2]string{Adam[rand.Intn(2)], Eve[rand.Intn(2)]}
		if kid[0] == "Y" {
			nMale++
		} else {
			newFem++
			nEve++
		}
	}

	//Next gen born from Lilith
	for i := 0.0; i < (float64(seedPop[2]) * birthRateELD[1]); i++ {
		kid := [2]string{Lilith[rand.Intn(2)], Adam[rand.Intn(2)]}
		if kid[0] == "Y" && kid[1] == "Y" {
			i-- //YY not viable, try again
		} else if kid[0] == "Y" && kid[1] == "X" {
			if viableY == "Y" || viableY == "y" {
				nMale++
			} else {
				i-- //Y egg not viable, by user input
			}
		} else if kid[0] == "Z" && kid[1] == "Y" {
			nLilith++
			newFem++
		} else {
			nDiana++
			newFem++
		}
	}

	//Next gen born from Diana
	for i := 0.0; i < (float64(seedPop[3]) * birthRateELD[2]); i++ {
		kid := [2]string{Diana[rand.Intn(2)], Adam[rand.Intn(2)]}
		if kid[0] == "Z" && kid[1] == "X" {
			newFem++
			nDiana++
		} else if kid[0] == "Z" && kid[1] == "Y" {
			newFem++
			nLilith++
		} else if kid[0] == "X" && kid[1] == "X" {
			newFem++
			nEve++
		} else {
			nMale++
		}
	}

	newPop := [4]int{nMale, nEve, nLilith, nDiana}
	total := nMale + newFem

	//series of return values that communicate what happened through unnatural output.
	if nMale == 0 && newFem == 0 {
		TotalExtinction++
		fmt.Println("EVERYONE DIED OUT!  Both Men AND Women are gone!")
	}
	if nMale == 0 {
		MaleExtinction++
		return [4]int{0, 0, 0, 1}
		//Male extiction is an existential issue for sexual reprodction.
		//Since a Diana and a male can sexualy reproduce all 4 types,
		//we need not worry if Eve or Lilith dies out for 1 generation. (only if Both Z women do)
	}
	if newFem == 0 {
		FemExtinction++
		return [4]int{0, 0, 0, 2} //unnatural output because if nMale == 0, the output would be [4]int{0,0,0,1}
	}
	if nLilith == 0 && nDiana == 0 {
		Zextinction++
		return [4]int{0, 0, 0, 3} //unnatural output because if nMale == 0, the output would be [4]int{0,0,0,1}
	}
	if total >= maxPopulation {
		MaxPopReached++
		return [4]int{0, 0, 0, 4}
		//unnstural output because if nMale == 0, the output would be [4]int{0,0,0,1}
		//This will make it possible to exit out of the tryFail loop early, even though we still have a viable generation.
	}

	fmt.Println(newPop)
	var malePercentage float64 = float64(nMale) / float64(total)
	var evePercentage float64 = float64(nEve) / float64(total)
	var lilithPercentage float64 = float64(nLilith) / float64(total)
	var dianaPercentage float64 = float64(nDiana) / float64(total)

	fmt.Printf("Male %.2f%%", (malePercentage * 100))
	fmt.Printf("\nEve %.2f%%", (evePercentage * 100))
	fmt.Printf("\nLilith %.2f%%", (lilithPercentage * 100))
	fmt.Printf("\nDiana %.2f%%", (dianaPercentage * 100))
	fmt.Println()

	return newPop
}

func GenTryFail(seedPop [4]int, birthRateELD [3]float64, viableY string, maxPopulation int, generations int) bool {
	fmt.Println("\nStarting Values for this Timeline:")
	fmt.Println(seedPop)
	series := nextGen(seedPop, birthRateELD, viableY, maxPopulation)

	for count := 1; count <= generations; count++ {

		if series == [4]int{0, 0, 0, 0} {
			fmt.Printf("Ended in %v generations because EVERYONE died out.\n", count)
			if LastGen <= count {
				LastGen = count
			}
			return false
		}
		if series == [4]int{0, 0, 0, 1} {
			fmt.Printf("Ended in %v generations because Men died out.\n", count)
			if LastGen <= count {
				LastGen = count
			}
			return false
		}
		if series == [4]int{0, 0, 0, 2} {
			fmt.Printf("Ended in %v generation because Women died out.\n", count)
			if LastGen <= count {
				LastGen = count
			}
			return false
		}
		if series == [4]int{0, 0, 0, 3} {
			fmt.Printf("Ended in %v generations because Z chromosome died out.\n", count)
			if LastGen <= count {
				LastGen = count
			}
			return false
		}
		if series == [4]int{0, 0, 0, 4} {
			fmt.Printf("Population size exceeded in %v generations!\n", count)
			if PopCapGen <= count {
				PopCapGen = count
			}
			return true
		}

		series = nextGen(series, birthRateELD, viableY, maxPopulation)
		if count == generations {
			return true
		}
	}
	fmt.Println("ERROR! genTryFail exited incorrectly!")
	return false //this should never be reached
}

func GenTryFailWithPop(seedPop [4]int, birthRateELD [3]float64, viableY string, maxPopulation int, generations int) (bool, [4]int) {
	series := nextGen(seedPop, birthRateELD, viableY, maxPopulation)

	for count := 1; count <= generations; count++ {
		if series == [4]int{0, 0, 0, 0} {
			return false, series
		}
		if series == [4]int{0, 0, 0, 1} {
			return false, series
		}
		if series == [4]int{0, 0, 0, 2} {
			return false, series
		}
		if series == [4]int{0, 0, 0, 3} {
			return false, series
		}
		if series == [4]int{0, 0, 0, 4} {
			return true, series
		}
		series = nextGen(series, birthRateELD, viableY, maxPopulation)
		if count == generations {
			return true, series
		}
	}
	return false, series // fallback
}
