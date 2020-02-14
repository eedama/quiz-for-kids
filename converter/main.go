package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/globalsign/mgo"
	// "github.com/k0kubun/pp"
	"github.com/karrick/godirwalk"
)

var (
	isMongoDB 		= false
	isExportJSON 	= true
	database 		= "QuizzForKids"
	collection 		= "Questions"
	session         *mgo.Session
)

func main() {

	optVerbose := flag.Bool("verbose", false, "Print file system entries.")
	flag.Parse()

	i := 1
	err := godirwalk.Walk("../dataset", &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if *optVerbose {
				fmt.Printf("%s %s\n", de.ModeType(), osPathname)
			}
			if de.ModeType() != os.ModeDir && !strings.HasSuffix(osPathname, "otdb") {
				convertQuizz(osPathname, i)
			}
			i++
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			if *optVerbose {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			}

			// For the purposes of this example, a simple SkipNode will suffice,
			// although in reality perhaps additional logic might be called for.
			return godirwalk.SkipNode
		},
		Unsorted: true, // set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func convertQuizz(filePath string, id int) {
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened file: ", filePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v\n", err)
	}

	if isMongoDB {
		session, err = mgo.Dial("mongodb://127.0.0.1:27017/" + database)
		if err != nil {
			log.Fatalf("cannot connect to mongodb host: %v\n", err)
		}
	}

	var otdb Opentdb
	for _, q := range data.Kahoot.Questions {

		otdbr := OpentdbResult{
			Category: "Environment",
			Difficulty: "easy",
			Question: q.Question,
		}

		if len(q.Choices) == 2 {
			otdbr.Type = "boolean"
		} else {
			otdbr.Type = "multiple"				
		}

		for _, c := range q.Choices {

			otdbr.IncorrectAnswers = append(otdbr.IncorrectAnswers, c.Answer)
			// question.Answers = append(question.Answers, c.Answer)
			if c.Correct {
				otdbr.CorrectAnswer = c.Answer
			}
		}
		otdbr.Image = q.Image

		// question.Time = q.Time
		// question.PointsMultiplier = q.PointsMultiplier
		otdb.Results = append(otdb.Results, otdbr)
	}

	if isMongoDB {
		coll := session.DB(database).C(collection)
		err = coll.Insert(otdb)
		if err != nil {
			log.Fatalln(err)
		}
		// Close session as normal
		session.Close()
	}

	// pp.Println(otdb)
    b, err := json.MarshalIndent(otdb, "", "\t")
    if err != nil {
        fmt.Println("error:", err)
    }
    ioutil.WriteFile(filePath+".otdb", b, os.ModePerm)


}
