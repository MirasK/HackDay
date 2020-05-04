'use strict'

let add = document.querySelector(".add-require");
let form = document.querySelectorAll('.form');
let action = "";
let data = new URLSearchParams();

// request btns
let accs = document.querySelectorAll('.accept span');
let refs = document.querySelectorAll('.refuse span');
let pop = document.querySelector('.employer-response');
// end

if (accs.length != 0) {
    accs.forEach(v => v.onclick = accept);
    refs.forEach(v => v.onclick = refuse);
    pop.onclick = (e) => {
        if (e.target.classList.length == 2) {
            pop.classList.toggle("open");
        }
    }
}

if (form != null) {
    form[0].onsubmit = (e) => {
        formF(0, e);
    };
    if (form[1] != null) {
        form[1].onsubmit = (e) => {
            formF(1, e);
        };
    }

}

add.onclick = () => {
    let newIn = document.createElement('input');
    newIn.classList.add('requirements');
    newIn.setAttribute("type", "text");
    newIn.setAttribute("required", "");
    newIn.setAttribute("placeholder", "C#");
    newIn.setAttribute("name", "requirements");
    add.before(newIn);
}

function reqArr() {
    let requires = document.querySelectorAll('.requirements');
    let arr = [];
    requires.forEach(v => {
        arr.push(v.value);
    });
    return arr
}

function checkPhone() {
    let phone = document.querySelector('.form-phone > input').value;
    if (phone.match(/^8[0-9]{10}$/) == null && phone.match(/^\+7[0-9]{10}$/) == null) {
        return "wrong number";
    }
    return "";
}

function doAccOrRef(e) {
    data = new URLSearchParams();
    let id;
    if (e.target.getAttribute('data-id') == null) {
        id = e.target.parentNode.getAttribute('data-id');
    } else {
        id = e.target.getAttribute('data-id');
    }
    data.append("id", id);
    data.append("type", "to student");
    form[0].setAttribute("action", `/works/${id}`);
    pop.classList.toggle("open");
    return data
}

function accept(e) {
    doAccOrRef(e);
    data.append("status", "true");
}

function refuse(e) {
    doAccOrRef(e);
    data.append("status", "false");
}

function formF(index, e) {
    action = form[index].getAttribute("action");
    e.preventDefault();

    if (action == "/create-work") {
        let err = checkPhone();
        if (err !== "") {
            alert(err);
            return
        }
    }

    let formData = new FormData(form[index]);


    if (action == "/create-work") {
        formData.set('requirements', reqArr());
    }

    for (let pair of formData) {
        data.append(pair[0], pair[1]);
    }

    fetching(data);
}

function fetching(data) {
    fetch(`${action}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: data
        })
        .then(res => res.json())
        .then(res => {
            if (res.err == "") {
                alert(res.msg);
                if (res.type == "to employer") {
                    window.location.reload();
                }
            } else {
                alert(res.err);
            }
        })
        .catch(err => console.log(err));
}