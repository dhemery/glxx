package describe

import "strings"

func printParticipation(p *Person, role string) {
	label := strings.ToUpper(role[:1]) + role[1:] + ":"
	printPersonReference(label, p)
}
