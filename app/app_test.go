package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler())
	defer ts.Close()
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}}) // 프론트에서 요청하는 역할(addTodoHandler가 Form으로 받기 때문에 go에서 지원하는 PostForm으로 요청해야한다.(나중에 JSON으로 보내고 받을 때 변경))
	// ts.URL = local:3000
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode) // todos경로에 POST로 보내는 것은 addTodoHandler인데 응답할 때 StatusCreated를 보낸다. 따라서 resp.StatusCode가 StatusCreated여야한다.
}
