'use strict'


let form = document.querySelector('.form');
let action = "";

if (form != null) {
    form.onsubmit = (e) => {
        formF(e);
    };
}

function formF(e) {
    e.preventDefault();
    let formData = new FormData(form[index]);
    action = form[index].getAttribute("action");
    fetchingFile(formData);
}

function fetchingFile(formData) {
    fetch(`${action}`, {
            method: 'POST',
            body: formData
        })
        .then(res => res.json())
        .then(res => {
            if (res.err == "") {
                alert(res.msg);
            } else {
                alert(res.err);
            }
        })
        .catch(err => console.log(err));
}