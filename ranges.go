package timez

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	defaultTimeZone = "Etc/UTC"

	dateFmt = "2006-01-02"
)

type TimeRange struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

type DateRange TimeRange

type dateRangeWithStrDate struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Timezone string `json:"timezone"`
}

func (d *DateRange) ToTimeRange() *TimeRange {
	tr := TimeRange(*d)
	return &tr
}

func (d *DateRange) MarshalJSON() ([]byte, error) {
	startTz, _ := d.Start.Zone()
	endTz, _ := d.End.Zone()

	if startTz != endTz {
		return nil, errors.New("range dates have different timezone")
	}

	return json.Marshal(&dateRangeWithStrDate{
		Start:    d.Start.Format(dateFmt),
		End:      d.End.Format(dateFmt),
		Timezone: startTz,
	})
}

func (d *DateRange) UnmarshalJSON(b []byte) error {
	dateRangeWithStrDate := &dateRangeWithStrDate{}
	err := json.Unmarshal(b, dateRangeWithStrDate)
	if err != nil {
		return err
	}

	tz := getTimezone(dateRangeWithStrDate.Timezone)
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}

	start, err := time.ParseInLocation(dateFmt, dateRangeWithStrDate.Start, location)
	if err != nil {
		return err
	}

	d.Start = &start

	end, err := time.ParseInLocation(dateFmt, dateRangeWithStrDate.End, location)
	if err != nil {
		return err
	}

	end = end.AddDate(0, 0, 1).Add(time.Nanosecond * -1)
	d.End = &end

	return nil
}

func getTimezone(srcTimezone string) string {
	if srcTimezone == "" {
		return defaultTimeZone
	}

	return srcTimezone
}
