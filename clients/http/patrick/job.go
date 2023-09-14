//Fork Author: Josh Durham
//Date Created: 09/14/2023
//Date Updated: 09/14/2023 - Josh Durham
//Purpose: provides a client for interacting with a job management service. 
//	   It allows fetching jobs, job details, job logs, canceling jobs, and retrying jobs by specific ID's. 

//Possible Changes to Code Base: I believe in order to create more maintainable, readable, and better reusable code, 
// 				 ensuring each file as a header which has the date the file was created on, the author, 
//				 date updated, and a description of what the file does.
//------------------------------------------------------------------------------------------------------------------------------


// The package "patrick" provides a client library for interacting with jobs.
package patrick

import (
	"fmt"
	"io"
	"net/http"
	patrickIface "github.com/taubyte/go-interfaces/services/patrick"
)

type data struct {
	ProjectId string
	JobIds    []string
}

// Fetches a list of job IDs for a given project.
func (c *Client) Jobs(projectId string) (jobList []string, err error) {
	var jobs data
	url := "/jobs/" + projectId
	// Makes an HTTP Get request to obtain the job ID for a given project.
	if err = c.http.Get(url, &jobs); err != nil {
		err = fmt.Errorf("failed getting jobs for project `%s` with: %w", projectId, err)
		return
	}

	return jobs.JobIds, nil
}

// Fetches details of a job by the ID.
func (c *Client) Job(jid string) (job *patrickIface.Job, err error) {
	receive := &struct {
		Job patrickIface.Job
	}{}
	url := "/job/" + jid
	// Makes an HTTP Get request to obtain job information.
	if err = c.http.Get(url, &receive); err != nil {
		err = fmt.Errorf("failed getting job `%s` with: %w", jid, err)
		return
	}

	return &receive.Job, nil
}

// Fetches the log information of the job.
func (c *Client) LogFile(jobId, resourceId string) (log io.ReadCloser, err error) {
	method := http.MethodGet
	path := "/logs" + "/" + resourceId

	// Makes an HTTP Get request to obtain log information.
	req, err := http.NewRequestWithContext(c.http.Context(), method, c.http.Url()+path, nil)
	if err != nil {
		err = fmt.Errorf("%s -- `%s` failed with %s", method, path, err.Error())
		return
	}

	
	//This seems to get some sort of HTTP authorization token.
	req.Header.Add("Authorization", c.http.AuthHeader())
	resp, err := c.http.Client().Do(req)
	if err != nil {
		err = fmt.Errorf("%s -- `%s` failed with %s", method, path, err.Error())
		return
	}

	// Closes when the operation is done
	go func() {
		<-c.http.Context().Done()
		resp.Body.Close()
	}()

	
	//Try catch to determine if the HTTP request was accepted
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s -- `%s` failed with status: %s", method, path, resp.Status)
		return
	}

	return resp.Body, nil
}

// Cancels a job by its ID
func (c *Client) Cancel(jid string) (response interface{}, err error) {
	url := "/cancel/" + jid
	// Make an HTTP POST request to cancel a job by ID
	if err = c.http.Post(url, nil, &response); err != nil {
		err = fmt.Errorf("failed getting job `%s` with: %w", jid, err)
		return
	}

	return
}

// Retries a job by its ID.
func (c *Client) Retry(jid string) (response interface{}, err error) {
	url := "/retry/" + jid
	// Make an HTTP POST request to retry a job
	if err = c.http.Post(url, nil, &response); err != nil {
		err = fmt.Errorf("failed getting job `%s` with: %w", jid, err)
		return
	}

	return
}
