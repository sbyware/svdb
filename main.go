package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Service represents a service entry in the database
type Service struct {
	Port        string `json:"port"`
	Description string `json:"description"`
	Tcp         bool   `json:"tcp"`
	Udp         bool   `json:"udp"`
	Status      string `json:"status"`
}

var (
	dbPath     = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".svdb")
	logPrefix  = "[svdb]"
	validFlags = []string{"p", "json", "m"}
)

// LoadDB loads the service database from JSON file
func LoadDB(path string) (map[string][]Service, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s database file does not exist", logPrefix)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s error opening database file: %s", logPrefix, err)
	}
	defer file.Close()

	var db map[string][]Service
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&db)
	if err != nil {
		return nil, fmt.Errorf("%s error decoding database file: %s", logPrefix, err)
	}

	return db, nil
}

func main() {
	db, err := LoadDB(dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	port := flag.String("p", "", "query by port number(s) (comma separated list of ports)")
	jsonOutput := flag.Bool("json", false, "output in json format")
	match := flag.String("m", "", "query service database by regular expression pattern matching")
	flag.Parse()

	if !validateFlags(flag.CommandLine) {
		flag.PrintDefaults()
		return
	}

	switch {
	case *port != "":
		queryByPort(db, *port, *jsonOutput)
	case *match != "":
		queryByMatch(db, *match, *jsonOutput)
	default:
		fmt.Println("svdb, the service database.")
		flag.PrintDefaults()
	}
}

// validateFlags checks if the provided flags are valid
func validateFlags(f *flag.FlagSet) bool {
	valid := true
	f.Visit(func(flag *flag.Flag) {
		if !contains(validFlags, flag.Name) {
			fmt.Printf("%s error: invalid flag %s\n", logPrefix, flag.Name)
			valid = false
		}
	})
	return valid
}

// contains checks if a string slice contains a specific string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// queryByPort queries the database by port number(s)
func queryByPort(db map[string][]Service, port string, jsonOutput bool) {
	ports := strings.Split(port, ",")
	var services []Service

	for _, p := range ports {
		foundServices, exists := db[p]
		if exists {
			services = append(services, foundServices...)
		} else {
			fmt.Printf("service with port '%s' not found in db\n\n", p)
		}
	}

	if jsonOutput {
		printJSON(services)
	} else {
		printServices(services)
	}
}

// printJSON prints services in JSON format
func printJSON(services []Service) {
	jsonData, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		fmt.Printf("%s error marshalling JSON: %s\n", logPrefix, err)
		return
	}
	fmt.Println(string(jsonData))
}

// printServices prints services in human-readable format
func printServices(services []Service) {
	for _, s := range services {
		proto := ""
		if s.Tcp && s.Udp {
			proto = "tcp/udp"
		} else if s.Tcp {
			proto = "tcp"
		} else if s.Udp {
			proto = "udp"
		}
		fmt.Printf("port\t\t: %s\ndescription\t: %s\nprotocol\t: %s\nstatus\t\t: %s\n\n", s.Port, s.Description, proto, s.Status)
	}
}

// queryByMatch queries the database by regular expression pattern matching
func queryByMatch(db map[string][]Service, pattern string, jsonOutput bool) {
	var services []Service

	r, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("%s error compiling regular expression: %s\n", logPrefix, err)
		return
	}

	for _, v := range db {
		for _, s := range v {
			descProto := fmt.Sprintf("%s %s", strings.ToLower(s.Description), getProtocolString(s.Tcp, s.Udp))
			if r.MatchString(descProto) {
				services = append(services, s)
			}
		}
	}

	if jsonOutput {
		printJSON(services)
	} else {
		printServices(services)
	}
}

// getProtocolString returns the protocol string based on TCP and UDP flags
func getProtocolString(tcp, udp bool) string {
	if tcp && udp {
		return "tcp/udp"
	} else if tcp {
		return "tcp"
	} else if udp {
		return "udp"
	}
	return ""
}
