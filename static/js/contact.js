'use strict'

// sing-in\up
let form = document.querySelector('.form');
let action = "/contact";


if (form != null) {
    form.onsubmit = (e) => {
        e.preventDefault();
        let formData = new FormData(form);
        const data = new URLSearchParams();

        for (const pair of formData) {
            data.append(pair[0], pair[1])
        }
        fetching(data)
    };
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
            console.log(res);
            if (res.msg == "Sended") {
                alert("Sended!");
            } else {
                alert(res.err);
            }
        })
        .catch(err => console.log(err));
}