package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	client     *http.Client
	baseURL    string
	credential string
	nonce      string
}

func (c *Client) Login(password string) error {
	arg := base64.URLEncoding.EncodeToString([]byte("admin:" + password))
	c.nonce = "1337"
	resp, err := c.client.Get(c.baseURL + "/login?arg=" + arg + "&_n=" + c.nonce)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.credential = string(body)
	if c.credential == "" {
		return fmt.Errorf("no credential cookie in the response")
	}

	return nil
}

func (c Client) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Cookie", "credential="+c.credential)
	resp, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func (c Client) SNMPGet(oids []string) ([]byte, error) {
	url := c.baseURL + "/snmpGet?oids="
	for _, oid := range oids {
		url += oid + ";"
	}
	url += "&_n=" + c.nonce
	return c.get(url)
}

func (c Client) SNMPSet(oid string, value string, dataType string) ([]byte, error) {
	url := c.baseURL + "/snmpSet?oids="
	url += oid + "=" + value + ";" + dataType + ";"
	url += "&_n=" + c.nonce
	return c.get(url)
}

func (c Client) SNMPWalk(oids []string) ([]byte, error) {
	url := c.baseURL + "/walk?oids="
	for _, oid := range oids {
		url += oid + ";"
	}
	url += "&_n=" + c.nonce
	return c.get(url)
}
