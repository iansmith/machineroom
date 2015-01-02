package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	s5 "github.com/seven5/seven5/client"

	"github.com/igneous-systems/beta/shared"
)

//typed constants
var (
	currentParent = s5.NewHtmlId("div", "current-parent")
	submit        = s5.NewHtmlId("button", "send-params")
	errText       = s5.NewHtmlId("h5", "err-text")
	errRegion     = s5.NewHtmlId("div", "err-region")
	usernameInput = s5.NewHtmlId("input", "username")
	passwordInput = s5.NewHtmlId("input", "pwd")

	row     = s5.NewCssClass("row")
	col2    = s5.NewCssClass("col-sm-2")
	col6    = s5.NewCssClass("col-sm-6")
	col4    = s5.NewCssClass("col-sm-4")
	offset2 = s5.NewCssClass("col-sm-offset-2")
	gray    = s5.NewCssClass("gray")
	bold    = s5.NewCssClass("bold")
	warning = s5.NewCssClass("bg-warning")
)

type mainPage struct {
	//page level state
	usernameCurrent s5.StringAttribute
	passwordCurrent s5.StringAttribute
	usernameDisplay s5.StringAttribute
	passwordDisplay s5.StringAttribute
}

func slideDown(message string, clazz string) {
	errText.Dom().SetText(message)

	//xxx should be doing this with a constant, not a string
	jquery.NewJQuery("div#err-region").Underlying().Call("slideDown", "slow")

	errRegion.Dom().AddClass(clazz /*why isn't this typed?*/)

	js.Global.Call("setTimeout", func() {
		jquery.NewJQuery("div#err-region").Underlying().Call("slideUp", "slow")
		errRegion.Dom().RemoveClass(clazz /*why isn't this typed?*/)
	}, 5000)
}

func (mp mainPage) sendParams(evt jquery.Event) {
	evt.StopPropagation()
	var data shared.ApiPayload
	data.Username = usernameInput.Dom().Val()
	data.Password = passwordInput.Dom().Val()

	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	if err := enc.Encode(data); err != nil {
		slideDown(fmt.Sprintf("json encoding: %v", err), "bg-error")
		return
	}
	payload := buff.String()
	jquery.Ajax(
		map[string]interface{}{
			"data":     payload,
			"dataType": "json",
			"type":     "POST",
			"url":      "/beta/api",
			"cache":    false,
		}).
		Then(func(ignored js.Object) {
		mp.updateCurrent()
	}).
		Fail(func(p js.Object) {
		mp.failSet(p.Get("status").Int(), p.Get("statusText").Str())
	})
}

func (mp mainPage) failSet(statusCode int, statusMsg string) {
	slideDown(fmt.Sprintf("unable to set (%d): %s", statusCode, statusMsg),
		"bg-danger")
}

func (mp mainPage) failGet(statusCode int, statusMsg string) {
	if statusCode == http.StatusNotFound {
		print("no keys found... not updating screen")
		slideDown("Did not find any configuration data.", "bg-warning")
		return
	}
	str := fmt.Sprintf("unexpected error on get (%d): %s", statusCode, statusMsg)
	slideDown(str, "bg-error")
}

func (mp mainPage) succeedGet(payload shared.ApiPayload) {
	mp.passwordCurrent.Set(payload.Password)
	mp.usernameCurrent.Set(payload.Username)
}

func (mp mainPage) succeedPost() {
}

func (mp mainPage) buildPrimary() {
	currentParent.Dom().Append(
		s5.DIV(
			s5.Class(col4),
			s5.Class(offset2),
			s5.DIV(
				s5.Class(row),
				s5.DIV(
					s5.Class(col6),
					s5.Class(gray),
					s5.SPAN(
						s5.Class(bold),
						s5.Text("Username"),
					),
				),
				s5.DIV(
					s5.Class(col6),
					s5.SPAN(
						s5.Class(bold),
						s5.Text("Password"),
					),
				),
			),
			s5.DIV(
				s5.Class(row),
				s5.DIV(
					s5.Class(col6),
					s5.Class(gray),
					s5.SPAN(
						s5.TextEqual(mp.usernameDisplay),
					),
				),
				s5.DIV(
					s5.Class(col6),
					s5.SPAN(
						s5.TextEqual(mp.passwordDisplay),
					),
				),
			),
		).Build(),
	)

	//not using input constraints for a form this simple
	submit.Dom().On(s5.CLICK, func(event jquery.Event) {
		///foo
		mp.sendParams(event)
	})
}

//called when the dom is ready
func (mp mainPage) Start() {

	mp.buildPrimary()
	mp.updateCurrent()
}

func (mp mainPage) updateCurrent() {

	jquery.Ajax(
		map[string]interface{}{
			"dataType": "json",
			"type":     "GET",
			"url":      "/beta/api",
			"cache":    false,
		}).
		Then(func(v js.Object) {
		var payload shared.ApiPayload
		if err := s5.UnpackJson(&payload, v); err != nil {
			slideDown(err.Error(), "bg-danger")
		}
		mp.succeedGet(payload)
	}).
		Fail(func(p js.Object) {
		mp.failGet(p.Get("status").Int(), p.Get("statusText").Str())
	})

}

func showOrSayUnset(values []s5.Equaler) s5.Equaler {
	v := values[0].(s5.StringEqualer).S
	v = strings.TrimSpace(v)
	if v == "" {
		return s5.StringEqualer{S: "not set"}
	}
	return s5.StringEqualer{S: v}
}

func newMainPage() mainPage {
	result := mainPage{
		usernameCurrent: s5.NewStringSimple(""),
		passwordCurrent: s5.NewStringSimple(""),
		usernameDisplay: s5.NewStringSimple(""),
		passwordDisplay: s5.NewStringSimple(""),
	}
	uCons := s5.NewSimpleConstraint(showOrSayUnset, result.usernameCurrent)
	pCons := s5.NewSimpleConstraint(showOrSayUnset, result.passwordCurrent)
	result.usernameDisplay.Attach(uCons)
	result.passwordDisplay.Attach(pCons)

	return result
}

func main() {
	s5.Main(newMainPage())
}
