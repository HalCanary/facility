package download // import "github.com/HalCanary/facility/download"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

FUNCTIONS

func Get(client *http.Client, requestUrl, referer, userAgent string) (io.ReadCloser, string, error)
    Call `client.Do(http.NewRequest())`, but with extra steps. Set `User-Agent`
    and `Referer` (if set). If Error status code is returned, close body and
    return `error`.

func GetUrl(url, ref string, force bool) (io.ReadCloser, error)
    Fetch the content of a URL, using a cache if possible and if force is false.

func Post(client *http.Client, data url.Values, requestUrl, referer, userAgent string) (io.ReadCloser, string, error)
    Call `client.Do(http.NewRequest())`, but with extra steps. Make sure
    `Content-Type` is `application/x-www-form-urlencoded` Set `User-Agent` and
    `Referer` (if set). If Error status code is returned, close body and return
    `error`.

