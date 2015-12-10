package guage

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui"
	DB "github.com/vrecan/FluxDash/influx"
)

type GaugeInfo struct {
	From  string
	Time  string
	Title string
	Where string
}

type Gauge struct {
	I  GaugeInfo
	G  *ui.Gauge
	db *DB.Influx
}

func NewGauge(barColor ui.Attribute, db *DB.Influx, info GaugeInfo) *Gauge {

	g := &Gauge{G: ui.NewGauge(), db: db, I: info}
	g.G.BarColor = barColor
	// g.G.PercentColor = ui.ColorRed
	// g.G.PercentColorHighlighted = ui.ColorMagenta
	g.G.Width = 50
	g.G.Height = 3
	return g
}

func (s *Gauge) Gauges() *ui.Gauge {
	return s.G
}

func (s *Gauge) Update() {
	s.SetData()
	s.SetTitle()
}

func (s *Gauge) SetData() {
	meanTotal := getData(s.db, buildQuery("mean(value)", s.I.From, s.I.Where, s.I.Time, ""))
	s.G.Percent = meanTotal[0]
}

func (s *Gauge) SetTitle() {
	meanTotal := getData(s.db, buildQuery("mean(value)", s.I.From, s.I.Where, s.I.Time, ""))
	maxTotal := getData(s.db, buildQuery("max(value)", s.I.From, s.I.Where, s.I.Time, ""))
	s.G.Label = fmt.Sprintf("%s mean:%v%% max:%v%%", s.I.Title, meanTotal[0], maxTotal[0])
}
func (s *Gauge) GetColumns() []*ui.Row {
	return []*ui.Row{ui.NewCol(12, 0, s.Gauges())}
}
func buildQuery(sel string, from string, where string, time string, groupBy string) string {
	if len(sel) == 0 || len(from) == 0 || len(time) == 0 {
		log.Fatal("invalid query string :", fmt.Sprintf("SELECT %s FROM %s WHERE %s AND time > %s %s", sel, from, where, groupBy))
	}
	if len(where) > 0 {
		return fmt.Sprintf("SELECT %s FROM %s WHERE %s AND time > %s %s", sel, from, where, time, groupBy)
	} else {
		return fmt.Sprintf("SELECT %s FROM %s WHERE time > %s %s", sel, from, time, groupBy)
	}
}
func getData(db *DB.Influx, q string) (data []int) {
	r, err := db.Query(q)
	if nil != err {
		log.Fatal(err)
	}
	if len(r) == 0 || len(r[0].Series) == 0 {
		log.Fatal(q)
	}

	for _, row := range r[0].Series[0].Values {
		_, err := time.Parse(time.RFC3339, row[0].(string))
		if err != nil {
			log.Fatal(err)
		}
		if len(row) > 1 {
			if nil != row[1] {
				val, err := row[1].(json.Number).Float64()
				if nil != err {
					fmt.Println("ERR: ", err)
				}
				data = append(data, int(val))
			}
		}

	}
	return data
}
