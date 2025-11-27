package main

import (
	"ZtestAssisted/ztestlogic"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func main() {
	fs := http.FileServer(http.Dir("./docs"))
	http.Handle("/", fs)

	http.HandleFunc("/api/simulate", handleSimulation)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSimulation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pop := [4]int{req.Adam, req.Eve, req.Lilith, req.Diana}
	birth := [3]float64{req.EveBirth, req.LilithBirth, req.DianaBirth}
	viableY := "N"
	if req.ViableY {
		viableY = "Y"
	}

	stats := &ztestlogic.SimulationStats{}
	success := 0

	for i := 0; i < req.Timelines; i++ {
		if ztestlogic.GenTryFail(pop, birth, viableY, req.MaxPop, req.Generations, stats) {
			success++
		}
	}

	// Build summary string (similar to GUI)
	summary := ""
	summary += strconv.Itoa(success) + " out of " + strconv.Itoa(req.Timelines) + " timelines still had the Z chromosome by the end.\n"
	if stats.TotalExtinction > 0 {
		summary += strconv.Itoa(stats.TotalExtinction) + " failed because EVERYONE died out.\n"
	}
	if stats.Zextinction > 0 {
		summary += strconv.Itoa(stats.Zextinction) + " failed because Lilith and Diana died out. There were still Women, but no more Z chromosomes.\n"
	}
	if stats.MaleExtinction > 0 {
		summary += strconv.Itoa(stats.MaleExtinction) + " failed because men died out. Usually because total population got too small.\n"
	}
	if stats.FemExtinction > 0 {
		summary += strconv.Itoa(stats.FemExtinction) + " failed because women died out completely. Usually because total population got too small.\n"
	}
	if stats.LastGen > 0 {
		summary += "If they ended without either men or a Z chromosome, they did so within " + strconv.Itoa(stats.LastGen) + " generations.\n"
	}
	if stats.MaxPopReached > 0 {
		summary += strconv.Itoa(stats.MaxPopReached) + " were cut off early because they reached a population size of " + strconv.Itoa(req.MaxPop) + "\n"
		summary += "They hit that population cap at or below " + strconv.Itoa(stats.PopCapGen) + " generations.\n"
	}

	// Run one last time with population details for the example
	ZisThere, finalPop := ztestlogic.GenTryFailWithPop(pop, birth, viableY, req.MaxPop, req.Generations, stats)

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
		SuccessCount:    success,
		TotalTimelines:  req.Timelines,
		TotalExtinction: stats.TotalExtinction,
		ZExtinction:     stats.Zextinction,
		MaleExtinction:  stats.MaleExtinction,
		FemExtinction:   stats.FemExtinction,
		LastGen:         stats.LastGen,
		MaxPopReached:   stats.MaxPopReached,
		PopCapGen:       stats.PopCapGen,
		Summary:         summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
