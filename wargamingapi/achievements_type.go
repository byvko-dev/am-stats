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

// Achievements - Frame for achievements on vehicle
type Achievements struct {
	ArmorPiercer                int `json:"armorPiercer,omitempty" bson:"armorPiercer,omitempty"`
	MedalFadin                  int `json:"medalFadin,omitempty" bson:"medalFadin,omitempty"`
	MedalCarius                 int `json:"medalCarius,omitempty" bson:"medalCarius,omitempty"`
	MedalEkins                  int `json:"medalEkins,omitempty" bson:"medalEkins,omitempty"`
	CollectorGuP                int `json:"collectorGuP,omitempty" bson:"collectorGuP,omitempty"`
	MedalHalonen                int `json:"medalHalonen,omitempty" bson:"medalHalonen,omitempty"`
	HeroesOfRassenay            int `json:"heroesOfRassenay,omitempty" bson:"heroesOfRassenay,omitempty"`
	FirstVictory                int `json:"firstVictory,omitempty" bson:"firstVictory,omitempty"`
	Defender                    int `json:"defender,omitempty" bson:"defender,omitempty"`
	Creative                    int `json:"creative,omitempty" bson:"creative,omitempty"`
	ESportFinal                 int `json:"eSportFinal,omitempty" bson:"eSportFinal,omitempty"`
	Supporter                   int `json:"supporter,omitempty" bson:"supporter,omitempty"`
	GoldClanRibbonSEA           int `json:"goldClanRibbonSEA,omitempty" bson:"goldClanRibbonSEA,omitempty"`
	PlatinumTwisterMedalSEA     int `json:"platinumTwisterMedalSEA,omitempty" bson:"platinumTwisterMedalSEA,omitempty"`
	MedalLehvaslaiho            int `json:"medalLehvaslaiho,omitempty" bson:"medalLehvaslaiho,omitempty"`
	TankExpert                  int `json:"tankExpert,omitempty" bson:"tankExpert,omitempty"`
	ESportQualification         int `json:"eSportQualification,omitempty" bson:"eSportQualification,omitempty"`
	MarkI                       int `json:"MarkI,omitempty" bson:"MarkI,omitempty"`
	MedalSupremacy              int `json:"medalSupremacy,omitempty" bson:"medalSupremacy,omitempty"`
	ParticipantofWGFest2017     int `json:"participantofWGFest2017,omitempty" bson:"participantofWGFest2017,omitempty"`
	MedalTournamentOffseason1   int `json:"medalTournamentOffseason1,omitempty" bson:"medalTournamentOffseason1,omitempty"`
	JointVictory                int `json:"jointVictory,omitempty" bson:"jointVictory,omitempty"`
	MedalTournamentOffseason2   int `json:"medalTournamentOffseason2,omitempty" bson:"medalTournamentOffseason2,omitempty"`
	MedalTournamentOffseason4   int `json:"medalTournamentOffseason4,omitempty" bson:"medalTournamentOffseason4,omitempty"`
	Sniper                      int `json:"sniper,omitempty" bson:"sniper,omitempty"`
	TitleSniper                 int `json:"titleSniper,omitempty" bson:"titleSniper,omitempty"`
	MedalCrucialContribution    int `json:"medalCrucialContribution,omitempty" bson:"medalCrucialContribution,omitempty"`
	Scout                       int `json:"scout,omitempty" bson:"scout,omitempty"`
	GoldTwisterMedalRU          int `json:"goldTwisterMedalRU,omitempty" bson:"goldTwisterMedalRU,omitempty"`
	TankExpert3                 int `json:"tankExpert3,omitempty" bson:"tankExpert3,omitempty"`
	TankExpert2                 int `json:"tankExpert2,omitempty" bson:"tankExpert2,omitempty"`
	TankExpert1                 int `json:"tankExpert1,omitempty" bson:"tankExpert1,omitempty"`
	TankExpert0                 int `json:"tankExpert0,omitempty" bson:"tankExpert0,omitempty"`
	MarkOfMastery               int `json:"markOfMastery,omitempty" bson:"markOfMastery,omitempty"`
	TankExpert6                 int `json:"tankExpert6,omitempty" bson:"tankExpert6,omitempty"`
	TankExpert5                 int `json:"tankExpert5,omitempty" bson:"tankExpert5,omitempty"`
	TankExpert4                 int `json:"tankExpert4,omitempty" bson:"tankExpert4,omitempty"`
	GoldTwisterMedalEU          int `json:"goldTwisterMedalEU,omitempty" bson:"goldTwisterMedalEU,omitempty"`
	ChristmasTreeLevelUpNY2019  int `json:"ChristmasTreeLevelUpNY2019,omitempty" bson:"ChristmasTreeLevelUpNY2019,omitempty"`
	MedalLavrinenko             int `json:"medalLavrinenko,omitempty" bson:"medalLavrinenko,omitempty"`
	MedalKolobanov              int `json:"medalKolobanov,omitempty" bson:"medalKolobanov,omitempty"`
	MedalLafayettePool          int `json:"medalLafayettePool,omitempty" bson:"medalLafayettePool,omitempty"`
	GoldClanRibbonEU            int `json:"goldClanRibbonEU,omitempty" bson:"goldClanRibbonEU,omitempty"`
	OlimpicGolden               int `json:"olimpicGolden,omitempty" bson:"olimpicGolden,omitempty"`
	MedalKnispel                int `json:"medalKnispel,omitempty" bson:"medalKnispel,omitempty"`
	Invader                     int `json:"invader,omitempty" bson:"invader,omitempty"`
	GoldTwisterMedalNA          int `json:"goldTwisterMedalNA,omitempty" bson:"goldTwisterMedalNA,omitempty"`
	MechanicEngineer            int `json:"mechanicEngineer,omitempty" bson:"mechanicEngineer,omitempty"`
	MarkOfMasteryII             int `json:"markOfMasteryII,omitempty" bson:"markOfMasteryII,omitempty"`
	FirstBlood                  int `json:"firstBlood,omitempty" bson:"firstBlood,omitempty"`
	MedalKay                    int `json:"medalKay,omitempty" bson:"medalKay,omitempty"`
	MedalOrlik                  int `json:"medalOrlik,omitempty" bson:"medalOrlik,omitempty"`
	MedalBrothersInArms         int `json:"medalBrothersInArms,omitempty" bson:"medalBrothersInArms,omitempty"`
	MedalAbrams                 int `json:"medalAbrams,omitempty" bson:"medalAbrams,omitempty"`
	MedalAtgm                   int `json:"medalAtgm,omitempty" bson:"medalAtgm,omitempty"`
	MainGun                     int `json:"mainGun,omitempty" bson:"mainGun,omitempty"`
	IronMan                     int `json:"ironMan,omitempty" bson:"ironMan,omitempty"`
	PlatinumClanRibbonEU        int `json:"platinumClanRibbonEU,omitempty" bson:"platinumClanRibbonEU,omitempty"`
	PlatinumClanRibbonSEA       int `json:"platinumClanRibbonSEA,omitempty" bson:"platinumClanRibbonSEA,omitempty"`
	Warrior                     int `json:"warrior,omitempty" bson:"warrior,omitempty"`
	GoldClanRibbonRU            int `json:"goldClanRibbonRU,omitempty" bson:"goldClanRibbonRU,omitempty"`
	MedalRadleyWalters          int `json:"medalRadleyWalters,omitempty" bson:"medalRadleyWalters,omitempty"`
	Raider                      int `json:"raider,omitempty" bson:"raider,omitempty"`
	ParticipantofNewStart       int `json:"participantofNewStart,omitempty" bson:"participantofNewStart,omitempty"`
	DiamondClanRibbon           int `json:"diamondClanRibbon,omitempty" bson:"diamondClanRibbon,omitempty"`
	MedalBillotte               int `json:"medalBillotte,omitempty" bson:"medalBillotte,omitempty"`
	PlatinumTwisterMedalEU      int `json:"platinumTwisterMedalEU,omitempty" bson:"platinumTwisterMedalEU,omitempty"`
	Diehard                     int `json:"diehard,omitempty" bson:"diehard,omitempty"`
	MasterofContinents          int `json:"masterofContinents,omitempty" bson:"masterofContinents,omitempty"`
	Evileye                     int `json:"evileye,omitempty" bson:"evileye,omitempty"`
	Cadet                       int `json:"cadet,omitempty" bson:"cadet,omitempty"`
	SupremacyHunter             int `json:"supremacyHunter,omitempty" bson:"supremacyHunter,omitempty"`
	ContinentalContender        int `json:"continentalContender,omitempty" bson:"continentalContender,omitempty"`
	Steelwall                   int `json:"steelwall,omitempty" bson:"steelwall,omitempty"`
	SupremacyLegend             int `json:"supremacyLegend,omitempty" bson:"supremacyLegend,omitempty"`
	Punisher                    int `json:"punisher,omitempty" bson:"punisher,omitempty"`
	ESport                      int `json:"eSport,omitempty" bson:"eSport,omitempty"`
	PlatinumTwisterMark         int `json:"platinumTwisterMark,omitempty" bson:"platinumTwisterMark,omitempty"`
	GoldClanRibbonNA            int `json:"goldClanRibbonNA,omitempty" bson:"goldClanRibbonNA,omitempty"`
	MedalPoppel                 int `json:"medalPoppel,omitempty" bson:"medalPoppel,omitempty"`
	MechanicEngineer6           int `json:"mechanicEngineer6,omitempty" bson:"mechanicEngineer6,omitempty"`
	MechanicEngineer4           int `json:"mechanicEngineer4,omitempty" bson:"mechanicEngineer4,omitempty"`
	GoldTwisterMedalSEA         int `json:"goldTwisterMedalSEA,omitempty" bson:"goldTwisterMedalSEA,omitempty"`
	MechanicEngineer2           int `json:"mechanicEngineer2,omitempty" bson:"mechanicEngineer2,omitempty"`
	MechanicEngineer3           int `json:"mechanicEngineer3,omitempty" bson:"mechanicEngineer3,omitempty"`
	MechanicEngineer0           int `json:"mechanicEngineer0,omitempty" bson:"mechanicEngineer0,omitempty"`
	MechanicEngineer1           int `json:"mechanicEngineer1,omitempty" bson:"mechanicEngineer1,omitempty"`
	MechanicEngineer5           int `json:"mechanicEngineer5,omitempty" bson:"mechanicEngineer5,omitempty"`
	MedalTarczay                int `json:"medalTarczay,omitempty" bson:"medalTarczay,omitempty"`
	Sinai                       int `json:"sinai,omitempty" bson:"sinai,omitempty"`
	PattonValley                int `json:"pattonValley,omitempty" bson:"pattonValley,omitempty"`
	MedalDeLanglade             int `json:"medalDeLanglade,omitempty" bson:"medalDeLanglade,omitempty"`
	DiamondTwisterMedal         int `json:"diamondTwisterMedal,omitempty" bson:"diamondTwisterMedal,omitempty"`
	Beasthunter                 int `json:"beasthunter,omitempty" bson:"beasthunter,omitempty"`
	SupremacyVeteran            int `json:"supremacyVeteran,omitempty" bson:"supremacyVeteran,omitempty"`
	Kamikaze                    int `json:"kamikaze,omitempty" bson:"kamikaze,omitempty"`
	OlimpicBronze               int `json:"olimpicBronze,omitempty" bson:"olimpicBronze,omitempty"`
	MedalTournamentOffseason3   int `json:"medalTournamentOffseason3,omitempty" bson:"medalTournamentOffseason3,omitempty"`
	PlatinumClanRibbonRU        int `json:"platinumClanRibbonRU,omitempty" bson:"platinumClanRibbonRU,omitempty"`
	MedalOskin                  int `json:"medalOskin,omitempty" bson:"medalOskin,omitempty"`
	Invincible                  int `json:"invincible,omitempty" bson:"invincible,omitempty"`
	PlatinumClanRibbonNA        int `json:"platinumClanRibbonNA,omitempty" bson:"platinumClanRibbonNA,omitempty"`
	PlatinumTwisterMedalRU      int `json:"platinumTwisterMedalRU,omitempty" bson:"platinumTwisterMedalRU,omitempty"`
	ContinentalViceChampion     int `json:"continentalViceChampion,omitempty" bson:"continentalViceChampion,omitempty"`
	OlimpicSilver               int `json:"olimpicSilver,omitempty" bson:"olimpicSilver,omitempty"`
	MarkOfMasteryI              int `json:"markOfMasteryI,omitempty" bson:"markOfMasteryI,omitempty"`
	ContinentalCompetitor       int `json:"continentalCompetitor,omitempty" bson:"continentalCompetitor,omitempty"`
	MedalTournamentSummerSeason int `json:"medalTournamentSummerSeason,omitempty" bson:"medalTournamentSummerSeason,omitempty"`
	Mousebane                   int `json:"mousebane,omitempty" bson:"mousebane,omitempty"`
	MedalBrunoPietro            int `json:"medalBrunoPietro,omitempty" bson:"medalBrunoPietro,omitempty"`
	MedalTournamentSpringSeason int `json:"medalTournamentSpringSeason,omitempty" bson:"medalTournamentSpringSeason,omitempty"`
	GoldTwisterMark             int `json:"goldTwisterMark,omitempty" bson:"goldTwisterMark,omitempty"`
	CollectorWarhammer          int `json:"collectorWarhammer,omitempty" bson:"collectorWarhammer,omitempty"`
	MarkOfMasteryIII            int `json:"markOfMasteryIII,omitempty" bson:"markOfMasteryIII,omitempty"`
	MedalLeClerc                int `json:"medalLeClerc,omitempty" bson:"medalLeClerc,omitempty"`
	MedalTournamentProfessional int `json:"medalTournamentProfessional,omitempty" bson:"medalTournamentProfessional,omitempty"`
	MedalCommunityChampion      int `json:"medalCommunityChampion,omitempty" bson:"medalCommunityChampion,omitempty"`
	DiamondTwisterMark          int `json:"diamondTwisterMark,omitempty" bson:"diamondTwisterMark,omitempty"`
	PlatinumTwisterMedalNA      int `json:"platinumTwisterMedalNA,omitempty" bson:"platinumTwisterMedalNA,omitempty"`
	HandOfDeath                 int `json:"handOfDeath,omitempty" bson:"handOfDeath,omitempty"`
	Huntsman                    int `json:"huntsman,omitempty" bson:"huntsman,omitempty"`
	Camper                      int `json:"camper,omitempty" bson:"camper,omitempty"`
	MedalNikolas                int `json:"medalNikolas,omitempty" bson:"medalNikolas,omitempty"`
	AndroidTest                 int `json:"androidTest,omitempty" bson:"androidTest,omitempty"`
	Sturdy                      int `json:"sturdy,omitempty" bson:"sturdy,omitempty"`
	MedalTwitch                 int `json:"medalTwitch,omitempty" bson:"medalTwitch,omitempty"`
	MedalWGfestTicket           int `json:"medalWGfestTicket,omitempty" bson:"medalWGfestTicket,omitempty"`
	ChampionofNewStart          int `json:"championofNewStart,omitempty" bson:"championofNewStart,omitempty"`
}

// AchievementsFrame -
type AchievementsFrame struct {
	Achievements Achievements `json:"achievements"`
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
