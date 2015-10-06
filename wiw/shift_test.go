package wiw

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSpec(t *testing.T) {

	Convey("Given an empty collection of shifts", t, func() {
		var shifts []Shift

		NoHoursWorked := map[int][]float64{}

		Convey("There should be not any hours worked", func() {
			summ := SummarizeShifts(shifts)
			So(summ, ShouldResemble, NoHoursWorked)
		})

		Convey("Add a shift with an entire work day: 8 hours", func() {
			ts := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
			tsPlusAWorkDay := time.Date(2015, 1, 1, 8, 0, 0, 0, time.UTC)

			shifts = append(shifts, Shift{StartTime: ts, EndTime: tsPlusAWorkDay})

			Convey("Hours worked should be 8", func() {
				summ := SummarizeShifts(shifts)

				year, weekNumber := ts.ISOWeek()
				weekNumber-- //iso week number goes from 1 to 52 or 53

				expected := map[int][]float64{
					year: make([]float64, ISOWeeksCount(ts)),
				}

				expected[year][weekNumber] = WorkHoursPerDay

				So(summ, ShouldResemble, expected)

				Convey("Add another shift with 8 hours, so total should be 16", func() {
					shifts = append(shifts, Shift{StartTime: ts, EndTime: tsPlusAWorkDay})
					summ = SummarizeShifts(shifts)

					expected[year][weekNumber] = WorkHoursPerDay * 2

					So(summ, ShouldResemble, expected)
				})
			})
		})

		Convey("A shift with 5 weeks, starting on 1st Jan. ", func() {
			st := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
			et := st.Add(time.Hour * HoursPerWeek * 5)

			shifts := append(shifts, Shift{StartTime: st, EndTime: et})

			summ := SummarizeShifts(shifts)

			t1 := time.Date(2015, time.December, 31, 0, 0, 0, 0, time.UTC)
			t2 := time.Date(2016, time.December, 31, 0, 0, 0, 0, time.UTC)

			expected := map[int][]float64{
				2015: make([]float64, ISOWeeksCount(t1)),
				2016: make([]float64, ISOWeeksCount(t2)),
			}

			expected[2015][ISOWeeksCount(t1)-1] = 72
			expected[2016][0] = HoursPerWeek
			expected[2016][1] = HoursPerWeek
			expected[2016][2] = HoursPerWeek
			expected[2016][3] = HoursPerWeek
			expected[2016][4] = HoursPerWeek - expected[2015][ISOWeeksCount(t1)-1]

			So(summ, ShouldResemble, expected)

		})

	})
}
