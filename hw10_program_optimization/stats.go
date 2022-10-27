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
	regexpForEmail := regexp.MustCompile(`\.` + domain)
	result := make(DomainStat)
	scan := bufio.NewScanner(r)

	for i := 0; scan.Scan(); i++ {
		email := fastjson.GetString(scan.Bytes(), "Email")

		if regexpForEmail.MatchString(email) {
			domain := strings.ToLower(regexpForDomain.FindString(email))
			result[domain]++
		}
	}

	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}
