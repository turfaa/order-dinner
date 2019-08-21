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

var (
	ErrNotOk  = errors.New("Not OK.")
	ErrNoFood = errors.New("No food to orders.")
)

func NewDinnerClient(ctx context.Context, ba, token string, mid int, fids []int) (Client, error) {
	c := &client{
		ctx:           ctx,
		baseApi:       ba,
		authorization: "Token " + token,
		menuID:        mid,
		foodIDs:       fids,
	}

	if err := c.HealthCheck(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*2)
	defer cancel()

	req, _ := http.NewRequest("GET", c.baseApi+"/menu/"+strconv.Itoa(c.menuID), nil)
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", c.authorization)

	hc := http.DefaultClient

	if res, err := hc.Do(req); err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		b := make([]byte, 100000)
		_, _ = res.Body.Read(b)

		return ErrNotOk
	}

	return nil
}

func (c *client) Order() error {
	var err error
	err = nil

	if len(c.foodIDs) == 0 {
		return ErrNoFood
	}

	url := c.baseApi + "/order/" + strconv.Itoa(c.menuID)
	reqd := orderRequest{}

	for _, fid := range c.foodIDs {
		ctxt, cancel := context.WithTimeout(c.ctx, time.Second*2)

		reqd.FoodID = fid
		json, err := jsoniter.Marshal(reqd)
		if err != nil {
			cancel()
			continue
		}

		req, _ := http.NewRequest("POST", url, bytes.NewReader(json))
		req = req.WithContext(ctxt)
		req.Header.Add("Authorization", c.authorization)
		req.Header.Add("Content-Type", "application/json")

		hc := http.DefaultClient

		if _, err = hc.Do(req); err == nil {
			cancel()
			break
		}

		cancel()
	}

	return err
}
