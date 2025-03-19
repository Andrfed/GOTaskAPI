package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bytes"

	"encoding/json"

	"net/http"

	"io/ioutil"

	"fmt"
)

func getAllTasks() []Task {
	url := "http://localhost:8000/tasks/"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gtasks []Task
	err := json.Unmarshal(body, &gtasks)
	if err != nil {
		panic(err)
	}
	return gtasks
}

func addTask(title string, text string) Task {
	var task Task
	task.Title = title
	task.Text = text
	data, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	url := "http://localhost:8000/tasks/"
	r := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", url, r)
	req.Header.Set("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gtask Task
	err = json.Unmarshal(body, &gtask)
	if err != nil {
		panic(err)
	}
	return gtask
}

func deleteTask(id int) Task {
	url := "http://localhost:8000/tasks/" + fmt.Sprint(id)
	req, _ := http.NewRequest("DELETE", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gtask Task
	err := json.Unmarshal(body, &gtask)
	if err != nil {
		panic(err)
	}
	return gtask
}

func getTask(id int) Task {
	url := "http://localhost:8000/tasks/" + fmt.Sprint(id)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gtask Task
	err := json.Unmarshal(body, &gtask)
	if err != nil {
		panic(err)
	}
	return gtask
}

func updateTask(id int, title string, text string) Task {
	var task Task
	task.Title = title
	task.Text = text
	data, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	url := "http://localhost:8000/tasks/" + fmt.Sprint(id)
	r := bytes.NewBuffer(data)
	req, _ := http.NewRequest("PUT", url, r)
	req.Header.Set("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var gtask Task
	err = json.Unmarshal(body, &gtask)
	if err != nil {
		panic(err)
	}
	return gtask
}
func TestAddTask(t *testing.T) {
	task := Task{}
	id := 0
	task = addTask("Test", "Test Task")
	assert.Equal(t, "Test", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	id = task.Id
	task = deleteTask(id)
	assert.Equal(t, id, task.Id)
}

func TestGetTask(t *testing.T) {
	task := Task{}
	id := 0
	task = addTask("Test", "Test Task")
	assert.Equal(t, "Test", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	id = task.Id
	task = getTask(id)
	assert.Equal(t, "Test", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	task = deleteTask(id)
	assert.Equal(t, id, task.Id)
}

func TestUpdateTask(t *testing.T) {
	task := Task{}
	id := 0
	task = addTask("Test", "Test Task")
	assert.Equal(t, "Test", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	id = task.Id
	task = updateTask(id, "Updated", "Updated Task")
	assert.Equal(t, "Updated", task.Title)
	assert.Equal(t, "Updated Task", task.Text)
	task = getTask(id)
	assert.Equal(t, "Updated", task.Title)
	assert.Equal(t, "Updated Task", task.Text)
	task = deleteTask(id)
	assert.Equal(t, id, task.Id)
}

func TestGetAllTasks(t *testing.T) {
	task := Task{}
	tasks := []Task{}
	var ids [3]int
	task = addTask("Test1", "Test Task")
	assert.Equal(t, "Test1", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	ids[0] = task.Id
	task = addTask("Test2", "Test Task")
	assert.Equal(t, "Test2", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	ids[1] = task.Id
	task = addTask("Test3", "Test Task")
	assert.Equal(t, "Test3", task.Title)
	assert.Equal(t, "Test Task", task.Text)
	ids[2] = task.Id
	tasks = getAllTasks()
	assert.Equal(t, 3, len(tasks))
	for i := 0; i < len(tasks); i++ {
		assert.Equal(t, ids[i], tasks[i].Id)
		assert.Equal(t, fmt.Sprint("Test", i+1), tasks[i].Title)
		assert.Equal(t, "Test Task", tasks[i].Text)
	}
	task = deleteTask(ids[0])
	assert.Equal(t, ids[0], task.Id)
	task = deleteTask(ids[1])
	assert.Equal(t, ids[1], task.Id)
	task = deleteTask(ids[2])
	assert.Equal(t, ids[2], task.Id)
	tasks = getAllTasks()
	assert.Equal(t, 0, len(tasks))
}
