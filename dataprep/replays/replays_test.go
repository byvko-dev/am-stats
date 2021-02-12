package dataprep

import (
	"log"
	"testing"
)

func TestProcessReplay(t *testing.T) {
	replay := "https://cdn.discordapp.com/attachments/719875141047418962/806183292293873674/20210202_1017___BuCKaC__IS-4_4279240738041098.wotbreplay" // NA
	// replay := "https://cdn.discordapp.com/attachments/719831153162321981/808352798285758464/20210208_1115__i_Yuki_S16_Kranvagn_2509572478643416.wotbreplay" // ASIA
	export, err := ProcessReplay(replay)
	if err != nil {
		log.Print(err)
		t.FailNow()
		return
	}
	log.Printf("%+v", export)
}
