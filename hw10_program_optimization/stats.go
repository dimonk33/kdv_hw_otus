package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/goccy/go-json"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	if err := countDomains(r, domain, &result); err != nil {
		return DomainStat{}, err
	}

	return result, nil
}

func countDomains(r io.Reader, domain string, result *DomainStat) error {
	var user User
	var index string
	var num int
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		if err := json.Unmarshal(fileScanner.Bytes(), &user); err != nil {
			return err
		}
		if strings.Contains(user.Email, "."+domain) {
			index = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num = (*result)[index]
			num++
			(*result)[index] = num
		}
	}
	return nil
}
