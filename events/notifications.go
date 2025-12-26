package events

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/views"
	"github.com/getsentry/sentry-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/palantir/stacktrace"
)

func (handler *EventHandler) RegisterAssignmentHandler(bot *discordgo.Session, signal chan os.Signal) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to start cron").Error())
		return
	}

	task := gocron.NewTask(handler.ProcessAssignments, bot)
	_, err = scheduler.NewJob(gocron.DurationJob(time.Minute), task)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to schedule new job").Error())
	}

	scheduler.Start()

	slog.Info("running assignment notification handler")
	<-signal

	err = scheduler.Shutdown()
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to shut down scheduler").Error())
	}

}

var AssignmentStatusDuration = map[string]time.Duration{
	"one hour": time.Hour,
	"one day":  24 * time.Hour,
}

func (handler *EventHandler) ProcessAssignments(bot *discordgo.Session) {
	transaction := sentry.StartTransaction(context.Background(), "processAssignments")
	defer transaction.Finish()

	for i, currentStatus := range clients.AssignmentStatus[:len(clients.AssignmentStatus)-1] {
		newStatus := clients.AssignmentStatus[i+1]
		assignments, err := handler.HakaseClient.SearchAssignments(transaction, clients.SearchAssignmentsQuery{
			Type: "and",
			Value: []clients.SearchAssignmentsQuery{
				{
					Type:  "lte",
					Field: "due",
					Value: time.Now().UTC().Add(AssignmentStatusDuration[newStatus]),
				},
				{
					Type:  "eq",
					Field: "status",
					Value: currentStatus,
				},
			},
		})
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to get assignments due in %s", newStatus).Error())
		} else {
			for _, assignment := range assignments {
				go handler.sendAssignmentNotification(transaction, bot, assignment, newStatus)
				go handler.updateAssignmentStatus(transaction, assignment, newStatus)
			}
			slog.Info("sent notifications for assignments due in one day", "count", len(assignments))
		}
	}

	toBeDeleted, err := handler.HakaseClient.SearchAssignments(transaction, clients.SearchAssignmentsQuery{
		Type: "and",
		Value: []clients.SearchAssignmentsQuery{
			{
				Type:  "lte",
				Field: "due",
				Value: time.Now().UTC(),
			},
			{
				Type:  "eq",
				Field: "status",
				Value: clients.AssignmentStatus[2],
			},
		},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get assignments to be deleted").Error())
	} else {
		for _, assignment := range toBeDeleted {
			err := handler.HakaseClient.DeleteAssignment(transaction, assignment.ID)
			if err != nil {
				slog.Error(stacktrace.Propagate(err, "failed to delete assignment with id: %s", assignment.ID).Error())
			}
		}
		slog.Info("deleted overdue assignments", "count", len(toBeDeleted))
	}

}

func (handler *EventHandler) sendAssignmentNotification(span *sentry.Span, bot *discordgo.Session, assignment clients.Assignment, newStatus string) {
	course, err := handler.HakaseClient.ReadCourse(span, assignment.CourseID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to read course of assignment: %s", assignment.ID).Error())
		return
	}

	guild, err := bot.Guild(course.CourseID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to retrieve guild system channel for %s", course.CourseID).Error())
		return
	}

	notifyChannel := course.NotifyChannel
	if notifyChannel == "" {
		notifyChannel = guild.SystemChannelID
	}

	notifyTarget := course.NotifyGroup
	if notifyTarget == "" {
		notifyTarget = fmt.Sprintf("<@%s>", guild.OwnerID)
	} else {
		notifyTarget = fmt.Sprintf("<@&%s>", notifyTarget)
	}

	_, err = bot.ChannelMessageSendComplex(notifyChannel, &discordgo.MessageSend{
		Content: notifyTarget,
		Embeds:  []*discordgo.MessageEmbed{views.NotificationView(assignment, newStatus)},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to send due date notification for assignment: %s", assignment.ID).Error())
	}
}

func (handler *EventHandler) updateAssignmentStatus(span *sentry.Span, assignment clients.Assignment, newStatus string) {
	assignment.Assignment = assignment.ID
	assignment.Status = newStatus
	err := handler.HakaseClient.UpdateAssignment(span, assignment)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update assignment status: %s", assignment.ID).Error())
	}

}
