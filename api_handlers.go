// GENERATED CODE
// DO NOT EDIT
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
)

// ProfileParams validator
func (out *ProfileParams) Decode(in map[string][]string) error {
	var val []string
	var exists bool
	val, exists = in["login"]
	if !exists {
		return fmt.Errorf("login must me not empty")
	}
	if exists {
		out.Login = val[0]
	}
	return nil
}

// CreateParams validator
func (out *CreateParams) Decode(in map[string][]string) error {
	var val []string
	var exists bool
	val, exists = in["login"]
	if !exists {
		return fmt.Errorf("login must me not empty")
	}
	if exists {
		out.Login = val[0]
		if len(out.Login) < 10 {
			return fmt.Errorf("login len must be >= 10")
		}
	}
	val, exists = in["full_name"]
	if exists {
		out.Name = val[0]
	}
	val, exists = in["status"]
	if exists {
		out.Status = val[0]
		if out.Status == "" {
			out.Status = "user"
		} else {
			opts := []string{"user", "moderator", "admin"}
			if !slices.Contains(opts, out.Status) {
				return fmt.Errorf("status must be one of [user, moderator, admin]")
			}
		}
	}
	val, exists = in["age"]
	if exists {
		age, err := strconv.Atoi(val[0])
		if err != nil {
			return fmt.Errorf("age must be int")
		}
		out.Age = age
		if out.Age < 0 {
			return fmt.Errorf("age must be >= 0")
		}
		if out.Age > 128 {
			return fmt.Errorf("age must be <= 128")
		}
	}
	return nil
}

// handler for Profile method
func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	ctx := r.Context()
	in := ProfileParams{}

	err := in.Decode(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(jsonStr))
		return
	}

	out := new(User)
	out, err = srv.Profile(ctx, in)
	if err != nil {
		if err, ok := err.(ApiError); ok {
			w.WriteHeader(err.HTTPStatus)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		resp["error"] = err.Error()
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = "internal server error"
		}
		fmt.Fprintf(w, string(jsonStr))
		return
	}
	resp["error"] = ""
	resp["response"] = out
	w.WriteHeader(http.StatusOK)
	jsonStr, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = "internal server error"
	}
	fmt.Fprintf(w, string(jsonStr))
}

// handler for Create method
func (srv *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	ctx := r.Context()

	in := CreateParams{}
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = "internal server error"
		jsonStr, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsonStr))
		return
	}
	req := url.URL{RawQuery: string(bs)}

	err = in.Decode(req.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(jsonStr))
		return
	}

	out := new(NewUser)
	out, err = srv.Create(ctx, in)
	if err != nil {
		if err, ok := err.(ApiError); ok {
			w.WriteHeader(err.HTTPStatus)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		resp["error"] = err.Error()
		jsonStr, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsonStr))
		return
	}

	resp["error"] = ""
	resp["response"] = out
	w.WriteHeader(http.StatusOK)
	jsonStr, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = "internal server error"
	}
	fmt.Fprintf(w, string(jsonStr))
}

// OtherCreateParams validator
func (out *OtherCreateParams) Decode(in map[string][]string) error {
	var val []string
	var exists bool
	val, exists = in["username"]
	if !exists {
		return fmt.Errorf("username must me not empty")
	}
	if exists {
		out.Username = val[0]
		if len(out.Username) < 3 {
			return fmt.Errorf("username len must be >= 3")
		}
	}
	val, exists = in["account_name"]
	if exists {
		out.Name = val[0]
	}
	val, exists = in["class"]
	if exists {
		out.Class = val[0]
		if out.Class == "" {
			out.Class = "warrior"
		} else {
			opts := []string{"warrior", "sorcerer", "rouge"}
			if !slices.Contains(opts, out.Class) {
				return fmt.Errorf("class must be one of [warrior, sorcerer, rouge]")
			}
		}
	}
	val, exists = in["level"]
	if exists {
		level, err := strconv.Atoi(val[0])
		if err != nil {
			return fmt.Errorf("level must be int")
		}
		out.Level = level
		if out.Level < 1 {
			return fmt.Errorf("level must be >= 1")
		}
		if out.Level > 50 {
			return fmt.Errorf("level must be <= 50")
		}
	}
	return nil
}

// handler for Create method
func (srv *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	ctx := r.Context()

	in := OtherCreateParams{}
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = "internal server error"
		jsonStr, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsonStr))
		return
	}
	req := url.URL{RawQuery: string(bs)}

	err = in.Decode(req.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(jsonStr))
		return
	}

	out := new(OtherUser)
	out, err = srv.Create(ctx, in)
	if err != nil {
		if err, ok := err.(ApiError); ok {
			w.WriteHeader(err.HTTPStatus)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		resp["error"] = err.Error()
		jsonStr, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(jsonStr))
		return
	}

	resp["error"] = ""
	resp["response"] = out
	w.WriteHeader(http.StatusOK)
	jsonStr, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = "internal server error"
	}
	fmt.Fprintf(w, string(jsonStr))
}

// multiplexer for MyApi
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	switch r.URL.Path { 
    case "/user/profile":
		if r.Method == http.MethodPost {
			bs, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp["error"] = "internal server error"
				jsonStr, _ := json.Marshal(resp)
				fmt.Fprintf(w, string(jsonStr))
				return
			}
			r.URL.RawQuery = string(bs)
		}
		srv.handlerProfile(w, r)
    case "/user/create":
		if r.Method != http.MethodPost {
			resp["error"] = "bad method"
			w.WriteHeader(http.StatusNotAcceptable)
			jsonStr, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp["error"] = "internal server error"
				jsonStr, _ := json.Marshal(resp)
				fmt.Fprintf(w, string(jsonStr))
				return
			}
			fmt.Fprintf(w, string(jsonStr))
			return
		}
		if r.Header.Get("X-Auth") != "100500" {
			resp["error"] = "unauthorized"
			w.WriteHeader(http.StatusForbidden)
			jsonStr, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp["error"] = "internal server error"
				jsonStr, _ := json.Marshal(resp)
				fmt.Fprintf(w, string(jsonStr))
				return
			}
			fmt.Fprintf(w, string(jsonStr))
			return
		}
		srv.handlerCreate(w, r)
	default:
		resp["error"] = "unknown method"
		w.WriteHeader(http.StatusNotFound)
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = "internal server error"
			jsonStr, _ := json.Marshal(resp)
			fmt.Fprintf(w, string(jsonStr))
			return
		}
		fmt.Fprintf(w, string(jsonStr))
		return
	}
}

// multiplexer for OtherApi
func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	switch r.URL.Path { 
    case "/user/create":
		if r.Method != http.MethodPost {
			resp["error"] = "bad method"
			w.WriteHeader(http.StatusNotAcceptable)
			jsonStr, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp["error"] = "internal server error"
				jsonStr, _ := json.Marshal(resp)
				fmt.Fprintf(w, string(jsonStr))
				return
			}
			fmt.Fprintf(w, string(jsonStr))
			return
		}
		if r.Header.Get("X-Auth") != "100500" {
			resp["error"] = "unauthorized"
			w.WriteHeader(http.StatusForbidden)
			jsonStr, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp["error"] = "internal server error"
				jsonStr, _ := json.Marshal(resp)
				fmt.Fprintf(w, string(jsonStr))
				return
			}
			fmt.Fprintf(w, string(jsonStr))
			return
		}
		srv.handlerCreate(w, r)
	default:
		resp["error"] = "unknown method"
		w.WriteHeader(http.StatusNotFound)
		jsonStr, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = "internal server error"
			jsonStr, _ := json.Marshal(resp)
			fmt.Fprintf(w, string(jsonStr))
			return
		}
		fmt.Fprintf(w, string(jsonStr))
		return
	}
}

