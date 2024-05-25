package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	DB_PATH    = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".svdb")
	LOG_PREFIX = "[svdb]"
)

type Service struct {
	Port        string `json:"port"`
	Description string `json:"description"`
	Tcp         bool   `json:"tcp"`
	Udp         bool   `json:"udp"`
	Status      string `json:"status"`
}

func LoadDB() (map[string][]Service, error) {
	// if DB_PATH aint there, copy ./db.json to DB_PATH and return parsed db
	if _, err := os.Stat(DB_PATH); os.IsNotExist(err) {
		if err := CopyDB(); err != nil {
			return nil, fmt.Errorf("%s error copying db: %s", LOG_PREFIX, err)
		}
	}

	db, err := ParseDB()
	if err != nil {
		return nil, fmt.Errorf("%s error parsing db: %s", LOG_PREFIX, err)
	}

	return db, nil
}

func CopyDB() error {
	src, err := os.Open("./db.json")
	if err != nil {
		return fmt.Errorf("%s error opening db file: %s", LOG_PREFIX, err)
	}
	defer src.Close()

	dst, err := os.Create(DB_PATH)
	if err != nil {
		return fmt.Errorf("%s error creating db file: %s", LOG_PREFIX, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("%s error copying db file: %s", LOG_PREFIX, err)
	}

	return nil
}

func ParseDB() (map[string][]Service, error) {
	var db map[string][]Service

	file, err := os.Open(DB_PATH)

	if err != nil {
		return nil, fmt.Errorf("%s error opening db file: %s", LOG_PREFIX, err)
	}
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&db); err != nil {
		return nil, fmt.Errorf("%s error decoding db file: %s", LOG_PREFIX, err)
	}

	defer file.Close()

	return db, nil
}

func main() {
	portFlag := flag.String("p", "", "query by (-p)ort number(s) (comma separated list of ports)")
	jsonOutputFlag := flag.Bool("j", false, "output in (-j)son format")
	matchFlag := flag.String("X", "", "query service database by regular e(-X)pression pattern matching")
	filterFlag := flag.String("select", "", "reduce output to (-select=fields,like,so)")
	flag.Parse()

	db, err := LoadDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	var matchedServices []Service

	switch {
	case *portFlag != "":
		services, err := queryByPort(db, *portFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		matchedServices = services
	case *matchFlag != "":
		services, err := queryByMatch(db, *matchFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		matchedServices = services
	default:
		fmt.Println("svdb, the service database.")
		flag.PrintDefaults()
		return
	}

	if *filterFlag != "" {
		matchedServices = selectFields(matchedServices, *filterFlag)
		printFiltered(matchedServices)
		return
	}

	if *jsonOutputFlag {
		printJSON(matchedServices)
	} else {
		printPlain(matchedServices)
	}
}

// selectFields filters services based on specified fields
func selectFields(services []Service, fields string) []Service {
	if fields == "" {
		return services
	}

	fieldsArr := strings.Split(fields, ",")

	filteredServices := make([]Service, len(services))
	for i, s := range services {
		filteredServices[i] = Service{}
		for _, f := range fieldsArr {
			switch f {
			case "port":
				filteredServices[i].Port = s.Port
			case "description":
				filteredServices[i].Description = s.Description
			case "tcp":
				filteredServices[i].Tcp = s.Tcp
			case "udp":
				filteredServices[i].Udp = s.Udp
			case "status":
				filteredServices[i].Status = s.Status
			}
		}
	}

	return filteredServices
}

// queryByPort queries the database by port number
func queryByPort(db map[string][]Service, port string) ([]Service, error) {
	var services []Service
	for _, p := range strings.Split(port, ",") {
		foundServices, exists := db[p]
		if !exists {
			fmt.Printf("%s service with port '%s' not found in db\n", LOG_PREFIX, p)
			continue
		}
		services = append(services, foundServices...)
	}
	return services, nil
}

// queryByMatch queries the database by regular expression pattern matching
func queryByMatch(db map[string][]Service, pattern string) ([]Service, error) {
	var services []Service

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%s error compiling regular expression: %s", LOG_PREFIX, err)
	}

	for _, v := range db {
		for _, s := range v {
			descProto := fmt.Sprintf("%s %s", strings.ToLower(s.Description), getProtocolString(s.Tcp, s.Udp))
			if r.MatchString(descProto) {
				services = append(services, s)
			}
		}
	}

	return services, nil
}

// printJSON prints services in JSON format
func printJSON(services []Service) {
	jsonData, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		fmt.Printf("%s error marshalling JSON: %s\n", LOG_PREFIX, err)
		return
	}
	fmt.Println(string(jsonData))
}

// printPlain prints services in human-readable format
func printPlain(services []Service) {
	for _, s := range services {
		proto := getProtocolString(s.Tcp, s.Udp)
		printKeyVal("port", s.Port)
		printKeyVal("description", s.Description)
		printKeyVal("protocol", proto)
		printKeyVal("status", s.Status)
		fmt.Println()
	}
}

func printFiltered(services []Service) {
	// if theres more than 1 service and there's more than 1 key in each service,
	// we call printPlain.
	// else, means we only have 1 key in a single service, so we print just the value.
	if len(services) > 1 && len(serviceToArr(services[0])) > 1 {
		printPlain(services)
	} else {
		for _, s := range services {
			fmt.Println(serviceToArr(s)[0])
		}
	}
}

func serviceToArr(s Service) []string {
	return []string{s.Port, s.Description, getProtocolString(s.Tcp, s.Udp), s.Status}
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

func getProtocolString(tcp, udp bool) string {
	switch {
	case tcp && udp:
		return "tcp/udp"
	case tcp:
		return "tcp"
	case udp:
		return "udp"
	default:
		return ""
	}
}

func printKeyVal(key, val string) {
	// needs to be even with the key. and the : needs to be in same col every time
	if val == "" {
		return
	}

	fmt.Printf("%-15s: %s\n", key, val)
}
