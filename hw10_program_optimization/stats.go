package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"sync"

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
	chUser := make(chan User, 1000)
	var errUser error

	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go getUsers(r, chUser, &errUser)
	go countDomains(chUser, &result, re, &wg)

	wg.Wait()

	if errUser != nil {
		return nil, errUser
	}

	return result, nil
}

func getUsers(r io.Reader, ch chan User, errOut *error) {
	defer close(ch)
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)
	var user User
	for fileScanner.Scan() {
		if err := json.Unmarshal(fileScanner.Bytes(), &user); err != nil {
			*errOut = err
			return
		}
		ch <- user
	}
}

func countDomains(ch chan User, result *DomainStat, re *regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	for user := range ch {
		matched := re.Match([]byte(user.Email))
		if matched {
			num := (*result)[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			(*result)[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
}
