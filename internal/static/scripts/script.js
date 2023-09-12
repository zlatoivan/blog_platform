function showLikes(articleId) {
    let data = {
        ArticleId: articleId,
    };
    fetch("/showLikes", {
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            let btn = document.getElementById('articleIdButton=' + articleId.toString())

            console.log(result["IsLiked"])
            if (result["IsLiked"]) {
                btn.classList.add("active")
            } else {
                btn.classList.remove("active")
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}

function likeOnClick(articleId) {
    let data = {
        ArticleId: articleId,
    };
    fetch("/someoneIsLiked", {
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            let label = document.getElementById('articleIdLabel=' + articleId.toString())
            label.textContent = result["LikesCnt"]

            let btn = document.getElementById('articleIdButton=' + articleId.toString())
            if (result["IsLiked"]) {
                btn.classList.add("active")
            } else {
                btn.classList.remove("active")
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}