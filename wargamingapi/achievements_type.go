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

// AchievementsFrame - Frame for achievements on vehicle
type AchievementsFrame struct {
	Achievements struct {
		ArmorPiercer                int `json:"armorPiercer"`
		MedalFadin                  int `json:"medalFadin"`
		MedalCarius                 int `json:"medalCarius"`
		MedalEkins                  int `json:"medalEkins"`
		CollectorGuP                int `json:"collectorGuP"`
		MedalHalonen                int `json:"medalHalonen"`
		HeroesOfRassenay            int `json:"heroesOfRassenay"`
		FirstVictory                int `json:"firstVictory"`
		Defender                    int `json:"defender"`
		Creative                    int `json:"creative"`
		ESportFinal                 int `json:"eSportFinal"`
		Supporter                   int `json:"supporter"`
		GoldClanRibbonSEA           int `json:"goldClanRibbonSEA"`
		PlatinumTwisterMedalSEA     int `json:"platinumTwisterMedalSEA"`
		MedalLehvaslaiho            int `json:"medalLehvaslaiho"`
		TankExpert                  int `json:"tankExpert"`
		ESportQualification         int `json:"eSportQualification"`
		MarkI                       int `json:"MarkI"`
		MedalSupremacy              int `json:"medalSupremacy"`
		ParticipantofWGFest2017     int `json:"participantofWGFest2017"`
		MedalTournamentOffseason1   int `json:"medalTournamentOffseason1"`
		JointVictory                int `json:"jointVictory"`
		MedalTournamentOffseason2   int `json:"medalTournamentOffseason2"`
		MedalTournamentOffseason4   int `json:"medalTournamentOffseason4"`
		Sniper                      int `json:"sniper"`
		TitleSniper                 int `json:"titleSniper"`
		MedalCrucialContribution    int `json:"medalCrucialContribution"`
		Scout                       int `json:"scout"`
		GoldTwisterMedalRU          int `json:"goldTwisterMedalRU"`
		TankExpert3                 int `json:"tankExpert3"`
		TankExpert2                 int `json:"tankExpert2"`
		TankExpert1                 int `json:"tankExpert1"`
		TankExpert0                 int `json:"tankExpert0"`
		MarkOfMastery               int `json:"markOfMastery"`
		TankExpert6                 int `json:"tankExpert6"`
		TankExpert5                 int `json:"tankExpert5"`
		TankExpert4                 int `json:"tankExpert4"`
		GoldTwisterMedalEU          int `json:"goldTwisterMedalEU"`
		ChristmasTreeLevelUpNY2019  int `json:"ChristmasTreeLevelUpNY2019"`
		MedalLavrinenko             int `json:"medalLavrinenko"`
		MedalKolobanov              int `json:"medalKolobanov"`
		MedalLafayettePool          int `json:"medalLafayettePool"`
		GoldClanRibbonEU            int `json:"goldClanRibbonEU"`
		OlimpicGolden               int `json:"olimpicGolden"`
		MedalKnispel                int `json:"medalKnispel"`
		Invader                     int `json:"invader"`
		GoldTwisterMedalNA          int `json:"goldTwisterMedalNA"`
		MechanicEngineer            int `json:"mechanicEngineer"`
		MarkOfMasteryII             int `json:"markOfMasteryII"`
		FirstBlood                  int `json:"firstBlood"`
		MedalKay                    int `json:"medalKay"`
		MedalOrlik                  int `json:"medalOrlik"`
		MedalBrothersInArms         int `json:"medalBrothersInArms"`
		MedalAbrams                 int `json:"medalAbrams"`
		MedalAtgm                   int `json:"medalAtgm"`
		MainGun                     int `json:"mainGun"`
		IronMan                     int `json:"ironMan"`
		PlatinumClanRibbonEU        int `json:"platinumClanRibbonEU"`
		PlatinumClanRibbonSEA       int `json:"platinumClanRibbonSEA"`
		Warrior                     int `json:"warrior"`
		GoldClanRibbonRU            int `json:"goldClanRibbonRU"`
		MedalRadleyWalters          int `json:"medalRadleyWalters"`
		Raider                      int `json:"raider"`
		ParticipantofNewStart       int `json:"participantofNewStart"`
		DiamondClanRibbon           int `json:"diamondClanRibbon"`
		MedalBillotte               int `json:"medalBillotte"`
		PlatinumTwisterMedalEU      int `json:"platinumTwisterMedalEU"`
		Diehard                     int `json:"diehard"`
		MasterofContinents          int `json:"masterofContinents"`
		Evileye                     int `json:"evileye"`
		Cadet                       int `json:"cadet"`
		SupremacyHunter             int `json:"supremacyHunter"`
		ContinentalContender        int `json:"continentalContender"`
		Steelwall                   int `json:"steelwall"`
		SupremacyLegend             int `json:"supremacyLegend"`
		Punisher                    int `json:"punisher"`
		ESport                      int `json:"eSport"`
		PlatinumTwisterMark         int `json:"platinumTwisterMark"`
		GoldClanRibbonNA            int `json:"goldClanRibbonNA"`
		MedalPoppel                 int `json:"medalPoppel"`
		MechanicEngineer6           int `json:"mechanicEngineer6"`
		MechanicEngineer4           int `json:"mechanicEngineer4"`
		GoldTwisterMedalSEA         int `json:"goldTwisterMedalSEA"`
		MechanicEngineer2           int `json:"mechanicEngineer2"`
		MechanicEngineer3           int `json:"mechanicEngineer3"`
		MechanicEngineer0           int `json:"mechanicEngineer0"`
		MechanicEngineer1           int `json:"mechanicEngineer1"`
		MechanicEngineer5           int `json:"mechanicEngineer5"`
		MedalTarczay                int `json:"medalTarczay"`
		Sinai                       int `json:"sinai"`
		PattonValley                int `json:"pattonValley"`
		MedalDeLanglade             int `json:"medalDeLanglade"`
		DiamondTwisterMedal         int `json:"diamondTwisterMedal"`
		Beasthunter                 int `json:"beasthunter"`
		SupremacyVeteran            int `json:"supremacyVeteran"`
		Kamikaze                    int `json:"kamikaze"`
		OlimpicBronze               int `json:"olimpicBronze"`
		MedalTournamentOffseason3   int `json:"medalTournamentOffseason3"`
		PlatinumClanRibbonRU        int `json:"platinumClanRibbonRU"`
		MedalOskin                  int `json:"medalOskin"`
		Invincible                  int `json:"invincible"`
		PlatinumClanRibbonNA        int `json:"platinumClanRibbonNA"`
		PlatinumTwisterMedalRU      int `json:"platinumTwisterMedalRU"`
		ContinentalViceChampion     int `json:"continentalViceChampion"`
		OlimpicSilver               int `json:"olimpicSilver"`
		MarkOfMasteryI              int `json:"markOfMasteryI"`
		ContinentalCompetitor       int `json:"continentalCompetitor"`
		MedalTournamentSummerSeason int `json:"medalTournamentSummerSeason"`
		Mousebane                   int `json:"mousebane"`
		MedalBrunoPietro            int `json:"medalBrunoPietro"`
		MedalTournamentSpringSeason int `json:"medalTournamentSpringSeason"`
		GoldTwisterMark             int `json:"goldTwisterMark"`
		CollectorWarhammer          int `json:"collectorWarhammer"`
		MarkOfMasteryIII            int `json:"markOfMasteryIII"`
		MedalLeClerc                int `json:"medalLeClerc"`
		MedalTournamentProfessional int `json:"medalTournamentProfessional"`
		MedalCommunityChampion      int `json:"medalCommunityChampion"`
		DiamondTwisterMark          int `json:"diamondTwisterMark"`
		PlatinumTwisterMedalNA      int `json:"platinumTwisterMedalNA"`
		HandOfDeath                 int `json:"handOfDeath"`
		Huntsman                    int `json:"huntsman"`
		Camper                      int `json:"camper"`
		MedalNikolas                int `json:"medalNikolas"`
		AndroidTest                 int `json:"androidTest"`
		Sturdy                      int `json:"sturdy"`
		MedalTwitch                 int `json:"medalTwitch"`
		MedalWGfestTicket           int `json:"medalWGfestTicket"`
		ChampionofNewStart          int `json:"championofNewStart"`
	} `json:"achievements"`
	MaxSeries struct {
		ArmorPiercer int `json:"armorPiercer"`
		Punisher     int `json:"punisher"`
		TitleSniper  int `json:"titleSniper"`
		Invincible   int `json:"invincible"`
		TankExpert   int `json:"tankExpert"`
		MedalKay     int `json:"medalKay"`
		Diehard      int `json:"diehard"`
		Beasthunter  int `json:"beasthunter"`
		HandOfDeath  int `json:"handOfDeath"`
		JointVictory int `json:"jointVictory"`
		Sinai        int `json:"sinai"`
		PattonValley int `json:"pattonValley"`
	} `json:"max_series"`
}

// Diff - Find changes in achievements
func (new AchievementsFrame) Diff(old AchievementsFrame) (result AchievementsFrame) {
	// MaxSeries
	result.MaxSeries = new.MaxSeries
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
