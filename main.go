package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(filename string)([]problem, error)  {
	//1. open the file
	if fObj,err := os.Open(filename); err == nil {
		//2. create a reader
		csvR := csv.NewReader(fObj)
		//3. read the file
		if clines, err := csvR.ReadAll(); err == nil {
			//4. call the problemparcer function
			return problemPracer(clines), nil
		}else{
			return nil, fmt.Errorf("failed to read in csv file" + "format from %s file: %s", filename , err.Error())
		}
	}else{
		return nil, fmt.Errorf("failed to open %s file: %s", filename, err.Error())
	}
}

func main(){
	fmt.Printf("Welcome to the Quiz.\nAnswer the following Questions:\n")

	//1. Input name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file")

	//2. set the duration of the timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()

	//3. pull the problem from the file using problemPuller
	problems, err := problemPuller(*fName)

	// 4. handle the error appropriately
		if err != nil {
			exit(fmt.Sprintf("Something went wrong%s", err.Error()))
		}
		
	// 5. create a variable to count the number of correct answers
	correctAns := 0

	// 6. using the duration of the timer, we would initalize the timer
		tObj := time.NewTimer(time.Duration(*timer) * time.Second)
		ansC := make(chan string)

	//7. loop through the problem and print the question, we'll accept the answer
	problemLoop :
		for i, p := range problems{
			var answer string
			fmt.Printf("Problem %d: %s=", i+1,p.q)

			go func ()  {
				fmt.Scanf("%s" , &answer)
				ansC <- answer
			}()
			select{
			case <- tObj.C:
				fmt.Println()
				break problemLoop
			case iAns := <- ansC:
				if iAns == p.a {
					correctAns++
				}	
			}
			if i == len(problems)-1 {
				close(ansC)
			}
			
		}	

	//8. we'll calculate the number of correct answers and print the result
	fmt.Printf("Your results is: %d out of %d\n",correctAns , len(problems))
	exit("Thanks for playing...")
	<- ansC
}

func problemPracer(lines [][]string)[]problem {
	//go through the lines and parse them with the problem struct 
	r := make([]problem, len(lines))
	for i:= 0; i < len(lines); i++{
		r[i] = problem{q: lines[i][0] , a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
    os.Exit(1)
}