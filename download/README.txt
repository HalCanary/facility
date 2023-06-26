package download // import "github.com/HalCanary/facility/download"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

FUNCTIONS

func Get(client *http.Client, requestUrl, referer, userAgent string) (io.ReadCloser, string, error)
func GetUrl(url, ref string, force bool) (io.ReadCloser, error)
    Fetch the content of a URL, using a cache if possible and if force is fakse.

func Post(client *http.Client, data url.Values, requestUrl, referer, userAgent string) (io.ReadCloser, string, error)
