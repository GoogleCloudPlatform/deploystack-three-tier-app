/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

var basepath = "//127.0.0.1:9000/api/v1/todo";

document.addEventListener('DOMContentLoaded', function(){
    listTodos();
});


function listTodos() {
    var xmlhttp = new XMLHttpRequest();

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
           if (xmlhttp.status == 200) {
               renderListTodos(xmlhttp.response);
           }
           else if (xmlhttp.status == 400) {
              alert('There was an error 400');
           }
           else {
               alert('something else other than 200 was returned');
           }
        }
    };

    xmlhttp.open("GET", basepath, true);
    xmlhttp.send();
}

function renderListTodos(resp){
    let todos = JSON.parse(resp);
    let content = document.querySelector(".content");
    content.innerHTML = "";

    let ul = document.createElement("ul");
    ul.classList.add("list")

    todos.forEach(todo => {
        let li = document.createElement("li");
        let el = renderTodo(todo);
        li.appendChild(el)
        ul.appendChild(li);
    });

    let li = document.createElement("li");
    let el = renderNewTodo()
    li.appendChild(el)
    ul.appendChild(li);


    content.appendChild(ul);

}

function renderNewTodo(){
    let div = document.createElement("div");
    div.classList.add("todo");

    let input = document.createElement("input");
    input.type = "checkbox";
    input.id = `todo-new-cb`;
    input.disabled = true;

    let editor = document.createElement("div");
    editor.classList.add("editor");
    editor.classList.add("editor-new");
    editor.contentEditable = true;
    editor.dataset.placeholder = "Type something here to add a new task. "
    editor.id = `todo-new`;
    editor.addEventListener("blur", createHandler);
    editor.addEventListener("keypress", catchEnter);
    editor.addEventListener("click", function(e){e.target.focus();e.target.innerHTML = "   "});


    let h1 = document.createElement("h1");
    h1.appendChild(input);
    h1.appendChild(editor);

    div.appendChild(h1);


    return div;

}


function renderTodo(todo){
    let div = document.createElement("div");
    div.classList.add("todo");
    if (todo.complete){
        div.classList.add("complete");
    }

    let input = document.createElement("input");
    input.type = "checkbox";
    input.id = `todo-${todo.id}-cb`;
    input.checked = todo.complete;
    input.addEventListener("change", checkHandler);

    let editor = document.createElement("span");
    editor.classList.add("editor");
    editor.contentEditable = true;
    editor.innerHTML = todo.title;
    editor.id = `todo-${todo.id}`;
    editor.addEventListener("blur", blurHandler);

    let icon = document.createElement("span");
    icon.classList.add("material-icons", "delete");
    icon.innerHTML = "delete";
    icon.id = `todo-${todo.id}-delete`;
    icon.addEventListener("click", deleteHandler);

    let h1 = document.createElement("h1");
    h1.appendChild(input);
    h1.appendChild(editor);
    h1.appendChild(icon);

    div.appendChild(h1);


    return div;

}

function blurHandler(e){
    let complete = e.target.parentElement.childNodes[0].checked;
    let title = e.target.innerHTML;
    let id = e.target.id.split("-")[1];
    updateTodo(id, title, complete);
}

function createHandler(e){
    let title = e.target.innerHTML;

    if (title.trim().length == 0){
        e.target.innerHTML= "";
        return
    }

    createTodo(title);
}

function catchEnter(e){
    if (e.key === "Enter") {
        e.preventDefault();
        e.target.blur();
      }
}


function checkHandler(e){
    let complete = e.target.checked;
    let title = e.target.parentElement.childNodes[1].innerHTML;
    let id = e.target.parentElement.childNodes[1].id.split("-")[1];


    if (complete){
        e.target.parentElement.parentElement.classList.add("complete");
    }  else{
        e.target.parentElement.parentElement.classList.remove("complete");
    }

    updateTodo(id, title, complete);
}

function deleteHandler(e){
    let id = e.target.id.split("-")[1];
    deleteTodo(id);
}

function updateTodo(id, title, complete){
    var xmlhttp = new XMLHttpRequest();
    let form  = new FormData();
    form.append("title", title);
    form.append("complete", complete);

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
           if (xmlhttp.status == 200) {
                listTodos();
           }
           else if (xmlhttp.status == 400) {
              alert('There was an error 400');
           }
           else {
               alert('something else other than 200 was returned');
           }
        }
    };

    xmlhttp.open("POST", basepath+"/"+ id, true);
    xmlhttp.send(form);
}

function createTodo(title){
    var xmlhttp = new XMLHttpRequest();
    let form  = new FormData();
    form.append("title", title);

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
           if (xmlhttp.status == 201) {
                listTodos();
           }
           else if (xmlhttp.status == 400) {
              alert('There was an error 400');
           }
           else {
               alert('something else other than 201 was returned');
               console.log(xmlhttp.status);
           }
        }
    };

    xmlhttp.open("POST", basepath, true);
    xmlhttp.send(form);
}

function deleteTodo(id){
    var xmlhttp = new XMLHttpRequest();

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
           if (xmlhttp.status == 204) {
            listTodos();
           }
           else if (xmlhttp.status == 400) {
              alert('There was an error 400');
           }
           else {
               alert('something else other than 204 was returned');
               console.log(xmlhttp.status);
           }
        }
    };

    xmlhttp.open("DELETE", basepath+"/"+ id, true);
    xmlhttp.send();
}
