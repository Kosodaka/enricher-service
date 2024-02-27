package enricher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kosodaka/enricher-service/internal/domain/ports/enricher"
	"net/http"
	"sync"
)

type Configer interface {
	GetAgeApiURL() string
	GetGenderApiURL() string
	GetNationalityApiURL() string
}

type Enricher struct {
	AgeUrl         string `json:"age_url"`
	GenderUrl      string `json:"gender_url"`
	NationalityUrl string `json:"nationality_url"`
	Client         *http.Client
}

type PersonNationality struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
type PersonNationalities struct {
	Country []PersonNationality `json:"country"`
}
type PersonAge struct {
	Age int `json:"age"`
}
type PersonGender struct {
	Gender string `json:"gender"`
}

func NewEnricher(cfg Configer) *Enricher {
	return &Enricher{
		AgeUrl:         cfg.GetAgeApiURL(),
		GenderUrl:      cfg.GetGenderApiURL(),
		NationalityUrl: cfg.GetNationalityApiURL(),
		Client:         &http.Client{},
	}
}
func (e Enricher) Enrich(ctx context.Context, name string) (*enricher.EnrichData, error) {
	errCh := make(chan error)
	resCh := make(chan *enricher.EnrichData)

	newCtx, cansel := context.WithCancel(ctx)
	defer cansel()

	go func() {
		defer close(errCh)
		defer close(resCh)

		var (
			age         *PersonAge
			gender      *PersonGender
			nationality *PersonNationality
			err         error
		)

		w := &sync.WaitGroup{}
		w.Add(3)

		go func() {
			defer w.Done()
			age, err = e.getAge(newCtx, name)
			if err != nil {
				errCh <- err
				return
			}
		}()

		go func() {
			defer w.Done()
			gender, err = e.getGender(newCtx, name)
			if err != nil {
				errCh <- err
				return
			}
		}()

		go func() {
			defer w.Done()
			nationalities, err := e.getNationality(newCtx, name)
			if err != nil {
				errCh <- err
				return
			}
			if len(nationalities.Country) == 0 {
				errCh <- fmt.Errorf("no nationality")
				return
			}
			// The first nationality from api url has the most probability
			nationality = &PersonNationality{CountryId: nationalities.Country[0].CountryId, Probability: nationalities.Country[0].Probability}
		}()

		w.Wait()
		if age != nil && gender != nil && nationality != nil {
			resCh <- &enricher.EnrichData{
				Age:         age.Age,
				Gender:      gender.Gender,
				Nationality: nationality.CountryId,
			}

		}

	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("error whith enricher")
	case err := <-errCh:
		return nil, err
	case res := <-resCh:
		return res, nil
	}

}

// Come to api with request on env:AGE_API_URL and get Age
func (e Enricher) getAge(ctx context.Context, name string) (*PersonAge, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", e.AgeUrl, name), nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("errror to get age")
	}
	defer resp.Body.Close()
	age := &PersonAge{}
	err = json.NewDecoder(resp.Body).Decode(&age)
	if err != nil {
		return nil, fmt.Errorf("%s: error to decode", err)
	}

	return age, nil
}

// Come to api with request on env:GENDER_API_URL and get gender
func (e Enricher) getGender(ctx context.Context, name string) (*PersonGender, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", e.GenderUrl, name), nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("errror to get gender")
	}
	defer resp.Body.Close()
	gender := &PersonGender{}
	err = json.NewDecoder(resp.Body).Decode(&gender)
	if err != nil {
		return nil, err
	}

	return gender, nil
}

// Come to api with request on env:NATIONALITY_API_URL get nationality
func (e Enricher) getNationality(ctx context.Context, name string) (*PersonNationalities, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s?name=%s", e.NationalityUrl, name), nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("errror to get nationality")
	}
	defer resp.Body.Close()
	nationalities := &PersonNationalities{}
	err = json.NewDecoder(resp.Body).Decode(&nationalities)
	if err != nil {
		return nil, err
	}

	return nationalities, nil
}
