package render_test

import (
	"log"
	"testing"
	"github.com/cufee/am-stats/stats"
	"github.com/cufee/am-stats/render"
)

func TestRealRender(t *testing.T){
	data, err := stats.ExportSessionAsStruct(1025273213, "NA", 0)
	log.Println(err)
	result, err := render.ImageFromStats(data)
	log.Println(err)
	log.Printf("%T", result)
}