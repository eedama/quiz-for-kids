package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"path/filepath"
	"path"

	// "github.com/k0kubun/pp"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/karrick/godirwalk"
	"github.com/spf13/pflag"
)

var (
	help 			bool
	debug 			bool
	dataDir 		string
	mongoHost  		string
	mongoPort  		string
	mongoDB 		string 
	mongoCollection string
	isMongoDB 		 = true
	isExportJSON 	 = true
	isDropCollection = false
	// database 		 = "QuizzForKids"
	// collection 		 = "Questions"
	session          *mgo.Session
)

func main() {

	pflag.StringVarP(&dataDir, "data-dir", "", "/opt/quiz-for-kids/data", "data directory path (with the kahoot json dumps).")

	pflag.StringVarP(&mongoHost, "mongo-host", "", "localhost", "mongodb host.")
	pflag.StringVarP(&mongoPort, "mongo-port", "", "27017", "mongodb port.")
	pflag.StringVarP(&mongoDB, "mongo-database", "", "QuizzForKids", "mongodb database name.")
	pflag.StringVarP(&mongoCollection, "mongo-collection", "", "Questions", "mongodb collection name.")

	pflag.BoolVarP(&debug, "debug", "d", false, "debug mode")
	pflag.BoolVarP(&help, "help", "h", false, "help info")
	pflag.Parse()
	if help {
		pflag.PrintDefaults()
		os.Exit(1)
	}

	optVerbose := flag.Bool("verbose", false, "Print file system entries.")
	flag.Parse()

	i := 1
	err := godirwalk.Walk("/opt/quiz-for-kids/data", &godirwalk.Options{
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

func filenameWithoutExtension(fn string) string {
      return strings.TrimSuffix(fn, path.Ext(fn))
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
		mongoAddr := fmt.Sprintf("mongodb://%s:%s/", mongoHost, mongoPort)
		session, err = mgo.Dial(mongoAddr + mongoDB)
		if err != nil {
			log.Fatalf("cannot connect to mongodb host: %v\n", err)
		}
		if isDropCollection {
			col := session.DB(mongoDB).C(mongoCollection)
			col.DropCollection()
		}

	}

	var otdb Opentdb
	otdb.Name = filenameWithoutExtension(filepath.Base(filePath))
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

		// check if exists
		var otdbExist *Opentdb
		c := session.DB(mongoDB).C(mongoCollection)
		err := c.Find(bson.M{"name": otdb.Name}).One(&otdbExist)
		if err != nil {
			coll := session.DB(mongoDB).C(mongoCollection)
			err = coll.Insert(otdb)
			if err != nil {
				log.Fatalln(err)
			}
			// Close session as normal
			session.Close()
		}
	}

	// pp.Println(otdb)
    b, err := json.MarshalIndent(otdb, "", "\t")
    if err != nil {
        log.Println("error:", err)
    }
    ioutil.WriteFile(filePath+".otdb", b, os.ModePerm)


}
