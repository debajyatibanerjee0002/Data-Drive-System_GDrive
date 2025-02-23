package utility

import "strings"

func GetUserName(email string) (string, int) {
	atIndex := strings.Index(email, "@") // Find '@' position
	if atIndex == -1 {
		return "", -1 // Return empty if '@' not found
	}
	return email[:atIndex], atIndex // Extract everything before '@'
}
