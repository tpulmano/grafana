package conditions

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/tsdb"
)

func TestSimpleReducer(t *testing.T) {
	Convey("Test simple reducer by calculating", t, func() {

		Convey("sum", func() {
			result := testReducer("sum", 1, 2, 3)
			So(result, ShouldEqual, float64(6))
		})

		Convey("min", func() {
			result := testReducer("min", 3, 2, 1)
			So(result, ShouldEqual, float64(1))
		})

		Convey("max", func() {
			result := testReducer("max", 1, 2, 3)
			So(result, ShouldEqual, float64(3))
		})

		Convey("count", func() {
			result := testReducer("count", 1, 2, 3000)
			So(result, ShouldEqual, float64(3))
		})

		Convey("last", func() {
			result := testReducer("last", 1, 2, 3000)
			So(result, ShouldEqual, float64(3000))
		})

		Convey("median odd amount of numbers", func() {
			result := testReducer("median", 1, 2, 3000)
			So(result, ShouldEqual, float64(2))
		})

		Convey("median even amount of numbers", func() {
			result := testReducer("median", 1, 2, 4, 3000)
			So(result, ShouldEqual, float64(3))
		})

		Convey("median with one values", func() {
			result := testReducer("median", 1)
			So(result, ShouldEqual, float64(1))
		})

		Convey("median should ignore null values", func() {
			reducer := newSimpleReducer("median")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 3))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(float64(1)), 4))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(float64(2)), 5))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(float64(3)), 6))

			result := reducer.Reduce(series)
			So(result.Valid, ShouldEqual, true)
			So(result.Float64, ShouldEqual, float64(2))
		})

		Convey("avg", func() {
			result := testReducer("avg", 1, 2, 3)
			So(result, ShouldEqual, float64(2))
		})

		Convey("avg with only nulls", func() {
			reducer := newSimpleReducer("avg")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
			So(reducer.Reduce(series).Valid, ShouldEqual, false)
		})

		Convey("count_non_null", func() {
			Convey("with null values and real values", func() {
				reducer := newSimpleReducer("count_non_null")
				series := &tsdb.TimeSeries{
					Name: "test time serie",
				}

				series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
				series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))
				series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(3), 3))
				series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(3), 4))

				So(reducer.Reduce(series).Valid, ShouldEqual, true)
				So(reducer.Reduce(series).Float64, ShouldEqual, 2)
			})

			Convey("with null values", func() {
				reducer := newSimpleReducer("count_non_null")
				series := &tsdb.TimeSeries{
					Name: "test time serie",
				}

				series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
				series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))

				So(reducer.Reduce(series).Valid, ShouldEqual, false)
			})
		})

		Convey("avg of number values and null values should ignore nulls", func() {
			reducer := newSimpleReducer("avg")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(3), 1))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 3))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(3), 4))

			So(reducer.Reduce(series).Float64, ShouldEqual, float64(3))
		})

		Convey("diff one point", func() {
			result := testReducer("diff", 30)
			So(result, ShouldEqual, float64(0))
		})

		Convey("diff two points", func() {
			result := testReducer("diff", 30, 40)
			So(result, ShouldEqual, float64(10))
		})

		Convey("diff three points", func() {
			result := testReducer("diff", 30, 40, 40)
			So(result, ShouldEqual, float64(10))
		})

		Convey("diff with only nulls", func() {
			reducer := newSimpleReducer("diff")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))

			So(reducer.Reduce(series).Valid, ShouldEqual, false)
		})

		Convey("diff with increase", func() {
			result := testReducer("diff", 11, 30)
			So(result, ShouldEqual, float64(19))
		})

		Convey("diff with decrease", func() {
			result := testReducer("diff", 30, 11)
			So(result, ShouldEqual, float64(19))
		})

		Convey("percent_diff one point", func() {
			result := testReducer("percent_diff", 40)
			So(result, ShouldEqual, float64(0))
		})

		Convey("percent_diff two points", func() {
			result := testReducer("percent_diff", 30, 40)
			So(result, ShouldEqual, float64(33.33333333333333))
		})

		Convey("percent_diff three points", func() {
			result := testReducer("percent_diff", 30, 40, 40)
			So(result, ShouldEqual, float64(33.33333333333333))
		})

		Convey("percent_diff with only nulls", func() {
			reducer := newSimpleReducer("percent_diff")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))

			So(reducer.Reduce(series).Valid, ShouldEqual, false)
		})

		Convey("percent_diff with increase", func() {
			result := testReducer("percent_diff", 20, 40)
			So(result, ShouldEqual, float64(100))
		})

		Convey("percent_diff with decrease", func() {
			result := testReducer("percent_diff", 40, 20)
			So(result, ShouldEqual, float64(50))
		})

		Convey("change with one point", func() {
			result := testReducer("change", 30)
			So(result, ShouldEqual, float64(0))
		})

		Convey("change with two points", func() {
			result := testReducer("change", 30, 50)
			So(result, ShouldEqual, float64(20))
		})

		Convey("change with three points", func() {
			result := testReducer("change", 30, 50, 45)
			So(result, ShouldEqual, float64(15))
		})

		Convey("change with only nulls", func() {
			reducer := newSimpleReducer("change")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))

			So(reducer.Reduce(series).Valid, ShouldEqual, false)
		})

		Convey("change with increase", func() {
			result := testReducer("change", 30, 60)
			So(result, ShouldEqual, float64(30))
		})

		Convey("change with decrease", func() {
			result := testReducer("change", 40, 10)
			So(result, ShouldEqual, float64(-30))
		})

		Convey("percent_change with one point", func() {
			result := testReducer("percent_change", 40)
			So(result, ShouldEqual, float64(0))
		})

		Convey("percent_change with two points", func() {
			result := testReducer("percent_change", 40, 46)
			So(result, ShouldEqual, float64(15))
		})

		Convey("percent_change with three points", func() {
			result := testReducer("percent_change", 40, 35, 52)
			So(result, ShouldEqual, float64(30))
		})

		Convey("percent_change with only nulls", func() {
			reducer := newSimpleReducer("percent_change")
			series := &tsdb.TimeSeries{
				Name: "test time serie",
			}

			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 1))
			series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFromPtr(nil), 2))

			So(reducer.Reduce(series).Valid, ShouldEqual, false)
		})

		Convey("percent_change with increase", func() {
			result := testReducer("percent_change", 40, 54)
			So(result, ShouldEqual, float64(35))
		})

		Convey("percent_change with decrease", func() {
			result := testReducer("percent_change", 40, 28)
			So(result, ShouldEqual, float64(-30))
		})

	})
}

func testReducer(reducerType string, datapoints ...float64) float64 {
	reducer := newSimpleReducer(reducerType)
	series := &tsdb.TimeSeries{
		Name: "test time serie",
	}

	for idx := range datapoints {
		series.Points = append(series.Points, tsdb.NewTimePoint(null.FloatFrom(datapoints[idx]), 1234134))
	}

	return reducer.Reduce(series).Float64
}
