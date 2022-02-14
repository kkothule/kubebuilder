package controllers

import (
	"time"

	kbatch "k8s.io/api/batch/v1"
	tzcronjobberv1 "kkothule.io/tzcronjobber/api/v1"
)

func getCurrentTimeInZone(sj *tzcronjobberv1.TZCronjobber) (time.Time, error) {
	if sj.Spec.TimeZone == "" {
		return time.Now(), nil
	}

	loc, err := time.LoadLocation(sj.Spec.TimeZone)
	if err != nil {
		return time.Now(), err
	}

	return time.Now().In(loc), nil
}

// byJobStartTime sorts a list of jobs by start timestamp, using their names as a tie breaker.
type byJobStartTime []kbatch.Job

func (o byJobStartTime) Len() int      { return len(o) }
func (o byJobStartTime) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o byJobStartTime) Less(i, j int) bool {
	if o[j].Status.StartTime == nil {
		return o[i].Status.StartTime != nil
	}

	if o[i].Status.StartTime.Equal(o[j].Status.StartTime) {
		return o[i].Name < o[j].Name
	}

	return o[i].Status.StartTime.Before(o[j].Status.StartTime)
}
