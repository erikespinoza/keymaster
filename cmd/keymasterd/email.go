package main

import (
	"bytes"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/Cloud-Foundations/golib/pkg/communications/configuredemail"
)

const emailAdminTemplateData = `
{{define "Bootstrap OTP Admin Email"}}
From: {{.InitiatorAddr}}
To: {{.AdminAddrs}}
Subject: Keymaster Bootstrap OTP generated for {{.Username}}

A Bootstrap OTP was generated by {{.InitiatorUser}} for user: {{.Username}}

The OTP fingerprint is: {{printf "%x" .Fingerprint}}

The OTP will expire in {{.Duration}}.

The user profile may be viewed at: {{.HostIdentity}}/profile/{{.Username}}
{{end}}
`

const emailUserTemplateData = `
{{define "Bootstrap OTP User Email"}}
From: {{.InitiatorAddr}}
To: {{.UserAddr}}
Subject: Welcome to Keymaster

Hi, {{.Username}}. Welcome to Keymaster. Please log in to:
{{.HostIdentity}}

with your username and password. After this step you will be asked to enter
your one-time passcode (Bootstrap OTP) which is:
{{.OTP}}

In case of debugging, your OTP fingerprint is: {{printf "%x" .Fingerprint}}

Please register your U2F security key after login.

You have {{.Duration}} to complete this operation before this passcode expires.
{{end}}
`

const emailTimeout = time.Second * 15

type bootstrapOtpEmailData struct {
	AdminAddrs    string
	Duration      time.Duration
	Fingerprint   [4]byte
	HostIdentity  string
	InitiatorAddr string
	InitiatorUser string
	OTP           string
	UserAddr      string
	Username      string
}

func (state *RuntimeState) initEmailDefaults() {
	state.Config.Email.AwsSecretLifetime = time.Minute * 5
}

func (state *RuntimeState) setupEmail() error {
	if state.Config.Email.Domain == "" {
		return nil
	}
	var err error
	state.emailManager, err = configuredemail.New(
		state.Config.Email.EmailConfig, logger)
	if err != nil {
		return err
	}
	return nil
}

func (state *RuntimeState) sendBootstrapOtpEmail(hash []byte, OTP string,
	duration time.Duration, initiatorUser, targetUser string) error {
	emailData := bootstrapOtpEmailData{
		Duration:      duration,
		HostIdentity:  state.Config.Base.HostIdentity,
		OTP:           OTP,
		InitiatorAddr: initiatorUser + "@" + state.Config.Email.Domain,
		InitiatorUser: initiatorUser,
		UserAddr:      targetUser + "@" + state.Config.Email.Domain,
		Username:      targetUser,
	}
	copy(emailData.Fingerprint[:], hash[:4])
	adminUsers := make(map[string]struct{})
	adminUsers[initiatorUser] = struct{}{}
	for _, user := range state.Config.Base.AdminUsers {
		adminUsers[user] = struct{}{}
	}
	// TODO(rgooch): refactor to use a UserInfo interface which is set up in
	//               config.go. This will require changes in several files.
	if state.gitDB != nil {
		for _, adminGroup := range state.Config.Base.AdminGroups {
			users, err := state.gitDB.GetUsersInGroup(adminGroup)
			if err != nil {
				return err
			}
			for _, user := range users {
				adminUsers[user] = struct{}{}
			}
		}
	}
	adminAddrs := make([]string, 0, len(adminUsers))
	for user := range adminUsers {
		adminAddrs = append(adminAddrs, user+"@"+state.Config.Email.Domain)
	}
	sort.Strings(adminAddrs)
	emailData.AdminAddrs = strings.Join(adminAddrs, ",")
	buffer := &bytes.Buffer{}
	err := state.textTemplates.ExecuteTemplate(buffer,
		"Bootstrap OTP Admin Email", emailData)
	if err != nil {
		return err
	}
	err = state.sendMail(emailData.InitiatorAddr, adminAddrs, buffer.Bytes(),
		emailTimeout)
	if err != nil {
		return err
	}
	buffer = &bytes.Buffer{}
	err = state.textTemplates.ExecuteTemplate(buffer,
		"Bootstrap OTP User Email", emailData)
	if err != nil {
		return err
	}
	err = state.sendMail(emailData.InitiatorAddr, []string{emailData.UserAddr},
		buffer.Bytes(), emailTimeout)
	if err != nil {
		return err
	}
	return nil
}

func (state *RuntimeState) sendMail(from string, to []string, msg []byte,
	timeout time.Duration) error {
	for msg[0] == '\n' { // Leading blank lines would prevent email being sent.
		msg = msg[1:]
	}
	timer := time.NewTimer(timeout)
	errorChan := make(chan error)
	go func() {
		err := state.emailManager.SendMail(from, to, msg)
		select {
		case errorChan <- err: // Status was consumed.
		default: // Status consumer has gone. Log it.
			if err != nil {
				state.logger.Printf("late failure sending email: %s\n", err)
			} else {
				state.logger.Println("late success sending email")
			}
		}
	}()
	select {
	case err := <-errorChan:
		return err
	case <-timer.C:
		return errors.New("timed out sending email")
	}
}
