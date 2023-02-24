package users

import (
    "fmt"
    "net/mail"
    "strings"
    "unicode"
)

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (u *RegisterRequest) Validate(usernameExists, emailExists bool) map[string]string {
    errors := make(map[string]string)

    var usernameErrors []string
    if len(u.Username) == 0 {
        usernameErrors = append(usernameErrors, "Required")
    } else if usernameExists {
        usernameErrors = append(usernameErrors, "Username already in use")
    } else {
        usernameErrors = append(usernameErrors, u.validateUsername()...)
    }

    if len(usernameErrors) > 0 {
        errors["username"] = strings.Join(usernameErrors, ", ")
    }

    var emailErrors []string
    if len(u.Email) == 0 {
        emailErrors = append(emailErrors, "Required")
    } else if emailExists {
        emailErrors = append(emailErrors, "Email already in use")
    } else {
        parser := mail.AddressParser{}
        _, err := parser.Parse(u.Email)
        if err != nil {
            emailErrors = append(emailErrors, "Invalid email address")
        }
    }

    if len(emailErrors) > 0 {
        errors["email"] = strings.Join(emailErrors, ", ")
    }

    var passwordErrors []string
    if len(u.Password) == 0 {
        passwordErrors = append(passwordErrors, "Required")
    } else {
        err := u.validatePassword()
        if err != nil {
            passwordErrors = append(passwordErrors, err.Error())
        }
    }

    if len(passwordErrors) > 0 {
        errors["password"] = strings.Join(passwordErrors, ", ")
    }

    return errors
}

func (u *RegisterRequest) validateUsername() []string {
    var errors []string

    const minLength = 6
    const maxLength = 30
    var usernameLen int

    for ind, ch := range u.Username {
        switch {
        case unicode.IsNumber(ch):
            if ind == 0 {
                errors = append(errors, "Cannot start with a number")
            }
            usernameLen++
        case unicode.IsPunct(ch) && string(ch) == "_":
            if ind == 0 {
                errors = append(errors, "Cannot start with an underscore")
            }
            usernameLen++
        case unicode.IsUpper(ch) || unicode.IsLower(ch):
            usernameLen++
        default:
            errors = append(errors, "Only alphanumeric characters and underscores (_) allowed")
            usernameLen++
        }
    }

    if usernameLen < minLength || usernameLen > maxLength {
        errors = append(errors, fmt.Sprintf("Must be between %d and %d characters long", minLength, maxLength))
    }

    return errors
}

func (u *RegisterRequest) validatePassword() error {
    const minLength = 6
    const maxLength = 64

    if len(u.Password) < minLength || len(u.Password) > maxLength {
        return fmt.Errorf("Must be between %d and %d characters long", minLength, maxLength)
    }

    return nil
}
