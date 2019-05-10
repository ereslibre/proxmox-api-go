package proxmox

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type StorageRef struct {
	Node    string
	Storage string
}

func (c *Client) Upload(storage StorageRef, filename, content string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("filename", filename)
	writer.WriteField("content", "snippets")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="filename"; filename="%s"`, filename))
	h.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(h)
	if err != nil {
		return err
	}

	part.Write([]byte(content))
	err = writer.Close()
	if err != nil {
		return err
	}
	headers := http.Header{}
	headers.Add("Content-Type", writer.FormDataContentType())
	request, err := c.session.NewRequest("POST", fmt.Sprintf("%s/nodes/%s/storage/%s/upload", c.ApiUrl, storage.Node, storage.Storage), &headers, body)
	if err != nil {
		return err
	}
	response, err := c.session.Do(request)
	if err != nil {
		return err
	}
	resp, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%v\n", string(resp))
	return err
}
