package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"

	"gitlab.com/schoentoon/rs-tools/lib/runemetrics"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Need 2 arguments, execute as %s [input] [output]", os.Args[0])
	}
	in := os.Args[1]
	out := os.Args[2]

	existing := readOutputFile(out)

	f, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	profile, err := runemetrics.ParseProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	newer := profile.Activities
	if len(existing) > 0 {
		newer = runemetrics.NewAchievementsSince(existing, profile.Activities)
	}

	fout, err := os.OpenFile(out, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	encoder := json.NewEncoder(fout)

	sort.Slice(newer, func(i, j int) bool { return newer[i].Date.Unix() < newer[j].Date.Unix() })

	for _, activity := range newer {
		err = encoder.Encode(activity)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readOutputFile(filename string) []runemetrics.Activity {
	out := []runemetrics.Activity{}

	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return out
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	for decoder.More() {
		activity := &runemetrics.Activity{}
		err := decoder.Decode(activity)
		if err != nil {
			panic(err)
			return out
		}
		out = append(out, *activity)
	}

	return out
}
