package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"spy-cat-agency/internal/store"
)

// go shows me store.Target is not a type compiler NotAType so i did this
type storeTarget = store.Target

func Seed(store store.Storage) {
	ctx := context.Background()

	cats := generateCats(100)
	for idx, cat := range cats {
		if err := store.Cat.Create(ctx, &cat); err != nil {
			log.Println("Error in seed creating cat: ", err.Error())
			return
		}
		cats[idx].ID = cat.ID
	}

	missions := generateMissions(55)
	for idx, mission := range missions {
		if err := store.Mission.Create(ctx, &mission); err != nil {
			log.Println("Error in seed creating mission: ", err.Error())
			return
		}
		missions[idx].ID = mission.ID
		missions[idx].Targets = mission.Targets
	}

	var targets []storeTarget = nil
	for _, mission := range missions {
		targets = append(targets, mission.Targets...)
	}

	notes := generateNotes(44, targets)
	for _, note := range notes {
		if err := store.Mission.AddNote(ctx, &note); err != nil {
			log.Println("Error in seed creating notes: ", err.Error())
			return
		}
	}

	missionsLen := len(missions)
	for idx, cat := range cats {
		missionIdx := idx % missionsLen
		mission := missions[missionIdx]

		if err := store.Mission.AssignCat(ctx, cat.ID, mission.ID); err != nil {
			log.Println("Error in seed assigning mission: ", err.Error())
			return
		}
	}
}

func generateCats(num int) []store.Cat {
	cats := make([]store.Cat, num)
	catNamesLen := len(catNames)
	catBreedLen := len(breedNames)
	for idx := range cats {
		catName := catNames[idx%catNamesLen] + fmt.Sprintf("_%d", idx)
		catBreed := breedNames[idx%catBreedLen]
		cats[idx] = store.Cat{
			Name:              catName,
			YearsOfExperience: rand.IntN(15),
			Breed:             catBreed,
			Salary:            rand.Float64() * 10000.0,
		}
	}

	return cats
}

func generateMissions(num int) []store.Mission {
	missions := make([]store.Mission, num)

	for idx := range missions {
		targs := generateTargets(1 + rand.IntN(2))
		missions[idx] = store.Mission{
			Targets: targs,
		}
	}

	return missions
}

func generateTargets(num int) []store.Target {
	targets := make([]store.Target, num)
	targetNamesLen := len(targetNames)
	countriesLen := len(countries)

	for idx := range targets {
		targetName := targetNames[rand.IntN(targetNamesLen)] + fmt.Sprintf("_%d", idx)
		country := countries[rand.IntN(countriesLen)]
		targets[idx] = store.Target{
			Name:    targetName,
			Country: country,
		}
	}

	return targets
}

func generateNotes(num int, targets []store.Target) []store.Note {
	notes := make([]store.Note, num)
	targetsLen := len(targets)
	peopleNotesLen := len(peopleNotes)

	for idx := range notes {
		target := targets[rand.IntN(targetsLen)]
		note := peopleNotes[idx%peopleNotesLen]
		notes[idx] = store.Note{
			TargetID:  target.ID,
			MissionID: target.MissionID,
			Note:      note,
		}
	}

	return notes
}
