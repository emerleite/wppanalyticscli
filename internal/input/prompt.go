package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// TokenPrompter defines the interface for prompting for tokens
type TokenPrompter interface {
	PromptForToken() (string, error)
}

// SecurePrompter implements TokenPrompter with secure input
type SecurePrompter struct{}

// NewSecurePrompter creates a new secure prompter
func NewSecurePrompter() *SecurePrompter {
	return &SecurePrompter{}
}

// PromptForToken prompts the user for an access token with hidden input
func (p *SecurePrompter) PromptForToken() (string, error) {
	fmt.Fprint(os.Stderr, "Enter Facebook Access Token: ")
	
	// Try to read from terminal with hidden input
	if term.IsTerminal(int(syscall.Stdin)) {
		token, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr) // Add newline after hidden input
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(token)), nil
	}
	
	// Fallback to regular input if not a terminal
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(token), nil
}