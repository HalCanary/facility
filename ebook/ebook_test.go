package ebook

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/HalCanary/facility/dom"
)

var testStrings = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Feugiat sed lectus vestibulum mattis ullamcorper velit sed. Et leo duis ut diam quam nulla porttitor massa. Aliquam malesuada bibendum arcu vitae elementum curabitur vitae nunc. Dictum at tempor commodo ullamcorper a lacus vestibulum sed. Suspendisse faucibus interdum posuere lorem ipsum dolor. Habitasse platea dictumst quisque sagittis. Semper risus in hendrerit gravida rutrum quisque non tellus. Sollicitudin tempor id eu nisl. Auctor augue mauris augue neque. Pulvinar etiam non quam lacus suspendisse. Duis ultricies lacus sed turpis tincidunt id aliquet risus. Congue quisque egestas diam in arcu cursus euismod quis viverra. At risus viverra adipiscing at in tellus integer feugiat scelerisque.",
	"Eu nisl nunc mi ipsum faucibus vitae aliquet. Donec enim diam vulputate ut pharetra sit. Amet consectetur adipiscing elit ut aliquam purus sit amet luctus. Amet tellus cras adipiscing enim eu. Vehicula ipsum a arcu cursus vitae congue mauris rhoncus. Nisl pretium fusce id velit ut tortor pretium. Diam maecenas ultricies mi eget mauris pharetra. Varius quam quisque id diam vel quam elementum pulvinar. Vitae proin sagittis nisl rhoncus mattis rhoncus urna. Tristique risus nec feugiat in fermentum posuere urna nec. Pharetra convallis posuere morbi leo urna molestie. Condimentum id venenatis a condimentum vitae sapien pellentesque. Sed augue lacus viverra vitae congue. Tellus at urna condimentum mattis pellentesque. Sit amet dictum sit amet justo donec enim diam vulputate. Fringilla urna porttitor rhoncus dolor purus non. Id interdum velit laoreet id donec ultrices tincidunt arcu. Ultricies tristique nulla aliquet enim.",
	"Congue quisque egestas diam in arcu. Sit amet purus gravida quis blandit turpis cursus. Elit at imperdiet dui accumsan. Odio euismod lacinia at quis risus. Potenti nullam ac tortor vitae purus faucibus ornare. Donec adipiscing tristique risus nec feugiat in fermentum posuere urna. Nunc faucibus a pellentesque sit amet porttitor eget. Egestas sed tempus urna et pharetra pharetra massa massa ultricies. Ipsum a arcu cursus vitae congue mauris rhoncus. Amet facilisis magna etiam tempor orci eu lobortis elementum nibh. Eget felis eget nunc lobortis mattis aliquam. Pellentesque habitant morbi tristique senectus et netus. Suspendisse potenti nullam ac tortor vitae purus faucibus ornare suspendisse. Nisi lacus sed viverra tellus in hac habitasse platea dictumst. Metus vulputate eu scelerisque felis imperdiet. Nibh sit amet commodo nulla facilisi nullam vehicula ipsum a.",
	"In est ante in nibh mauris cursus. Etiam erat velit scelerisque in dictum non consectetur a. Magna ac placerat vestibulum lectus mauris ultrices eros in cursus. Diam maecenas ultricies mi eget mauris. Sit amet justo donec enim diam vulputate ut pharetra sit. Eget gravida cum sociis natoque penatibus et. Ut faucibus pulvinar elementum integer enim neque volutpat ac tincidunt. Enim ut sem viverra aliquet. Fermentum odio eu feugiat pretium. Purus faucibus ornare suspendisse sed nisi lacus sed. Adipiscing diam donec adipiscing tristique risus nec feugiat in fermentum. Odio pellentesque diam volutpat commodo sed egestas egestas fringilla. Vulputate enim nulla aliquet porttitor lacus luctus accumsan. Urna neque viverra justo nec ultrices dui sapien eget mi. Nisl nunc mi ipsum faucibus vitae aliquet nec. Lacus vel facilisis volutpat est. Porttitor eget dolor morbi non arcu risus quis varius quam.",
	"Sit amet nulla facilisi morbi tempus iaculis urna. Nisi est sit amet facilisis magna etiam tempor orci. A condimentum vitae sapien pellentesque habitant morbi. Suspendisse interdum consectetur libero id faucibus nisl tincidunt eget nullam. Magna fringilla urna porttitor rhoncus dolor purus non. Ullamcorper eget nulla facilisi etiam dignissim diam quis enim. Velit ut tortor pretium viverra suspendisse potenti nullam. Ac placerat vestibulum lectus mauris. Tincidunt dui ut ornare lectus sit. Fermentum iaculis eu non diam phasellus vestibulum. Tempor orci dapibus ultrices in iaculis nunc sed. Nunc sed augue lacus viverra. Vitae tempus quam pellentesque nec nam. Sit amet consectetur adipiscing elit.",
	"Amet cursus sit amet dictum sit amet justo. Pretium vulputate sapien nec sagittis aliquam malesuada bibendum. Mi ipsum faucibus vitae aliquet nec ullamcorper. Et malesuada fames ac turpis. Arcu bibendum at varius vel pharetra vel. Lobortis feugiat vivamus at augue eget. Scelerisque eleifend donec pretium vulputate sapien nec. Arcu odio ut sem nulla pharetra. Dolor sed viverra ipsum nunc aliquet bibendum enim facilisis gravida. Et leo duis ut diam. Pulvinar sapien et ligula ullamcorper malesuada proin. Nulla facilisi nullam vehicula ipsum a.",
	"A pellentesque sit amet porttitor eget. Nec ullamcorper sit amet risus nullam eget felis eget. Est ante in nibh mauris. Integer eget aliquet nibh praesent tristique. Nisl purus in mollis nunc. Fringilla est ullamcorper eget nulla facilisi. Sit amet mauris commodo quis imperdiet massa tincidunt nunc pulvinar. Elit scelerisque mauris pellentesque pulvinar. Urna porttitor rhoncus dolor purus non enim praesent. Bibendum ut tristique et egestas quis. Eu tincidunt tortor aliquam nulla facilisi cras fermentum odio. Arcu risus quis varius quam quisque id diam vel. Proin sed libero enim sed. Faucibus turpis in eu mi bibendum. Pellentesque adipiscing commodo elit at imperdiet dui accumsan. Morbi tempus iaculis urna id volutpat lacus. Pellentesque eu tincidunt tortor aliquam. Viverra suspendisse potenti nullam ac tortor vitae purus.",
	"Lacinia at quis risus sed vulputate. Vel pretium lectus quam id leo in vitae. Nulla facilisi cras fermentum odio eu. Bibendum neque egestas congue quisque egestas diam in arcu cursus. Sit amet mauris commodo quis imperdiet massa tincidunt nunc pulvinar. Imperdiet proin fermentum leo vel orci porta non. Sed adipiscing diam donec adipiscing tristique risus nec feugiat in. Sodales ut etiam sit amet nisl purus. Rutrum tellus pellentesque eu tincidunt. Laoreet non curabitur gravida arcu ac.",
	"Pellentesque nec nam aliquam sem et tortor consequat id. Pretium vulputate sapien nec sagittis aliquam malesuada bibendum arcu vitae. Id volutpat lacus laoreet non curabitur gravida arcu ac tortor. Enim blandit volutpat maecenas volutpat blandit aliquam etiam erat. Semper risus in hendrerit gravida rutrum quisque non. Vulputate enim nulla aliquet porttitor lacus luctus accumsan tortor. Nisi quis eleifend quam adipiscing. Et ultrices neque ornare aenean euismod elementum nisi quis. Mauris sit amet massa vitae tortor condimentum lacinia quis vel. Vitae congue eu consequat ac felis donec et odio. Nibh mauris cursus mattis molestie. At elementum eu facilisis sed odio. Ultrices neque ornare aenean euismod elementum nisi quis eleifend quam. Nunc sed augue lacus viverra vitae.",
	"In eu mi bibendum neque egestas congue quisque. Laoreet sit amet cursus sit amet dictum sit amet justo. Dignissim enim sit amet venenatis urna cursus eget nunc scelerisque. Ut placerat orci nulla pellentesque dignissim enim. Praesent elementum facilisis leo vel. Penatibus et magnis dis parturient montes. Nunc mattis enim ut tellus elementum. Condimentum mattis pellentesque id nibh tortor id aliquet lectus. Morbi tincidunt augue interdum velit euismod in pellentesque. Malesuada proin libero nunc consequat. Facilisis volutpat est velit egestas dui id ornare. Ipsum dolor sit amet consectetur adipiscing elit ut aliquam purus.",
}

func TestEbook(t *testing.T) {
	now := time.Now()
	div := dom.Elem("div")
	for _, s := range testStrings {
		dom.Append(div, dom.Elem("p", dom.Text(s)))
	}
	ebook := EbookInfo{
		Authors:  "The Author",
		Comments: "Some Comments",
		Title:    "the Title",
		Source:   "Lorem Ipsum Generator",
		Language: "en",
		Chapters: []Chapter{
			Chapter{
				Title:    "Chapter One",
				Content:  div,
				Modified: now,
			},
		},
		Modified: now,
		Cover:    nil,
	}
	var buffer bytes.Buffer
	err := ebook.Write(&buffer)
	if err != nil {
		t.Error(err)
	} else {
		os.WriteFile("/tmp/ebooktest.epub", buffer.Bytes(), 0o644)
	}

}
