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

func (handler *EventHandler) RegisterNotificationHandler(bot *discordgo.Session, signal chan os.Signal) {
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
	defer scheduler.Shutdown()

	slog.Info("running assignment notification handler")
	<-signal

}

func (handler *EventHandler) ProcessAssignments(bot *discordgo.Session) {
	transaction := sentry.StartTransaction(context.Background(), "sendAssignmentNotifications")
	defer transaction.Finish()

	dueInOneDay, err := handler.HakaseClient.SearchAssignments(transaction, clients.SearchAssignmentsQuery{
		Type:  "eq",
		Field: "status",
		Value: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get assignments due in one day").Error())
	} else {
		for _, assignment := range dueInOneDay {
			assignmentStatus := clients.AssignmentStatus[1]
			go handler.SendAssignmentNotification(transaction, bot, assignment, assignmentStatus)
			go handler.UpdateAssignmentStatus(transaction, assignment, assignmentStatus)
		}
		slog.Info(fmt.Sprintf("sent notifications for %d assignments due in one day.", len(dueInOneDay)))
	}

	dueInOneHour, err := handler.HakaseClient.SearchAssignments(transaction, clients.SearchAssignmentsQuery{
		Type:  "lte",
		Field: "due",
		Value: time.Now().Add(time.Hour).Format(time.RFC3339),
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get assignments due in one hour").Error())
	} else {
		for _, assignment := range dueInOneHour {
			assignmentStatus := clients.AssignmentStatus[2]
			go handler.SendAssignmentNotification(transaction, bot, assignment, assignmentStatus)
			go handler.UpdateAssignmentStatus(transaction, assignment, assignmentStatus)
		}
		slog.Info(fmt.Sprintf("sent notifications for %d assignments due in one hour.", len(dueInOneHour)))
	}

}

func (handler *EventHandler) SendAssignmentNotification(span *sentry.Span, bot *discordgo.Session, assignment clients.Assignment, timeBuffer string) {
	course, err := handler.HakaseClient.ReadCourse(span, assignment.CourseID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to read course of assignment: %s", assignment.ID).Error())
		return
	}

	notifyChannel := course.NotifyChannel
	if notifyChannel == "" {
		guild, err := bot.Guild(course.CourseID)
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to retrieve guild system channel for %s", course.CourseID).Error())
			return
		}
		notifyChannel = guild.SystemChannelID
	}

	_, err = bot.ChannelMessageSendEmbed(notifyChannel, views.NotificationView(assignment, timeBuffer))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to send due date notification for assignment: %s", assignment.ID).Error())
	}
}

func (handler *EventHandler) UpdateAssignmentStatus(span *sentry.Span, assignment clients.Assignment, assignmentStatus string) {
	assignment.Status = assignmentStatus
	err := handler.HakaseClient.UpdateAssignment(span, assignment)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to update assignment status: %s", assignment.ID).Error())
	}
}
