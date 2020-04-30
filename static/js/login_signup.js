const info_btn = document.getElementsByClassName("info-btn")
for (let i = 0; i < info_btn.length; i++) {
    info_btn[i].onclick = () => {
        document.querySelector(".container").classList.toggle("log-in");
    };
}

// sing-in\up
let form = document.querySelectorAll('.form');
let action = "/";


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
                if (res.type == "up") {
                    window.location.replace("/profile/settings");
                } else {
                    window.location.replace("/profile");
                }
            } else {
                alert(res.err);
            }
        })
        .catch(err => console.log(err));
}