// Package views provides Discord message embeds and components for assignments.
package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

// AssignmentView returns a Discord message embed for the given assignment and member.
// It displays assignment details such as name, due date, author, and link.
func AssignmentView(member *discordgo.Member, assignment clients.Assignment) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       assignment.Name,
		Description: fmt.Sprintf("due %s", assignment.Due.Format(time.RFC1123)),
		Author:      &discordgo.MessageEmbedAuthor{Name: member.User.Username, IconURL: member.User.AvatarURL("")},
		URL:         assignment.URL,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("id %s", assignment.ID)},
	}
}

// AssignmentActions returns action buttons for editing or removing the given assignment.
func AssignmentActions(assignment clients.Assignment) *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "📝",
				},
				Label:    "edit",
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("updateAssignmentAction_%s", assignment.ID),
			},
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "🗑️",
				},
				Label:    "remove",
				Style:    discordgo.SecondaryButton,
				CustomID: fmt.Sprintf("deleteAssignmentAction_%s", assignment.ID),
			},
		},
	}
}

// AssignmentModal returns modal components for creating or updating an assignment.
// If assignment is nil, it creates a new assignment modal.
func AssignmentModal(assignment *clients.Assignment) []discordgo.MessageComponent {

	newAssignment := assignment == nil

	now := time.Now()
	if newAssignment {
		assignment = &clients.Assignment{
			// placeholder data
			Name: "Assignment 1",
			Due:  &now,
			URL:  "https://canvas.instructure.com",
		}
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentName",
					Label:       "assignment name:",
					Style:       discordgo.TextInputShort,
					Placeholder: assignment.Name,
					Required:    newAssignment,
					MaxLength:   50,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentDue",
					Label:       "due date:",
					Style:       discordgo.TextInputShort,
					Placeholder: assignment.Due.Format(time.RFC1123),
					Required:    newAssignment,
					MaxLength:   50,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "assignmentLink",
					Label:       "link:",
					Style:       discordgo.TextInputShort,
					Placeholder: assignment.URL,
					Required:    false,
					MaxLength:   50,
				},
			},
		},
	}
}
