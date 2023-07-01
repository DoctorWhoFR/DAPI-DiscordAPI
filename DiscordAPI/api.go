/*
Package DiscordAPI

Package used to manage to Discord REST API part.

In fact most of the function is linked to a specific Type,
there is not so much function that you are going to call directly.

If you want to directly call discord API, you can use the in-build function RequestDiscord, like :

	response := RequestDiscord("/channels/"+channelId, http.MethodGet, "channels", nil, true)

Or:

	answer := RequestDiscord("/channels/"+content.ID+"/messages", http.MethodPost, "channels", body, true)
*/
package DiscordAPI

import (
	"azginfr/dapi/DiscordInternal"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var timeout int

/*
fakeHTTPResponse

Basic in-build BotInternal function to simply generate a new fake http.Response for error handling purpose in other function.
*/
func fakeHTTPResponse() *http.Response {
	var response *http.Response      // Declare a pointer to http.Response
	fakeResponse := &http.Response{} // Create a new http.Response and take its address
	response = fakeResponse
	return response
}

/*
RequestDiscordForm

Function used to send a new request to the Discord Go channel managing the discord call.

You should use this function instead of directly send message to channel.

Example for FormData payload:

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part1, _ := writer.CreateFormFile("files[0]", filepath.Base(fileContent.Name()))
	_, _ = io.Copy(part1, fileContent)
	_ = writer.WriteField("payload_json", string(messageJson))
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}
*/
func RequestDiscordForm(url, methode, bucket string, payload []byte, wantAnswer bool, buffer *bytes.Buffer, contentType string) BucketRequestAnswer {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue:   answer,
		BucketName:    bucket,
		Url:           url,
		Methode:       methode,
		Payload:       payload,
		WantAnswer:    wantAnswer,
		FormWriter:    buffer,
		FormHTTP:      true,
		ContentHeader: contentType,
	}

	addRequest(request)

	response := <-answer
	close(answer)
	return response
}

/*
RequestDiscord

Function used to send a new request to the Discord Go channel managing the discord call.

You should use this function instead of directly send message to channel.

Example :

	test := struct {
		WebhookChannelId string `json:"webhook_channel_id"`
	}{WebhookChannelId: channelId}

	body, err := json.Marshal(test)

	if err != nil {
		panic(err)
	}

	answer := RequestDiscord(fmt.Sprintf("/channels/%s/followers", content.ID), http.MethodPost, "channels", body, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
*/
func RequestDiscord(url, methode, bucket string, payload []byte, wantAnswer bool) BucketRequestAnswer {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  bucket,
		Url:         url,
		Methode:     methode,
		Payload:     payload,
		WantAnswer:  wantAnswer,
		FormHTTP:    false,
	}

	addRequest(request)

	response := <-answer
	close(answer)
	return response
}

/*
handleRateLimit
*/
func handleRateLimit(res *http.Response, b *Bucket) {
	Reset := res.Header.Get("X-RateLimit-Reset")
	DiscordInternal.LogTrace(Reset)

	resetNumber, err := strconv.ParseFloat(Reset, 64)
	if err != nil {
		log.Println("error parsing ratelimit")
		return
	}
	tm := time.Unix(int64(resetNumber), 0)

	Remaining := res.Header.Get("X-RateLimit-Remaining")
	BucketID := res.Header.Get("X-RateLimit-Bucket")
	BucketGlobal := res.Header.Get("X-RateLimit-Global")
	BucketScope := res.Header.Get("X-RateLimit-Scope")

	DiscordInternal.LogDebug("GLOGAL", BucketGlobal)
	DiscordInternal.LogDebug("SCOPE", BucketScope)

	remaining_n, err := strconv.ParseInt(Remaining, 10, 64)

	if err != nil {
		panic(err)
	}

	if remaining_n == 0 {
		b.lockBucket()
		DiscordInternal.LogTrace("not remaining, should wait next time for reset.")
	}

	if b.BucketID != BucketID {
		DiscordInternal.LogInfo("Bucket id changed for", b.Route, "from", b.BucketID, "to", BucketID)
	}

	b.Remaining = remaining_n
	b.ResetTime = tm
	b.BucketID = BucketID

	DiscordInternal.LogTrace(tm)
	DiscordInternal.LogTrace("rm", Remaining)
}

/*
httpDiscordCallFormData

Internal function to call Discord, with some in-build features like rate-limit handling.

There is two version of this function, httpDiscordCallJson and httpDiscordCallFormData one.

This is FormData version, where DiscordPayload is a bytes.Buffer representation of a http FormWriter payload.

You can use directly this function (or FormData) one, but you can also (and should) use RequestDiscord function
to use in-build rate-limit and other security features.
*/
func httpDiscordCallFormData(url, method string, wantAnswer bool, body *bytes.Buffer, contentHeader string) ([]byte, **http.Response, error) {
	urlFinal := "https://discord.com/api/" + url

	DiscordInternal.LogTrace(urlFinal, url, method, body)

	splitted := strings.Split(url, "/")

	_bucket := FindBucket(splitted[1], method)

	if _bucket.Blocked {
		until := time.Until(_bucket.ResetTime)
		DiscordInternal.LogTrace("blocked, waiting")

		DiscordInternal.SimpleSleep() // wait a little bit more than discord asked time value
		<-time.After(until)

		_bucket.unLockBucket()
		DiscordInternal.LogTrace("waited")
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, urlFinal, body)

	if err != nil {
		response := fakeHTTPResponse()
		return nil, &response, err
	}

	req.Header.Add("Authorization", "Bot "+os.Getenv("TOKEN"))
	req.Header.Add("Content-Type", contentHeader)

	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing")
		}
	}(res.Body)

	if res.StatusCode == 429 {
		timeout++
		_bucket.lockBucket()
		SaveBucket(_bucket)
		return httpDiscordCallFormData(url, method, wantAnswer, body, contentHeader)
	}

	DiscordInternal.HandlingTimeDiscord = append(DiscordInternal.HandlingTimeDiscord, time.Since(start))
	DiscordInternal.LogInfo("IsHandled discord request in", time.Since(start).Milliseconds(), "medium: ", DiscordInternal.MediumValueAPI(), "ms")

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		response := fakeHTTPResponse()
		return nil, &response, err
	}
	DiscordInternal.LogDebug(wantAnswer)

	DiscordInternal.LogDebug("HTTP CODE", res.StatusCode)
	DiscordInternal.LogDebug("BUCKET", _bucket.Route, _bucket.Key)
	DiscordInternal.LogDebug("REMAINING", _bucket.Remaining)

	handleRateLimit(res, &_bucket)

	DiscordInternal.LogInfo("TIMEOUT >", timeout)

	SaveBucket(_bucket)

	return responseBody, &res, err
}

/*
httpDiscordCallJson

Internal function to call Discord, with some in-build features like rate-limit handling.

There is two version of this function, httpDiscordCallJson and httpDiscordCallFormData one.

This is Json version, where DiscordPayload is a []byte representation of a discord json payload.

You can use directly this function (or FormData) one, but you can also (and should) use RequestDiscord function
to use in-build rate-limit and other security features.
*/
func httpDiscordCallJson(DiscordEndpoint, DiscordMethod string, DiscordPayload []byte, wantAnswer bool) ([]byte, **http.Response, error) {

	urlFinal := "https://discord.com/api/" + DiscordEndpoint

	DiscordInternal.LogTrace("HTTP CALL", urlFinal, DiscordEndpoint, DiscordMethod, string(DiscordPayload))

	splitted := strings.Split(DiscordEndpoint, "/")

	payload := bytes.NewReader(DiscordPayload)

	_bucket := FindBucket(splitted[1], DiscordMethod)

	if _bucket.Blocked {
		until := time.Until(_bucket.ResetTime)
		DiscordInternal.LogTrace("blocked, waiting")

		DiscordInternal.SimpleSleep() // wait a little bit more than discord asked time value
		<-time.After(until)

		_bucket.unLockBucket()
		DiscordInternal.LogTrace("waited")
	}

	client := &http.Client{}

	req, err := http.NewRequest(DiscordMethod, urlFinal, payload)

	if err != nil {
		response := fakeHTTPResponse()
		return nil, &response, err
	}

	req.Header.Add("Authorization", "Bot "+os.Getenv("TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing")
		}
	}(res.Body)

	if res.StatusCode == 429 {
		timeout++
		_bucket.lockBucket()
		SaveBucket(_bucket)
		return httpDiscordCallJson(DiscordEndpoint, DiscordMethod, DiscordPayload, wantAnswer)
	}

	DiscordInternal.HandlingTimeDiscord = append(DiscordInternal.HandlingTimeDiscord, time.Since(start))
	DiscordInternal.LogInfo("IsHandled discord request in", time.Since(start).Milliseconds(), "medium: ", DiscordInternal.MediumValueAPI(), "ms")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		response := fakeHTTPResponse()
		return nil, &response, err
	}

	DiscordInternal.LogDebug(wantAnswer)

	DiscordInternal.LogDebug("HTTP CODE", res.StatusCode)
	DiscordInternal.LogDebug("BUCKET", _bucket.Route, _bucket.Key)
	DiscordInternal.LogDebug("REMAINING", _bucket.Remaining)

	handleRateLimit(res, &_bucket)

	DiscordInternal.LogInfo("TIMEOUT >", timeout)

	SaveBucket(_bucket)

	return body, &res, err
}

/*
NewBucketHandler

Internal function used to start the goroutine for each new discord Bucket.

If no goroutine is find wait a simple time, using DiscordInternal.SimpleSleep
*/
func NewBucketHandler() ([]byte, *http.Response) {
	for {
		GlobalDiscordBuckets.Range(func(key, value any) bool {
			v := value.(Bucket)

			if !v.IsHandled {

				go HandleBucketRequest(v)
				v.handled()
				SaveBucket(v)

				DiscordInternal.LogDebug("Handing new bucket", v.Index, v.Method)
			}
			return false
		})
		DiscordInternal.SimpleSleep()
	}
}
