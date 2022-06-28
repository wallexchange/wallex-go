package wallex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.wallex.ir"

// Error is a service error.
type Error struct {
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("wallex: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("wallex: %s", e.Message)
}

// Client provides idiomatic methods to call Wallex API.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// ClientOptions customizes client's properties.
type ClientOptions struct {

	// APIKey is optional, but necessary to call private API.
	// If empty, defaults to WALLEX_API_KEY environment variable.
	APIKey string

	// HTTPClient is used to establish connection and to send HTTP requests.
	// If nil, it defaults to http.DefaultClient.
	HTTPClient *http.Client
}

// New instantiates a new Client.
func New(opt ClientOptions) *Client {
	c := &Client{}
	if opt.APIKey == "" {
		c.apiKey = os.Getenv("WALLEX_API_KEY")
	}
	if opt.HTTPClient != nil {
		c.httpClient = opt.HTTPClient
	} else {
		c.httpClient = http.DefaultClient
	}
	return c
}

// -----------------------------------------------------------------------------
// Markets
// -----------------------------------------------------------------------------

// Market represents a market information.
type Market struct {
	Symbol             string `json:"symbol"`
	BaseAsset          string `json:"baseAsset"`
	BaseAssetPrecision int    `json:"baseAssetPrecision"`
	QuoteAsset         string `json:"quoteAsset"`
	QuotePrecision     int    `json:"quotePrecision"`
	FarsiName          string `json:"faName"`
	FarsiBaseAsset     string `json:"faBaseAsset"`
	FarsiQuoteAsset    string `json:"faQuoteAsset"`
	StepSize           int    `json:"stepSize"`
	TickSize           int    `json:"tickSize"`
	MinQty             Number `json:"minQty"`
	MinNotional        Number `json:"minNotional"`
	Stats              struct {
		BidPrice       Number `json:"bidPrice"`
		AskPrice       Number `json:"askPrice"`
		Change24H      Number `json:"24h_ch"`
		Change7D       Number `json:"7d_ch"`
		Volume24H      Number `json:"24h_volume"`
		Volume7D       Number `json:"7d_volume"`
		QuoteVolume24H Number `json:"24h_quoteVolume"`
		HighPrice24H   Number `json:"24h_highPrice"`
		LowPrice24H    Number `json:"24h_lowPrice"`
		LastPrice      Number `json:"lastPrice"`
		LastQty        Number `json:"lastQty"`
		LastTradeSide  string `json:"lastTradeSide"`
		BidVolume      Number `json:"bidVolume"`
		AskVolume      Number `json:"askVolume"`
		BidCount       Number `json:"bidCount"`
		AskCount       Number `json:"askCount"`
		Direction      struct {
			Sell int `json:"SELL"`
			Buy  int `json:"BUY"`
		} `json:"direction"`
	} `json:"stats"`
	CreatedAt time.Time `json:"createdAt"`
}

// Markets retrieves a list of all available markets and their stats.
func (c *Client) Markets() ([]*Market, error) {
	resp, err := c.httpClient.Get(baseURL + "/v1/markets")
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result struct {
			Symbols map[string]*Market `json:"symbols"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	markets := make([]*Market, 0, len(result.Result.Symbols))
	for _, m := range result.Result.Symbols {
		markets = append(markets, m)
	}
	return markets, nil
}

// Currency represents a crypto-currency information.
type Currency struct {
	Key                 string    `json:"key"`
	Name                string    `json:"name"`
	NameEn              string    `json:"name_en"`
	Rank                int       `json:"rank"`
	Dominance           Number    `json:"dominance"`
	Volume24H           Number    `json:"volume_24h"`
	MarketCap           Number    `json:"market_cap"`
	ATH                 Number    `json:"ath"`
	ATHChangePercentage Number    `json:"ath_change_percentage"`
	ATHDate             time.Time `json:"ath_date"`
	Price               Number    `json:"price"`
	DailyHighPrice      Number    `json:"daily_high_price"`
	DailyLowPrice       Number    `json:"daily_low_price"`
	WeeklyHighPrice     Number    `json:"weekly_high_price"`
	WeeklyLowPrice      Number    `json:"weekly_low_price"`
	PercentChange1H     Number    `json:"percent_change_1h"`
	PercentChange24H    Number    `json:"percent_change_24h"`
	PercentChange7D     Number    `json:"percent_change_7d"`
	PercentChange14D    Number    `json:"percent_change_14d"`
	PercentChange30D    Number    `json:"percent_change_30d"`
	PercentChange60D    Number    `json:"percent_change_60d"`
	PercentChange200D   Number    `json:"percent_change_200d"`
	PercentChange1Y     Number    `json:"percent_change_1y"`
	PriceChange24H      Number    `json:"price_change_24h"`
	PriceChange7D       Number    `json:"price_change_7d"`
	PriceChange14D      Number    `json:"price_change_14d"`
	PriceChange30D      Number    `json:"price_change_30d"`
	PriceChange60D      Number    `json:"price_change_60d"`
	PriceChange200D     Number    `json:"price_change_200d"`
	PriceChange1Y       Number    `json:"price_change_1y"`
	MaxSupply           Number    `json:"max_supply"`
	TotalSupply         Number    `json:"total_supply"`
	CirculatingSupply   Number    `json:"circulating_supply"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Currencies retrieves a list of all available crypto-currencies and their stats.
func (c *Client) Currencies() ([]*Currency, error) {
	resp, err := c.httpClient.Get(baseURL + "/v1/currencies/stats")
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result []*Currency `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// MarketOrder represents an open order in the market.
type MarketOrder struct {
	Price    Number `json:"price"`
	Quantity Number `json:"quantity"`
	Sum      Number `json:"sum"`
}

// OpenOrders retrieves list of open-orders in a market.
func (c *Client) OpenOrders(symbol string) (ask []*MarketOrder, bid []*MarketOrder, _ error) {
	query := url.Values{}
	query.Add("symbol", symbol)

	resp, err := c.httpClient.Get(baseURL + "/v1/depth?" + query.Encode())
	if err != nil {
		return nil, nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result struct {
			Ask []*MarketOrder `json:"ask"`
			Bid []*MarketOrder `json:"bid"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, wrapRequestError(err)
	}

	return result.Result.Ask, result.Result.Bid, nil
}

// MarketOrder represents an trade in the market.
type MarketTrade struct {
	Symbol    string    `json:"symbol"`
	Quantity  Number    `json:"quantity"`
	Price     Number    `json:"price"`
	Sum       Number    `json:"sum"`
	Timestamp time.Time `json:"timestamp"`
}

// OpenOrders retrieves list of most recent trades in a market.
func (c *Client) Trades(symbol string) ([]*MarketTrade, error) {
	query := url.Values{}
	query.Add("symbol", symbol)

	resp, err := c.httpClient.Get(baseURL + "/v1/trades?" + query.Encode())
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result struct {
			LatestTrades []*MarketTrade `json:"latestTrades"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result.LatestTrades, nil
}

// Candle is an OHLCV candle.
type Candle struct {
	Timestamp time.Time `json:"timestamp"`
	Open      Number    `json:"open"`
	High      Number    `json:"high"`
	Low       Number    `json:"low"`
	Close     Number    `json:"close"`
	Volume    Number    `json:"volume"`
}

// List of available candle resolutions.
const (
	Minute     = "1"
	Hour       = "60"
	ThreeHour  = "180"
	SixHour    = "360"
	TwelveHour = "720"
	Day        = "1D"
)

// Candles retrieves OHLCV candles for the given time interval.
func (c *Client) Candles(symbol, resolution string, from, to time.Time) ([]*Candle, error) {
	query := url.Values{}
	query.Add("symbol", symbol)
	query.Add("resolution", resolution)
	query.Add("from", strconv.FormatInt(from.Unix(), 10))
	query.Add("to", strconv.FormatInt(to.Unix(), 10))

	resp, err := http.Get(baseURL + "/v1/udf/history?" + query.Encode())
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		T []int64  `json:"t"`
		O []string `json:"o"`
		H []string `json:"h"`
		L []string `json:"l"`
		C []string `json:"c"`
		V []string `json:"v"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	candles := make([]*Candle, 0, len(result.T))
	for i, t := range result.T {
		candles = append(candles, &Candle{
			Timestamp: time.Unix(t, 0),
			Open:      Number(result.O[i]),
			High:      Number(result.H[i]),
			Low:       Number(result.L[i]),
			Close:     Number(result.C[i]),
			Volume:    Number(result.V[i]),
		})
	}
	return candles, nil
}

// -----------------------------------------------------------------------------
// Error handling
// -----------------------------------------------------------------------------

// List of common service errors.
var (
	ErrBadRequest   = &Error{Message: "bad request"}
	ErrUnauthorized = &Error{Message: "unauthorized"}
	ErrForbidden    = &Error{Message: "access forbidden"}
	ErrNotFound     = &Error{Message: "resource not found"}
	ErrUnknown      = &Error{Message: "unknown error"}
)

func wrapRequestError(err error) error {
	return &Error{
		Message: "request failed",
		Cause:   err,
	}
}

func errNonOKResponse(code int) error {
	switch code {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return ErrUnknown
	}
}
