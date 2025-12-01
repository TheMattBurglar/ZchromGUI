package main

import (
	"ZtestAssisted/ztestlogic"
	"encoding/json"
	"strconv"
	"syscall/js"
)

type SimulationRequest struct {
	Adam        int     `json:"adam"`
	Eve         int     `json:"eve"`
	Lilith      int     `json:"lilith"`
	Diana       int     `json:"diana"`
	EveBirth    float64 `json:"eveBirth"`
	LilithBirth float64 `json:"lilithBirth"`
	DianaBirth  float64 `json:"dianaBirth"`
	MaxPop      int     `json:"maxPop"`
	Generations int     `json:"generations"`
	Timelines   int     `json:"timelines"`
	ViableY     bool    `json:"viableY"`
}

type SimulationResponse struct {
	SuccessCount    int    `json:"successCount"`
	TotalTimelines  int    `json:"totalTimelines"`
	TotalExtinction int    `json:"totalExtinction"`
	ZExtinction     int    `json:"zExtinction"`
	MaleExtinction  int    `json:"maleExtinction"`
	FemExtinction   int    `json:"femExtinction"`
	LastGen         int    `json:"lastGen"`
	MaxPopReached   int    `json:"maxPopReached"`
	PopCapGen       int    `json:"popCapGen"`
	Summary         string `json:"summary"`
}

// Global state for the current simulation
var (
	currentReq         SimulationRequest
	currentStats       *ztestlogic.SimulationStats
	timelinesCompleted int
	successCount       int
)

func initSimulation(this js.Value, args []js.Value) interface{} {
	jsonInput := args[0].String()
	err := json.Unmarshal([]byte(jsonInput), &currentReq)
	if err != nil {
		return "Error parsing JSON: " + err.Error()
	}

	// Reset state
	currentStats = &ztestlogic.SimulationStats{}
	timelinesCompleted = 0
	successCount = 0

	return nil
}

func runBatch(this js.Value, args []js.Value) interface{} {
	batchSize := args[0].Int()

	pop := [4]int{currentReq.Adam, currentReq.Eve, currentReq.Lilith, currentReq.Diana}
	birth := [3]float64{currentReq.EveBirth, currentReq.LilithBirth, currentReq.DianaBirth}
	viableY := "N"
	if currentReq.ViableY {
		viableY = "Y"
	}

	for i := 0; i < batchSize; i++ {
		if timelinesCompleted >= currentReq.Timelines {
			break
		}
		if ztestlogic.GenTryFail(pop, birth, viableY, currentReq.MaxPop, currentReq.Generations, currentStats) {
			successCount++
		}
		timelinesCompleted++
	}

	return map[string]interface{}{
		"finished":  timelinesCompleted >= currentReq.Timelines,
		"completed": timelinesCompleted,
		"total":     currentReq.Timelines,
	}
}

func getResults(this js.Value, args []js.Value) interface{} {
	// Build summary string
	summary := ""
	summary += strconv.Itoa(successCount) + " out of " + strconv.Itoa(currentReq.Timelines) + " timelines still had the Z chromosome by the end.\n"
	if currentStats.TotalExtinction > 0 {
		summary += strconv.Itoa(currentStats.TotalExtinction) + " failed because EVERYONE died out.\n"
	}
	if currentStats.Zextinction > 0 {
		summary += strconv.Itoa(currentStats.Zextinction) + " failed because Lilith and Diana died out. There were still Women, but no more Z chromosomes.\n"
	}
	if currentStats.MaleExtinction > 0 {
		summary += strconv.Itoa(currentStats.MaleExtinction) + " failed because men died out. Usually because total population got too small.\n"
	}
	if currentStats.FemExtinction > 0 {
		summary += strconv.Itoa(currentStats.FemExtinction) + " failed because women died out completely. Usually because total population got too small.\n"
	}
	if currentStats.LastGen > 0 {
		summary += "If they ended without either men or a Z chromosome, they did so within " + strconv.Itoa(currentStats.LastGen) + " generations.\n"
	}
	if currentStats.MaxPopReached > 0 {
		summary += strconv.Itoa(currentStats.MaxPopReached) + " were cut off early because they reached a population size of " + strconv.Itoa(currentReq.MaxPop) + "\n"
		summary += "They hit that population cap at or below " + strconv.Itoa(currentStats.PopCapGen) + " generations.\n"
	}

	// Run one last time with population details for the example
	pop := [4]int{currentReq.Adam, currentReq.Eve, currentReq.Lilith, currentReq.Diana}
	birth := [3]float64{currentReq.EveBirth, currentReq.LilithBirth, currentReq.DianaBirth}
	viableY := "N"
	if currentReq.ViableY {
		viableY = "Y"
	}

	ZisThere, finalPop := ztestlogic.GenTryFailWithPop(pop, birth, viableY, currentReq.MaxPop, currentReq.Generations, currentStats)

	isMarker := finalPop[0] == 0 && finalPop[1] == 0 && finalPop[2] == 0 && finalPop[3] > 0
	total := finalPop[0] + finalPop[1] + finalPop[2] + finalPop[3]

	if !isMarker && total > 0 {
		summary += "\nExample timeline final population totals:\n"
		summary += "Adam: " + strconv.Itoa(finalPop[0]) + "\n"
		summary += "Eve: " + strconv.Itoa(finalPop[1]) + "\n"
		summary += "Lilith: " + strconv.Itoa(finalPop[2]) + "\n"
		summary += "Diana: " + strconv.Itoa(finalPop[3]) + "\n"
		summary += "\nPercentages:\n"
		summary += "Adam: " + strconv.FormatFloat(100*float64(finalPop[0])/float64(total), 'f', 2, 64) + "%\n"
		summary += "Eve: " + strconv.FormatFloat(100*float64(finalPop[1])/float64(total), 'f', 2, 64) + "%\n"
		summary += "Lilith: " + strconv.FormatFloat(100*float64(finalPop[2])/float64(total), 'f', 2, 64) + "%\n"
		summary += "Diana: " + strconv.FormatFloat(100*float64(finalPop[3])/float64(total), 'f', 2, 64) + "%\n"
	}
	if ZisThere && isMarker {
		summary += "\nExample timeline ended successfully by reaching the population cap.\n"
	}

	if isMarker && !ZisThere {
		summary += "\nExample timeline ended with a marker indicating Men, Women, or the Z chromosome died out.\n"
	}

	resp := SimulationResponse{
		SuccessCount:    successCount,
		TotalTimelines:  currentReq.Timelines,
		TotalExtinction: currentStats.TotalExtinction,
		ZExtinction:     currentStats.Zextinction,
		MaleExtinction:  currentStats.MaleExtinction,
		FemExtinction:   currentStats.FemExtinction,
		LastGen:         currentStats.LastGen,
		MaxPopReached:   currentStats.MaxPopReached,
		PopCapGen:       currentStats.PopCapGen,
		Summary:         summary,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return "Error marshalling response: " + err.Error()
	}

	return string(jsonResp)
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("initSimulation", js.FuncOf(initSimulation))
	js.Global().Set("runBatch", js.FuncOf(runBatch))
	js.Global().Set("getResults", js.FuncOf(getResults))
	<-c
}
