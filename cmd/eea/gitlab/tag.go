package gitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/l10178/eea-cli/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AddTagReq struct {
	ProjectId string
	Ref       string
	TagName   string
	Message   string
}

type AddTagResponse struct {
	Name   string
	Target string
}

type RepositoryTag struct {
	Name    string
	Message string
	Target  string
	Commit  Commit `json:"commit"`
}

type Commit struct {
	Id      string
	ShortId string `json:"short_id"`
	Message string
}

func AddTag(req *AddTagReq) (AddTagResponse, error) {
	var result AddTagResponse
	gitlab := config.GetConfig().GitLab

	id := url.QueryEscape(req.ProjectId)
	gitApiTags := gitlab.ApiRoot + "/projects/" + id + "/repository/tags"

	httpReq, err := http.NewRequest("POST", gitApiTags, nil)
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	query := httpReq.URL.Query()
	query.Add("tag_name", req.TagName)
	query.Add("ref", req.Ref)
	query.Add("message", req.Message)
	httpReq.URL.RawQuery = query.Encode()
	httpReq.Header.Add("PRIVATE-TOKEN", gitlab.PrivateToken)

	var response *http.Response
	response, err = http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		err = json.Unmarshal(body, &result)
	} else {
		err = errors.New(response.Status)
		log.Printf("[%s] %s", req.ProjectId, string(body))
	}
	return result, err
}

func BatchTag(file string, tag string, ref string) error {

	// read all projects from the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	gitlab := config.GetConfig().GitLab

	// one projects a line
	slices := strings.Split(string(data), "\n")

	// count errors
	errs := 0

	for _, line := range slices {

		//ignore empty line or commented
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// the line is `group:project-name:master`
		sls := strings.Split(line, ":")

		// id = group + name or a number id
		id := strings.TrimSpace(sls[1])

		if strings.TrimSpace(sls[0]) != "" {
			id = strings.TrimSpace(sls[0]) + "/" + strings.TrimSpace(sls[1])
		}

		//use file ref if the `ref` is empty
		if ref == "" && len(sls) > 2 {
			ref = strings.TrimSpace(sls[2])
		}

		//yeah, use tag as message
		msg := tag

		req := &AddTagReq{
			ProjectId: id,
			TagName:   tag,
			Ref:       ref,
			Message:   msg,
		}
		_, err = AddTag(req)

		if err != nil {
			errs += 1
			continue
		}
		//wait for tag web hook
		duration, _ := time.ParseDuration(gitlab.TagSleepSeconds)
		time.Sleep(duration)
	}

	if errs > 0 {
		return errors.New("error")
	}
	return nil
}

func GetRepositoryTag(projectId string, tag string) (RepositoryTag, error) {
	var result RepositoryTag
	gitlab := config.GetConfig().GitLab

	id := url.QueryEscape(projectId)

	//GET /projects/:id/repository/tags/:tag_name
	gitApiTag := gitlab.ApiRoot + "/projects/" + id + "/repository/tags/" + tag

	httpReq, err := http.NewRequest("GET", gitApiTag, nil)
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	query := httpReq.URL.Query()
	httpReq.URL.RawQuery = query.Encode()
	httpReq.Header.Add("PRIVATE-TOKEN", gitlab.PrivateToken)

	var response *http.Response
	response, err = http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		err = json.Unmarshal(body, &result)
	} else {
		err = errors.New(response.Status)
		log.Printf("[%s] %s", projectId, string(body))
	}
	return result, err
}

func BatchCommit(file string, tag string) error {
	gitlab := config.GetConfig().GitLab
	// read all projects from the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// one projects a line
	slices := strings.Split(string(data), "\n")

	// count errors
	errs := 0

	for _, line := range slices {

		//ignore empty line or commented
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// the line is `group:project-name:master`
		sls := strings.Split(line, ":")

		// id = group + name or a number id
		name := strings.TrimSpace(sls[1])
		id := name

		if strings.TrimSpace(sls[0]) != "" {
			id = strings.TrimSpace(sls[0]) + "/" + name
		}

		resp, err := GetRepositoryTag(id, tag)

		if err != nil {
			errs += 1
			continue
		}
		// full git url
		gitUrl := gitlab.GitRoot + "/" + id + ".git"
		fmt.Printf("%s,%s,%s\n", name, gitUrl, resp.Commit.Id)
	}

	if errs > 0 {
		return errors.New("error")
	}
	return nil
}
