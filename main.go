// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"strings"
// )

// func main() {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("Enter domains to check (press Ctrl+D to exit):")
// 	for scanner.Scan() {
// 		domain := strings.TrimSpace(scanner.Text())
// 		if domain != "" {
// 			checkDomain(domain)
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatalf("Error reading input: %v\n", err)
// 	}
// }

// func checkDomain(domain string) {
// 	var hasMx, hasSPF, hasDMARC bool
// 	var spfRecord, DMARCRecord string

// 	mxRecords, err := net.LookupMX(domain)
// 	if err != nil {
// 		log.Printf("Error looking up MX records for %s: %v\n", domain, err)
// 	} else {
// 		if len(mxRecords) > 0 {
// 			hasMx = true
// 		}
// 	}

// 	txtRecords, err := net.LookupTXT(domain)
// 	if err != nil {
// 		log.Printf("Error looking up TXT records for %s: %v\n", domain, err)
// 	} else {
// 		for _, record := range txtRecords {
// 			if strings.HasPrefix(record, "v=spf1") {
// 				hasSPF = true
// 				spfRecord = record
// 				break
// 			}
// 		}
// 	}

// 	dmarcDomain := "_dmarc." + domain
// 	DMARCRecords, err := net.LookupTXT(dmarcDomain)
// 	if err != nil {
// 		log.Printf("Error looking up TXT records for %s: %v\n", dmarcDomain, err)
// 	} else {
// 		for _, record := range DMARCRecords {
// 			if strings.HasPrefix(record, "v=DMARC1") {
// 				hasDMARC = true
// 				DMARCRecord = record
// 				break
// 			}
// 		}
// 	}

// 	fmt.Printf("Domain: %s\nMX Records: %v\nSPF: %v\nSPF Record: %s\nDMARC: %v\nDMARC Record: %s\n\n",
// 		domain, hasMx, hasSPF, spfRecord, hasDMARC, DMARCRecord)
// }


package main

func main() {
    verify()
}

//If we run the main.go program with the go run main.go command, 
//you will get an error: ./main.go:4:2: undefined: verify. 
//This is because you have not specified all the main package files that are needed to run the application. 
//If you run this program with the go run main.go verify.go command, then it will run without errors:

//If we run the verify.go program with the go run verify.go command with the main pkg same as in main will through error 
//pkg main not def but if we call it in main and then call go run main.go verify.go will run very well.
