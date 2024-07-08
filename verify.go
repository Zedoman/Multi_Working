package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/gocolly/colly/v2"
)

func verify() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Verify Email")
		fmt.Println("2. Verify Domain")
		fmt.Println("3. Find Emails from Website")
		fmt.Println("4. Exit")
		fmt.Print("Enter your choice: ")
		var choice int
		if !scanner.Scan() {
			log.Fatalf("Error reading input: %v\n", scanner.Err())
		}
		choice = getInt(scanner.Text())

		switch choice {
		case 1:
			fmt.Print("Enter email to verify: ")
			if !scanner.Scan() {
				log.Fatalf("Error reading input: %v\n", scanner.Err())
			}
			email := strings.TrimSpace(scanner.Text())
			if email != "" {
				verifyEmail(email)
			}

		case 2:
			fmt.Print("Enter domain to verify: ")
			if !scanner.Scan() {
				log.Fatalf("Error reading input: %v\n", scanner.Err())
			}
			domain := strings.TrimSpace(scanner.Text())
			if domain != "" {
				verifyDomain(domain)
			}

		case 3:
			fmt.Print("Enter website to find emails (e.g., https://example.com): ")
			if !scanner.Scan() {
				log.Fatalf("Error reading input: %v\n", scanner.Err())
			}
			website := strings.TrimSpace(scanner.Text())
			if website != "" {
				findEmailsFromWebsite(website)
			}

		case 4:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Please choose again.")
		}
	}
}

func verifyEmail(email string) {
	fmt.Printf("Verifying email: %s\n", email)

	// Step 1: Format validation using regular expression
	if !isValidEmailFormat(email) {
		fmt.Printf("Invalid email format: %s\n", email)
		return
	}

	// Step 2: Domain validation by checking MX records
	parts := strings.Split(email, "@")
	domain := parts[1]

	mxRecords, err := lookupMXRecords(domain)
	if err != nil {
		log.Printf("Error looking up MX records for %s: %v\n", domain, err)
		fmt.Printf("Domain %s does not have valid MX records\n", domain)
		return
	}

	if len(mxRecords) == 0 {
		fmt.Printf("Domain %s does not have valid MX records\n", domain)
		return
	}

	// Step 3: SMTP validation
	if !isValidSMTP(email, mxRecords) {
		fmt.Printf("Email %s does not exist\n", email)
		return
	}

	fmt.Printf("Email %s is valid\n", email)
}

func isValidSMTP(email string, mxRecords []*net.MX) bool {
	// Try each MX record until a valid one is found
	for _, mx := range mxRecords {
		client, err := smtp.Dial(mx.Host + ":25")
		if err != nil {
			log.Printf("Error connecting to SMTP server %s: %v\n", mx.Host, err)
			continue // Try the next MX record
		}
		defer client.Close()

		// Set the HELO message using your own domain
		err = client.Hello("gmail.com")
		if err != nil {
			log.Printf("Error sending HELO command: %v\n", err)
			continue // Try the next MX record
		}

		// Set the sender using a valid email address from your domain
		err = client.Mail("avra6269@gmail.com")
		if err != nil {
			log.Printf("Error setting sender email: %v\n", err)
			continue // Try the next MX record
		}

		// Set the recipient
		err = client.Rcpt(email)
		if err != nil {
			log.Printf("SMTP server %s rejected recipient: %v\n", mx.Host, err)
			continue // Try the next MX record
		}

		// If no errors, the email is valid
		return true
	}

	// If all MX records failed, the email is not valid
	return false
}


func lookupMXRecords(domain string) ([]*net.MX, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second,
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}

	var mxRecords []*net.MX
	var err error

	for i := 0; i < 3; i++ { // Retry up to 3 times
		mxRecords, err = resolver.LookupMX(context.Background(), domain)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return mxRecords, err
}

func isValidEmailFormat(email string) bool {
	// Regex pattern for basic email format validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func verifyDomain(domain string) {
	fmt.Printf("Verifying domain: %s\n", domain)
	// Implement domain verification logic here
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, DMARCRecord string

	mxRecords, err := lookupMXRecords(domain)
	if err != nil {
		log.Printf("Error looking up MX records for %s: %v\n", domain, err)
	} else {
		if len(mxRecords) > 0 {
			hasMx = true
		}
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v\n", domain, err)
	} else {
		for _, record := range txtRecords {
			if strings.HasPrefix(record, "v=spf1") {
				hasSPF = true
				spfRecord = record
				break
			}
		}
	}

	dmarcDomain := "_dmarc." + domain
	DMARCRecords, err := net.LookupTXT(dmarcDomain)
	if err != nil {
		log.Printf("Error looking up TXT records for %s: %v\n", dmarcDomain, err)
	} else {
		for _, record := range DMARCRecords {
			if strings.HasPrefix(record, "v=DMARC1") {
				hasDMARC = true
				DMARCRecord = record
				break
			}
		}
	}

	fmt.Printf("Domain: %s\nMX Records: %v\nSPF: %v\nSPF Record: %s\nDMARC: %v\nDMARC Record: %s\n\n",
		domain, hasMx, hasSPF, spfRecord, hasDMARC, DMARCRecord)
}



//Data Scrapping kr rha ha ye 
func findEmailsFromWebsite(website string) {
	fmt.Printf("Finding emails from website: %s\n", website)

	// Create a new collector
	c := colly.NewCollector(
		colly.Async(true), // Use asynchronous mode
	)

	// Set a limit for the number of visited links
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",    // Apply to all domains
		Parallelism: 2,      // Number of parallel requests
		RandomDelay: 2 * time.Second, // Random delay between requests
	})

	// Regular expression to match email addresses
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// On every HTML element call callback
	c.OnHTML("html", func(e *colly.HTMLElement) {
		emails := emailRegex.FindAllString(e.Text, -1)
		for _, email := range emails {
			fmt.Println("Found email:", email)
		}
	})

	// On every link found, visit the link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		//fmt.Println("Visiting link:", link)
		e.Request.Visit(link)
	})

	// Start scraping the website
	err := c.Visit(website)
	if err != nil {
		log.Printf("Error visiting website %s: %v\n", website, err)
	}

	// Wait until all asynchronous tasks are finished
	c.Wait()
}

func getInt(input string) int {
	var num int
	_, err := fmt.Sscanf(input, "%d", &num)
	if err != nil {
		log.Println("Error converting input to integer:", err)
	}
	return num
}