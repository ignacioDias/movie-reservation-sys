import { isUserLoggedIn, logout } from "./utils.js";
const $div = document.querySelector(".profile");
const $currentMoviesBtn = document.querySelector('.current-movies-btn');
const $soonMoviesBtn = document.querySelector('.soon-movies-btn');
const $signUpBtn = document.querySelector('.signup-btn');
const $loginBtn = document.querySelector('.login-btn');
const $profileBtn = document.querySelector('.profile-btn');
const $logoutBtn = document.querySelector('.logout-btn');
const $loggedDivBtns = document.querySelector('.logged')
const $notLoggedDivBtns = document.querySelector('.not-logged')
const $goBackBtn = document.querySelector('.go-back-btn')

window.addEventListener("load", async () => {
    if (!await isUserLoggedIn()) {
        window.location.href = "/";
        return;
    }
    $loggedDivBtns.style.display = "block";
    $notLoggedDivBtns.style.display = "none";
    try {
        const response = await fetch("/api/v1/users/me");
        if (!response.ok) {
            throw new Error(`${response.status} getting profile`);
        }
        const result = await response.json();
        renderProfile(result);
   } catch (error) {
        console.error(error);
    }
});

function renderProfile(result) {
    const $card = document.createElement("section");
    $card.classList.add("profile-card");

    const $profilePicture = document.createElement("img");
    $profilePicture.src = result.profilePicture;
    $profilePicture.alt = "Profile avatar";
    $profilePicture.classList.add("profile-avatar");
    $card.appendChild($profilePicture);

    const $documentNumber = document.createElement("p");
    $documentNumber.classList.add("profile-line");
    $documentNumber.textContent = `Document: ${result.documentNumber}`;
    $card.appendChild($documentNumber);

    const $email = document.createElement("p");
    $email.classList.add("profile-line");
    $email.textContent = `Email: ${result.email}`;
    $card.appendChild($email);

    const $deleteDivSection = document.createElement("div");
    $deleteDivSection.classList.add("danger-zone");

    const $deleteP = document.createElement("p");
    $deleteP.textContent = "Do you want to delete your account?";
    $deleteDivSection.appendChild($deleteP);

    const $deleteButton = document.createElement("button");
    $deleteButton.textContent = "Delete";
    $deleteButton.classList.add("delete-btn");
    $deleteDivSection.appendChild($deleteButton);
    $card.appendChild($deleteDivSection);

    const $confirmDiv = document.createElement("div");
    $confirmDiv.classList.add("confirm-actions", "hidden");
    
    const $confirmText = document.createElement("p");
    $confirmText.textContent = "Are you sure?";
    $confirmDiv.appendChild($confirmText);

    const $acceptButton = document.createElement("button");
    $acceptButton.textContent = "Yes";
    $acceptButton.classList.add("accept-btn");
    $acceptButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/api/v1/users/me", { method: "DELETE" });
            if (!response.ok) {
                throw new Error(`Error deleting: ${response.status}`);
            }
            window.location.href = "/";
        } catch (error) {
            console.error(error);
        }
    });
    $confirmDiv.appendChild($acceptButton);

    const $revertButton = document.createElement("button");
    $revertButton.textContent = "Cancel";
    $revertButton.classList.add("cancel-btn");
    $revertButton.addEventListener("click", () => {
        $confirmDiv.classList.add("hidden");
    });
    $confirmDiv.appendChild($revertButton);

    $card.appendChild($confirmDiv);
    $div.appendChild($card);

    $deleteButton.addEventListener("click", () => {
        $confirmDiv.classList.remove("hidden");
    });
}

$currentMoviesBtn.addEventListener("click", () => {
    window.location.href = "/movies/projecting";
})
$soonMoviesBtn.addEventListener("click", () => {
    window.location.href = "/movies/soon";
})
$signUpBtn.addEventListener("click", () => {
    window.location.href = "/register";
})
$loginBtn.addEventListener("click", () => {
    window.location.href = "/login";
})
$logoutBtn.addEventListener("click", async () => {
    await logout();
})
$profileBtn.addEventListener("click", () => {
    window.location.href = "/me";
})

$goBackBtn.addEventListener("click", () => {
    window.location.href = "/";
})