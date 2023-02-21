package main

import (
	"fmt"
    "os"
    "bufio"
    "regexp"
    "time"
    "bytes"
    "strings"
    "io/ioutil"
    "math/rand"
    "path/filepath"

    parser "github.com/GJSBRT/FiveMCompiler/luasyntax/parser"
    format "github.com/GJSBRT/FiveMCompiler/luasyntax/format"
)

var (
    eventTypes = []string{
        "AddEventHandler", 
        "RegisterNetEvent",
        "RegisterServerEvent", 
        "TriggerEvent",
        "TriggerClientEvent",
        "TriggerServerEvent",
    }

    events = make(map[string]string)
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func main() {
	fmt.Println("Running")
	
    /* Read all events */
    err := filepath.Walk("./test", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            readFile(path)
        }

        return nil
    })

    if err != nil {
        fmt.Println(err)
    }
    
    for k, v := range events {
        fmt.Println(k, " ", v)
    }

    /* Rename all events */
    fmt.Println("Renaming events")
    err = filepath.Walk("./test", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            matched, err := filepath.Match("*.lua", info.Name())
            if err != nil {
                return err
            }

            if info.Name() == "__resource.lua" || info.Name() == "fxmanifest.lua" {
                matched = false
            }

            if matched {
                read, err := ioutil.ReadFile(path)
                if err != nil {
                    return err
                }

                for oldEvent, newEvent := range events {
                    fmt.Println("Renaming " + oldEvent + " to " + newEvent + " in " + path)

                    read = bytes.Replace(read, []byte(oldEvent), []byte(newEvent), -1)
                }
                
                err = ioutil.WriteFile(path, read, 0666);  
                if err != nil {
                    return err
                }

                lua, err := parser.ParseFile(path, nil)
                if err != nil {
                    return err
                }

                format.Minify(lua)

                var output bytes.Buffer
                lua.WriteTo(&output)

                err = ioutil.WriteFile(path, []byte(output.String()), 0666);  
                if err != nil {
                    return err
                }


                return nil
            }
        }

        return nil
    })

	if err != nil {
		panic(err)
	}
}

func readFile(path string) {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    scanner := bufio.NewScanner(f)

    line := 1
    for scanner.Scan() {
        for _, eventType := range eventTypes {
            re := regexp.MustCompile("(?i)"+ eventType +"\\((\\n[\"'](.*?)[\"']|\\n\\s+[\"'](.*?)[\"']|.+[\"'](.*?)[\"']|[\"'](.*?)[\"'])")
            parameters := strings.Split(strings.Replace(string(re.Find(scanner.Bytes())), eventType + "(", "", -1), ",")
            event := strings.Replace(strings.Replace(parameters[0], "\"", "'", -1), "'", "", -1)

            if (re.MatchString(scanner.Text())) {
                events[event] = "druif_" + fmt.Sprint(rand.Int())
            }
        }

        line++
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
}