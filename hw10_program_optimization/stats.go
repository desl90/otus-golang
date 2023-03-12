package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
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
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	s := bufio.NewScanner(r)
	user := &User{}
	result := make(DomainStat)

	reg, _ := regexp.Compile("(.*)@[^.]+\\." + domain)

	for s.Scan() {
		if err := json.Unmarshal(s.Bytes(), user); err != nil {
			return nil, fmt.Errorf("error unmarshal user: %w", err)
		}

		if reg.MatchString(user.Email) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
