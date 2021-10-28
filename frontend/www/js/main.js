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

    let ul = document.createElement("ul");
    ul.classList.add("list")

    todos.forEach(todo => {
        let li = document.createElement("li");
        let el = renderTodo(todo);
        li.appendChild(el)
        ul.appendChild(li);
    });
    content.appendChild(ul);

    console.log(todos)
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
    console.log("delete");
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
               console.log("success")
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

function deleteTodo(id){
    var xmlhttp = new XMLHttpRequest();

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
           if (xmlhttp.status == 204) {
               console.log("success")
           }
           else if (xmlhttp.status == 400) {
              alert('There was an error 400');
           }
           else {
               alert('something else other than 204 was returned');
           }
        }
    };

    xmlhttp.open("DELETE", basepath+"/"+ id, true);
    xmlhttp.send();
}

// TODO: Add Create new todo interface
// TODO: Add Delete todo interface  