package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kayn-form/cmd/models"
	"log"
	"net/http"
	"os"
	"strings"

	"kayn_form/models"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/sirupsen/logrus"
)

type KaynFormAdapter struct {
	*log.Logger
}

func NewKaynFormAdapter(logger *log.Logger) *KaynFormAdapter {
	return &KaynFormAdapter{
		logger,
	}
}

// GetSummonerInfo is an endpoint that accepts only POST requests.
// Most of this is stolen from my Help Pix project...
func (adapter *KaynFormAdapter) GetSummonerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		adapter.Logger.Printf("Not a valid POST request!")
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		adapter.Logger.Printf(`whoops there's an error while reading request body`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// read what the user has sent us in the format {"username": "", "region": ""}
	formResp := &models.UserData{}
	json.Unmarshal(body, formResp)
	// Make sure that the region is a valid Riot region.
	// TODO: Add new regions from South East Asia.
	var region api.Region
	switch strings.ToUpper(formResp.Region) {
	case "EUW":
		region = api.RegionEuropeWest
	case "EUNE":
		region = api.RegionEuropeNorthEast
	case "NA":
		region = api.RegionNorthAmerica
	case "BR":
		region = api.RegionBrasil
	case "KR":
		region = api.RegionKorea
	case "JP":
		region = api.RegionJapan
	case "LAN":
		region = api.RegionLatinAmericaNorth
	case "LAS":
		region = api.RegionLatinAmericaSouth
	case "OCE":
		region = api.RegionOceania
	case "RU":
		region = api.RegionRussia
	default:
		adapter.Logger.Print(`Invalid region provided.`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return

	}

	api_key := os.Getenv(`RIOT_API_KEY`)
	client := golio.NewClient(api_key, golio.WithRegion(region), golio.WithLogger(logrus.New()))
	summoner, err := client.Riot.LoL.Summoner.GetByName(formResp.Username)
	if err != nil {
		adapter.Logger.Printf(`error while getting summoner: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// the following code is a product of both me operating on zero sleep and some mank from accessing the data from riots side. anyway to make this a lot nicer would be gr8
	matchHistoryStream := client.Riot.LoL.Match.ListStream(summoner.PUUID)
	var matchIDs []string
	for k := range matchHistoryStream {
		matchIDs = append(matchIDs, k.MatchID)
	}
	var matchHistoryActual []*lol.Match
	for _, v := range matchIDs {
		match, _ := client.Riot.LoL.Match.Get(v)

		matchHistoryActual = append(matchHistoryActual, match)
	}
	var kaynData *models.KaynData
	// Loop through match history. We only care about Kayn games played here.
	for _, v := range matchHistoryActual {
		if v.Info.QueueID == 420 { // Ranked games only for now
			for _, v := range v.Info.Participants {
				if v.PUUID == summoner.PUUID {
					if v.ChampionName == "Kayn" {
						var items []string
						i1, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item1))
						i2, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item2))
						i3, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item3))
						i4, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item4))
						i5, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item5))
						i6, _ := client.DataDragon.GetItem(fmt.Sprint(v.Item6))
						items = append(items, i1.Name)
						items = append(items, i2.Name)
						items = append(items, i3.Name)
						items = append(items, i4.Name)
						items = append(items, i5.Name)
						items = append(items, i6.Name)
						// Now here, we check the transformation...
						// 0 is base form, 1 is Rhaast, 2 is Shadow Assassin.
						// should we give a shit about base form Kayn?
						if v.ChampionTransform == 1 {
							kaynData = adapter.MakeKaynData("Rhaast", v, items)

						} else if v.ChampionTransform == 2 {
							kaynData = adapter.MakeKaynData("Shadow Assassin", v, items)
						} else {
							adapter.Logger.Printf(`no kayn data found`)
						}
					}
				}
			}

		}
	}
	// Marshall the data we get from the Riot API into a nice JSON blob
	kaynDataJSON, err := json.Marshal(kaynData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	writeResponse(w, kaynDataJSON, http.StatusOK)

}

func (adapter *KaynFormAdapter) MakeKaynData(transformation string, p *lol.Participant, items []string) *models.KaynData {
	ratio := (p.Kills + p.Assists) / p.Deaths
	return &models.KaynData{
		Transformation:      transformation,
		Level:               p.ChampLevel,
		Deaths:              p.Deaths,
		FirstBlood:          p.FirstBloodKill,
		Gold:                p.GoldEarned,
		Victory:             p.Win,
		Items:               items,
		Kills:               p.Kills,
		Pentakills:          p.PentaKills,
		LargestKillingSpree: p.LargestKillingSpree,
		ObjectivesStolen:    p.ObjectivesStolen,
		TimePlayed:          p.TimePlayed,
		TotalDamage:         p.TotalDamageDealt,
		TotalDamageTaken:    p.TotalDamageTaken,
		DoubleKills:         p.DoubleKills,
		TripleKills:         p.TripleKills,
		QuadraKills:         p.QuadraKills,
		TrueDamage:          p.TrueDamageDealt,
		VisionScore:         p.VisionScore,
		BaronKills:          p.BaronKills,
		DragonKills:         p.DragonKills,
		Assists:             p.Assists,
		Ratio:               ratio,
	}
}
