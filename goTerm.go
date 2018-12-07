package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kniren/gota/dataframe"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// locate

var execCommands []string

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\t Go Shell")
	fmt.Println("---------------------")
	shellReader()

here:
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		// extracting the function called
		mainFn := strings.Split(text, " ")

		switch mainFn[0] {

		case "hi", "hey", "hello":
			greeter()
		case "pwd":
			pwd()
		case "ls":
			readCurrentDir(mainFn)
		case "cat":
			catty(mainFn)
		case "cd":
			changeDir(mainFn)
			pwd()
		case "sitwell":
			sitwell(mainFn)
		case "rm","rmdir":
			blackbuck(mainFn)
		case "touch":
			touch(mainFn)
		case "mkdir":
			mkdir(mainFn)
		case "clear":
			fmt.Println("\033[H\033[2J")
		case "mv":
			moveAndRename(mainFn)
		case "analyze":
			dataReader(mainFn)
		case "cp":
			Copy(mainFn)
		case "exit", "q":
			break here
		case "rapid":
			runGo(mainFn)
		default:
			quantifier(mainFn)

		}
	}
}

// prints current directory
func pwd() {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
}

// greets you
func greeter() {
	userDet, _ := user.Current()

	fields := *userDet
	print("Hey, ")
	fmt.Println(fields.Name)

}

// ls command
func readCurrentDir(dir []string) {
	var set string
	if len(dir) == 1 {
		set = "."
	} else {
		set = dir[1]
	}
	file, err := os.Open(set)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		fmt.Println(name)
	}
}

// cat elaboration
func catty(filename []string) {
	stream, err := ioutil.ReadFile(filename[1])

	checkErr(err)

	// Convert into a string
	readString := string(stream)

	fmt.Println(readString)

}

func changeDir(direc []string) {
	home := os.Getenv("HOME")
	set := home
	if len(direc) != 1 {
		set = direc[1]
	}
	err := os.Chdir(set)
	checkErr(err)
}

func sitwell(parsing []string) {
	if parsing[1] == "help" {
		fmt.Println("sitwell alias=yourAlias script-executable")
		return
	}

	alias := strings.Split(parsing[1], "=")[1]
	execName := parsing[2]
	fmt.Println(alias)
	fmt.Println(execName)
	execCommands = append(execCommands, alias)

	usr, err := user.Current()
	//mkdir([]string{"mkdir", usr.HomeDir+"/go-shell"})
	f, err := os.OpenFile(usr.HomeDir+"/go-shell/record.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	//mdir, err := os.Getwd()
	moveAndRename([]string{"mv",execName,usr.HomeDir+"/go-shell/"+alias})

	defer f.Close()

	if _, err = f.WriteString(alias); err != nil {
		panic(err)
	}
	if _, err = f.WriteString("\n"); err != nil {
		panic(err)
	}
	//if _, err = f.WriteString(execName); err != nil {
	//	panic(err)
	//}
	//
	//if _, err = f.WriteString("\n"); err != nil {
	//	panic(err)
	//}
}

func shellReader() {
	usr, err := user.Current()

	stream, err := ioutil.ReadFile(usr.HomeDir+"/go-shell/record.txt")

	checkErr(err)
	if err!=nil {
		crapper := []string {"touch",usr.HomeDir+"/go-shell/record.txt"}
		touch(crapper)
		return
	}
	// Convert into a string
	readString := string(stream)
	record := strings.Split(readString, "\n")
	fmt.Println(record)
	//var keys []string
	for k := range record {
		execCommands = append(execCommands,record[k])
		//if k%2==0 {
		//	keys = append(keys, record[k])
		//}else{
		//	execCommands[keys[k-1]] = record[k]
		//}
	}

}

func quantifier(cmd []string) {
		order := cmd[0]
		//var keys []string
		//for k := range execCommands {
		//	keys = append(keys, k)
		//}
		usr, err := user.Current()
		checkErr(err)
		if contains(execCommands, order) {
			curdir, err := os.Getwd()
			checkErr(err)
			changeDir([]string{"cd",usr.HomeDir+"/go-shell/"})
			//cmd := exec.Command("./"+order)
			//out, err1 := cmd.Output()
			//checkErr(err1)
			//cmd := exec.Command("python", "-c", `import subprocess as sb; sb.call(["` + usr.HomeDir + "/go-shell/" + order +`"])`)
			//out,err := cmd.CombinedOutput()
			//checkErr(err)
			//fmt.Println(string(out))


			command:= exec.Command("./"+order)

			// set var to get the output
			var out bytes.Buffer

			// set the output to our variable
			command.Stdout = &out
			err = command.Run()
			if err != nil {
				log.Println(err)
			}

			fmt.Println(out.String())


			changeDir([]string{"cd",curdir})
		}else{
			fmt.Println("Command not found!")
		}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
	if a == str {
		return true
		}
	}
	return false
}

// rm command
func blackbuck(instance []string){
	err := os.RemoveAll(instance[1])
	checkErr(err)
}

// touch
func touch(file []string){
	_, err := os.Create(file[1])
	checkErr(err)
	//fmt.Println(*result)
	fmt.Println("file created!")
}

func mkdir(msk []string){
	os.MkdirAll(msk[1], os.ModePerm)
}

func moveAndRename(dirs []string){
	err := os.Rename(dirs[1],dirs[2])
	checkErr(err)
}

func dataReader(fileInfo []string){
	stream, err := ioutil.ReadFile(fileInfo[1])

	checkErr(err)

	// Convert into a string
	dataStr := string(stream)

	if fileInfo[2] == "csv"{
		df := dataframe.ReadCSV(strings.NewReader(dataStr))
		fmt.Println(df)
		fmt.Println(df.Describe())
		fmt.Println("Number of rows: ",df.Nrow())
		fmt.Println("Number of rows: ",df.Ncol())
	}
	if fileInfo[2] == "json"{
		df := dataframe.ReadJSON(strings.NewReader(dataStr))
		fmt.Println(df)
		fmt.Println(df.Describe())
		fmt.Println("Number of rows: ",df.Nrow())
		fmt.Println("Number of rows: ",df.Ncol())
	}
	// display nrow and ncol
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func Copy(srcToDst []string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(srcToDst[1])
	checkErr(err)
	// Write data to dst
	err = ioutil.WriteFile(srcToDst[2], data, os.ModePerm)
	checkErr(err)
}


func runGo(rapid []string){
	output, err := exec.Command("go",  "run", rapid[1]).CombinedOutput()
	if err == nil {
		fmt.Println(string(output)) // write the output with ResponseWriter
	}
}