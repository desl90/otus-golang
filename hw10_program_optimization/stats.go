package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
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
	s := bufio.NewScanner(r)
	user := &User{}
	result := make(DomainStat)

	reg, _ := regexp.Compile("(?m)(@)[^.]+\\." + domain)

	for s.Scan() {
		if err := easyjson.Unmarshal(s.Bytes(), user); err != nil {
			return nil, fmt.Errorf("error unmarshal user: %w", err)
		}

		if reg.MatchString(user.Email) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
