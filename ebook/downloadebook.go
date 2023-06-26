package ebook

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"errors"
	"sync"
)

// A function that generates an ebook from a url.
// @param url - the URL of the title page of the book.
// @param doPopulate - if true, download and populate the entire EbookInfo data structure, not just its metadata.
type EbookGeneratorFunction func(url string, doPopulate bool) (EbookInfo, error)

// Returned by a EbookGeneratorFunction when the URL can not be handled.
var UnsupportedUrlError = errors.New("unsupported url")

var (
	registerdFunctions      []EbookGeneratorFunction
	registerdFunctionsMutex sync.Mutex
)

// Register the given function.
func RegisterEbookGenerator(downloadFunction EbookGeneratorFunction) {
	registerdFunctionsMutex.Lock()
	registerdFunctions = append(registerdFunctions, downloadFunction)
	registerdFunctionsMutex.Unlock()
}

// Return the result of the first registered download function that does not return UnsupportedUrlError.
// @param url - the URL of the title page of the book.
// @param doPopulate - if true, download and populate the entire EbookInfo data structure, not just its metadata.
func DownloadEbook(url string, doPopulate bool) (EbookInfo, error) {
	registerdFunctionsMutex.Lock()
	fns := registerdFunctions
	registerdFunctionsMutex.Unlock()
	for _, fn := range fns {
		info, err := fn(url, doPopulate)
		if err != UnsupportedUrlError {
			return info, err
		}
	}
	return EbookInfo{}, UnsupportedUrlError
}
