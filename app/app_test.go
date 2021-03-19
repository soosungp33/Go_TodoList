package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/soosungp33/Go_TodoList/model"
	"github.com/stretchr/testify/assert"
)

func TestTodo(t *testing.T) {
	os.Remove("./test.db") // 테스트하기 전에 테스트파일을 없애준다.
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close() // db를 종료해줘야함

	ts := httptest.NewServer(ah) // 테스트 서버 설정
	defer ts.Close()

	var todo model.Todo

	// add(POST)로 리스트 2개를 추가하는 코드
	// ts.URL = local:3000
	// Values = {name:item}
	// 요청한 응답이 resp에 저장된다.
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}}) // 프론트에서 add를 요청하는 역할(addTodoHandler가 Form으로 받기 때문에 go에서 지원하는 PostForm으로 요청해야한다.)
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode) // todos경로에 POST로 요청하는 것은 addTodoHandler인데 응답할 때 StatusCreated를 보낸다. 따라서 resp.StatusCode가 StatusCreated여야한다.
	err = json.NewDecoder(resp.Body).Decode(&todo)    // resp.body에 있는 JSON 데이터를 해독해서 todo객체에 저장한다.
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo") // 값이 맞는지 비교
	id1 := todo.ID
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo) // resp.body에 있는 JSON 데이터를 해독해서 todo객체에 저장한다.
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo2") // 값이 맞는지 비교
	id2 := todo.ID

	// 리스트를 받아오는 코드(GET)
	resp, err = http.Get(ts.URL + "/todos") // todos경로에 GET으로 요청하는 것은 getTodoListHandler 이다.
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2) // 위에서 2개 add했으므로 todos의 크기는 2여야 한다.
	for _, t := range todos {   // 받아온 Name이 각각 맞는지 검사
		if t.ID == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	// 리스트의 토글박스 부분
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true") // id1의 토글박스가 true가 되어 있음
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos") // GET으로 리스트를 가져와서
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos) // todos에 담고
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.ID == id1 { // id1일 때
			assert.True(t.Completed) // complete=true로 요청했으니까 t.Completed도 true여야 한다.
		}
	}

	// 리스트를 DELETE 하기(go의 http는 DELETE를 지원하지 않음)
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil) // DELETE request를 만들어줘야함(이 경로는 removeTodoHandler)
	resp, err = http.DefaultClient.Do(req)                                       // req를 요청함(req는 DELETE)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	resp, err = http.Get(ts.URL + "/todos") // GET으로 리스트를 가져와서
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos) // todos에 담고
	assert.NoError(err)
	assert.Equal(len(todos), 1) // 지웠으니까 사이즈가 1이어야 한다.
	for _, t := range todos {
		assert.Equal(t.ID, id2) // id1이 지워졌으니까 나오는 게 id2여야 한다.
	}
}
