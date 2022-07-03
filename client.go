package wallex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	baseURL      = "https://api.wallex.ir"
	apiKeyHeader = "x-api-key"
)

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

// MarketOrders retrieves list of active orders in a market.
func (c *Client) MarketOrders(symbol string) (ask []*MarketOrder, bid []*MarketOrder, _ error) {
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

// MarketTrade represents an trade in the market.
type MarketTrade struct {
	Symbol    string    `json:"symbol"`
	Quantity  Number    `json:"quantity"`
	Price     Number    `json:"price"`
	Sum       Number    `json:"sum"`
	Timestamp time.Time `json:"timestamp"`
}

// MarketTrades retrieves list of most recent trades in a market.
func (c *Client) MarketTrades(symbol string) ([]*MarketTrade, error) {
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
// Account
// -----------------------------------------------------------------------------

// Profile represents a Wallex account.
type Profile struct {
	TrackingID   int       `json:"tracking_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	NationalCode string    `json:"national_code"`
	FaceImage    string    `json:"face_image"`
	Birthday     time.Time `json:"birthday"`
	Address      struct {
		Country     *string `json:"country"`
		Province    *string `json:"province"`
		City        *string `json:"city"`
		Location    *string `json:"location"`
		PostalCode  *string `json:"postal_code"`
		HouseNumber *string `json:"house_number"`
	} `json:"address"`
	PhoneNumber struct {
		AreaCode   string `json:"area_code"`
		MainNumber string `json:"main_number"`
	} `json:"phone_number"`
	MobileNumber string  `json:"mobile_number"`
	Verification string  `json:"verification"`
	Email        string  `json:"email"`
	InviteCode   string  `json:"invite_code"`
	Avatar       *string `json:"avatar"`
	Commission   int     `json:"commission"`
	Settings     struct {
		Theme              string   `json:"theme"`
		Mode               string   `json:"mode"`
		OrderSubmitConfirm bool     `json:"order_submit_confirm"`
		OrderDeleteConfirm bool     `json:"order_delete_confirm"`
		DefaultMode        bool     `json:"default_mode"`
		FavoriteMarkets    []string `json:"favorite_markets"`
		ChooseTradingType  bool     `json:"choose_trading_type"`
		CoinDeposit        bool     `json:"coin_deposit"`
		CoinWithdraw       bool     `json:"coin_withdraw"`
		MoneyDeposit       bool     `json:"money_deposit"`
		MoneyWithdraw      bool     `json:"money_withdraw"`
		Logins             bool     `json:"logins"`
		Trade              bool     `json:"trade"`
		APIKeyExpiration   bool     `json:"api_key_expiration"`
		Notification       struct {
			Email struct {
				IsEnable bool `json:"is_enable"`
				Actions  struct {
					CoinDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"coin_deposit"`
					CoinWithdraw struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"coin_withdraw"`
					MoneyDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"money_deposit"`
					MoneyWithdraw struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"money_withdraw"`
					Logins struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"logins"`
					APIKeyExpiration struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"api_key_expiration"`
					ManualDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"manual_deposit"`
				} `json:"actions"`
				Label string `json:"label"`
			} `json:"email"`
			Announcement struct {
				IsEnable bool `json:"is_enable"`
				Actions  struct {
					CoinDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"coin_deposit"`
					CoinWithdraw struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"coin_withdraw"`
					MoneyDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"money_deposit"`
					MoneyWithdraw struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"money_withdraw"`
					Logins struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"logins"`
					Trade struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"trade"`
					APIKeyExpiration struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"api_key_expiration"`
					ManualDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"manual_deposit"`
					PriceAlert struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"price_alert"`
				} `json:"actions"`
				Label string `json:"label"`
			} `json:"announcement"`
			Push struct {
				IsEnable bool `json:"is_enable"`
				Actions  struct {
					CoinDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"coin_deposit"`
					CoinWithdraw struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"coin_withdraw"`
					MoneyDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"money_deposit"`
					MoneyWithdraw struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"money_withdraw"`
					Logins struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"logins"`
					Trade struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"trade"`
					APIKeyExpiration struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"api_key_expiration"`
					ManualDeposit struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"manual_deposit"`
					PriceAlert struct {
						IsEnable bool   `json:"is_enable"`
						Label    string `json:"label"`
					} `json:"price_alert"`
				} `json:"actions"`
				Label string `json:"label"`
			} `json:"push"`
		} `json:"notification"`
	} `json:"settings"`
	Status struct {
		FirstName         string `json:"first_name"`
		LastName          string `json:"last_name"`
		NationalCode      string `json:"national_code"`
		NationalCardImage string `json:"national_card_image"`
		FaceImage         string `json:"face_image"`
		Birthday          string `json:"birthday"`
		Address           string `json:"address"`
		PhoneNumber       string `json:"phone_number"`
		MobileNumber      string `json:"mobile_number"`
		Email             string `json:"email"`
	} `json:"status"`
	KycInfo struct {
		Details struct {
			MobileActivation bool `json:"mobile_activation"`
			PersonalInfo     bool `json:"personal_info"`
			FinancialInfo    bool `json:"financial_info"`
			PhoneNumber      bool `json:"phone_number"`
			NationalCard     bool `json:"national_card"`
			FaceRecognition  bool `json:"face_recognition"`
			AdminApproval    bool `json:"admin_approval"`
		} `json:"details"`
		Level int `json:"level"`
	} `json:"kyc_info"`
	Meta struct {
		DisabledFeatures []string `json:"disabled_features"`
	} `json:"meta"`
}

// Profile retrieves account profile.
func (c *Client) Profile() (*Profile, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/profile", nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result *Profile `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// Balance represents holdings for an asset.
type Balance struct {
	Asset  string `json:"asset"`
	FaName string `json:"faName"`
	Fiat   bool   `json:"fiat"`
	Value  Number `json:"value"`
	Locked Number `json:"locked"`
}

// Balances retrieves a mapping between assets and their holdings.
func (c *Client) Balances() (map[string]*Balance, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/balances", nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result struct {
			Balances map[string]*Balance `json:"balances"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result.Balances, nil
}

// NamedFeeLevel represents a certain level on fees.
type NamedFeeLevel struct {
	MakerFee Number `json:"maker_fee"`
	TakerFee Number `json:"taker_fee"`
	Name     Number `json:"name"`
}

// FeeLevel contains the information for all fee levels and the current level.
type FeeLevel struct {
	Levels        map[Number]*NamedFeeLevel `json:"levels"`
	RecentDays    int                       `json:"recent_days"`
	RecentDaysSum Number                    `json:"recent_days_sum"`
	MakerFee      Number                    `json:"maker_fee"`
	TakerFee      Number                    `json:"taker_fee"`
	IsFixed       bool                      `json:"is_fixed"`
}

// FeeLevels retrieves a mapping between symbols and fee levels.
func (c *Client) FeeLevels() (map[string]*FeeLevel, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/fee", nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result map[string]*FeeLevel `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// BankingCard represents a banking card.
type BankingCard struct {
	ID         int      `json:"id"`
	CardNumber string   `json:"card_number"`
	Owners     []string `json:"owners"`
	Status     string   `json:"status"`
	IsDefault  int      `json:"is_default"`
}

// BankingCards retrieves a list of all user's banking cards.
func (c *Client) BankingCards() ([]*BankingCard, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/card-numbers", nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result []*BankingCard `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// BankAccount represents a bank account.
type BankAccount struct {
	ID          int      `json:"id"`
	IBAN        string   `json:"iban"`
	Owners      []string `json:"owners"`
	BankName    string   `json:"bank_name"`
	Status      string   `json:"status"`
	IsDefault   int      `json:"is_default"`
	BankDetails struct {
		Code  string `json:"code"`
		Label string `json:"label"`
	} `json:"bank_details"`
}

// BankAccounts retrieves a list of all user's bank accounts.
func (c *Client) BankAccounts() ([]*BankAccount, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/ibans", nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result []*BankAccount `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// -----------------------------------------------------------------------------
// Orders and trades
// -----------------------------------------------------------------------------

// List of order types.
const (
	OrderTypeLimit  = "LIMIT"
	OrderTypeMarket = "MARKET"
)

// List of order sides.
const (
	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"
)

// OrderParams is the request params to place an order.
type OrderParams struct {
	Symbol   string `json:"symbol"`
	Type     string `json:"type"`
	Side     string `json:"side"`
	Price    Number `json:"price"`
	Quantity Number `json:"quantity"`
	ClientID string `json:"client_id,omitempty"`
}

// Order represents a placed order.
type Order struct {
	Symbol          string    `json:"symbol"`
	Type            string    `json:"type"`
	Side            string    `json:"side"`
	Price           Number    `json:"price"`
	OrigQty         Number    `json:"origQty"`
	OrigSum         Number    `json:"origSum"`
	ExecutedPrice   *Number   `json:"executedPrice"`
	ExecutedQty     *Number   `json:"executedQty"`
	ExecutedSum     *Number   `json:"executedSum"`
	ExecutedPercent *Number   `json:"executedPercent"`
	Status          string    `json:"status"`
	Active          bool      `json:"active"`
	ClientOrderID   string    `json:"clientOrderId"`
	CreatedAt       time.Time `json:"created_at"`
}

// PlaceOrder places a new order.
func (c *Client) PlaceOrder(p *OrderParams) (*Order, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	body, _ := json.Marshal(p)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/v1/account/orders", bytes.NewReader(body))
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result *Order `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// CancelOrder cancels a placed order.
func (c *Client) CancelOrder(clientOrderID string) error {
	if c.apiKey == "" {
		return ErrMissingAPIKey
	}

	query := url.Values{}
	query.Add("clientOrderId", clientOrderID)

	req, err := http.NewRequest(http.MethodDelete, baseURL+"/v1/account/orders?"+query.Encode(), nil)
	if err != nil {
		return wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errNonOKResponse(resp.StatusCode)
	}

	return nil
}

// Order retrieves details for a placed order.
func (c *Client) Order(clientOrderID string) (*Order, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/orders/"+clientOrderID, nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result *Order `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result, nil
}

// OpenOrders retrievs a list of user's active orders.
// If symbol is empty, it retrieves active orders for all markets.
func (c *Client) OpenOrders(symbol string) ([]*Order, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	query := url.Values{}
	if symbol != "" {
		query.Add("symbol", symbol)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/openOrders?"+query.Encode(), nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result struct {
			Orders []*Order `json:"orders"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result.Orders, nil
}

// Trade represents a fulfilled trade.
type Trade struct {
	Symbol         string    `json:"symbol"`
	Quantity       Number    `json:"quantity"`
	Price          Number    `json:"price"`
	Sum            Number    `json:"sum"`
	Fee            Number    `json:"fee"`
	FeeCoefficient Number    `json:"feeCoefficient"`
	FeeAsset       string    `json:"feeAsset"`
	IsBuyer        bool      `json:"isBuyer"`
	Timestamp      time.Time `json:"timestamp"`
}

// Trades retrieves list of most recent user's trades.
// If symbol is empty, it retrieves trades for all markets.
// If side is empty, it retrieves trades for both sides.
func (c *Client) Trades(symbol, side string) ([]*Trade, error) {
	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	query := url.Values{}
	if symbol != "" {
		query.Add("symbol", symbol)
	}
	if side != "" {
		query.Add("side", side)
	}

	req, err := http.NewRequest(http.MethodGet, baseURL+"/v1/account/trades?"+query.Encode(), nil)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	req.Header.Add(apiKeyHeader, c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapRequestError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errNonOKResponse(resp.StatusCode)
	}

	result := struct {
		Result struct {
			AccountLatestTrades []*Trade `json:"AccountLatestTrades"`
		} `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, wrapRequestError(err)
	}

	return result.Result.AccountLatestTrades, nil
}

// -----------------------------------------------------------------------------
// Error handling
// -----------------------------------------------------------------------------

// List of common service errors.
var (
	ErrMissingAPIKey = &Error{Message: "missing api key"}
	ErrBadRequest    = &Error{Message: "bad request"}
	ErrUnauthorized  = &Error{Message: "unauthorized"}
	ErrForbidden     = &Error{Message: "access forbidden"}
	ErrNotFound      = &Error{Message: "resource not found"}
	ErrUnknown       = &Error{Message: "unknown error"}
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
