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
    let checked = "";
    if (todo.complete){
        div.classList.add("complete");
        checked = "checked"
    }


    div.innerHTML = `
        <h1 contenteditable="true">
            <input type="checkbox" ${checked} id="todo-${todo.id}" name="todo-${todo.id}" value="">
            ${todo.title}</h1>
    `;
    return div;

}