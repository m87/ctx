package core_test

// func TestCreateContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	state := cs.Load()
// 	assert.Len(t, state.Contexts, 1)
// 	createdContext := state.Contexts[test.TestId]
// 	assert.Equal(t, createdContext.Id, test.TestId)
// 	assert.Equal(t, createdContext.Description, test.TestDescription)
// 	assert.Equal(t, createdContext.State, core.ACTIVE)
// 	assert.Equal(t, createdContext.Duration, time.Duration(0))
// 	assert.Len(t, createdContext.Intervals, 0)
// 	assert.Len(t, createdContext.Comments, 0)
// }

// func TestCreateExistingId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)
// 	err := cm.CreateContext(test.TestId, test.TestDescription)

// 	assert.Error(t, err, errors.New("context already exists"))
// }

// func TestDontCreateContextWithEmptyDescription(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	err := cm.CreateContext(test.TestId, "  \t")

// 	assert.Error(t, err, errors.New("empty description"))
// }

// func TestDontCreateContextWithEmptyId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	err := cm.CreateContext(" \t", test.TestDescription)

// 	assert.Error(t, err, errors.New("empty id"))
// }


// func TestSwitchContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	err := cm.Switch(test.TestId)

// 	state := cs.Load()
// 	assert.NoError(t, err)
// 	assert.Equal(t, test.TestId, state.CurrentId)

// }

// func TestDontSwitchWithEmptyId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	err := cm.Switch("\t")

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("empty id"))
// 	assert.Equal(t, "", state.CurrentId)
// }

// func TestDontSwitchIfDoesNotExists(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)
// 	cm.Switch(test.TestId)
// 	err := cm.Switch("test")

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("context does not exist"))
// 	assert.Equal(t, test.TestId, state.CurrentId)
// }

// func TestSwitchNotExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.PrevTestId, test.TestDescription)

// 	cm.Switch(test.PrevTestId)
// 	err := cm.Switch(test.TestId)

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("context does not exist"))
// 	assert.Equal(t, test.PrevTestId, state.CurrentId)

// }
// func TestSwitchCreateIfNotExists(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())

// 	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	state := cs.Load()
// 	assert.NoError(t, err)
// 	assert.Equal(t, test.TestId, state.CurrentId)
// 	assert.NotNil(t, state.Contexts[test.TestId])
// 	assert.Len(t, state.Contexts, 1)
// 	createdContext := state.Contexts[test.TestId]
// 	assert.Equal(t, createdContext.Id, test.TestId)
// 	assert.Equal(t, createdContext.Description, test.TestDescription)
// 	assert.Equal(t, createdContext.State, core.ACTIVE)
// 	assert.Equal(t, createdContext.Duration, time.Duration(0))
// 	assert.Len(t, createdContext.Intervals, 1)
// 	assert.Len(t, createdContext.Comments, 0)

// }

// func TestDontSwitchOrCreateWithEmptyId(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())

// 	err := cm.CreateIfNotExistsAndSwitch("\t", test.TestDescription)

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("empty id"))
// 	assert.Equal(t, "", state.CurrentId)
// }

// func TestDontSwitchOrCreateWithEmptyDescription(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())

// 	err := cm.CreateIfNotExistsAndSwitch(test.TestId, " \t ")

// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("empty id"))
// 	assert.Equal(t, "", state.CurrentId)
// }

// func TestSwitchCreateIfNotExistsOnExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	assert.Equal(t, cs.Load().CurrentId, "")
// 	err := cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	assert.Equal(t, cs.Load().CurrentId, test.TestId)
// 	assert.NoError(t, err)
// }

// func TestSwitchAlreadyActiveContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	err := cm.Switch(test.TestId)
// 	assert.NoError(t, err)

// 	err = cm.Switch(test.TestId)
// 	state := cs.Load()
// 	assert.Error(t, err, errors.New("context already active"))
// 	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

// 	err = cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	state = cs.Load()
// 	assert.Error(t, err, errors.New("context already active"))
// 	assert.Len(t, state.Contexts[test.TestId].Intervals, 1)

// }

// func TestIntervals(t *testing.T) {
// 	cs := NewTestContextStore()
// 	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), tp)
// 	cm.CreateContext(test.TestId, test.TestDescription)

// 	tp.Current = ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	state := cs.Load()
// 	assert.Equal(t, test.TestId, state.CurrentId)
// 	prevCtx := state.Contexts[state.CurrentId]
// 	assert.Len(t, prevCtx.Intervals, 1)
// 	assert.Equal(t, prevCtx.Intervals[0].Start, tp.Current)
// 	assert.True(t, prevCtx.Intervals[0].End.Time.IsZero())

// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state = cs.Load()
// 	prevCtx = state.Contexts[test.TestId]
// 	assert.Equal(t, prevCtx.Intervals[0].Start, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// 	assert.Equal(t, prevCtx.Intervals[0].End, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, test.PrevTestId, state.CurrentId)
// 	nextCtx := state.Contexts[state.CurrentId]
// 	assert.Len(t, nextCtx.Intervals, 1)
// 	assert.Equal(t, nextCtx.Intervals[0].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.True(t, nextCtx.Intervals[0].End.Time.IsZero())

// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	state = cs.Load()
// 	nextCtx = state.Contexts[test.PrevTestId]
// 	assert.Equal(t, nextCtx.Intervals[0].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, nextCtx.Intervals[0].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, test.TestId, state.CurrentId)
// 	prevCtx = state.Contexts[state.CurrentId]
// 	assert.Len(t, prevCtx.Intervals, 2)
// 	assert.Equal(t, prevCtx.Intervals[1].Start, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.True(t, prevCtx.Intervals[1].End.Time.IsZero())

// }


// func TestFree(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), NewTestTimerProvider("2025-03-13 13:00:00"))
// 	dt, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	assert.Equal(t, test.TestId, cs.Load().CurrentId)

// 	cm.Free()
// 	state := cs.Load()
// 	assert.Equal(t, "", state.CurrentId)
// 	assert.Equal(t, ctxtime.ZonedTime{Time: dt, Timezone: time.UTC.String()}, state.Contexts[test.TestId].Intervals[0].Start)
// 	assert.Equal(t, ctxtime.ZonedTime{Time: dt, Timezone: time.UTC.String()}, state.Contexts[test.TestId].Intervals[0].End)

// }

// func TestFreeWithNowCurrentContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateContext(test.TestId, test.TestDescription)
// 	assert.Equal(t, "", cs.Load().CurrentId)
// 	err := cm.Free()
// 	assert.Error(t, err, errors.New("no active context"))
// }

// func TestDeleteContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	cm.Free()
// 	assert.Len(t, cs.Load().Contexts, 1)
// 	assert.Len(t, es.Load().Events, 4)
// 	err := cm.Delete(test.TestId)
// 	assert.NoError(t, err)
// 	assert.Len(t, es.Load().Events, 5)
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.DELETE_CTX)
// 	assert.Len(t, cs.Load().Contexts, 0)
// }

// func TestDeleteActiveContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	assert.Len(t, cs.Load().Contexts, 1)
// 	err := cm.Delete(test.TestId)
// 	assert.Error(t, err, errors.New("context is active"))
// 	assert.Len(t, cs.Load().Contexts, 1)
// }

// func TestDeleteNotExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.Delete(test.PrevTestId)
// 	assert.Error(t, err, errors.New("context does not exist"))
// 	assert.Len(t, cs.Load().Contexts, 1)
// }


// func TestMergeContexts(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.TestIdExtra, test.PrevDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.Switch(test.TestId)
// 	cm.Switch(test.PrevTestId)
// 	cm.Free()
// 	assert.Len(t, es.Load().Events, 21)
// 	assert.Len(t, cs.Load().Contexts, 3)
// 	assert.Equal(t, cs.Load().Contexts[test.PrevTestId].Duration, dt3.Sub(dt2))

// 	err := cm.MergeContext(test.TestId, test.PrevTestId)

// 	assert.NoError(t, err)
// 	assert.Equal(t, cs.Load().Contexts[test.PrevTestId].Duration, dt3.Sub(dt2)*2)
// 	events := es.Load().Events
// 	assert.Len(t, events, 23)
// 	assert.Len(t, cs.Load().Contexts, 2)
// 	for _, event := range events {
// 		if event.Type == core.SWITCH_CTX {
// 			if v, ok := event.Data["from"]; ok && v != "" && event.CtxId == test.TestIdExtra {
// 				assert.Equal(t, v, test.TestId)
// 			}
// 		}
// 	}
// }


// func TestErrorOnEditCurrentContextInterval(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)

// 	id := cs.Load().Contexts[test.TestId].Intervals[0].Id
// 	err := cm.EditContextInterval(test.TestId, id, ctxtime.ZonedTime{Time: time.Now().Local()}, ctxtime.ZonedTime{Time: time.Now().Local()})

// 	assert.Error(t, errors.New("context is active"), err)

// }

// func TestEditContextInterval(t *testing.T) {
// 	as := NewTestArchiveStore()
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, as, tp)
// 	dt1, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.TestDescription)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.TestDescription)

// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt2.Sub(dt1))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)

// 	id := cs.Load().Contexts[test.TestId].Intervals[0].Id
// 	err := cm.EditContextInterval(test.TestId, id, ctxtime.ZonedTime{Time: dt1}, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt1, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt3.Sub(dt1))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)
// 	assert.NoError(t, err, errors.New("context is active"))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.EDIT_CTX_INTERVAL)
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.start"], dt1.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.end"], dt2.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.start"], dt1.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.end"], dt3.Format(time.RFC3339))

// 	id = cs.Load().Contexts[test.TestId].Intervals[0].Id
// 	err = cm.EditContextInterval(test.TestId, id, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Start, ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].End, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[1].Duration, dt3.Sub(dt2))
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration, cs.Load().Contexts[test.TestId].Intervals[0].Duration+cs.Load().Contexts[test.TestId].Intervals[1].Duration)
// 	assert.NoError(t, err, errors.New("context is active"))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.EDIT_CTX_INTERVAL)
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.start"], dt1.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["old.end"], dt3.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.start"], dt2.Format(time.RFC3339))
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Data["new.end"], dt3.Format(time.RFC3339))
// }

// func TestGetIntervalDurationForDateInBetween(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:05:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-14", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 24*time.Hour, duration)
// }

// func TestGetIntervalDurationForDateOutOfBounds(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:05:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-16", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, time.Duration(0), duration)
// 	date, _ = time.ParseInLocation(time.DateOnly, "2025-03-12", time.UTC)
// 	duration, err = cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, time.Duration(0), duration)
// }

// func TestGetIntervalDurationForDateBefore(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 10:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:00:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-15", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 13*time.Hour, duration)
// }

// func TestGetIntervalDurationForDateAfter(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 10:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-15 13:00:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-13", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 14*time.Hour, duration)
// }

// func TestGetIntervalDurationForDateEqual(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 10:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:00:00", time.UTC)
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	state := cs.Load()
// 	date, _ := time.ParseInLocation(time.DateOnly, "2025-03-13", time.UTC)
// 	duration, err := cm.GetIntervalDurationsByDate(&state, test.TestId, ctxtime.ZonedTime{Time: date, Timezone: time.UTC.String()})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 3*time.Hour, duration)
// }

// func TestDeleteINtervalActiveContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.Error(t, err, errors.New("context is active"))
// }

// func TestDeleteIntervalNotExistingContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.Error(t, err, errors.New("context does not exist"))
// }

// func TestDeleteIntervalOutOfBounds(t *testing.T) {
// 	cs := NewTestContextStore()
// 	cm := core.NewContextManager(cs, NewTestEventsStore(), NewTestArchiveStore(), ctxtime.NewTimer())
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.Error(t, err, errors.New("interval out of bounds"))
// }

// func TestDeleteInterval(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)
// 	dt2, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:05:00", time.UTC)
// 	dt3, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:10:00", time.UTC)
// 	dt4, _ := time.ParseInLocation(time.DateTime, "2025-03-13 13:15:00", time.UTC)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt2, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	tp.Current = ctxtime.ZonedTime{Time: dt4, Timezone: time.UTC.String()}
// 	cm.CreateIfNotExistsAndSwitch(test.PrevTestId, test.PrevDescription)

// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 2)
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration.Seconds(), time.Duration(600000000000).Seconds())

// 	err := cm.DeleteIntervalByIndex(test.TestId, 0)

// 	assert.NoError(t, err)
// 	assert.Len(t, cs.Load().Contexts[test.TestId].Intervals, 1)
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].Start, ctxtime.ZonedTime{Time: dt3, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Intervals[0].End, ctxtime.ZonedTime{Time: dt4, Timezone: time.UTC.String()})
// 	assert.Equal(t, cs.Load().Contexts[test.TestId].Duration.Seconds(), time.Duration(300000000000).Seconds())
// 	// assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.DELETE_CTX_INTERVAL)
// }

// func TestAddLabelToContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.LabelContext(test.TestId, "test-label")

// 	assert.NoError(t, err)
// 	assert.Contains(t, cs.Load().Contexts[test.TestId].Labels, "test-label")
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.LABEL_CTX)
// }

// func TestRemoveLabelFromContext(t *testing.T) {
// 	cs := NewTestContextStore()
// 	tp := NewTestTimerProvider("2025-03-13 13:00:00")
// 	es := NewTestEventsStore()
// 	cm := core.NewContextManager(cs, es, NewTestArchiveStore(), tp)

// 	cm.CreateIfNotExistsAndSwitch(test.TestId, test.TestDescription)
// 	err := cm.LabelContext(test.TestId, "test-label")
// 	assert.NoError(t, err)

// 	err = cm.DeleteLabelContext(test.TestId, "test-label")
// 	assert.NoError(t, err)
// 	assert.NotContains(t, cs.Load().Contexts[test.TestId].Labels, "test-label")
// 	assert.Equal(t, es.Load().Events[len(es.Load().Events)-1].Type, core.DELETE_CTX_LABEL)
// }

