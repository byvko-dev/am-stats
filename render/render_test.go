package render_test

import (
	"os"
	"log"
	"testing"
	 "image/png" 
	"github.com/cufee/am-stats/stats"
	"github.com/cufee/am-stats/render"
)

func TestRealRender(t *testing.T){
	data, err := stats.ExportSessionAsStruct(1025273213, "NA", 0)
	log.Println(err)
	result, err := render.ImageFromStats(data, "-battles", 10)
	log.Println(err)
	
	f, err := os.Create("../am-stats/render/rendercache/out.jpg")
	if err != nil {
		// Handle error
	}
	defer f.Close()
	err = png.Encode(f, result)
	if err != nil {
		// Handle error
	}	
}