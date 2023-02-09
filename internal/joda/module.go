package joda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type ModuleSummary struct {
	Name     string `json:"Name"`
	Path     string `json:"Path"`
	Language string `json:"Language"`
	Type     string `json:"Type"`
}

type moduleResponse struct {
	Status *string `json:"status"`
	Error  *string `json:"error"`
}

func (j Joda) UploadModule(f multipart.File, h *multipart.FileHeader) error {

	// Buffer file
	var b bytes.Buffer
	// File writer
	w := multipart.NewWriter(&b)
	defer w.Close()
	// Add module file
	var fw io.Writer
	var err error
	if fw, err = w.CreateFormFile("module", h.Filename); err != nil {
		return err
	}
	// Write file
	if _, err := io.Copy(fw, f); err != nil {
		return err
	}
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	url := j.host + api + "/module"
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Could not read module response. Parse Error: " + err.Error())
		return errors.New("Could not read module response. Parse Error: " + err.Error())
	}

	response_message := new(moduleResponse)
	err = json.Unmarshal(body, &response_message)
	if err != nil {
		log.Println("Could not parse module response. Parse Error: " + err.Error())
		return errors.New("Could not parse module response. Parse Error: " + err.Error())
	}

	if response_message.Error != nil {
		return fmt.Errorf("error uploading module: %s", *response_message.Error)
	}

	return nil
}

func (j Joda) GetModules() ([]ModuleSummary, error) {
	endpoint := j.host + api + "/module"
	log.Println("Getting Modules (" + endpoint + ")")
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Println("Could not get modules. HTTP Error: " + err.Error())
		return nil, errors.New("Could not get modules. HTTP Error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Could not read modules. Parse Error: " + err.Error())
		return nil, errors.New("Could not read modules. Parse Error: " + err.Error())
	}

	modules := make([]ModuleSummary, 0)
	err = json.Unmarshal(body, &modules)
	if err != nil {
		log.Println("Could not parse modules. Parse Error: " + err.Error())
		return nil, errors.New("Could not parse modules. Parse Error: " + err.Error())
	}

	return modules, nil
}
