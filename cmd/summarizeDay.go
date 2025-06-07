package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/m87/ctx/ctx"
	"github.com/m87/ctx/ctx_model"
	"github.com/m87/ctx/util"
	"github.com/spf13/cobra"
)

func roundDuration(d time.Duration, unit string) time.Duration {
	switch unit {
	case "nanosecond":
		return d.Round(time.Nanosecond)
	case "microsecond":
		return d.Round(time.Microsecond)
	case "millisecond":
		return d.Round(time.Millisecond)
	case "second":
		return d.Round(time.Second)
	case "minute":
		return d.Round(time.Minute)
	case "hour":
		return d.Round(time.Hour)
	default:
		return d.Round(time.Nanosecond)
	}
}

var jiraLineRegex = regexp.MustCompile(`^([A-Z][A-Z0-9]+-\d+)\s*(.*)$`)

func GenerateJiraCurlCommand(entry string, duration time.Duration, started time.Time, jiraBaseURL, authToken string) string {
	matches := jiraLineRegex.FindStringSubmatch(entry)
	if matches == nil {
		return ""
	}

	issueKey := matches[1]
	comment := strings.TrimSpace(matches[2])

	minutes := int(duration.Minutes())
	if minutes <= 0 {
		minutes = 1
	}
	timeSpent := fmt.Sprintf("%dm", minutes)
	startedFormatted := started.Format("2006-01-02T15:04:05.000-0700")

	jsonBody := fmt.Sprintf(`{
  "timeSpent": "%s",
  "started": "%s",
  "comment": "%s"
}`, timeSpent, startedFormatted, comment)

	curl := fmt.Sprintf(`curl --fail -s -o /dev/null -w "%%{http_code}" \
  -X POST \
  -H "Authorization: Basic %s" \
  -H "Content-Type: application/json" \
  --data '%s' \
  %s/rest/api/3/issue/%s/worklog \
  || echo "❌ Nie udało się zalogować czasu do zadania %s"`,
		authToken,
		jsonBody,
		jiraBaseURL,
		issueKey,
		issueKey,
	)

	return curl
}

var summarizeDayCmd = &cobra.Command{
	Use:     "day",
	Aliases: []string{"d", "day"},
	Short:   "Summarize day",
	Run: func(cmd *cobra.Command, args []string) {
		roundUnit, _ := cmd.Flags().GetString("round")
		loc, err := time.LoadLocation(ctx_model.DetectTimezoneName())
		if err != nil {
			loc = time.UTC
		}
		mgr := ctx.CreateManager()
		date := mgr.TimeProvider.Now().Time.In(loc)
		if len(args) > 0 {
			rawDate := strings.TrimSpace(args[0])

			if rawDate != "" {
				date, _ = time.ParseInLocation(time.DateOnly, rawDate, loc)
			}
		}

		durations := map[string]time.Duration{}
		overallDuration := time.Duration(0)

		mgr.ContextStore.Read(func(s *ctx_model.State) error {
			for ctxId, _ := range s.Contexts {
				d, err := mgr.GetIntervalDurationsByDate(s, ctxId, ctx_model.ZonedTime{Time: date, Timezone: loc.String()})
				util.Checkm(err, "Unable to get interval durations for context "+ctxId)
				durations[ctxId] = roundDuration(d, roundUnit)
			}
			return nil
		})

		sortedIds := make([]string, 0, len(durations))
		for k := range durations {
			sortedIds = append(sortedIds, k)
		}
		sort.Strings(sortedIds)

		for _, c := range sortedIds {
			d := durations[c]
			ctx, _ := mgr.Ctx(c)
			if d > 0 {
				fmt.Printf("- %s: %s\n", ctx.Description, d)
				overallDuration += d
				if f, _ := cmd.Flags().GetBool("verbose"); f {
					mgr.ContextStore.Read(func(s *ctx_model.State) error {
						for _, interval := range mgr.GetIntervalsByDate(s, c, ctx_model.ZonedTime{Time: date, Timezone: loc.String()}) {
							fmt.Printf("\t- %s - %s\n", interval.Start.Time.Format(time.RFC3339), interval.End.Time.Format(time.RFC3339))
						}
						return nil
					})
				}
			}
		}
		mgr.ContextStore.Read(func(s *ctx_model.State) error {
			if f, _ := cmd.Flags().GetBool("jira"); f {
				fmt.Println("\nJira curl commands:")
				jiraBaseURL := "https://your-jira-instance.atlassian.net"
				authToken := "your_base64_encoded_auth_token"
				for _, c := range sortedIds {
					d := durations[c]
					ctx, _ := mgr.Ctx(c)
					if d > 0 {
						for _, interval := range mgr.GetIntervalsByDate(s, c, ctx_model.ZonedTime{Time: date, Timezone: loc.String()}) {
							curlCommand := GenerateJiraCurlCommand(ctx.Description, interval.End.Time.Sub(interval.Start.Time), interval.Start.Time, jiraBaseURL, authToken)
							fmt.Println(curlCommand)
						}
					}
				}

			}
			return nil
		})
		fmt.Printf("Overall: %s\n", overallDuration)
	},
}

func init() {
	summarizeCmd.AddCommand(summarizeDayCmd)
	summarizeDayCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	summarizeDayCmd.Flags().StringP("round", "r", "nanosecond", "Round to the nearest nanosecond, microsecond, millisecond, second, minute, hour")
	summarizeDayCmd.Flags().BoolP("jira", "j", false, "Generate curl commands for Jira time tracking")
}
