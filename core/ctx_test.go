package core_test






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

