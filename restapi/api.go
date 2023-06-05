package restapi

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"test/dapi/internal"
	"time"
)

var timeout int

func fakeHTTPResponse() *http.Response {
	var response *http.Response      // Declare a pointer to http.Response
	fakeResponse := &http.Response{} // Create a new http.Response and take its address
	response = fakeResponse
	return response
}

func RequestDiscord(url, methode, bucket string, payload []byte, wantAnswer bool) BucketRequestAnswer {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  bucket,
		Url:         url,
		Methode:     methode,
		Payload:     payload,
		WantAnswer:  wantAnswer,
	}

	addRequest(request)

	response := <-answer
	close(answer)
	return response
}

func handleRateLimit(res *http.Response, b *Bucket) {
	Reset := res.Header.Get("X-RateLimit-Reset")
	internal.LogTrace(Reset)

	reset_number, err := strconv.ParseFloat(Reset, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(int64(reset_number), 0)

	Remaining := res.Header.Get("X-RateLimit-Remaining")
	BucketID := res.Header.Get("X-RateLimit-Bucket")
	BucketGlobal := res.Header.Get("X-RateLimit-Global")
	BucketScope := res.Header.Get("X-RateLimit-Scope")

	internal.LogDebug("GLOGAL", BucketGlobal)
	internal.LogDebug("SCOPE", BucketScope)

	remaining_n, err := strconv.ParseInt(Remaining, 10, 64)

	if err != nil {
		panic(err)
	}

	if remaining_n == 0 {
		b.lockBucket()
		internal.LogTrace("not remaining, should wait next time for reset.")
	}

	if b.BucketID != BucketID {
		internal.LogInfo("Bucket id changed for", b.Route, "from", b.BucketID, "to", BucketID)
	}

	b.Remaining = remaining_n
	b.ResetTime = tm
	b.BucketID = BucketID

	internal.LogTrace(tm)
	internal.LogTrace("rm", Remaining)
}

func httpDiscordCall(url, method string, bodyR []byte, wantAnswer bool) ([]byte, **http.Response, error) {

	urlFinal := "https://discord.com/api/" + url

	internal.LogTrace(urlFinal, url, method, string(bodyR))

	splitted := strings.Split(url, "/")

	payload := bytes.NewReader(bodyR)

	_bucket := Buckets.findBucket(splitted[1], method)

	if _bucket.Blocked {
		until := time.Until(_bucket.ResetTime)
		internal.LogTrace("blocked, waiting")

		sleep() // wait a little bit more than discord asked time value
		<-time.After(until)

		_bucket.unLockBucket()
		internal.LogTrace("waited")
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, urlFinal, payload)

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
		Buckets.saveBucket(_bucket)
		return httpDiscordCall(url, method, bodyR, wantAnswer)
	}

	handlingTimeDiscord = append(handlingTimeDiscord, time.Since(start))
	internal.LogInfo("Handled discord request in", time.Since(start).Milliseconds(), "medium: ", MediumValueAPI(), "ms")

	// everything was OK
	if !wantAnswer {
		return nil, &res, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		response := fakeHTTPResponse()
		return nil, &response, err
	}

	internal.LogDebug("HTTP CODE", res.StatusCode)
	internal.LogDebug("BUCKET", _bucket.Route, _bucket.Key)
	internal.LogDebug("REMAINING", _bucket.Remaining)

	handleRateLimit(res, &_bucket)

	internal.LogInfo("TIMEOUT >", timeout)

	Buckets.saveBucket(_bucket)

	return body, &res, err
}

func sleep() {
	time.Sleep(time.Millisecond * 10)
}

var handlingTimeDiscord = []time.Duration{0}
var handlingTimeBot = []time.Duration{0}

// MediumValueAPI return value in ms
func MediumValueAPI() int64 {
	var total int64
	for _, timed := range handlingTimeDiscord {
		total = total + timed.Milliseconds()
	}

	return total / int64(len(handlingTimeDiscord))
}

// MediumValueBOT return value in ms
func MediumValueBOT() int64 {
	var total int64

	for _, timed := range handlingTimeBot {
		total = total + timed.Milliseconds()
	}

	return total / int64(len(handlingTimeBot))
}

func MediumValueWithoutAPI() int64 {
	mediumApi := MediumValueAPI()
	mediumBot := MediumValueBOT()

	return mediumBot - mediumApi
}

func HandleBucket(b Bucket) {
	for {
		for c, v := range b.Requests {
			start := time.Now()

			body, res, err := httpDiscordCall(v.Url, v.Methode, v.Payload, v.WantAnswer)
			answer := BucketRequestAnswer{Body: body, Res: **res, Err: err}
			v.AnswerQueue <- answer
			delete(b.Requests, c)

			handlingTimeBot = append(handlingTimeBot, time.Since(start))
			internal.LogInfo("handled total request and answered in", time.Since(start).Milliseconds(), "ms, medium: ", MediumValueBOT(), "ms")
			internal.LogInfo("handled without discord medium:", MediumValueWithoutAPI(), "ms")
		}
		sleep()
	}

}

func DiscordAPIHandler() ([]byte, *http.Response) {
	for {
		if len(Buckets) == 0 {
			sleep()
			continue
		}

		for _, v := range Buckets {
			if !v.Handled {
				go HandleBucket(v)
				v.handled()
				v.save(Buckets)
				internal.LogDebug("Handing new bucket", v.Index, v.Method)
			}
		}
		sleep()
	}
}
