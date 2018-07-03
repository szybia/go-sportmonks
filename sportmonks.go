package sportmonks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antonholmquist/jason"
)

var apiToken = ""
var apiURL = "https://soccer.sportmonks.com/api/v2.0/"
var dataNotFound = "key 'data' not found"
var tooManyAttemptsError = "invalid character 'T' looking for beginning of value"
var tooManyAttemptsMessage = "Too Many Attempts."

//Logger is a dedicated sportmonks logger
var Logger = log.Logger{}

type paginatedRequest struct {
	pageNumber int64
	data       []*jason.Object
}

//NoIncludes specifies the default of no includes
var NoIncludes = ""

//FirstOrAllPages specifies the default when a specific page is not requested
//this default should be used when requesting the first page or all pages
var FirstOrAllPages = 0

//AllPages specifies to fetch all pages available for that request
var AllPages = true

//SinglePage specifies to only fetch a single page for the request
var SinglePage = false

//SetAPIToken sets the API token for sportmonks
func SetAPIToken(s string) {
	apiToken = s
}

//Get API request
func Get(endpoint string, include string, page int, allPages bool) ([]byte, error) {
	if endpoint == "" {
		return []byte{}, errors.New("no endpoint provided")
	} else if apiToken == "" {
		return []byte{}, errors.New("apiToken has not been set")
	}

	requestURL := apiURL + endpoint
	r, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return []byte{}, err
	}

	q := r.URL.Query()
	q.Add("api_token", apiToken)
	if include != "" {
		q.Add("include", include)
	}
	if page != FirstOrAllPages {
		q.Add("page", strconv.Itoa(page))
		allPages = false
	}
	r.URL.RawQuery = q.Encode()

	resp, err := http.Get(r.URL.String())
	if err != nil {
		return []byte{}, err
	}

	body, err := jason.NewObjectFromReader(resp.Body)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), tooManyAttemptsError) {
			err = errors.New(tooManyAttemptsMessage)
		}
		return []byte{}, err
	}

	var isObject bool
	var dataObject *jason.Object
	var dataObjectArray []*jason.Object

	dataObjectArray, err = body.GetObjectArray("data")
	if err != nil {
		dataObject, err = body.GetObject("data")
		if err != nil {
			if strings.Contains(fmt.Sprint(err), dataNotFound) {
				//No errors. Empty response
				return []byte{}, nil
			}
			return []byte{}, err
		}
		isObject = true
	}
	if allPages {
		pages, err := body.GetInt64("meta", "pagination", "total_pages")
		//	No error means endpoint is paginated
		if err == nil {
			if pages > 1 {
				c := make(chan paginatedRequest)
				requests := make([][]*jason.Object, pages)
				for i := int64(2); i <= pages; i++ {
					go getRequest(r.URL.String(), i, c)
				}

				for i := int64(2); i <= pages; i++ {
					g := <-c
					requests[g.pageNumber-1] = g.data
				}

				for i := int64(1); i < pages; i++ {
					dataObjectArray = append(dataObjectArray, requests[i]...)
				}
			}
		}
	}
	var m []byte

	if isObject {
		m, err = json.Marshal(dataObject)
		if err != nil {
			return []byte{}, err
		}
	} else {
		m, err = json.Marshal(dataObjectArray)
		if err != nil {
			return []byte{}, err
		}
	}
	return m, nil
}

//Goroutine which fetches a specified page
func getRequest(requestURL string, pageNumber int64, c chan paginatedRequest) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		Logger.Println(err)
		c <- paginatedRequest{pageNumber, []*jason.Object{}}
		return
	}

	q := req.URL.Query()
	q.Add("page", strconv.FormatInt(pageNumber, 10))
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		Logger.Println(err)
		c <- paginatedRequest{pageNumber, []*jason.Object{}}
		return
	}

	body, err := jason.NewObjectFromReader(resp.Body)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), tooManyAttemptsError) {
			err = errors.New(tooManyAttemptsMessage)
		}

		Logger.Println(err)
		c <- paginatedRequest{pageNumber, []*jason.Object{}}
		return
	}

	var isObject bool
	var dataObject *jason.Object
	var dataObjectArray []*jason.Object

	dataObjectArray, err = body.GetObjectArray("data")
	if err != nil {
		dataObject, err = body.GetObject("data")
		if err != nil {
			if strings.Contains(fmt.Sprint(err), dataNotFound) {
				//No errors. Empty response
				c <- paginatedRequest{pageNumber, []*jason.Object{}}
				return
			}
			Logger.Println(err)
			c <- paginatedRequest{pageNumber, []*jason.Object{}}
			return
		}
		isObject = true
	}

	if isObject {
		c <- paginatedRequest{pageNumber, []*jason.Object{dataObject}}
	} else {
		c <- paginatedRequest{pageNumber, dataObjectArray}
	}
}

//Continents request for all continents
func Continents(include string) ([]byte, error) {
	return Get("continents", include, FirstOrAllPages, SinglePage)
}

//Continent request for continent identified by supplied ID
func Continent(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("continents/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Countries request for all countries
func Countries(include string, page int, allPages bool) ([]byte, error) {
	return Get("countries", include, page, allPages)
}

//Country request for specific country identified by ID
func Country(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("countries/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Leagues request for all leagues
func Leagues(include string, page int, allPages bool) ([]byte, error) {
	return Get("leagues", include, page, allPages)
}

//League request for specific league identified by ID
func League(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("leagues/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Seasons request for all seasons
func Seasons(include string, page int, allPages bool) ([]byte, error) {
	return Get("seasons", include, page, allPages)
}

//Season request for specific season identified by ID
func Season(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("seasons/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Fixture request for specific fixture identified by ID
func Fixture(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("fixtures/%v", ID), include, FirstOrAllPages, SinglePage)
}

//FixturesFromToDateTeam fetches all fixtures between the specified dates for the team identified by teamID
func FixturesFromToDateTeam(fromDate, toDate string, teamID int, include string, page int, allPages bool) ([]byte, error) {
	return Get(fmt.Sprintf("fixtures/between/%v/%v/%v", fromDate, toDate, teamID), include, page, allPages)
}

//FixturesDate fetches fixtures played at a specific date
func FixturesDate(date string, include string, page int, allPages bool) ([]byte, error) {
	return Get(fmt.Sprintf("fixtures/date/%v", date), include, page, allPages)
}

//FixturesFromToDate fetches fixtures played between the specified dates
func FixturesFromToDate(fromDate, toDate string, include string, page int, allPages bool) ([]byte, error) {
	return Get(fmt.Sprintf("fixtures/between/%v/%v", fromDate, toDate), include, page, allPages)
}

//FixturesMultipleList fetches all fixtures with IDs contained in the supplied comma-seperated int string
func FixturesMultipleList(IDs, include string) ([]byte, error) {
	return Get(fmt.Sprintf("fixtures/multi/%v", IDs), include, FirstOrAllPages, SinglePage)
}

//FixturesMultipleIntList fetches all fixtures with IDs contained in the supplied int slice
func FixturesMultipleIntList(IDs []int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("fixtures/multi/%v", IntSliceToSepString(IDs, ",")), include, FirstOrAllPages, SinglePage)
}

//StagesSeason fetches all stages for season specified by ID
func StagesSeason(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("stages/season/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Stage fetches stage identified by ID
func Stage(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("stages/%v", ID), include, FirstOrAllPages, SinglePage)
}

//LivescoresNow fetches live fixtures
func LivescoresNow(include string) ([]byte, error) {
	return Get("livescores/now", include, FirstOrAllPages, SinglePage)
}

//Livescores fetches fixtures which are played today
func Livescores(include string, page int, allPages bool) ([]byte, error) {
	return Get("livescores", include, page, allPages)
}

//CommentariesFixture fetches textual representations of events within a fixture
func CommentariesFixture(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("commentaries/fixture/%v", ID), "", FirstOrAllPages, SinglePage)
}

//VideoHighlights fetches links to videos posted on social media (Community feature)
func VideoHighlights(include string, page int, allPages bool) ([]byte, error) {
	return Get("highlights", include, page, allPages)
}

//Head2Head fetches all fixtures involving two teams
func Head2Head(team1ID, team2ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("head2head/%v/%v", team1ID, team2ID), include, FirstOrAllPages, SinglePage)
}

//TvStationsFixture fetches all Tv stations which broadcast the specified fixture
func TvStationsFixture(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("tvstations/fixture/%v", ID), "", FirstOrAllPages, SinglePage)
}

//StandingsSeason fetches standings for the specified season
func StandingsSeason(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("standings/season/%v", ID), include, FirstOrAllPages, SinglePage)
}

//LiveStandingsSeason fetches live standings for the specified season
func LiveStandingsSeason(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("standings/season/live/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Team fetches a team identified by specified ID
func Team(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("teams/%v", ID), include, FirstOrAllPages, SinglePage)
}

//SeasonTeams fetches all teams from a specified season
func SeasonTeams(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("teams/season/%v", ID), include, FirstOrAllPages, SinglePage)
}

//SeasonTopScorer fetches the top goal scorers for a specified season
func SeasonTopScorer(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("topscorers/season/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Venue fetches a venue specified by ID
func Venue(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("venues/%v", ID), "", FirstOrAllPages, SinglePage)
}

//SeasonVenues fetches venues for specified season
func SeasonVenues(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("venues/season/%v", ID), "", FirstOrAllPages, SinglePage)
}

//SeasonRounds fetches all rounds for specified season
func SeasonRounds(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("rounds/season/%v", ID), include, FirstOrAllPages, SinglePage)
}

//Round fetches round specified by ID
func Round(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("rounds/%v", ID), include, FirstOrAllPages, SinglePage)
}

//OddsFixtureBookmaker fetches betting information for specified fixture for specified bookmaker
func OddsFixtureBookmaker(fixtureID, bookmakerID int) ([]byte, error) {
	return Get(fmt.Sprintf("odds/fixture/%v/bookmaker/%v", fixtureID, bookmakerID), "", FirstOrAllPages, SinglePage)
}

//OddsFixture fetches betting odds for specified fixture
func OddsFixture(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("odds/fixture/%v", ID), "", FirstOrAllPages, SinglePage)
}

//OddsFixtureMarket fetches betting odds for specified fixture for specified market
func OddsFixtureMarket(fixtureID, marketID int) ([]byte, error) {
	return Get(fmt.Sprintf("odds/fixture/%v/market/%v", fixtureID, marketID), "", FirstOrAllPages, SinglePage)
}

//OddsInPlayFixture fetches in play odds for specified fixture
func OddsInPlayFixture(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("odds/inplay/fixture/%v", ID), "", FirstOrAllPages, SinglePage)
}

//Bookmakers fetches all bookmakers
func Bookmakers() ([]byte, error) {
	return Get("bookmakers", "", FirstOrAllPages, SinglePage)
}

//Bookmaker fetches bookmaker specified by ID
func Bookmaker(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("bookmakers/%v", ID), "", FirstOrAllPages, SinglePage)
}

//Markets fetches all markets
func Markets() ([]byte, error) {
	return Get("markets", "", FirstOrAllPages, SinglePage)
}

//Market fetches a market specified by ID
func Market(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("markets/%v", ID), "", FirstOrAllPages, SinglePage)
}

//Player fetches a player specified by ID
func Player(ID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("players/%v", ID), include, FirstOrAllPages, SinglePage)
}

//SeasonTeamSquad fetches a squad for a specified team for a specified season
func SeasonTeamSquad(seasonID, teamID int, include string) ([]byte, error) {
	return Get(fmt.Sprintf("squad/season/%v/team/%v", seasonID, teamID), include, FirstOrAllPages, SinglePage)
}

//Coach fetches a squad specified by ID
func Coach(ID int) ([]byte, error) {
	return Get(fmt.Sprintf("coaches/%v", ID), "", FirstOrAllPages, SinglePage)
}

//IntSliceToSepString takes a int slice and generates a string of values seperated by specified seperator
func IntSliceToSepString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
