package externalapis

// dataToPIDres - JSON response from WG API
type vehiclesAchievmentsRes struct {
	Data   map[string]AchievementsFrame `json:"data"`
	Status string                       `json:"status"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}

// AchievementsFrame -
type AchievementsFrame struct {
	Achievements struct {
		ArmorPiercer                int `json:"armorpiercer,omitempty" bson:"armorpiercer,omitempty"`
		MedalFadin                  int `json:"medalfadin,omitempty" bson:"medalfadin,omitempty"`
		MedalCarius                 int `json:"medalcarius,omitempty" bson:"medalcarius,omitempty"`
		MedalEkins                  int `json:"medalekins,omitempty" bson:"medalekins,omitempty"`
		CollectorGuP                int `json:"collectorgup,omitempty" bson:"collectorgup,omitempty"`
		MedalHalonen                int `json:"medalhalonen,omitempty" bson:"medalhalonen,omitempty"`
		HeroesOfRassenay            int `json:"heroesofrassenay,omitempty" bson:"heroesofrassenay,omitempty"`
		FirstVictory                int `json:"firstvictory,omitempty" bson:"firstvictory,omitempty"`
		Defender                    int `json:"defender,omitempty" bson:"defender,omitempty"`
		Creative                    int `json:"creative,omitempty" bson:"creative,omitempty"`
		ESportFinal                 int `json:"esportfinal,omitempty" bson:"esportfinal,omitempty"`
		Supporter                   int `json:"supporter,omitempty" bson:"supporter,omitempty"`
		GoldClanRibbonSEA           int `json:"goldclanribbonsea,omitempty" bson:"goldclanribbonsea,omitempty"`
		PlatinumTwisterMedalSEA     int `json:"platinumtwistermedalsea,omitempty" bson:"platinumtwistermedalsea,omitempty"`
		MedalLehvaslaiho            int `json:"medallehvaslaiho,omitempty" bson:"medallehvaslaiho,omitempty"`
		TankExpert                  int `json:"tankexpert,omitempty" bson:"tankexpert,omitempty"`
		ESportQualification         int `json:"esportqualification,omitempty" bson:"esportqualification,omitempty"`
		MarkI                       int `json:"marki,omitempty" bson:"marki,omitempty"`
		MedalSupremacy              int `json:"medalsupremacy,omitempty" bson:"medalsupremacy,omitempty"`
		ParticipantofWGFest2017     int `json:"participantofwgfest2017,omitempty" bson:"participantofwgfest2017,omitempty"`
		MedalTournamentOffseason1   int `json:"medaltournamentoffseason1,omitempty" bson:"medaltournamentoffseason1,omitempty"`
		JointVictory                int `json:"jointvictory,omitempty" bson:"jointvictory,omitempty"`
		MedalTournamentOffseason2   int `json:"medaltournamentoffseason2,omitempty" bson:"medaltournamentoffseason2,omitempty"`
		MedalTournamentOffseason4   int `json:"medaltournamentoffseason4,omitempty" bson:"medaltournamentoffseason4,omitempty"`
		Sniper                      int `json:"sniper,omitempty" bson:"sniper,omitempty"`
		TitleSniper                 int `json:"titlesniper,omitempty" bson:"titlesniper,omitempty"`
		MedalCrucialContribution    int `json:"medalcrucialcontribution,omitempty" bson:"medalcrucialcontribution,omitempty"`
		Scout                       int `json:"scout,omitempty" bson:"scout,omitempty"`
		GoldTwisterMedalRU          int `json:"goldtwistermedalru,omitempty" bson:"goldtwistermedalru,omitempty"`
		TankExpert3                 int `json:"tankexpert3,omitempty" bson:"tankexpert3,omitempty"`
		TankExpert2                 int `json:"tankexpert2,omitempty" bson:"tankexpert2,omitempty"`
		TankExpert1                 int `json:"tankexpert1,omitempty" bson:"tankexpert1,omitempty"`
		TankExpert0                 int `json:"tankexpert0,omitempty" bson:"tankexpert0,omitempty"`
		MarkOfMastery               int `json:"markofmastery,omitempty" bson:"markofmastery,omitempty"`
		TankExpert6                 int `json:"tankexpert6,omitempty" bson:"tankexpert6,omitempty"`
		TankExpert5                 int `json:"tankexpert5,omitempty" bson:"tankexpert5,omitempty"`
		TankExpert4                 int `json:"tankexpert4,omitempty" bson:"tankexpert4,omitempty"`
		GoldTwisterMedalEU          int `json:"goldtwistermedaleu,omitempty" bson:"goldtwistermedaleu,omitempty"`
		ChristmasTreeLevelUpNY2019  int `json:"christmastreelevelupny2019,omitempty" bson:"christmastreelevelupny2019,omitempty"`
		MedalLavrinenko             int `json:"medallavrinenko,omitempty" bson:"medallavrinenko,omitempty"`
		MedalKolobanov              int `json:"medalkolobanov,omitempty" bson:"medalkolobanov,omitempty"`
		MedalLafayettePool          int `json:"medallafayettepool,omitempty" bson:"medallafayettepool,omitempty"`
		GoldClanRibbonEU            int `json:"goldclanribboneu,omitempty" bson:"goldclanribboneu,omitempty"`
		OlimpicGolden               int `json:"olimpicgolden,omitempty" bson:"olimpicgolden,omitempty"`
		MedalKnispel                int `json:"medalknispel,omitempty" bson:"medalknispel,omitempty"`
		Invader                     int `json:"invader,omitempty" bson:"invader,omitempty"`
		GoldTwisterMedalNA          int `json:"goldtwistermedalna,omitempty" bson:"goldtwistermedalna,omitempty"`
		MechanicEngineer            int `json:"mechanicengineer,omitempty" bson:"mechanicengineer,omitempty"`
		MarkOfMasteryII             int `json:"markofmasteryii,omitempty" bson:"markofmasteryii,omitempty"`
		FirstBlood                  int `json:"firstblood,omitempty" bson:"firstblood,omitempty"`
		MedalKay                    int `json:"medalkay,omitempty" bson:"medalkay,omitempty"`
		MedalOrlik                  int `json:"medalorlik,omitempty" bson:"medalorlik,omitempty"`
		MedalBrothersInArms         int `json:"medalbrothersinarms,omitempty" bson:"medalbrothersinarms,omitempty"`
		MedalAbrams                 int `json:"medalabrams,omitempty" bson:"medalabrams,omitempty"`
		MedalAtgm                   int `json:"medalatgm,omitempty" bson:"medalatgm,omitempty"`
		MainGun                     int `json:"maingun,omitempty" bson:"maingun,omitempty"`
		IronMan                     int `json:"ironman,omitempty" bson:"ironman,omitempty"`
		PlatinumClanRibbonEU        int `json:"platinumclanribboneu,omitempty" bson:"platinumclanribboneu,omitempty"`
		PlatinumClanRibbonSEA       int `json:"platinumclanribbonsea,omitempty" bson:"platinumclanribbonsea,omitempty"`
		Warrior                     int `json:"warrior,omitempty" bson:"warrior,omitempty"`
		GoldClanRibbonRU            int `json:"goldclanribbonru,omitempty" bson:"goldclanribbonru,omitempty"`
		MedalRadleyWalters          int `json:"medalradleywalters,omitempty" bson:"medalradleywalters,omitempty"`
		Raider                      int `json:"raider,omitempty" bson:"raider,omitempty"`
		ParticipantofNewStart       int `json:"participantofnewstart,omitempty" bson:"participantofnewstart,omitempty"`
		DiamondClanRibbon           int `json:"diamondclanribbon,omitempty" bson:"diamondclanribbon,omitempty"`
		MedalBillotte               int `json:"medalbillotte,omitempty" bson:"medalbillotte,omitempty"`
		PlatinumTwisterMedalEU      int `json:"platinumtwistermedaleu,omitempty" bson:"platinumtwistermedaleu,omitempty"`
		Diehard                     int `json:"diehard,omitempty" bson:"diehard,omitempty"`
		MasterofContinents          int `json:"masterofcontinents,omitempty" bson:"masterofcontinents,omitempty"`
		Evileye                     int `json:"evileye,omitempty" bson:"evileye,omitempty"`
		Cadet                       int `json:"cadet,omitempty" bson:"cadet,omitempty"`
		SupremacyHunter             int `json:"supremacyhunter,omitempty" bson:"supremacyhunter,omitempty"`
		ContinentalContender        int `json:"continentalcontender,omitempty" bson:"continentalcontender,omitempty"`
		Steelwall                   int `json:"steelwall,omitempty" bson:"steelwall,omitempty"`
		SupremacyLegend             int `json:"supremacylegend,omitempty" bson:"supremacylegend,omitempty"`
		Punisher                    int `json:"punisher,omitempty" bson:"punisher,omitempty"`
		ESport                      int `json:"esport,omitempty" bson:"esport,omitempty"`
		PlatinumTwisterMark         int `json:"platinumtwistermark,omitempty" bson:"platinumtwistermark,omitempty"`
		GoldClanRibbonNA            int `json:"goldclanribbonna,omitempty" bson:"goldclanribbonna,omitempty"`
		MedalPoppel                 int `json:"medalpoppel,omitempty" bson:"medalpoppel,omitempty"`
		MechanicEngineer6           int `json:"mechanicengineer6,omitempty" bson:"mechanicengineer6,omitempty"`
		MechanicEngineer4           int `json:"mechanicengineer4,omitempty" bson:"mechanicengineer4,omitempty"`
		GoldTwisterMedalSEA         int `json:"goldtwistermedalsea,omitempty" bson:"goldtwistermedalsea,omitempty"`
		MechanicEngineer2           int `json:"mechanicengineer2,omitempty" bson:"mechanicengineer2,omitempty"`
		MechanicEngineer3           int `json:"mechanicengineer3,omitempty" bson:"mechanicengineer3,omitempty"`
		MechanicEngineer0           int `json:"mechanicengineer0,omitempty" bson:"mechanicengineer0,omitempty"`
		MechanicEngineer1           int `json:"mechanicengineer1,omitempty" bson:"mechanicengineer1,omitempty"`
		MechanicEngineer5           int `json:"mechanicengineer5,omitempty" bson:"mechanicengineer5,omitempty"`
		MedalTarczay                int `json:"medaltarczay,omitempty" bson:"medaltarczay,omitempty"`
		Sinai                       int `json:"sinai,omitempty" bson:"sinai,omitempty"`
		PattonValley                int `json:"pattonvalley,omitempty" bson:"pattonvalley,omitempty"`
		MedalDeLanglade             int `json:"medaldelanglade,omitempty" bson:"medaldelanglade,omitempty"`
		DiamondTwisterMedal         int `json:"diamondtwistermedal,omitempty" bson:"diamondtwistermedal,omitempty"`
		Beasthunter                 int `json:"beasthunter,omitempty" bson:"beasthunter,omitempty"`
		SupremacyVeteran            int `json:"supremacyveteran,omitempty" bson:"supremacyveteran,omitempty"`
		Kamikaze                    int `json:"kamikaze,omitempty" bson:"kamikaze,omitempty"`
		OlimpicBronze               int `json:"olimpicbronze,omitempty" bson:"olimpicbronze,omitempty"`
		MedalTournamentOffseason3   int `json:"medaltournamentoffseason3,omitempty" bson:"medaltournamentoffseason3,omitempty"`
		PlatinumClanRibbonRU        int `json:"platinumclanribbonru,omitempty" bson:"platinumclanribbonru,omitempty"`
		MedalOskin                  int `json:"medaloskin,omitempty" bson:"medaloskin,omitempty"`
		Invincible                  int `json:"invincible,omitempty" bson:"invincible,omitempty"`
		PlatinumClanRibbonNA        int `json:"platinumclanribbonna,omitempty" bson:"platinumclanribbonna,omitempty"`
		PlatinumTwisterMedalRU      int `json:"platinumtwistermedalru,omitempty" bson:"platinumtwistermedalru,omitempty"`
		ContinentalViceChampion     int `json:"continentalvicechampion,omitempty" bson:"continentalvicechampion,omitempty"`
		OlimpicSilver               int `json:"olimpicsilver,omitempty" bson:"olimpicsilver,omitempty"`
		MarkOfMasteryI              int `json:"markofmasteryi,omitempty" bson:"markofmasteryi,omitempty"`
		ContinentalCompetitor       int `json:"continentalcompetitor,omitempty" bson:"continentalcompetitor,omitempty"`
		MedalTournamentSummerSeason int `json:"medaltournamentsummerseason,omitempty" bson:"medaltournamentsummerseason,omitempty"`
		Mousebane                   int `json:"mousebane,omitempty" bson:"mousebane,omitempty"`
		MedalBrunoPietro            int `json:"medalbrunopietro,omitempty" bson:"medalbrunopietro,omitempty"`
		MedalTournamentSpringSeason int `json:"medaltournamentspringseason,omitempty" bson:"medaltournamentspringseason,omitempty"`
		GoldTwisterMark             int `json:"goldtwistermark,omitempty" bson:"goldtwistermark,omitempty"`
		CollectorWarhammer          int `json:"collectorwarhammer,omitempty" bson:"collectorwarhammer,omitempty"`
		MarkOfMasteryIII            int `json:"markofmasteryiii,omitempty" bson:"markofmasteryiii,omitempty"`
		MedalLeClerc                int `json:"medalleclerc,omitempty" bson:"medalleclerc,omitempty"`
		MedalTournamentProfessional int `json:"medaltournamentprofessional,omitempty" bson:"medaltournamentprofessional,omitempty"`
		MedalCommunityChampion      int `json:"medalcommunitychampion,omitempty" bson:"medalcommunitychampion,omitempty"`
		DiamondTwisterMark          int `json:"diamondtwistermark,omitempty" bson:"diamondtwistermark,omitempty"`
		PlatinumTwisterMedalNA      int `json:"platinumtwistermedalna,omitempty" bson:"platinumtwistermedalna,omitempty"`
		HandOfDeath                 int `json:"handofdeath,omitempty" bson:"handofdeath,omitempty"`
		Huntsman                    int `json:"huntsman,omitempty" bson:"huntsman,omitempty"`
		Camper                      int `json:"camper,omitempty" bson:"camper,omitempty"`
		MedalNikolas                int `json:"medalnikolas,omitempty" bson:"medalnikolas,omitempty"`
		AndroidTest                 int `json:"androidtest,omitempty" bson:"androidtest,omitempty"`
		Sturdy                      int `json:"sturdy,omitempty" bson:"sturdy,omitempty"`
		MedalTwitch                 int `json:"medaltwitch,omitempty" bson:"medaltwitch,omitempty"`
		MedalWGfestTicket           int `json:"medalwgfestticket,omitempty" bson:"medalwgfestticket,omitempty"`
		ChampionofNewStart          int `json:"championofnewstart,omitempty" bson:"championofnewstart,omitempty"`
	} `json:"achievements"`
}

// Add - Add two Achievement structs together
func (a AchievementsFrame) Add(more AchievementsFrame) (result AchievementsFrame) {
	// Achievements
	result.Achievements.ArmorPiercer = a.Achievements.ArmorPiercer + more.Achievements.ArmorPiercer
	result.Achievements.MedalFadin = a.Achievements.MedalFadin + more.Achievements.MedalFadin
	result.Achievements.MedalCarius = a.Achievements.MedalCarius + more.Achievements.MedalCarius
	result.Achievements.MedalEkins = a.Achievements.MedalEkins + more.Achievements.MedalEkins
	result.Achievements.CollectorGuP = a.Achievements.CollectorGuP + more.Achievements.CollectorGuP
	result.Achievements.MedalHalonen = a.Achievements.MedalHalonen + more.Achievements.MedalHalonen
	result.Achievements.HeroesOfRassenay = a.Achievements.HeroesOfRassenay + more.Achievements.HeroesOfRassenay
	result.Achievements.FirstVictory = a.Achievements.FirstVictory + more.Achievements.FirstVictory
	result.Achievements.Defender = a.Achievements.Defender + more.Achievements.Defender
	result.Achievements.Creative = a.Achievements.Creative + more.Achievements.Creative
	result.Achievements.ESportFinal = a.Achievements.ESportFinal + more.Achievements.ESportFinal
	result.Achievements.Supporter = a.Achievements.Supporter + more.Achievements.Supporter
	result.Achievements.GoldClanRibbonSEA = a.Achievements.GoldClanRibbonSEA + more.Achievements.GoldClanRibbonSEA
	result.Achievements.PlatinumTwisterMedalSEA = a.Achievements.PlatinumTwisterMedalSEA + more.Achievements.PlatinumTwisterMedalSEA
	result.Achievements.MedalLehvaslaiho = a.Achievements.MedalLehvaslaiho + more.Achievements.MedalLehvaslaiho
	result.Achievements.TankExpert = a.Achievements.TankExpert + more.Achievements.TankExpert
	result.Achievements.ESportQualification = a.Achievements.ESportQualification + more.Achievements.ESportQualification
	result.Achievements.MarkI = a.Achievements.MarkI + more.Achievements.MarkI
	result.Achievements.MedalSupremacy = a.Achievements.MedalSupremacy + more.Achievements.MedalSupremacy
	result.Achievements.ParticipantofWGFest2017 = a.Achievements.ParticipantofWGFest2017 + more.Achievements.ParticipantofWGFest2017
	result.Achievements.MedalTournamentOffseason1 = a.Achievements.MedalTournamentOffseason1 + more.Achievements.MedalTournamentOffseason1
	result.Achievements.JointVictory = a.Achievements.JointVictory + more.Achievements.JointVictory
	result.Achievements.MedalTournamentOffseason2 = a.Achievements.MedalTournamentOffseason2 + more.Achievements.MedalTournamentOffseason2
	result.Achievements.MedalTournamentOffseason4 = a.Achievements.MedalTournamentOffseason4 + more.Achievements.MedalTournamentOffseason4
	result.Achievements.Sniper = a.Achievements.Sniper + more.Achievements.Sniper
	result.Achievements.TitleSniper = a.Achievements.TitleSniper + more.Achievements.TitleSniper
	result.Achievements.MedalCrucialContribution = a.Achievements.MedalCrucialContribution + more.Achievements.MedalCrucialContribution
	result.Achievements.Scout = a.Achievements.Scout + more.Achievements.Scout
	result.Achievements.GoldTwisterMedalRU = a.Achievements.GoldTwisterMedalRU + more.Achievements.GoldTwisterMedalRU
	result.Achievements.TankExpert3 = a.Achievements.TankExpert3 + more.Achievements.TankExpert3
	result.Achievements.TankExpert2 = a.Achievements.TankExpert2 + more.Achievements.TankExpert2
	result.Achievements.TankExpert1 = a.Achievements.TankExpert1 + more.Achievements.TankExpert1
	result.Achievements.TankExpert0 = a.Achievements.TankExpert0 + more.Achievements.TankExpert0
	result.Achievements.MarkOfMastery = a.Achievements.MarkOfMastery + more.Achievements.MarkOfMastery
	result.Achievements.TankExpert6 = a.Achievements.TankExpert6 + more.Achievements.TankExpert6
	result.Achievements.TankExpert5 = a.Achievements.TankExpert5 + more.Achievements.TankExpert5
	result.Achievements.TankExpert4 = a.Achievements.TankExpert4 + more.Achievements.TankExpert4
	result.Achievements.GoldTwisterMedalEU = a.Achievements.GoldTwisterMedalEU + more.Achievements.GoldTwisterMedalEU
	result.Achievements.ChristmasTreeLevelUpNY2019 = a.Achievements.ChristmasTreeLevelUpNY2019 + more.Achievements.ChristmasTreeLevelUpNY2019
	result.Achievements.MedalLavrinenko = a.Achievements.MedalLavrinenko + more.Achievements.MedalLavrinenko
	result.Achievements.MedalKolobanov = a.Achievements.MedalKolobanov + more.Achievements.MedalKolobanov
	result.Achievements.MedalLafayettePool = a.Achievements.MedalLafayettePool + more.Achievements.MedalLafayettePool
	result.Achievements.GoldClanRibbonEU = a.Achievements.GoldClanRibbonEU + more.Achievements.GoldClanRibbonEU
	result.Achievements.OlimpicGolden = a.Achievements.OlimpicGolden + more.Achievements.OlimpicGolden
	result.Achievements.MedalKnispel = a.Achievements.MedalKnispel + more.Achievements.MedalKnispel
	result.Achievements.Invader = a.Achievements.Invader + more.Achievements.Invader
	result.Achievements.GoldTwisterMedalNA = a.Achievements.GoldTwisterMedalNA + more.Achievements.GoldTwisterMedalNA
	result.Achievements.MechanicEngineer = a.Achievements.MechanicEngineer + more.Achievements.MechanicEngineer
	result.Achievements.MarkOfMasteryII = a.Achievements.MarkOfMasteryII + more.Achievements.MarkOfMasteryII
	result.Achievements.FirstBlood = a.Achievements.FirstBlood + more.Achievements.FirstBlood
	result.Achievements.MedalKay = a.Achievements.MedalKay + more.Achievements.MedalKay
	result.Achievements.MedalOrlik = a.Achievements.MedalOrlik + more.Achievements.MedalOrlik
	result.Achievements.MedalBrothersInArms = a.Achievements.MedalBrothersInArms + more.Achievements.MedalBrothersInArms
	result.Achievements.MedalAbrams = a.Achievements.MedalAbrams + more.Achievements.MedalAbrams
	result.Achievements.MedalAtgm = a.Achievements.MedalAtgm + more.Achievements.MedalAtgm
	result.Achievements.MainGun = a.Achievements.MainGun + more.Achievements.MainGun
	result.Achievements.IronMan = a.Achievements.IronMan + more.Achievements.IronMan
	result.Achievements.PlatinumClanRibbonEU = a.Achievements.PlatinumClanRibbonEU + more.Achievements.PlatinumClanRibbonEU
	result.Achievements.PlatinumClanRibbonSEA = a.Achievements.PlatinumClanRibbonSEA + more.Achievements.PlatinumClanRibbonSEA
	result.Achievements.Warrior = a.Achievements.Warrior + more.Achievements.Warrior
	result.Achievements.GoldClanRibbonRU = a.Achievements.GoldClanRibbonRU + more.Achievements.GoldClanRibbonRU
	result.Achievements.MedalRadleyWalters = a.Achievements.MedalRadleyWalters + more.Achievements.MedalRadleyWalters
	result.Achievements.Raider = a.Achievements.Raider + more.Achievements.Raider
	result.Achievements.ParticipantofNewStart = a.Achievements.ParticipantofNewStart + more.Achievements.ParticipantofNewStart
	result.Achievements.DiamondClanRibbon = a.Achievements.DiamondClanRibbon + more.Achievements.DiamondClanRibbon
	result.Achievements.MedalBillotte = a.Achievements.MedalBillotte + more.Achievements.MedalBillotte
	result.Achievements.PlatinumTwisterMedalEU = a.Achievements.PlatinumTwisterMedalEU + more.Achievements.PlatinumTwisterMedalEU
	result.Achievements.Diehard = a.Achievements.Diehard + more.Achievements.Diehard
	result.Achievements.MasterofContinents = a.Achievements.MasterofContinents + more.Achievements.MasterofContinents
	result.Achievements.Evileye = a.Achievements.Evileye + more.Achievements.Evileye
	result.Achievements.Cadet = a.Achievements.Cadet + more.Achievements.Cadet
	result.Achievements.SupremacyHunter = a.Achievements.SupremacyHunter + more.Achievements.SupremacyHunter
	result.Achievements.ContinentalContender = a.Achievements.ContinentalContender + more.Achievements.ContinentalContender
	result.Achievements.Steelwall = a.Achievements.Steelwall + more.Achievements.Steelwall
	result.Achievements.SupremacyLegend = a.Achievements.SupremacyLegend + more.Achievements.SupremacyLegend
	result.Achievements.Punisher = a.Achievements.Punisher + more.Achievements.Punisher
	result.Achievements.ESport = a.Achievements.ESport + more.Achievements.ESport
	result.Achievements.PlatinumTwisterMark = a.Achievements.PlatinumTwisterMark + more.Achievements.PlatinumTwisterMark
	result.Achievements.GoldClanRibbonNA = a.Achievements.GoldClanRibbonNA + more.Achievements.GoldClanRibbonNA
	result.Achievements.MedalPoppel = a.Achievements.MedalPoppel + more.Achievements.MedalPoppel
	result.Achievements.MechanicEngineer6 = a.Achievements.MechanicEngineer6 + more.Achievements.MechanicEngineer6
	result.Achievements.MechanicEngineer4 = a.Achievements.MechanicEngineer4 + more.Achievements.MechanicEngineer4
	result.Achievements.GoldTwisterMedalSEA = a.Achievements.GoldTwisterMedalSEA + more.Achievements.GoldTwisterMedalSEA
	result.Achievements.MechanicEngineer2 = a.Achievements.MechanicEngineer2 + more.Achievements.MechanicEngineer2
	result.Achievements.MechanicEngineer3 = a.Achievements.MechanicEngineer3 + more.Achievements.MechanicEngineer3
	result.Achievements.MechanicEngineer0 = a.Achievements.MechanicEngineer0 + more.Achievements.MechanicEngineer0
	result.Achievements.MechanicEngineer1 = a.Achievements.MechanicEngineer1 + more.Achievements.MechanicEngineer1
	result.Achievements.MechanicEngineer5 = a.Achievements.MechanicEngineer5 + more.Achievements.MechanicEngineer5
	result.Achievements.MedalTarczay = a.Achievements.MedalTarczay + more.Achievements.MedalTarczay
	result.Achievements.Sinai = a.Achievements.Sinai + more.Achievements.Sinai
	result.Achievements.PattonValley = a.Achievements.PattonValley + more.Achievements.PattonValley
	result.Achievements.MedalDeLanglade = a.Achievements.MedalDeLanglade + more.Achievements.MedalDeLanglade
	result.Achievements.DiamondTwisterMedal = a.Achievements.DiamondTwisterMedal + more.Achievements.DiamondTwisterMedal
	result.Achievements.Beasthunter = a.Achievements.Beasthunter + more.Achievements.Beasthunter
	result.Achievements.SupremacyVeteran = a.Achievements.SupremacyVeteran + more.Achievements.SupremacyVeteran
	result.Achievements.Kamikaze = a.Achievements.Kamikaze + more.Achievements.Kamikaze
	result.Achievements.OlimpicBronze = a.Achievements.OlimpicBronze + more.Achievements.OlimpicBronze
	result.Achievements.MedalTournamentOffseason3 = a.Achievements.MedalTournamentOffseason3 + more.Achievements.MedalTournamentOffseason3
	result.Achievements.PlatinumClanRibbonRU = a.Achievements.PlatinumClanRibbonRU + more.Achievements.PlatinumClanRibbonRU
	result.Achievements.MedalOskin = a.Achievements.MedalOskin + more.Achievements.MedalOskin
	result.Achievements.Invincible = a.Achievements.Invincible + more.Achievements.Invincible
	result.Achievements.PlatinumClanRibbonNA = a.Achievements.PlatinumClanRibbonNA + more.Achievements.PlatinumClanRibbonNA
	result.Achievements.PlatinumTwisterMedalRU = a.Achievements.PlatinumTwisterMedalRU + more.Achievements.PlatinumTwisterMedalRU
	result.Achievements.ContinentalViceChampion = a.Achievements.ContinentalViceChampion + more.Achievements.ContinentalViceChampion
	result.Achievements.OlimpicSilver = a.Achievements.OlimpicSilver + more.Achievements.OlimpicSilver
	result.Achievements.MarkOfMasteryI = a.Achievements.MarkOfMasteryI + more.Achievements.MarkOfMasteryI
	result.Achievements.ContinentalCompetitor = a.Achievements.ContinentalCompetitor + more.Achievements.ContinentalCompetitor
	result.Achievements.MedalTournamentSummerSeason = a.Achievements.MedalTournamentSummerSeason + more.Achievements.MedalTournamentSummerSeason
	result.Achievements.Mousebane = a.Achievements.Mousebane + more.Achievements.Mousebane
	result.Achievements.MedalBrunoPietro = a.Achievements.MedalBrunoPietro + more.Achievements.MedalBrunoPietro
	result.Achievements.MedalTournamentSpringSeason = a.Achievements.MedalTournamentSpringSeason + more.Achievements.MedalTournamentSpringSeason
	result.Achievements.GoldTwisterMark = a.Achievements.GoldTwisterMark + more.Achievements.GoldTwisterMark
	result.Achievements.CollectorWarhammer = a.Achievements.CollectorWarhammer + more.Achievements.CollectorWarhammer
	result.Achievements.MarkOfMasteryIII = a.Achievements.MarkOfMasteryIII + more.Achievements.MarkOfMasteryIII
	result.Achievements.MedalLeClerc = a.Achievements.MedalLeClerc + more.Achievements.MedalLeClerc
	result.Achievements.MedalTournamentProfessional = a.Achievements.MedalTournamentProfessional + more.Achievements.MedalTournamentProfessional
	result.Achievements.MedalCommunityChampion = a.Achievements.MedalCommunityChampion + more.Achievements.MedalCommunityChampion
	result.Achievements.DiamondTwisterMark = a.Achievements.DiamondTwisterMark + more.Achievements.DiamondTwisterMark
	result.Achievements.PlatinumTwisterMedalNA = a.Achievements.PlatinumTwisterMedalNA + more.Achievements.PlatinumTwisterMedalNA
	result.Achievements.HandOfDeath = a.Achievements.HandOfDeath + more.Achievements.HandOfDeath
	result.Achievements.Huntsman = a.Achievements.Huntsman + more.Achievements.Huntsman
	result.Achievements.Camper = a.Achievements.Camper + more.Achievements.Camper
	result.Achievements.MedalNikolas = a.Achievements.MedalNikolas + more.Achievements.MedalNikolas
	result.Achievements.AndroidTest = a.Achievements.AndroidTest + more.Achievements.AndroidTest
	result.Achievements.Sturdy = a.Achievements.Sturdy + more.Achievements.Sturdy
	result.Achievements.MedalTwitch = a.Achievements.MedalTwitch + more.Achievements.MedalTwitch
	result.Achievements.MedalWGfestTicket = a.Achievements.MedalWGfestTicket + more.Achievements.MedalWGfestTicket
	result.Achievements.ChampionofNewStart = a.Achievements.ChampionofNewStart + more.Achievements.ChampionofNewStart

	return result
}

// Diff - Find changes in achievements
func (new AchievementsFrame) Diff(old AchievementsFrame) (result AchievementsFrame) {
	// Achievements
	result.Achievements.ArmorPiercer = new.Achievements.ArmorPiercer - old.Achievements.ArmorPiercer
	result.Achievements.MedalFadin = new.Achievements.MedalFadin - old.Achievements.MedalFadin
	result.Achievements.MedalCarius = new.Achievements.MedalCarius - old.Achievements.MedalCarius
	result.Achievements.MedalEkins = new.Achievements.MedalEkins - old.Achievements.MedalEkins
	result.Achievements.CollectorGuP = new.Achievements.CollectorGuP - old.Achievements.CollectorGuP
	result.Achievements.MedalHalonen = new.Achievements.MedalHalonen - old.Achievements.MedalHalonen
	result.Achievements.HeroesOfRassenay = new.Achievements.HeroesOfRassenay - old.Achievements.HeroesOfRassenay
	result.Achievements.FirstVictory = new.Achievements.FirstVictory - old.Achievements.FirstVictory
	result.Achievements.Defender = new.Achievements.Defender - old.Achievements.Defender
	result.Achievements.Creative = new.Achievements.Creative - old.Achievements.Creative
	result.Achievements.ESportFinal = new.Achievements.ESportFinal - old.Achievements.ESportFinal
	result.Achievements.Supporter = new.Achievements.Supporter - old.Achievements.Supporter
	result.Achievements.GoldClanRibbonSEA = new.Achievements.GoldClanRibbonSEA - old.Achievements.GoldClanRibbonSEA
	result.Achievements.PlatinumTwisterMedalSEA = new.Achievements.PlatinumTwisterMedalSEA - old.Achievements.PlatinumTwisterMedalSEA
	result.Achievements.MedalLehvaslaiho = new.Achievements.MedalLehvaslaiho - old.Achievements.MedalLehvaslaiho
	result.Achievements.TankExpert = new.Achievements.TankExpert - old.Achievements.TankExpert
	result.Achievements.ESportQualification = new.Achievements.ESportQualification - old.Achievements.ESportQualification
	result.Achievements.MarkI = new.Achievements.MarkI - old.Achievements.MarkI
	result.Achievements.MedalSupremacy = new.Achievements.MedalSupremacy - old.Achievements.MedalSupremacy
	result.Achievements.ParticipantofWGFest2017 = new.Achievements.ParticipantofWGFest2017 - old.Achievements.ParticipantofWGFest2017
	result.Achievements.MedalTournamentOffseason1 = new.Achievements.MedalTournamentOffseason1 - old.Achievements.MedalTournamentOffseason1
	result.Achievements.JointVictory = new.Achievements.JointVictory - old.Achievements.JointVictory
	result.Achievements.MedalTournamentOffseason2 = new.Achievements.MedalTournamentOffseason2 - old.Achievements.MedalTournamentOffseason2
	result.Achievements.MedalTournamentOffseason4 = new.Achievements.MedalTournamentOffseason4 - old.Achievements.MedalTournamentOffseason4
	result.Achievements.Sniper = new.Achievements.Sniper - old.Achievements.Sniper
	result.Achievements.TitleSniper = new.Achievements.TitleSniper - old.Achievements.TitleSniper
	result.Achievements.MedalCrucialContribution = new.Achievements.MedalCrucialContribution - old.Achievements.MedalCrucialContribution
	result.Achievements.Scout = new.Achievements.Scout - old.Achievements.Scout
	result.Achievements.GoldTwisterMedalRU = new.Achievements.GoldTwisterMedalRU - old.Achievements.GoldTwisterMedalRU
	result.Achievements.TankExpert3 = new.Achievements.TankExpert3 - old.Achievements.TankExpert3
	result.Achievements.TankExpert2 = new.Achievements.TankExpert2 - old.Achievements.TankExpert2
	result.Achievements.TankExpert1 = new.Achievements.TankExpert1 - old.Achievements.TankExpert1
	result.Achievements.TankExpert0 = new.Achievements.TankExpert0 - old.Achievements.TankExpert0
	result.Achievements.MarkOfMastery = new.Achievements.MarkOfMastery - old.Achievements.MarkOfMastery
	result.Achievements.TankExpert6 = new.Achievements.TankExpert6 - old.Achievements.TankExpert6
	result.Achievements.TankExpert5 = new.Achievements.TankExpert5 - old.Achievements.TankExpert5
	result.Achievements.TankExpert4 = new.Achievements.TankExpert4 - old.Achievements.TankExpert4
	result.Achievements.GoldTwisterMedalEU = new.Achievements.GoldTwisterMedalEU - old.Achievements.GoldTwisterMedalEU
	result.Achievements.ChristmasTreeLevelUpNY2019 = new.Achievements.ChristmasTreeLevelUpNY2019 - old.Achievements.ChristmasTreeLevelUpNY2019
	result.Achievements.MedalLavrinenko = new.Achievements.MedalLavrinenko - old.Achievements.MedalLavrinenko
	result.Achievements.MedalKolobanov = new.Achievements.MedalKolobanov - old.Achievements.MedalKolobanov
	result.Achievements.MedalLafayettePool = new.Achievements.MedalLafayettePool - old.Achievements.MedalLafayettePool
	result.Achievements.GoldClanRibbonEU = new.Achievements.GoldClanRibbonEU - old.Achievements.GoldClanRibbonEU
	result.Achievements.OlimpicGolden = new.Achievements.OlimpicGolden - old.Achievements.OlimpicGolden
	result.Achievements.MedalKnispel = new.Achievements.MedalKnispel - old.Achievements.MedalKnispel
	result.Achievements.Invader = new.Achievements.Invader - old.Achievements.Invader
	result.Achievements.GoldTwisterMedalNA = new.Achievements.GoldTwisterMedalNA - old.Achievements.GoldTwisterMedalNA
	result.Achievements.MechanicEngineer = new.Achievements.MechanicEngineer - old.Achievements.MechanicEngineer
	result.Achievements.MarkOfMasteryII = new.Achievements.MarkOfMasteryII - old.Achievements.MarkOfMasteryII
	result.Achievements.FirstBlood = new.Achievements.FirstBlood - old.Achievements.FirstBlood
	result.Achievements.MedalKay = new.Achievements.MedalKay - old.Achievements.MedalKay
	result.Achievements.MedalOrlik = new.Achievements.MedalOrlik - old.Achievements.MedalOrlik
	result.Achievements.MedalBrothersInArms = new.Achievements.MedalBrothersInArms - old.Achievements.MedalBrothersInArms
	result.Achievements.MedalAbrams = new.Achievements.MedalAbrams - old.Achievements.MedalAbrams
	result.Achievements.MedalAtgm = new.Achievements.MedalAtgm - old.Achievements.MedalAtgm
	result.Achievements.MainGun = new.Achievements.MainGun - old.Achievements.MainGun
	result.Achievements.IronMan = new.Achievements.IronMan - old.Achievements.IronMan
	result.Achievements.PlatinumClanRibbonEU = new.Achievements.PlatinumClanRibbonEU - old.Achievements.PlatinumClanRibbonEU
	result.Achievements.PlatinumClanRibbonSEA = new.Achievements.PlatinumClanRibbonSEA - old.Achievements.PlatinumClanRibbonSEA
	result.Achievements.Warrior = new.Achievements.Warrior - old.Achievements.Warrior
	result.Achievements.GoldClanRibbonRU = new.Achievements.GoldClanRibbonRU - old.Achievements.GoldClanRibbonRU
	result.Achievements.MedalRadleyWalters = new.Achievements.MedalRadleyWalters - old.Achievements.MedalRadleyWalters
	result.Achievements.Raider = new.Achievements.Raider - old.Achievements.Raider
	result.Achievements.ParticipantofNewStart = new.Achievements.ParticipantofNewStart - old.Achievements.ParticipantofNewStart
	result.Achievements.DiamondClanRibbon = new.Achievements.DiamondClanRibbon - old.Achievements.DiamondClanRibbon
	result.Achievements.MedalBillotte = new.Achievements.MedalBillotte - old.Achievements.MedalBillotte
	result.Achievements.PlatinumTwisterMedalEU = new.Achievements.PlatinumTwisterMedalEU - old.Achievements.PlatinumTwisterMedalEU
	result.Achievements.Diehard = new.Achievements.Diehard - old.Achievements.Diehard
	result.Achievements.MasterofContinents = new.Achievements.MasterofContinents - old.Achievements.MasterofContinents
	result.Achievements.Evileye = new.Achievements.Evileye - old.Achievements.Evileye
	result.Achievements.Cadet = new.Achievements.Cadet - old.Achievements.Cadet
	result.Achievements.SupremacyHunter = new.Achievements.SupremacyHunter - old.Achievements.SupremacyHunter
	result.Achievements.ContinentalContender = new.Achievements.ContinentalContender - old.Achievements.ContinentalContender
	result.Achievements.Steelwall = new.Achievements.Steelwall - old.Achievements.Steelwall
	result.Achievements.SupremacyLegend = new.Achievements.SupremacyLegend - old.Achievements.SupremacyLegend
	result.Achievements.Punisher = new.Achievements.Punisher - old.Achievements.Punisher
	result.Achievements.ESport = new.Achievements.ESport - old.Achievements.ESport
	result.Achievements.PlatinumTwisterMark = new.Achievements.PlatinumTwisterMark - old.Achievements.PlatinumTwisterMark
	result.Achievements.GoldClanRibbonNA = new.Achievements.GoldClanRibbonNA - old.Achievements.GoldClanRibbonNA
	result.Achievements.MedalPoppel = new.Achievements.MedalPoppel - old.Achievements.MedalPoppel
	result.Achievements.MechanicEngineer6 = new.Achievements.MechanicEngineer6 - old.Achievements.MechanicEngineer6
	result.Achievements.MechanicEngineer4 = new.Achievements.MechanicEngineer4 - old.Achievements.MechanicEngineer4
	result.Achievements.GoldTwisterMedalSEA = new.Achievements.GoldTwisterMedalSEA - old.Achievements.GoldTwisterMedalSEA
	result.Achievements.MechanicEngineer2 = new.Achievements.MechanicEngineer2 - old.Achievements.MechanicEngineer2
	result.Achievements.MechanicEngineer3 = new.Achievements.MechanicEngineer3 - old.Achievements.MechanicEngineer3
	result.Achievements.MechanicEngineer0 = new.Achievements.MechanicEngineer0 - old.Achievements.MechanicEngineer0
	result.Achievements.MechanicEngineer1 = new.Achievements.MechanicEngineer1 - old.Achievements.MechanicEngineer1
	result.Achievements.MechanicEngineer5 = new.Achievements.MechanicEngineer5 - old.Achievements.MechanicEngineer5
	result.Achievements.MedalTarczay = new.Achievements.MedalTarczay - old.Achievements.MedalTarczay
	result.Achievements.Sinai = new.Achievements.Sinai - old.Achievements.Sinai
	result.Achievements.PattonValley = new.Achievements.PattonValley - old.Achievements.PattonValley
	result.Achievements.MedalDeLanglade = new.Achievements.MedalDeLanglade - old.Achievements.MedalDeLanglade
	result.Achievements.DiamondTwisterMedal = new.Achievements.DiamondTwisterMedal - old.Achievements.DiamondTwisterMedal
	result.Achievements.Beasthunter = new.Achievements.Beasthunter - old.Achievements.Beasthunter
	result.Achievements.SupremacyVeteran = new.Achievements.SupremacyVeteran - old.Achievements.SupremacyVeteran
	result.Achievements.Kamikaze = new.Achievements.Kamikaze - old.Achievements.Kamikaze
	result.Achievements.OlimpicBronze = new.Achievements.OlimpicBronze - old.Achievements.OlimpicBronze
	result.Achievements.MedalTournamentOffseason3 = new.Achievements.MedalTournamentOffseason3 - old.Achievements.MedalTournamentOffseason3
	result.Achievements.PlatinumClanRibbonRU = new.Achievements.PlatinumClanRibbonRU - old.Achievements.PlatinumClanRibbonRU
	result.Achievements.MedalOskin = new.Achievements.MedalOskin - old.Achievements.MedalOskin
	result.Achievements.Invincible = new.Achievements.Invincible - old.Achievements.Invincible
	result.Achievements.PlatinumClanRibbonNA = new.Achievements.PlatinumClanRibbonNA - old.Achievements.PlatinumClanRibbonNA
	result.Achievements.PlatinumTwisterMedalRU = new.Achievements.PlatinumTwisterMedalRU - old.Achievements.PlatinumTwisterMedalRU
	result.Achievements.ContinentalViceChampion = new.Achievements.ContinentalViceChampion - old.Achievements.ContinentalViceChampion
	result.Achievements.OlimpicSilver = new.Achievements.OlimpicSilver - old.Achievements.OlimpicSilver
	result.Achievements.MarkOfMasteryI = new.Achievements.MarkOfMasteryI - old.Achievements.MarkOfMasteryI
	result.Achievements.ContinentalCompetitor = new.Achievements.ContinentalCompetitor - old.Achievements.ContinentalCompetitor
	result.Achievements.MedalTournamentSummerSeason = new.Achievements.MedalTournamentSummerSeason - old.Achievements.MedalTournamentSummerSeason
	result.Achievements.Mousebane = new.Achievements.Mousebane - old.Achievements.Mousebane
	result.Achievements.MedalBrunoPietro = new.Achievements.MedalBrunoPietro - old.Achievements.MedalBrunoPietro
	result.Achievements.MedalTournamentSpringSeason = new.Achievements.MedalTournamentSpringSeason - old.Achievements.MedalTournamentSpringSeason
	result.Achievements.GoldTwisterMark = new.Achievements.GoldTwisterMark - old.Achievements.GoldTwisterMark
	result.Achievements.CollectorWarhammer = new.Achievements.CollectorWarhammer - old.Achievements.CollectorWarhammer
	result.Achievements.MarkOfMasteryIII = new.Achievements.MarkOfMasteryIII - old.Achievements.MarkOfMasteryIII
	result.Achievements.MedalLeClerc = new.Achievements.MedalLeClerc - old.Achievements.MedalLeClerc
	result.Achievements.MedalTournamentProfessional = new.Achievements.MedalTournamentProfessional - old.Achievements.MedalTournamentProfessional
	result.Achievements.MedalCommunityChampion = new.Achievements.MedalCommunityChampion - old.Achievements.MedalCommunityChampion
	result.Achievements.DiamondTwisterMark = new.Achievements.DiamondTwisterMark - old.Achievements.DiamondTwisterMark
	result.Achievements.PlatinumTwisterMedalNA = new.Achievements.PlatinumTwisterMedalNA - old.Achievements.PlatinumTwisterMedalNA
	result.Achievements.HandOfDeath = new.Achievements.HandOfDeath - old.Achievements.HandOfDeath
	result.Achievements.Huntsman = new.Achievements.Huntsman - old.Achievements.Huntsman
	result.Achievements.Camper = new.Achievements.Camper - old.Achievements.Camper
	result.Achievements.MedalNikolas = new.Achievements.MedalNikolas - old.Achievements.MedalNikolas
	result.Achievements.AndroidTest = new.Achievements.AndroidTest - old.Achievements.AndroidTest
	result.Achievements.Sturdy = new.Achievements.Sturdy - old.Achievements.Sturdy
	result.Achievements.MedalTwitch = new.Achievements.MedalTwitch - old.Achievements.MedalTwitch
	result.Achievements.MedalWGfestTicket = new.Achievements.MedalWGfestTicket - old.Achievements.MedalWGfestTicket
	result.Achievements.ChampionofNewStart = new.Achievements.ChampionofNewStart - old.Achievements.ChampionofNewStart
	return result
}
