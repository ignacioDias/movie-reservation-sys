import { renderComingSoonMovies, isUserLoggedIn } from "./utils.js";

const $currentMoviesBtn = document.querySelector('.current-movies-btn');
const $signUpBtn = document.querySelector('.signup-btn');
const $loginBtn = document.querySelector('.login-btn');
const $profileBtn = document.querySelector('.profile-btn');
const $logoutBtn = document.querySelector('.logout-btn');
const $loggedDivBtns = document.querySelector('.logged')
const $notLoggedDivBtns = document.querySelector('.not-logged')
const $goBackBtn = document.querySelector('.go-back-btn')

const $futureMoviesDiv = document.querySelector(".movies")

window.addEventListener("load", async () => {
    if (await isUserLoggedIn()) {
        $loggedDivBtns.style.display = "block";
        $notLoggedDivBtns.style.display = "none";
    } 
    try {
        const response = await fetch("/api/v1/movies/soon")
        if(!response.ok) {
            throw new Error(`error response: ${response.status}`)
        }
        const result = await response.json()
        renderComingSoonMovies($futureMoviesDiv, result)
    } catch(error) {
        console.error(error)
    }
})

$currentMoviesBtn.addEventListener("click", () => {
    window.location.href = "/movies/projecting";
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