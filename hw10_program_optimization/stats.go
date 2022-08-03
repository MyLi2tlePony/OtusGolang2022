package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type (
	DomainStat map[string]int
)

var regexpForDomain = regexp.MustCompile(`[^@]+$`)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	email, err := getUsersEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return countDomains(email, domain)
}

func getUsersEmails(r io.Reader) (<-chan string, error) {
	emails := make(chan string, 1000)
	scanner := bufio.NewScanner(r)

	go func() {
		defer close(emails)

		for i := 0; scanner.Scan(); i++ {
			emails <- fastjson.GetString(scanner.Bytes(), "Email")
		}
	}()

	return emails, scanner.Err()
}

func countDomains(emails <-chan string, domainPart string) (DomainStat, error) {
	result := make(DomainStat)
	regexpForEmail := regexp.MustCompile(`\.` + domainPart)

	for email := range emails {
		if regexpForEmail.MatchString(email) {
			domain := strings.ToLower(regexpForDomain.FindString(email))
			result[domain]++
		}
	}

	return result, nil
}
