const info_btn = document.getElementsByClassName("info-btn")
for (let i = 0; i < info_btn.length; i++) {
    info_btn[i].onclick = () => {
        document.querySelector(".container").classList.toggle("log-in");
    };
}

// sing-in\up
let form = document.querySelectorAll('.form');
let action = "";


if (form != null) {
    form[0].onsubmit = (e) => {
        formF(0, e);
    };
    form[1].onsubmit = (e) => {
        formF(1, e);
    };
}

function formF(index, e) {
    e.preventDefault();
    let formData = new FormData(form[index]);
    const data = new URLSearchParams();
    action = form[index].getAttribute("action");

    for (const pair of formData) {
        data.append(pair[0], pair[1])
    }
    fetching(data)
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
            if (res.msg == "redirect") {
                if (res.type == "in" || res.type == "restore") {
                    window.location.replace("/profile");
                } else if (res.type == "forgot") {
                    window.location.replace("/verification");
                } else if (res.type == "verification") {
                    window.location.replace("/restore");
                }
            } else if (res.type == "up" && res.err == "") {
                alert(res.msg);
            } else {
                alert(res.err);
            }
        })
        .catch(err => console.log(err));
}