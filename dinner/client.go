package dinner

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"bytes"

	"github.com/json-iterator/go"
)

type client struct {
	ctx           context.Context
	baseApi       string
	authorization string

	menuID  int
	foodIDs []int
}

type orderRequest struct {
	FoodID int `json:"food_id"`
}

type menu struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type menuResponse struct {
	response
	Menu menu `json:"menu"`
}

var (
	ErrNoFood = errors.New("No food to orders.")
)

func NewDinnerClient(ctx context.Context, ba, token string, fids []int) (Client, error) {
	c := &client{
		ctx:           ctx,
		baseApi:       ba,
		authorization: "Token " + token,
		menuID:        -1,
		foodIDs:       fids,
	}

	if err := c.UpdateMenu(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) UpdateMenu() error {
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*2)
	defer cancel()

	req, _ := http.NewRequest("GET", c.baseApi+"/current", nil)
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", c.authorization)

	hc := http.DefaultClient

	if res, err := hc.Do(req); err != nil {
		return err
	} else {
		var mres menuResponse

		b := make([]byte, 100000)
		_, _ = res.Body.Read(b)

		if err := jsoniter.Unmarshal(b, &mres); err != nil {
			return err
		} else if mres.Status == "success" {
			c.menuID = mres.Menu.ID
		} else {
			return errors.New(mres.Error)
		}
	}

	return nil
}

func (c *client) IsReady() bool {
	return c.menuID != -1
}

func (c *client) Order() error {
	var err error
	var res *http.Response
	err = nil

	if len(c.foodIDs) == 0 {
		return ErrNoFood
	}

	url := c.baseApi + "/order/" + strconv.Itoa(c.menuID)
	reqd := orderRequest{}

	for _, fid := range c.foodIDs {
		var json []byte
		ctxt, cancel := context.WithTimeout(c.ctx, time.Second*2)

		reqd.FoodID = fid
		json, err = jsoniter.Marshal(reqd)
		if err != nil {
			cancel()
			continue
		}

		req, _ := http.NewRequest("POST", url, bytes.NewReader(json))
		req = req.WithContext(ctxt)
		req.Header.Add("Authorization", c.authorization)
		req.Header.Add("Content-Type", "application/json")

		hc := http.DefaultClient

		res, err = hc.Do(req)
		if err == nil {
			if res.StatusCode == http.StatusOK {
				cancel()
				break
			} else {
				var r response

				b := make([]byte, 100000)
				_, _ = res.Body.Read(b)

				err = jsoniter.Unmarshal(b, &r)
				if err == nil && r.Status != "success" {
					err = errors.New(r.Error)
				}
			}
		}

		cancel()
	}

	return err
}
