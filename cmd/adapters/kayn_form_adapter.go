package adapters

import (
	"encoding/json"
	"io/ioutil"
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
	formResp := &models.FormResponse{}
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
	summoner, err := client.Riot.LoL.Summoner.GetByName(formResp.SummonerName)
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

	// Loop through match history. We only care about Kayn games played here.
	for _, v := range matchHistoryActual {
		if v.Info.QueueID == 420 { // Ranked games only for now
			for _, v := range v.Info.Participants {
				if v.PUUID == summoner.PUUID {
					if v.ChampionName == "Kayn" {
						// Now here, we check the transformation...
						// 0 is base form, 1 is Rhaast, 2 is Shadow Assassin.
						// should we give a shit about base form Kayn?
						if v.ChampionTransform == 1 {
							// Fill in Rhaast data here.
						} else if v.ChampionTransform == 2 {
							// Shadow Assasin data goes here
						} else {
							// No data found, display nothing.
						}
					}
				}
			}

		}
	}
}
