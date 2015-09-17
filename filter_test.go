// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beegae

import (
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beegae/context"

	"google.golang.org/appengine/aetest"
)

var FilterUser = func(ctx *context.Context) {
	ctx.Output.Body([]byte("i am " + ctx.Input.Params[":last"] + ctx.Input.Params[":first"]))
}

func TestFilter(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	r, _ := inst.NewRequest("GET", "/person/asta/Xie", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/person/:last/:first", BeforeRouter, FilterUser)
	handler.Add("/person/:last/:first", &TestController{})
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am astaXie" {
		t.Errorf("user define func can't run")
	}
}

var FilterAdminUser = func(ctx *context.Context) {
	ctx.Output.Body([]byte("i am admin"))
}

// Filter pattern /admin/:all
// all url like    /admin/    /admin/xie    will all get filter

func TestPatternTwo(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	r, _ := inst.NewRequest("GET", "/admin/", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/admin/?:all", BeforeRouter, FilterAdminUser)
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am admin" {
		t.Errorf("filter /admin/ can't run")
	}
}

func TestPatternThree(t *testing.T) {
	inst, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer inst.Close()

	r, _ := inst.NewRequest("GET", "/admin/astaxie", nil)
	w := httptest.NewRecorder()
	handler := NewControllerRegister()
	handler.InsertFilter("/admin/:all", BeforeRouter, FilterAdminUser)
	handler.ServeHTTP(w, r)
	if w.Body.String() != "i am admin" {
		t.Errorf("filter /admin/astaxie can't run")
	}
}
