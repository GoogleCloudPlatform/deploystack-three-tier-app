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

    let h1 = document.createElement("h1");
    h1.contentEditable = true;
    h1.id = `todo-${todo.id}`;
    h1.appendChild(input);
    h1.insertAdjacentHTML('beforeend', todo.title);
    h1.addEventListener("blur", blurHandler);

    div.appendChild(h1);


    return div;

}

function blurHandler(e){
    let complete = e.target.childNodes[0].checked;
    let title = e.target.childNodes[1].data;
    let id = e.target.id.split("-")[1];
    updateTodo(id, title, complete);
}

function checkHandler(e){
    let complete = e.target.checked;
    let title = e.target.parentElement.childNodes[1].data;
    let id = e.target.parentElement.id.split("-")[1];

    if (complete){
        e.target.parentElement.parentElement.classList.add("complete");
    }  else{
        e.target.parentElement.parentElement.classList.remove("complete");
    }

    updateTodo(id, title, complete);
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