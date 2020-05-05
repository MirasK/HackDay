'use strict'


const tabBtn = document.querySelectorAll(".tab");
const tab = document.querySelectorAll(".tabShow");

let form = document.querySelectorAll('.form');
let action = "";

function tabs(panelIndex) {
    tab.forEach(function(node) {
        node.style.display = "none";
    });
    tab[panelIndex].style.display = "block";
}

let doctab = document.querySelector('.tab').innerHTML;
tabs(doctab - 1);

// $(".tab").click(function() {
//     $(this).addClass("active").siblings().removeClass("active");
// })

if (form != null) {
    form[0].onsubmit = (e) => {
        formF(0, e);
    };
    form[1].onsubmit = (e) => {
        formF(1, e);
    };
    form[2].onsubmit = (e) => {
        formF(2, e);
    };
    form[3].onsubmit = (e) => {
        formF(3, e);
    };
    form[4].onsubmit = (e) => {
        formF(4, e);
    };
    form[5].onsubmit = (e) => {
        formF(5, e);
    };
}

function formF(index, e) {
    e.preventDefault();
    if (index == 0) {
        let err = checkPhone()
        if (err != "") {
            alert(err);
            return
        }
    }

    let formData = new FormData(form[index]);
    const data = new URLSearchParams();
    action = form[index].getAttribute("action");

    for (const pair of formData) {
        data.append(pair[0], pair[1])
    }
    if (index == 3 || index == 4) {
        fetchingFile(formData)
    } else {
        fetching(data);
    }

}

function checkPhone() {
    let phone = document.querySelector('.phone').value;
    if (phone.match(/^8[0-9]{10}$/) == null && phone.match(/^\+7[0-9]{10}$/) == null) {
        return "wrong number";
    }
    return "";
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
            } else {
                alert(res.err);
            }
        })
        .catch(err => console.log(err));
}