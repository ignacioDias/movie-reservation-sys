const $form = document.querySelector(".login-form");
const $registerBtn = document.querySelector(".register-btn")
$form.addEventListener("submit", async function (event) {
    event.preventDefault();
    const formData = new FormData($form)

    const email = formData.get("email");
    const password = formData.get("password");

    const loginReq = {
        "email": email,
        "password": password
    }

    try {
        const url = "/api/v1/auth/login"
        const response = await fetch(url, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(loginReq)
        })
        if(!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }
        window.location.href = "/"
    } catch (error) {
        console.error(error)
    }
})

$registerBtn.addEventListener("click", () => {
    window.location.href = "/register"
})